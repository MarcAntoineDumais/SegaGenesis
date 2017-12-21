package processor

import (
    "fmt"
)

//constants
//bitmasks
const (
    bit0 = uint8(1)
    bit1 = uint8(1 << uint(1))
    bit2 = uint8(1 << uint(2))
    bit3 = uint8(1 << uint(3))
    bit4 = uint8(1 << uint(4))
    bit5 = uint8(1 << uint(5))
    bit6 = uint8(1 << uint(6))
    bit7 = uint8(1 << uint(7))
)

// supposed to be running at 8MHz
type cpu struct {
    d [8]uint32         // data registers d0, d1, ..., d7
    a [8]uint32         // adress registers a0, a1, ..., a7
    sr uint16           // status register TT, S, M, 0, III, 0, 0, 0, X, N, Z, V, C
    rom [1048576]uint8  // 1MB ROM area
    ram [65536]uint8    // 64KB RAM area
    pc uint32           // program counter
}

func Create() (*cpu) {
    c := new(cpu)
    return c
}

func (c *cpu) String() string {
    return fmt.Sprintf("d: %v\na: %v\nsr: %v\npc: %v\nnext instruction: %v", c.d, c.a, c.sr, c.pc, c.rom[c.pc])
}

func (c *cpu) Step() {
    //switch rom[pc] & ()
    //fmt.Printf("%T: %b\n%T: %b\n%T: %b", bit0, bit0, bit2, bit2, bit7, bit7)
}

func (c *cpu) Run() {
    fmt.Println("running")
}