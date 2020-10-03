package main

import "fmt"
import "os"
import "strconv"

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "引数の個数が正しくありません\n")
        os.Exit(1)
    }

    a, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Fprintf(os.Stderr, "引数が整数ではありません\n")
        os.Exit(1)
    }

    fmt.Printf(".intel_syntax noprefix\n")
    fmt.Printf(".globl main\n")
    fmt.Printf("main:\n")
    fmt.Printf("  mov rax, %d\n", a)
    fmt.Printf("  ret\n")
}
