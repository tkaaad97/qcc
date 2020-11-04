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
    TokenChar
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
    NodeAdd NodeKind = 1
    NodeSub = 2
    NodeMul = 3
    NodeDiv = 4
    NodeNum = 5
    NodeAssign = 6
    NodeEq = 7
    NodeNeq = 8
    NodeLt = 9
    NodeLe = 10
    NodeGt = 11
    NodeGe = 12
    NodeLVar = 13
    NodeReturn = 14
    NodeIf = 15
    NodeEither = 16
    NodeFor = 17
    NodeForFirst = 18
    NodeForSecond = 19
    NodeWhile = 20
    NodeBlock = 21
    NodeBlockChild = 22
    NodeFuncCall = 23
    NodeFuncArg = 24
    NodeFuncDef = 25
    NodeAddr = 26
    NodeDeref = 27
    NodeGVar = 28
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
    Globals map[string]*Node
}

type NodeAndLocalSize struct {
    Node *Node
    LocalSize int
}

type CTypeKind int

const (
    CTypeInt CTypeKind = iota
    CTypeChar
    CTypePointer
    CTypeArray
    CTypeFunction
)

type Parameter struct {
    Name string
    Type *CType
}

type CType struct {
    Kind CTypeKind
    PointerTo *CType
    ArraySize int
    ReturnType *CType
    Parameters []Parameter
}

type Register64 int
type Register32 int
type Register16 int
type Register8 int

const (
    Rax Register64 = 0
    Rdi = 1
    Rsi = 2
    Rdx = 3
    Rcx = 4
    Rbp = 5
    Rsp = 6
    Rbx = 7
    R8 = 8
    R9 = 9
    R10 = 10
    R11 = 11
    R12 = 12
    R13 = 13
    R14 = 14
    R15 = 15
)

const (
    Eax Register32 = 0
    Edi = 1
    Esi = 2
    Edx = 3
    Ecx = 4
    Ebp = 5
    Esp = 6
    Ebx = 7
    R8d = 8
    R9d = 9
    R10d = 10
    R11d = 11
    R12d = 12
    R13d = 13
    R14d = 14
    R15d = 15
)

const (
    Ax Register16 = 0
    Di = 1
    Si = 2
    Dx = 3
    Cx = 4
    Bp = 5
    Sp = 6
    Bx = 7
    R8w = 8
    R9w = 9
    R10w = 10
    R11w = 11
    R12w = 12
    R13w = 13
    R14w = 14
    R15w = 15
)

const (
    Al Register8 = 0
    Dil = 1
    Sil = 2
    Dl = 3
    Cl = 4
    Bpl = 5
    Spl = 6
    Bl = 7
    R8b = 8
    R9b = 9
    R10b = 10
    R11b = 11
    R12b = 12
    R13b = 13
    R14b = 14
    R15b = 15
)

func PrintErrorAt(input string, pos int, err string) {
    fmt.Fprintf(os.Stderr, "%s\n", input)
    format := fmt.Sprintf("%%%ds", pos)
    fmt.Fprintf(os.Stderr, format, "")
    fmt.Fprintf(os.Stderr, "^ %s\n", err)
}

func Char() *CType {
    a := CType { CTypeChar, nil, 0, nil, nil }
    return &a
}

func Int() *CType {
    a := CType { CTypeInt, nil, 0, nil, nil }
    return &a
}

func Array(baseType *CType, size int) *CType {
    a := CType { CTypeArray, baseType, size, nil, nil }
    return &a
}

func PointerTo(base *CType) *CType {
    a := CType { CTypePointer, base, 0, nil, nil }
    return &a
}

func Function(returnType *CType, parameters []Parameter) *CType {
    a := CType { CTypeFunction, nil, 0, returnType, parameters }
    return &a
}

func SizeOf(t *CType) int {
    switch (*t).Kind {
    case CTypeChar:
        return 1
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

func IsExpr(node *Node) bool {
    if node == nil {
        return false
    }

    switch ((*node).Kind) {
    case NodeAdd:
        return true
    case NodeSub:
        return true
    case NodeMul:
        return true
    case NodeDiv:
        return true
    case NodeNum:
        return true
    case NodeAssign:
        return true
    case NodeEq:
        return true
    case NodeNeq:
        return true
    case NodeLt:
        return true
    case NodeLe:
        return true
    case NodeGt:
        return true
    case NodeGe:
        return true
    case NodeLVar:
        return true
    case NodeFuncCall:
        return true
    case NodeAddr:
        return true
    case NodeDeref:
        return true
    case NodeGVar:
        return true
    }

    return false
}

func ShowRegister64(r Register64) string {
    switch r {
    case Rax:
        return "rax"
    case Rdi:
        return "rdi"
    case Rsi:
        return "rsi"
    case Rdx:
        return "rdx"
    case Rcx:
        return "rcx"
    case Rbp:
        return "rbp"
    case Rsp:
        return "rsp"
    case Rbx:
        return "rbx"
    case R8:
        return "r8"
    case R9:
        return "r9"
    case R10:
        return "r10"
    case R11:
        return "r11"
    case R12:
        return "r12"
    case R13:
        return "r13"
    case R14:
        return "r14"
    case R15:
        return "r15"
    }

    return "unknown64"
}

func ShowRegister32(r Register32) string {
    switch r {
    case Eax:
        return "eax"
    case Edi:
        return "edi"
    case Esi:
        return "esi"
    case Edx:
        return "edx"
    case Ecx:
        return "ecx"
    case Ebp:
        return "ebp"
    case Esp:
        return "esp"
    case Ebx:
        return "ebx"
    case R8d:
        return "r8d"
    case R9d:
        return "r9d"
    case R10d:
        return "r10d"
    case R11d:
        return "r11d"
    case R12d:
        return "r12d"
    case R13d:
        return "r13d"
    case R14d:
        return "r14d"
    case R15d:
        return "r15d"
    }

    return "unknown32"
}

func ShowRegister16(r Register16) string {
    switch r {
    case Ax:
        return "ax"
    case Di:
        return "di"
    case Si:
        return "si"
    case Dx:
        return "dx"
    case Cx:
        return "cx"
    case Bp:
        return "bp"
    case Sp:
        return "sp"
    case Bx:
        return "bx"
    case R8w:
        return "r8w"
    case R9w:
        return "r9w"
    case R10w:
        return "r10w"
    case R11w:
        return "r11w"
    case R12w:
        return "r12w"
    case R13w:
        return "r13w"
    case R14w:
        return "r14w"
    case R15w:
        return "r15w"
    }

    return "unknown16"
}

func ShowRegister8(r Register8) string {
    switch r {
    case Al:
        return "ax"
    case Dil:
        return "di"
    case Sil:
        return "si"
    case Dl:
        return "dx"
    case Cl:
        return "cx"
    case Bpl:
        return "bp"
    case Spl:
        return "sp"
    case Bl:
        return "bx"
    case R8b:
        return "r8b"
    case R9b:
        return "r9b"
    case R10b:
        return "r10b"
    case R11b:
        return "r11b"
    case R12b:
        return "r12b"
    case R13b:
        return "r13b"
    case R14b:
        return "r14b"
    case R15b:
        return "r15b"
    }

    return "unknown8"
}
