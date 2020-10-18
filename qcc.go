package main

import (
    "fmt"
    "os"
)

type TokenKind int

const (
    TokenReserved TokenKind = iota
    TokenLeftBracket
    TokenRightBracket
    TokenNum
    TokenIdent
    TokenReturn
    TokenIf
    TokenElse
    TokenEof
)

type Token struct {
    Kind TokenKind
    Val int
    Str string
    Pos int
}

type NodeKind int

const (
    NodeAdd NodeKind = iota
    NodeSub
    NodeMul
    NodeDiv
    NodeNum
    NodeAssign
    NodeEq
    NodeNeq
    NodeLt
    NodeLe
    NodeGt
    NodeGe
    NodeLVar
    NodeReturn
    NodeIf
    NodeElse
)

type Node struct {
    Kind NodeKind
    Lhs *Node
    Rhs *Node
    Val int
    Offset int
}

type ParserState struct {
    Tokens []Token
    Offset int
    Locals *map[string]*Node
}

func PrintErrorAt(input string, pos int, err string) {
    fmt.Fprintf(os.Stderr, "%s\n", input)
    format := fmt.Sprintf("%%%ds", pos)
    fmt.Fprintf(os.Stderr, format, "")
    fmt.Fprintf(os.Stderr, "^ %s\n", err)
}
