package processor

import (
    "binary"
    "fmt"
)

// supposed to be running at 8MHz
type cpu struct {
    d [8]uint32         // data registers d0, d1, ..., d7
    a [8]uint32         // adress registers a0, a1, ..., a7
    sr uint16           // status register TT, S, M, 0, III, 0, 0, 0, X, N, Z, V, C
    rom [1048576]byte  // 1MB ROM area
    ram [65536]byte    // 64KB RAM area
    pc uint32           // program counter
}

func Create() (*cpu) {
    c := new(cpu)
    return c
}

func (c *cpu) String() string {
    return fmt.Sprintf("d: %v\na: %v\nsr: %v\npc: %v\nnext instruction: %v", c.d, c.a, c.sr, c.pc, c.rom[c.pc])
}

func (c *cpu) Step() bool {
    i1 := c.rom[c.pc]
    i2 := c.rom[c.pc+1]
    switch i1 & (bit7 | bit6) {
    case 0:
        if (i1 & (bit5 | bit4)) == 0 {

        } else { //MOVE instruction
            size := opsize(2, (i1 & bit5) != 0, (i1 & bit4) != 0)
            dm, dr = addressingmode((i1 & bit0) != 0, (i2 & bit7) != 0, (i2 & bit6) != 0,
                                    (i1 & bit3) != 0, (i1 & bit2) != 0, (i1 & bit1) != 0)
            sm, sr = addressingmode((i2 & bit5) != 0, (i2 & bit4) != 0, (i2 & bit3) != 0,
                                    (i2 & bit2) != 0, (i2 & bit1) != 0, (i2 & bit0) != 0)

            tmp := make([]byte, 4)
            switch sm {
            case 0:
                binary.BigEndian.PutUint32(tmp, c.d[sr])
            case 1:
                binary.BigEndian.PutUint32(tmp, c.a[sr])
            }

            switch dm {
            case 0:
                binary.BigEndian.PutUint32(c.d[dr], tmp)
            case 1:
                binary.BigEndian.PutUint32(c.a[dr], tmp)
            }

            c.sr &= ^wbit0
            c.sr &= ^wbit1
            var val uint32
            binary.BigEndian.PutUint32(val, tmp)
            if val == 0 {
                c.sr &= ^wbit3
                c.sr |= wbit2
            } else if (val & lbit31) != 0 {
                c.sr &= ^wbit2
                c.sr |= wbit3
            }
        }
    case bit6:

    case bit7:

    case bit6 | bit7:
        if (i1 & bit5) != 0 {

        } else {
            if (i1 & bit4) != 0 { //ADD instruction


            } else {

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
    default:
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

    return 1
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
        mode = 0
        if x1 {
            register += 4
        }
        if x2 {
            register += 2
        }
        if x3 {
            register += 1
        }
        if m1 {
            mode += 4
        }
        if m2 {
            mode += 2
        }
        if m3 {
            mode += 1
        }
    }
    return
}
