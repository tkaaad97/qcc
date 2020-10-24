package main

import (
    "fmt"
    "os"
)

type TokenKind int

const (
    TokenReserved TokenKind = iota
    TokenLeftParenthesis
    TokenRightParenthesis
    TokenLeftBrace
    TokenRightBrace
    TokenNum
    TokenIdent
    TokenReturn
    TokenIf
    TokenElse
    TokenFor
    TokenWhile
    TokenComma
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
    NodeEither
    NodeFor
    NodeForFirst
    NodeForSecond
    NodeWhile
    NodeBlock
    NodeBlockChild
    NodeFuncCall
    NodeFuncArg
    NodeFuncDef
    NodeAddr
    NodeDeref
)

type Node struct {
    Kind NodeKind
    Lhs *Node
    Rhs *Node
    Val int
    Offset int
    Ident string
}

type ParserState struct {
    Tokens []Token
    Offset int
    Locals map[string]*Node
}

type NodeAndLocals struct {
    Node *Node
    Locals map[string]*Node
}

func PrintErrorAt(input string, pos int, err string) {
    fmt.Fprintf(os.Stderr, "%s\n", input)
    format := fmt.Sprintf("%%%ds", pos)
    fmt.Fprintf(os.Stderr, format, "")
    fmt.Fprintf(os.Stderr, "^ %s\n", err)
}
