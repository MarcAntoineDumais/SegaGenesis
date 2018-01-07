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
    pc int           // program counter
}

func Create() (*cpu) {
    c := new(cpu)
    return c
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
                            } else {
                                //ADD
                                dm, dr := addressingmode(i2[5], i2[4], i2[3],
                                                         i2[2], i2[1], i2[0])
                                
                                switch dm {
                                //
                                default:
                                    err = c.error("ADD Unexpected addressing mode")
                                    return
                                }
                            }
                        } else { // Direction to Dn
                            //ADD
                            sm, sr := addressingmode(i2[5], i2[4], i2[3],
                                                     i2[2], i2[1], i2[0])
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
                            default:
                             err = c.error("ADD Unexpected addressing mode")
                             return
                            }
                            
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
                return
            }
        } else {
            if !i1[5] && !i1[4] {
                return
            } else { // MOVE instruction
                size := opsize(2, i1[5], i1[4])
                dm, dr := addressingmode(i1[0], i2[7], i2[6],
                                         i1[3], i1[2], i1[1])
                sm, sr := addressingmode(i2[5], i2[4], i2[3],
                                         i2[2], i2[1], i2[0])

                extraBytes := bytesUsedByAddressing(dm, size)
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
                case 2:
                    address := readBytes(c.a[sr][:], 4)
                    for i := 0; i < size; i++ {
                        tmp[i] = c.ram[address + i]
                    }
                case 3:
                    address := readBytes(c.a[sr][:], 4)
                    for i := 0; i < size; i++ {
                        tmp[i] = c.ram[address + i]
                    }
                    increment(c.a[sr][:], size)
                case 4:
                    increment(c.a[sr][:], -size)
                    address := readBytes(c.a[sr][:], 4)
                    for i := 0; i < size; i++ {
                        tmp[i] = c.ram[address + i]
                    }
                case 5:
                    address := readBytes(c.a[sr][:], 4)
                    address += binary.Uint16(c.rom[c.pc+extraBytes+2:c.pc+extraBytes+4])
                    for i := 0; i < size; i++ {
                        tmp[i] = c.ram[address + i]
                    }
                case 6:
                    //address with index
                case 7:
                    inc := readBytes(c.rom[c.pc+extraBytes:c.pc+extraBytes+2], 2)
                    address := c.pc + int(signExtend2to4(inc))
                    for i := 0; i < size; i++ {
                        tmp[i] = c.ram[address + i]
                    }
                case 8:
                    //PC with index
                case 9:
                    //absolute short
                case 10:
                    //absolute long
                case 11: //immediate
                    if size == 1 {
                        tmp[0] = c.rom[c.pc+extraBytes+3]
                    } else {
                        for i := 0; i < size; i++ {
                            tmp[i] = c.rom[c.pc+extraBytes+2+i]
                        }
                    }
                }
                

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
                default:
                    err = c.error("MOVE Unexpected addressing mode")
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

func (c *cpu) Run() error {
    var err error
    b := false
    for !b {
        b, err = c.Step()
    }
    return err
}
