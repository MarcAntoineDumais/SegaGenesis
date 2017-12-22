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
    i1 := c.rom[c.pc]
    i2 := c.rom[c.pc+1]
    switch i1 & (bit7 | bit6) {
    case 0:
        if (i1 & (bit5 | bit4)) == 0 {
            return
        } else { // MOVE instruction
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
            case 11: //immediate
                if size == 1 {
                    tmp[0] = c.rom[c.pc+3]
                    c.pc += 2
                } else {
                    for i := 0; i < size; i++ {
                        tmp[i] = c.rom[c.pc+2+i]
                    }
                    c.pc += size
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
            }

            c.sr[1] &= ^bit0
            c.sr[1] &= ^bit1
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
                c.sr[1] &= ^bit3
                c.sr[1] |= bit2
            } else if negative {
                c.sr[1] &= ^bit2
                c.sr[1] |= bit3
            }
        }
    case bit6:
        return
    case bit7:
        return
    case bit6 | bit7:
        if (i1 & bit5) != 0 {
            return
        } else {
            if (i1 & bit4) != 0 { //ADD instruction
                return

            } else {
                return
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
