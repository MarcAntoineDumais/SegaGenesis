package processor

import (
    
)

// Supposed to be running at 8MHz
type cpu struct {
    d [8][4]byte        // data registers d0, d1, ..., d7
    a [8][4]byte        // adress registers a0, a1, ..., a7
    sr [2]byte          // status register TT, S, M, 0, III, 0, 0, 0, X, N, Z, V, C
    rom [1048576]byte   // 1MB ROM area
    ram [65536]byte     // 64KB RAM area
    pc int           // program counter
}

func Create() (*cpu) {
    c := new(cpu)
    return c
}

func (c *cpu) Run() error {
    var err error
    b := false
    for !b {
        b, err = c.Step()
    }
    return err
}

func (c *cpu) Step() (b bool, err error) {
    b = true
    i1 := parseByte(c.rom[c.pc])
    i2 := parseByte(c.rom[c.pc+1])
    if i1[7] {
        if i1[6] {
            if i1[5] {
                return
            } else {
                if i1[4] { //ADD instruction
                    register := bits3ToInt(i1[3], i1[2], i1[1])
                    if i2[7] && i2[6] { //ADDA
                        //size := opsize(1, i1[0], false)
                        //TODO
                    } else {
                        size := opsize(0, i2[7], i2[6])
                        if i1[0] { // Direction to <ea>
                            if !i2[5] && !i2[4] {
                                //ADDX
                                //TODO
                            } else {
                                //ADD
                                dm, dr := addressingmode(i2[5], i2[4], i2[3],
                                                         i2[2], i2[1], i2[0])
                                var overflow, carry bool
                                var result []byte
                                switch dm {
                                case 2:
                                    address := int(read4Bytes(c.a[dr][:]))
                                    overflow, carry = addTo(c.ram[address:address+size], c.d[register][:], size, false)
                                    result = c.ram[address:address+size]
                                case 3:
                                    address := int(read4Bytes(c.a[dr][:]))
                                    overflow, carry = addTo(c.ram[address:address+size], c.d[register][:], size, false)
                                    result = c.ram[address:address+size]
                                    increment(c.a[dr][:], size)
                                case 4:
                                    increment(c.a[dr][:], -size)
                                    address := int(read4Bytes(c.a[dr][:]))
                                    overflow, carry = addTo(c.ram[address:address+size], c.d[register][:], size, false)
                                    result = c.ram[address:address+size]
                                case 5:
                                    address := int(read4Bytes(c.a[dr][:]))
                                    address += int(read2Bytes(c.rom[c.pc+2:c.pc+4]))
                                    overflow, carry = addTo(c.ram[address:address+size], c.d[register][:], size, false)
                                    result = c.ram[address:address+size]
                                case 6:
                                    data, reg, word := parse8bitDisplacement(c.rom[c.pc+2])
                                    address := int(read4Bytes(c.a[dr][:]))
                                    if data {
                                        if word {
                                            address += int(signExtend2to4(read2Bytes(c.d[reg][2:4])))
                                        } else {
                                            address += int(read4Bytes(c.d[reg][:]))
                                        }
                                    } else {
                                        if word {
                                            address += int(signExtend2to4(read2Bytes(c.a[reg][2:4])))
                                        } else {
                                            address += int(read4Bytes(c.a[reg][:]))
                                        }
                                    }
                                    address += int(c.rom[c.pc+3])
                                    overflow, carry = addTo(c.ram[address:address+size], c.d[register][:], size, false)
                                    result = c.ram[address:address+size]
                                case 9:
                                    address := int(signExtend2to4(read2Bytes(c.rom[c.pc:c.pc+2])))
                                    overflow, carry = addTo(c.ram[address:address+size], c.d[register][:], size, false)
                                    result = c.ram[address:address+size]
                                case 10:
                                    address := int(read4Bytes(c.rom[c.pc:c.pc+4]))
                                    overflow, carry = addTo(c.ram[address:address+size], c.d[register][:], size, false)
                                    result = c.ram[address:address+size]
                                default:
                                    err = c.error("Can only ADD to memory alterable addressing mode")
                                    return
                                }
                                
                                if carry {
                                    c.sr[1] |= bit0
                                    c.sr[1] |= bit4
                                } else {
                                    c.sr[1] &= ^bit0
                                    c.sr[1] &= ^bit4
                                }
                                if overflow {
                                    c.sr[1] |= bit1
                                } else {
                                    c.sr[1] &= ^bit1
                                }
                                if isZero(result) {
                                    c.sr[1] &= ^bit3
                                    c.sr[1] |= bit2
                                } else if isNegative(result) {
                                    c.sr[1] &= ^bit2
                                    c.sr[1] |= bit3
                                }
                            }
                        } else { // Direction to Dn
                            //ADD
                            sm, sr := addressingmode(i2[5], i2[4], i2[3],
                                                     i2[2], i2[1], i2[0])
                            tmp := c.loadByAddressing(sm, sr, size, 0)
                            
                            overflow, carry := addTo(c.d[register][:], tmp, size, false)
                            if carry {
                                c.sr[1] |= bit0
                                c.sr[1] |= bit4
                            } else {
                                c.sr[1] &= ^bit0
                                c.sr[1] &= ^bit4
                            }
                            if overflow {
                                c.sr[1] |= bit1
                            } else {
                                c.sr[1] &= ^bit1
                            }
                            if isZero(tmp) {
                                c.sr[1] &= ^bit3
                                c.sr[1] |= bit2
                            } else if isNegative(tmp) {
                                c.sr[1] &= ^bit2
                                c.sr[1] |= bit3
                            }
                        }
                    }
                } else {
                    return
                }
            }
        } else {
            return
        }
    } else {
        if i1[6] {
            if c.rom[c.pc] == 0x4e && c.rom[c.pc+1] == 0x71 {
                //NOP
            } else {
                //ADDQ somewhere here
                return
            }
        } else {
            if !i1[5] && !i1[4] {
                //ADDI somewhere here
                return
            } else { // MOVE instruction
                size := opsize(2, i1[5], i1[4])
                dm, dr := addressingmode(i1[0], i2[7], i2[6],
                                         i1[3], i1[2], i1[1])
                sm, sr := addressingmode(i2[5], i2[4], i2[3],
                                         i2[2], i2[1], i2[0])

                extraBytes := bytesUsedByAddressing(dm, size)
                tmp := c.loadByAddressing(sm, sr, size, extraBytes)

                switch dm {
                case 0:
                    for i := 0; i < size; i++ {
                        c.d[dr][4-size+i] = tmp[i]
                    }
                case 1: // MOVEA instruction
                    if size == 1 {
                        err = c.error("MOVEA cannot handle size 1")
                        return
                    }
                    for i := 0; i < size; i++ {
                        c.a[dr][4-size+i] = tmp[i]
                    }
                case 2:
                    address := int(read4Bytes(c.a[dr][:]))
                    for i := 0; i < size; i++ {
                        c.ram[address + i] = tmp[i]
                    }
                case 3:
                    address := int(read4Bytes(c.a[dr][:]))
                    for i := 0; i < size; i++ {
                        c.ram[address + i] = tmp[i]
                    }
                    increment(c.a[dr][:], size)
                case 4:
                    increment(c.a[dr][:], -size)
                    address := int(read4Bytes(c.a[dr][:]))
                    for i := 0; i < size; i++ {
                        c.ram[address + i] = tmp[i]
                    }
                case 5:
                    address := int(read4Bytes(c.a[dr][:]))
                    address += int(read2Bytes(c.rom[c.pc+2:c.pc+4]))
                    for i := 0; i < size; i++ {
                        c.ram[address + i] = tmp[i]
                    }
                case 6:
                    data, register, word := parse8bitDisplacement(c.rom[c.pc+2])
                    address := int(read4Bytes(c.a[dr][:]))
                    if data {
                        if word {
                            address += int(signExtend2to4(read2Bytes(c.d[register][2:4])))
                        } else {
                            address += int(read4Bytes(c.d[register][:]))
                        }
                    } else {
                        if word {
                            address += int(signExtend2to4(read2Bytes(c.a[register][2:4])))
                        } else {
                            address += int(read4Bytes(c.a[register][:]))
                        }
                    }
                    address += int(c.rom[c.pc+3])
                    for i := 0; i < size; i++ {
                        c.ram[address + i] = tmp[i]
                    }
                case 9:
                    address := int(signExtend2to4(read2Bytes(c.rom[c.pc:c.pc+2])))
                    for i := 0; i < size; i++ {
                        c.ram[address + i] = tmp[i]
                    }
                case 10:
                    address := int(read4Bytes(c.rom[c.pc:c.pc+4]))
                    for i := 0; i < size; i++ {
                        c.ram[address + i] = tmp[i]
                    }
                default:
                    err = c.error("Can only MOVE to data alterable addressing modes")
                    return
                }

                c.sr[1] &= ^bit0
                c.sr[1] &= ^bit1
                if isZero(tmp) {
                    c.sr[1] &= ^bit3
                    c.sr[1] |= bit2
                } else if isNegative(tmp) {
                    c.sr[1] &= ^bit2
                    c.sr[1] |= bit3
                }
                
                c.pc += extraBytes + bytesUsedByAddressing(sm, size)
            }
        }
    }

    c.pc += 2
    b = false
    return
}

func (c *cpu) loadByAddressing(mode, register, size, extraBytes int) []byte {
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
        address := int(read4Bytes(c.a[register][:]))
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 3:
        address := int(read4Bytes(c.a[register][:]))
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
        increment(c.a[register][:], size)
    case 4:
        increment(c.a[register][:], -size)
        address := int(read4Bytes(c.a[register][:]))
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 5:
        address := int(read4Bytes(c.a[register][:]))
        address += int(read2Bytes(c.rom[c.pc+extraBytes+2:c.pc+extraBytes+4]))
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 6:
        data, reg, word := parse8bitDisplacement(c.rom[c.pc+extraBytes+2])
        address := int(read4Bytes(c.a[register][:]))
        if data {
            if word {
                address += int(signExtend2to4(read2Bytes(c.d[reg][2:4])))
            } else {
                address += int(read4Bytes(c.d[reg][:]))
            }
        } else {
            if word {
                address += int(signExtend2to4(read2Bytes(c.a[reg][2:4])))
            } else {
                address += int(read4Bytes(c.a[reg][:]))
            }
        }
        address += int(c.rom[c.pc+extraBytes+3])
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 7:
        inc := int(signExtend2to4(read2Bytes(c.rom[c.pc+extraBytes:c.pc+extraBytes+2])))
        address := c.pc + inc
        for i := 0; i < size; i++ {
            tmp[i] = c.rom[address + i]
        }
    case 8:
        data, reg, word := parse8bitDisplacement(c.rom[c.pc+extraBytes+2])
        address := c.pc
        if data {
            if word {
                address += int(signExtend2to4(read2Bytes(c.d[reg][2:4])))
            } else {
                address += int(read4Bytes(c.d[reg][:]))
            }
        } else {
            if word {
                address += int(signExtend2to4(read2Bytes(c.a[reg][2:4])))
            } else {
                address += int(read4Bytes(c.a[reg][:]))
            }
        }
        address += int(c.rom[c.pc+extraBytes+3])
        for i := 0; i < size; i++ {
            tmp[i] = c.rom[address + i]
        }
    case 9:
        address := int(signExtend2to4(read2Bytes(c.rom[c.pc+extraBytes:c.pc+extraBytes+2])))
        for i := 0; i < size; i++ {
            tmp[i] = c.ram[address + i]
        }
    case 10:
        address := int(read4Bytes(c.rom[c.pc+extraBytes:c.pc+extraBytes+4]))
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