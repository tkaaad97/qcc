package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "引数の個数が正しくありません\n")
        os.Exit(1)
    }

    // トークナイズする
    inputFile := os.Args[1]
    var input []rune
    if bs, err := ioutil.ReadFile(inputFile); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err.Error())
        os.Exit(1)
    } else {
        input = bytes.Runes(bs)
    }
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
