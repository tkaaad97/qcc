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
    TokenLeftBracket
    TokenRightBracket
    TokenNum
    TokenIdent
    TokenReturn
    TokenIf
    TokenElse
    TokenFor
    TokenWhile
    TokenComma
    TokenInt
    TokenSizeOf
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
    NodeDecl
)

type Node struct {
    Kind NodeKind
    Lhs *Node
    Rhs *Node
    Val int
    Offset int
    Ident string
    Type *CType
}

type ParserState struct {
    Tokens []Token
    Offset int
    Locals map[string]*Node
    LocalOffset int
    Funcs map[string]*CType
}

type NodeAndLocalSize struct {
    Node *Node
    LocalSize int
}

type CTypeKind int

const (
    CTypeInt CTypeKind = iota
    CTypePointer
    CTypeArray
)

type CType struct {
    Kind CTypeKind
    PointerTo *CType
    ArraySize int
}

func PrintErrorAt(input string, pos int, err string) {
    fmt.Fprintf(os.Stderr, "%s\n", input)
    format := fmt.Sprintf("%%%ds", pos)
    fmt.Fprintf(os.Stderr, format, "")
    fmt.Fprintf(os.Stderr, "^ %s\n", err)
}

func Int() *CType {
    a := CType { CTypeInt, nil, 0 }
    return &a
}

func Array(baseType *CType, size int) *CType {
    a := CType { CTypeArray, baseType, size }
    return &a
}

func PointerTo(base *CType) *CType {
    a := CType { CTypePointer, base, 0 }
    return &a
}

func SizeOf(t *CType) int {
    switch (*t).Kind {
    case CTypeInt:
        return 4
    case CTypePointer:
        return 8
    case CTypeArray:
        return (*t).ArraySize * SizeOf((*t).PointerTo)
    }
    return -1
}

func DerefType(t *CType) (*CType, bool) {
    if t == nil {
        return nil, false
    }

    if (*t).Kind != CTypePointer {
        return nil, false
    }

    return (*t).PointerTo, true
}

func Gcd(a, b int) int {
    if b == 0 {
        return a
    }
    return Gcd(b, a % b)
}

func Lcm(a, b int) int {
    return a * b / Gcd(a, b)
}
