package processor

import (
    "encoding/binary"
    "errors"
    "fmt"
    "io/ioutil"
)

func (c *cpu) error(s string) error {
    s = fmt.Sprintf("CPU error at address %d: %s", c.pc, s)
    return errors.New(s)
}

func (c *cpu) LoadFile(filename string) error {
    b, e := ioutil.ReadFile(filename)
    for i := range b {
        c.rom[i] = b[i]
    }

    return e
}

func (c *cpu) String() string {
    return fmt.Sprintf("d: %s\na: %s\nsr: %s\npc: %v\nnext instruction: %4x",
                        formatRegisters(c.d), formatRegisters(c.a), formatBytesBin(c.sr[:]), c.pc, c.rom[c.pc:c.pc+2])
}

func formatBytesHex(b []byte) string {
    s := "["
    for i := range b {
        s += fmt.Sprintf("%02x ", b[i])
    }
    return s[:len(s)-1] + "]"
}

func formatBytesBin(b []byte) string {
    s := "["
    for i := range b {
        s += fmt.Sprintf("%08b ", b[i])
    }
    return s[:len(s)-1] + "]"
}

func formatRegisters(r [8][4]byte) string {
    s := "["
    for i := 0; i < 8; i++ {
        s += fmt.Sprintf("%d", i) + formatBytesHex(r[i][:]) + " "
    }
    return s[:len(s)-1] + "]"
}

func (c *cpu) PrintRom(n int) string {
    s := fmt.Sprintf("Reading %d bytes from ROM:\n", n)
    for i := 0; i < n; i++ {
        s += fmt.Sprintf("%2x ", c.rom[i])
        if i % 16 == 15{
            s += "\n"
        }
    }
    return s + "\n"
}

// Returns byte=1, word=2, long=4
func opsize(mode int, b1, b2 bool) int {
    switch mode {
    case 0:
        if !b1 {
            if b2 {
                return 2
            } else {
                return 1
            }
        } else if !b2 {
            return 4
        }
    case 1:
        if b1 {
            return 4
        } else {
            return 2
        }
    case 2:
        if b1 {
            if b2 {
                return 2
            } else {
                return 4
            }
        } else if b2 {
            return 1
        }
    }

    return -1
}

/*
    Mode 0 = data register
    Mode 1 = address register
    Mode 2 = address
    Mode 3 = address with postincrement
    Mode 4 = address with predecrement
    Mode 5 = address with displacement
    Mode 6 = address with index
    Mode 7 = PC with displacement
    Mode 8 = PC with index
    Mode 9 = Absolute short
    Mode 10 = Absolute long
    Mode 11 = Immediate
*/
func addressingmode(m1, m2, m3, x1, x2, x3 bool) (mode, register int) {
    register = 0
    mode = -1
    if m1 && m2 && m3 {
        if x1 && !x2 && !x3 {
            mode = 11
        } else if !x1 {
            if x2 {
                if x3 {
                    mode = 8
                } else {
                    mode = 7
                }
            } else {
                if x3 {
                    mode = 10
                } else {
                    mode = 9
                }
            }
        }
    } else {
        mode = bits3ToInt(m1, m2, m3)
        register = bits3ToInt(x1, x2, x3)
    }
    return
}

/* Converts 3 bits to int
   Used for getting a register index
*/
func bits3ToInt(b1, b2, b3 bool) int {
    register := 0
    if b1 {
        register += 4
    }
    if b2 {
        register += 2
    }
    if b3 {
        register += 1
    }
    return register
}

/* Converts 4 bits to int
   Used for getting a condition mode
    Condition 0: T True
    Condition 1: F False
    Condition 2: HI Higher
    Condition 3: LS Lower or Same
    Condition 4: CC Carry Clear
    Condition 5: CS Carry Set
    Condition 6: NE Not Equal
    Condition 7: EQ Equal
    Condition 8: VC Overflow Clear
    Condition 9: VS Overflow Set
    Condition 10: PL Plus
    Condition 11: MI Minus
    Condition 12: GE Greater or Equal
    Condition 13: LT Less Than
    Condition 14: GT Greater Than
    Condition 15: LE Less or Equal
*/
func bits4ToInt(c1, c2, c3, c4 bool) int {
    condition := 0
    if c1 {
        condition += 8
    }
    if c2 {
        condition += 4
    }
    if c3 {
        condition += 2
    }
    if c4 {
        condition += 1
    }
    return condition
}

func readBytes(b []byte, n int) interface{} {
    switch n {
        case 1:
            return uint8(b[0])
        case 2:
            return binary.BigEndian.Uint16(b[:n])
        case 4:
            return binary.BigEndian.Uint32(b[:n])
    }
    return nil
}

func isbitset(i, mask byte) bool {
    return (i & mask) != 0
}

func parseByte(b byte) []bool {
    bits := make([]bool, 8)
    for i := range bits {
        bits[i] = (b & uint8(1 << uint(i))) != 0
    }
    return bits
}

func addByte(a, b byte, x bool) (result byte, overflow, carry bool) {
    result = a+b
    if x {
        result++
    }
    signA := isbitset(a, bit7)
    signB := isbitset(b, bit7)
    signResult := isbitset(result, bit7)
    overflow =  (signA && signB && !signResult) || (!signA && !signB && signResult)
    carry = (result < a || result < b)

    return
}

func addTo(d, s []byte, n int, x bool) (overflow, carry bool) {
    ld := len(d)
    ls := len(s)
    d[ld-1], overflow, carry = addByte(d[ld-1], s[ls-1], x)
    for i := 0; i < n-1; i++ {
        d[ld-2-i], overflow, carry = addByte(d[ld-2-i], s[ls-2-i], carry)
    }
    // test behavior on simulator. What happens if carry bit creates overflow in out-of-scope byte
    if carry && n < ld {
        d[ld-1-n], _, carry = addByte(d[ld-1-n], 0, true)
    }
    return
}

func increment(b []byte, n int) {
    switch len(b) {
    case 1:
        b[0] += byte(n)
    case 2:
        tmp := binary.Uint16(b)
        tmp += uint16(n)
        binary.PutUint16(b, tmp)
    case 4:
        tmp := binary.Uint32(b)
        tmp += uint32(n)
        binary.PutUint32(b, tmp)
    }
}

func isZero(b []byte) {
    combined := byte(0)
    for _, v := range b {
        combined |= v
    }
    return combined == 0
}

func isNegative(b []byte) {
    return isbitset(b[0], bit7)
}

func bytesUsedByAddressing(mode, size int) int {
    if (mode >= 5 && mode <= 9) {
        return 2
    } else if mode == 10 {
        return 4
    } else if mode == 11 {
        if size == 4 {
            return 4
        } else {
            return 2
        }
    } else {
        return 0
    }
}

func signExtend2to4(x uint16) uint32 {
    if isbitset(x, wbit15) {
        return uint32(x) + uint32(0xffff0000)
    } else {
        return uint32(x)
    }
}

func signExtend1to2(x uint8) uint16 {
    if isbitset(x, bit7) {
        return uint16(x) + uint16(0xff00)
    } else {
        return uint16(x)
    }
}

func signExtend1to4(x uint8) uint32 {
    if isbitset(x, bit7) {
        return uint32(x) + uint32(0xffffff00)
    } else {
        return uint32(x)
    }
}

func parse8bitDisplacement(b byte) (data bool, register int, word bool) {
    data = !isbitset(b, bit7)
    word = !isbitset(b, bit3)
    register = bits3ToInt(isbitset(b, bit6), isbitset(b, bit5), isbitset(b, bit4))
    return
}

func loadByAddressing(mode, register, size, extraBytes int) []byte {
    tmp := make([]byte, size)
    switch mode {
    case 0:
        for i := 0; i < size; i++ {
            tmp[i] = c.d[register][4-size+i]
        }
    case 1:
        for i := 0; i < size; i++ {
            tmp[i] = c.a[register][4-size+i]
        }
    case 2:
        address := readBytes(c.a[register][:], 4)
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 3:
        address := readBytes(c.a[register][:], 4)
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
        increment(c.a[register][:], size)
    case 4:
        increment(c.a[register][:], -size)
        address := readBytes(c.a[register][:], 4)
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 5:
        address := readBytes(c.a[register][:], 4)
        address += binary.Uint16(c.rom[c.pc+extraBytes+2:c.pc+extraBytes+4])
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 6:
        data, reg, word := parse8bitDisplacement(c.rom[c.pc+extraBytes+2])
        address := readBytes(c.a[register][:], 4)
        if data {
            if word {
                address += signExtend2to4(readBytes(c.d[reg][2:4], 2))
            } else {
                address += readBytes(c.d[reg][:], 4)
            }
        } else {
            if word {
                address += signExtend2to4(readBytes(c.a[reg][2:4], 2))
            } else {
                address += readBytes(c.a[reg][:], 4)
            }
        }
        address += int(c.rom[c.pc+extraBytes+3])
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 7:
        inc := readBytes(c.rom[c.pc+extraBytes:c.pc+extraBytes+2], 2)
        address := c.pc + int(signExtend2to4(inc))
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 8:
        data, reg, word := parse8bitDisplacement(c.rom[c.pc+extraBytes+2])
        address := c.pc
        if data {
            if word {
                address += signExtend2to4(readBytes(c.d[reg][2:4], 2))
            } else {
                address += readBytes(c.d[reg][:], 4)
            }
        } else {
            if word {
                address += signExtend2to4(readBytes(c.a[reg][2:4], 2))
            } else {
                address += readBytes(c.a[reg][:], 4)
            }
        }
        address += int(c.rom[c.pc+extraBytes+3])
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 9:
        address := int(signExtend2to4(readBytes(c.rom[c.pc+extraBytes:c.pc+extraBytes+2], 2)))
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 10:
        address := int(readBytes(c.rom[c.pc+extraBytes:c.pc+extraBytes+4], 4))
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 11:
        if size == 1 {
            tmp[0] = c.rom[c.pc+extraBytes+3]
        } else {
            for i := 0; i < size; i++ {
                tmp[i] = c.rom[c.pc+extraBytes+2+i]
            }
        }
    }
    return tmp
}
