package main

import (
    "fmt"
    "errors"
    "strconv"
    "unicode"
)

func IsAlpha(a rune) bool {
    return (a >= 'a' && a <= 'z') || (a >= 'A' && a <= 'Z') || (a == '_')
}

func IsAlnum(a rune) bool {
    return unicode.IsDigit(a) || IsAlpha(a)
}

func IsIdent(str string) bool {
    for i, c := range([]rune(str)) {
        if i == 0 {
            if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
                continue
            }
        } else {
            if IsAlnum(c) {
                continue
            }
        }
        return false
    }
    return true;
}

func Tokenize(input []rune) ([]Token, error) {
    l := len(input)
    off := 0
    tokens := make([]Token, 0, 100)

    for {
        if (off >= l) {
            break
        }

        if (unicode.IsSpace(input[off])) {
            off++
            continue
        }

        if IsAlpha(input[off]) {
            ident := make([]rune, 1, 10)
            ident[0] = input[off]
            off++
            for {
                if off >= l {
                    break
                }
                if !IsAlnum(input[off]) {
                    break
                }
                ident = append(ident, input[off])
                off++
            }
            token := Token {
                Kind: TokenReserved,
                Str: string(ident),
                Pos: off,
            }
            if token.Str == "return" {
                token.Kind = TokenReturn
            }
            tokens = append(tokens, token)
            continue
        }

        if (input[off] == ';') {
            token := Token {
                Kind: TokenReserved,
                Str: ";",
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if (input[off] == '=') {
            s := "="
            if (off + 1 < l && input[off + 1] == '=') {
                s = "=="
            }
            token := Token {
                Kind: TokenReserved,
                Str: s,
                Pos: off,
            }
            tokens = append(tokens, token)
            off += len(s)
            continue
        }

        if (input[off] == '!') {
            if (off + 1 >= l || input[off + 1] != '=') {
                return tokens, errors.New("トークナイズ失敗しました。")
            }
            token := Token {
                Kind: TokenReserved,
                Str: "!=",
                Pos: off,
            }
            tokens = append(tokens, token)
            off += 2
            continue
        }

        if (input[off] == '<') {
            s := "<"
            if (off + 1 < l && input[off + 1] == '=') {
                s = "<="
            }
            token := Token {
                Kind: TokenReserved,
                Str: s,
                Pos: off,
            }
            tokens = append(tokens, token)
            off += len(s)
            continue
        }

        if (input[off] == '>') {
            s := ">"
            if (off + 1 < l && input[off + 1] == '=') {
                s = ">="
            }
            token := Token {
                Kind: TokenReserved,
                Str: s,
                Pos: off,
            }
            tokens = append(tokens, token)
            off += len(s)
            continue
        }

        if (input[off] == '+' || input[off] == '-' || input[off] == '*' || input[off] == '/') {
            token := Token {
                Kind: TokenReserved,
                Str: string([]rune{input[off]}),
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if input[off] == '(' {
            token := Token {
                Kind: TokenLeftBracket,
                Str: "(",
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if input[off] == ')' {
            token := Token {
                Kind: TokenRightBracket,
                Str: ")",
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if unicode.IsDigit(input[off]) {
            if token, remaining, err := ParseNum(input, off); err != nil {
                return tokens, errors.New("tokenizeに失敗しました。")
            } else {
                tokens = append(tokens, token)
                off = remaining
            }
            continue
        }

        return tokens, errors.New("tokenizeに失敗しました。")
    }

    return tokens, nil
}

func ParseNum(input []rune, offset int) (Token, int, error) {
    l := len(input)
    a := offset
    for {
        if a >= l {
            break
        }

        if unicode.IsDigit(input[a]) {
            a++
        } else {
            break
        }
    }

    if a == offset {
        return Token{}, offset, errors.New("ParseNum失敗")
    }

    str := string(input[offset:a])
    token := Token {
        Kind: TokenNum,
        Str: str,
        Pos: a,
    }
    if result, err := strconv.Atoi(str); err != nil {
        return Token{}, offset, errors.New("ParseNum失敗")
    } else {
        token.Val = result
    }

    return token, a, nil
}

func ConsumeLeftBracket(state *ParserState) bool {
    if (*state).Offset >= len((*state).Tokens) {
        return false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenLeftBracket {
        (*state).Offset++
        return true
    }
    return false
}

func ConsumeRightBracket(state *ParserState) bool {
    if (*state).Offset >= len((*state).Tokens) {
        return false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenRightBracket {
        (*state).Offset++
        return true
    }
    return false
}

func ConsumeOp(state *ParserState, op string) bool {
    if (*state).Offset >= len((*state).Tokens) {
        return false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenReserved && token.Str == op {
        (*state).Offset++
        return true
    }
    return false
}

func ConsumeNum(state *ParserState) (int, bool) {
    if (*state).Offset >= len((*state).Tokens) {
        return 0, false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenNum {
        (*state).Offset++
        return token.Val, true
    }
    return 0, false
}

func ConsumeIdent(state *ParserState) (string, bool) {
    if (*state).Offset >= len((*state).Tokens) {
        return "", false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenReserved && IsIdent(token.Str) {
        (*state).Offset++
        return token.Str, true
    }
    return "", false
}

func NewNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
    node := Node { kind, lhs, rhs, 0, 0 }
    return &node
}

func NewNodeNum(val int) *Node {
    p := NewNode(NodeNum, nil, nil)
    (*p).Val = val
    return p
}

func NewNodeLVar(state *ParserState, name string) *Node {
    locals := *(*state).Locals
    if n, ok := locals[name]; ok {
        return n
    }
    o := (len(locals) + 1) * 8
    p := NewNode(NodeLVar, nil, nil)
    (*p).Offset = o
    locals[name] = p
    return p
}

func Program(state *ParserState) ([]*Node, error) {
    nodes := make([]*Node, 0, 10)
    for {
        if (*state).Offset >= len((*state).Tokens) {
            break
        }

        if node, err := Stmt(state); err != nil {
            return []*Node{}, err
        } else {
            nodes = append(nodes, node)
        }
    }

    return nodes, nil
}

func Stmt(state *ParserState) (*Node, error) {
    if (*state).Offset >= len((*state).Tokens) {
        return nil, errors.New("Stmtパース失敗")
    }

    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenReturn {
        return Return(state)
    }

    var node *Node
    if expr, err := Expr(state); err != nil {
        return nil, err
    } else {
        node = expr
    }

    if !ConsumeOp(state, ";") {
        return nil, errors.New("Stmtパース失敗")
    }
    return node, nil
}

func Return(state *ParserState) (*Node, error) {
    if (*state).Offset >= len((*state).Tokens) {
        return nil, errors.New("Returnパース失敗")
    }

    token := (*state).Tokens[(*state).Offset]
    if token.Kind != TokenReturn {
        return nil, errors.New("Returnパース失敗")
    }
    (*state).Offset++

    if e, err := Expr(state); err != nil {
        return nil, errors.New("Returnパース失敗")
    } else {
        if !ConsumeOp(state, ";") {
            return nil, errors.New("Returnパース失敗")
        }
        return NewNode(NodeReturn, e, nil), nil
    }
}

func Primary(state *ParserState) (*Node, error) {
    if v, consumed := ConsumeNum(state); consumed {
        return NewNodeNum(v), nil
    }

    if ident, consumed := ConsumeIdent(state); consumed {
        return NewNodeLVar(state, ident), nil
    }

    if ConsumeLeftBracket(state) {
        if n, err := Expr(state); err != nil {
            return nil, err
        } else {
            if ConsumeRightBracket(state) {
                return n, nil
            } else {
                return nil, errors.New("右括弧が不足しています")
            }
        }
    }

    return nil, fmt.Errorf("Primaryのパースに失敗しました。state: %v", *state)
}

func Unary(state *ParserState) (*Node, error) {
    if ConsumeOp(state, "+") {
        return Primary(state)
    } else if ConsumeOp(state, "-") {
        if a, err := Primary(state); err != nil {
            return nil, err
        } else {
            return NewNode(NodeSub, NewNodeNum(0), a), nil
        }
    }

    return Primary(state)
}

func Mul(state *ParserState) (*Node, error) {
    var node *Node
    if lhs, err := Unary(state); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if ConsumeOp(state, "*") {
            if rhs, err := Unary(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeMul, node, rhs)
            }
        } else if ConsumeOp(state, "/") {
            if rhs, err := Unary(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeDiv, node, rhs)
            }
        } else {
            break
        }
    }
    return node, nil
}

func Expr(state *ParserState) (*Node, error) {
    return Assign(state)
}

func Add(state *ParserState) (*Node, error) {
    var node *Node
    if lhs, err := Mul(state); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if ConsumeOp(state, "+") {
            if rhs, err := Mul(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeAdd, node, rhs)
            }
        } else if ConsumeOp(state, "-") {
            if rhs, err := Mul(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeSub, node, rhs)
            }
        } else {
            break
        }
    }
    return node, nil
}

func Assign(state *ParserState) (*Node, error) {
    var node *Node
    if lhs, err := Equality(state); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    if ConsumeOp(state, "=") {
        if rhs, err := Assign(state); err != nil {
            return nil, err
        } else {
            node = NewNode(NodeAssign, node, rhs)
        }
    }
    return node, nil
}

func Equality(state *ParserState) (*Node, error) {
    var node *Node
    if lhs, err := Relational(state); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if ConsumeOp(state, "==") {
            if rhs, err := Relational(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeEq, node, rhs)
            }
        } else if ConsumeOp(state, "!=") {
            if rhs, err := Relational(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeNeq, node, rhs)
            }
        } else {
            break
        }
    }
    return node, nil
}

func Relational(state *ParserState) (*Node, error) {
    var node *Node
    if lhs, err := Add(state); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if ConsumeOp(state, "<") {
            if rhs, err := Add(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeLt, node, rhs)
            }
        } else if ConsumeOp(state, "<=") {
            if rhs, err := Add(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeLe, node, rhs)
            }
        } else if ConsumeOp(state, ">") {
            if rhs, err := Add(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeGt, node, rhs)
            }
        } else if ConsumeOp(state, ">=") {
            if rhs, err := Add(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeGe, node, rhs)
            }
        } else {
            break
        }
    }
    return node, nil
}
