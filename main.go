package main

import (
    "fmt"
    //"io/ioutil"
    //"os"
    //"strconv"
    //"strings"
    //"time"
    //"unsafe"
    
    "github.com/marcantoinedumais/SegaGenesis/processor"
)

func main() {
    /*if len(os.Args) == 1 {
        fmt.Println("Missing sudoku file name. \nUsage: SudokuGo filename")
        return
    }

    g := loadGrid(os.Args[1])
    fmt.Println(g.String())
    if g.solve() {
        fmt.Println(g.String())
    } else {
        fmt.Println("Could not find a solution for this sudoku.")
    }*/
    cpu := processor.Create()
    fmt.Println(cpu.String())
    cpu.Step()
}