package main

import (
    "errors"
    "strconv"
    "unicode"
)

func IsIdent(str string) bool {
    return len(str) == 1 && str[0] >= 'a' && str[0] <= 'z'
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

        if (input[off] >= 'a' && input[off] <= 'z') {
            token := Token {
                Kind: TokenReserved,
                Str: string([]rune{input[off]}),
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
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
            off += 2
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

func ConsumeLeftBracket(tokens []Token, offset *int) bool {
    if *offset >= len(tokens) {
        return false
    }
    token := tokens[*offset]
    if token.Kind == TokenLeftBracket {
        (*offset)++
        return true
    }
    return false
}

func ConsumeRightBracket(tokens []Token, offset *int) bool {
    if *offset >= len(tokens) {
        return false
    }
    token := tokens[*offset]
    if token.Kind == TokenRightBracket {
        (*offset)++
        return true
    }
    return false
}

func ConsumeOp(tokens []Token, offset *int, op string) bool {
    if *offset >= len(tokens) {
        return false
    }
    token := tokens[*offset]
    if token.Kind == TokenReserved && token.Str == op {
        (*offset)++
        return true
    }
    return false
}

func ConsumeNum(tokens []Token, offset *int) (int, bool) {
    if *offset >= len(tokens) {
        return 0, false
    }
    token := tokens[*offset]
    if token.Kind == TokenNum {
        (*offset)++
        return token.Val, true
    }
    return 0, false
}

func ConsumeIdent(tokens []Token, offset *int) (string, bool) {
    if *offset >= len(tokens) {
        return "", false
    }
    token := tokens[*offset]
    if token.Kind == TokenReserved && IsIdent(token.Str) {
        (*offset)++
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

func NewNodeLVar(offset int) *Node {
    p := NewNode(NodeNum, nil, nil)
    (*p).Offset = offset
    return p
}

func Program(tokens []Token, offset *int) ([]*Node, error) {
    nodes := make([]*Node, 0, 10)
    for {
        if *offset >= len(tokens) {
            break
        }

        if node, err := Stmt(tokens, offset); err != nil {
            return []*Node{}, err
        } else {
            nodes = append(nodes, node)
        }
    }

    return nodes, nil
}

func Stmt(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if expr, err := Expr(tokens, offset); err != nil {
        return nil, err
    } else {
        node = expr
    }

    if !ConsumeOp(tokens, offset, ";") {
        return nil, errors.New("Stmtパース失敗")
    }
    return node, nil
}

func Primary(tokens []Token, offset *int) (*Node, error) {
    if v, consumed := ConsumeNum(tokens, offset); consumed {
        return NewNodeNum(v), nil
    }

    if v, consumed := ConsumeIdent(tokens, offset); consumed {
        offset := int(v[0]) - 'a'
        return NewNodeLVar(offset), nil
    }

    if ConsumeLeftBracket(tokens, offset) {
        if n, err := Expr(tokens, offset); err != nil {
            return nil, err
        } else {
            if ConsumeRightBracket(tokens, offset) {
                return n, nil
            } else {
                return nil, errors.New("右括弧が不足しています")
            }
        }
    }

    return nil, errors.New("Primaryのパースに失敗しました。")
}

func Unary(tokens []Token, offset *int) (*Node, error) {
    if ConsumeOp(tokens, offset, "+") {
        return Primary(tokens, offset)
    } else if ConsumeOp(tokens, offset, "-") {
        if a, err := Primary(tokens, offset); err != nil {
            return nil, err
        } else {
            return NewNode(NodeSub, NewNodeNum(0), a), nil
        }
    }

    return Primary(tokens, offset)
}

func Mul(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := Unary(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if ConsumeOp(tokens, offset, "*") {
            if rhs, err := Unary(tokens, offset); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeMul, node, rhs)
            }
        } else if ConsumeOp(tokens, offset, "/") {
            if rhs, err := Unary(tokens, offset); err != nil {
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

func Expr(tokens []Token, offset *int) (*Node, error) {
    return Assign(tokens, offset)
}

func Add(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := Mul(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if ConsumeOp(tokens, offset, "+") {
            if rhs, err := Mul(tokens, offset); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeAdd, node, rhs)
            }
        } else if ConsumeOp(tokens, offset, "-") {
            if rhs, err := Mul(tokens, offset); err != nil {
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

func Assign(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := Equality(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    if ConsumeOp(tokens, offset, "=") {
        if rhs, err := Assign(tokens, offset); err != nil {
            return nil, err
        } else {
            node = NewNode(NodeAssign, node, rhs)
        }
    }
    return node, nil
}

func Equality(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := Relational(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if ConsumeOp(tokens, offset, "==") {
            if rhs, err := Relational(tokens, offset); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeEq, node, rhs)
            }
        } else if ConsumeOp(tokens, offset, "!=") {
            if rhs, err := Relational(tokens, offset); err != nil {
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

func Relational(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := Add(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if ConsumeOp(tokens, offset, "<") {
            if rhs, err := Add(tokens, offset); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeLt, node, rhs)
            }
        } else if ConsumeOp(tokens, offset, "<=") {
            if rhs, err := Add(tokens, offset); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeLe, node, rhs)
            }
        } else if ConsumeOp(tokens, offset, ">") {
            if rhs, err := Add(tokens, offset); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeGt, node, rhs)
            }
        } else if ConsumeOp(tokens, offset, ">=") {
            if rhs, err := Add(tokens, offset); err != nil {
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
