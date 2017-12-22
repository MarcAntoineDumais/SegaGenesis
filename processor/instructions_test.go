package processor

import (
    "testing"
)

func TestMove(t *testing.T) {
    c := Create()
    /*
        MOVE.L #$ff123456,D0
        MOVE.W #$5b2a,A0
        MOVE.B #$3f,D1
        MOVE.L D0,A1
        MOVE.W D0,D2
        MOVE.W A0,A2
        MOVE.W #$0,D3
    */
    c.LoadFile("tests/move.bin")
    e := c.Run()
    if e != nil {
        t.Errorf("Error in test: %v", e)
    }
    
    expected := cpu{
        d: [8][4]byte{
            [4]byte{0xff, 0x12, 0x34, 0x56},
            [4]byte{0x00, 0x00, 0x00, 0x3f},
            [4]byte{0x00, 0x00, 0x34, 0x56},
            [4]byte{0x00, 0x00, 0x00, 0x00},
            [4]byte{0x00, 0x00, 0x00, 0x00},
            [4]byte{0x00, 0x00, 0x00, 0x00},
            [4]byte{0x00, 0x00, 0x00, 0x00},
            [4]byte{0x00, 0x00, 0x00, 0x00},
        }, 
        a: [8][4]byte{
            [4]byte{0x00, 0x00, 0x5b, 0x2a},
            [4]byte{0xff, 0x12, 0x34, 0x56},
            [4]byte{0x00, 0x00, 0x5b, 0x2a},
            [4]byte{0x00, 0x00, 0x00, 0x00},
            [4]byte{0x00, 0x00, 0x00, 0x00},
            [4]byte{0x00, 0x00, 0x00, 0x00},
            [4]byte{0x00, 0x00, 0x00, 0x00},
            [4]byte{0x00, 0x00, 0x00, 0x00},
        }, 
        sr: [2]byte{0x00, 0x04},
    }
    
    compareCPU(*c, expected, t)
}

func compareCPU(a, b cpu, t *testing.T) {
    err := false
    if a.d != b.d {
        err = true
        t.Log("Data registers different from expected")
    }
    if a.a != b.a {
        err = true
        t.Log("Address registers different from expected")
    }
    if a.sr != b.sr {
        err = true
        t.Log("Status register is different from expected")
    }
    if err {
        t.Errorf("Actual:\n%s\nExpected:\n%s", a.String(), b.String())
    }
}

func TestAdd(t *testing.T) {

}