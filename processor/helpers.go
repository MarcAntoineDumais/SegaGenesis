package processor

import (
    "fmt"
    "io/ioutil"
)

func (c *cpu) LoadFile(filename string) error {
    /*f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()
    */
    b, e := ioutil.ReadFile(filename)
    for i := range b {
        c.rom[i] = b[i]
    }
    
    return e
}

func (c *cpu) String() string {
    return fmt.Sprintf("d: %s\na: %s\nsr: %s\npc: %v\nnext instruction: %4x",
                        formatRegisters(c.d), formatRegisters(c.a), formatBytes(c.sr[:]), c.pc, c.rom[c.pc:c.pc+2])
}

func formatBytes(b []byte) string {
    s := "["
    for i := range b {
        s += fmt.Sprintf("%2x ", b[i])
    }
    return s[:len(s)-1] + "]"
}

func formatRegisters(r [8][4]byte) string {
    s := "["
    for i := 0; i < 8; i++ {
        s += fmt.Sprintf("%d", i) + formatBytes(r[i][:]) + " "
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