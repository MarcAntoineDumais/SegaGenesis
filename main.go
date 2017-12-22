package main

import (
    "fmt"
    //"io/ioutil"
    "os"
    //"strconv"
    //"strings"
    //"time"
    //"unsafe"

    "github.com/marcantoinedumais/SegaGenesis/processor"
)

func main() {
    if len(os.Args) == 1 {
        fmt.Println("Missing binary file name. \nUsage: SegaGenesis filename")
        return
    }

    cpu := processor.Create()
    cpu.LoadFile(os.Args[1])    
    fmt.Println("State of processor before execution")
    fmt.Println(cpu.PrintRom(256))
    fmt.Println(cpu.String())
    cpu.Run()
    fmt.Println("State of processor after execution")
    fmt.Println(cpu.String())
}