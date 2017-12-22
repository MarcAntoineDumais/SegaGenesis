package processor

import (
    "encoding/binary"
)

// Supposed to be running at 8MHz
type cpu struct {
    d [8][4]byte        // data registers d0, d1, ..., d7
    a [8][4]byte        // adress registers a0, a1, ..., a7
    sr [2]byte          // status register TT, S, M, 0, III, 0, 0, 0, X, N, Z, V, C
    rom [1048576]byte   // 1MB ROM area
    ram [65536]byte     // 64KB RAM area
    pc uint32           // program counter
}

func Create() (*cpu) {
    c := new(cpu)
    return c
}

func (c *cpu) Step() bool {
    i1 := c.rom[c.pc]
    i2 := c.rom[c.pc+1]
    switch i1 & (bit7 | bit6) {
    case 0:
        if (i1 & (bit5 | bit4)) == 0 {
            return true
        } else { //MOVE instruction
            size := opsize(2, (i1 & bit5) != 0, (i1 & bit4) != 0)
            dm, dr := addressingmode((i1 & bit0) != 0, (i2 & bit7) != 0, (i2 & bit6) != 0,
                                    (i1 & bit3) != 0, (i1 & bit2) != 0, (i1 & bit1) != 0)
            sm, sr := addressingmode((i2 & bit5) != 0, (i2 & bit4) != 0, (i2 & bit3) != 0,
                                    (i2 & bit2) != 0, (i2 & bit1) != 0, (i2 & bit0) != 0)

            tmp := make([]byte, size)
            switch sm {
            case 0:
                for i := 0; i < size; i++ {
                    tmp[i] = c.d[sr][4-size+i]
                }
            case 1:
                for i := 0; i < size; i++ {
                    tmp[i] = c.a[sr][4-size+i]
                }
            }

            switch dm {
            case 0:
                for i := 0; i < size; i++ {
                    c.d[dr][4-size+i] = tmp[i]
                }
            case 1:
                for i := 0; i < size; i++ {
                    c.a[dr][4-size+i] = tmp[i]
                }
            }

            c.sr[1] &= ^bit7
            c.sr[1] &= ^bit6
            var zero, negative bool
            switch size {
            case 1:
                val := tmp[0]
                zero = val == 0
                negative = int8(val) < 0
            case 2:
                val := binary.BigEndian.Uint16(tmp)
                zero = val == 0
                negative = int16(val) < 0
            case 4:
                val := binary.BigEndian.Uint32(tmp)
                zero = val == 0
                negative = int32(val) < 0
            }
            if zero {
                c.sr[1] &= ^bit4
                c.sr[1] |= bit5
            } else if negative {
                c.sr[1] &= ^bit5
                c.sr[1] |= bit4
            }
        }
    case bit6:
        return true
    case bit7:
        return true
    case bit6 | bit7:
        if (i1 & bit5) != 0 {
            return true
        } else {
            if (i1 & bit4) != 0 { //ADD instruction
                return true

            } else {
                return true
            }
        }
    }

    c.pc += 2
    return false
    //fmt.Printf("%T: %b\n%T: %b\n%T: %b", bit0, bit0, bit2, bit2, bit7, bit7)
}

func (c *cpu) Run() {
    for !c.Step(){}
}
