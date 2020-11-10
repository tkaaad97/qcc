package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "引数の個数が正しくありません\n")
        os.Exit(1)
    }

    // トークナイズする
    input := []rune(os.Args[1])
    var tokens []Token
    if tokenized, err := Tokenize(input); err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
        os.Exit(1)
    } else {
        tokens = tokenized
    }

    // プログラムパース
    state := ParserState { tokens, 0, make(map[string]*Node), 0, make(map[string]*CType), make(map[string]*Node), make([]string, 0, 10) }
    if globals, defs, err := Program(&state); err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
        os.Exit(1)
    } else {
        // アセンブラ生成
        GenProgram(globals, state.StringLiterals, defs);
    }
}
