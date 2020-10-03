package main

import "errors"
import "fmt"
import "os"
import "strconv"

func isDigit(a int) bool  {
    return '0' <= a && a <= '9'
}

func parseNum(input string, offset int) (result int, remaining int, err error) {
    l := len(input)
    a := offset
    for {
        if isDigit(int(input[a])) {
            a++
        } else {
            break
        }

        if a == l {
            break
        }
    }

    if a == offset {
        result = 0
        remaining = offset
        err = errors.New("parseNum失敗")
        return
    }

    result, err = strconv.Atoi(input[offset:a])
    if err != nil {
        result = 0
        remaining = offset
        return
    }

    remaining = a
    return
}

func parseOp(input string, offset int) (result int, remaining int, err error) {
    if input[offset] == '+' || input[offset] == '-' {
        result = int(input[offset])
        remaining = offset + 1
        return
    }
    return 0, offset, errors.New("parseOp failed")
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "引数の個数が正しくありません\n")
        os.Exit(1)
    }

    input := os.Args[1]
    l := len(input)

    fmt.Printf(".intel_syntax noprefix\n")
    fmt.Printf(".globl main\n")
    fmt.Printf("main:\n")

    // 最初の数
    a0, off, err := parseNum(input, 0)
    if err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
        os.Exit(1)
    }
    fmt.Printf("  mov rax, %d\n", a0)

    for {
        if (off != l) {
            // +か-
            op, off1, err1 := parseOp(input, off)
            if err1 != nil {
                fmt.Fprintf(os.Stderr, "エラー: %s", err1.Error())
                os.Exit(1)
            }
            off = off1

            // 次の数
            a, off2, err2 := parseNum(input, off)
            if err2 != nil {
                fmt.Fprintf(os.Stderr, "エラー: %s", err2.Error())
                os.Exit(1)
            }
            off = off2

            if op == '+' {
                fmt.Printf("  add rax, %d\n", a)
            } else {
                fmt.Printf("  sub rax, %d\n", a)
            }
        } else {
            break
        }
    }

    fmt.Printf("  ret\n")
}
