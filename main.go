package main

import (
    "errors"
    "fmt"
    "os"
    "strconv"
    "unicode"
)

type TokenKind int

const (
    TokenReserved TokenKind = iota
    TokenLeftBracket
    TokenRightBracket
    TokenNum
    TokenEof
)

type Token struct {
    kind TokenKind
    val int
    str string
    pos int
}

type NodeKind int

const (
    NodeAdd NodeKind = iota
    NodeSub
    NodeMul
    NodeDiv
    NodeNum
    NodeEq
    NodeNeq
    NodeLt
    NodeLe
    NodeGt
    NodeGe
)

type Node struct {
    kind NodeKind
    lhs *Node
    rhs *Node
    val int
}

func tokenize(input []rune) ([]Token, error) {
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

        if (input[off] == '=') {
            if (off + 1 >= l || input[off + 1] != '=') {
                return tokens, errors.New("トークナイズ失敗しました。")
            }
            token := Token {
                kind: TokenReserved,
                str: "==",
                pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if (input[off] == '!') {
            if (off + 1 >= l || input[off + 1] != '=') {
                return tokens, errors.New("トークナイズ失敗しました。")
            }
            token := Token {
                kind: TokenReserved,
                str: "!=",
                pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if (input[off] == '<') {
            s := "<"
            if (off + 1 < l && input[off + 1] == '=') {
                s = "<="
            }
            token := Token {
                kind: TokenReserved,
                str: s,
                pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if (input[off] == '>') {
            s := ">"
            if (off + 1 < l && input[off + 1] == '=') {
                s = ">="
            }
            token := Token {
                kind: TokenReserved,
                str: s,
                pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if (input[off] == '+' || input[off] == '-' || input[off] == '*' || input[off] == '/') {
            token := Token {
                kind: TokenReserved,
                str: string([]rune{input[off]}),
                pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if input[off] == '(' {
            token := Token {
                kind: TokenLeftBracket,
                str: "(",
                pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if input[off] == ')' {
            token := Token {
                kind: TokenRightBracket,
                str: ")",
                pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if unicode.IsDigit(input[off]) {
            if token, remaining, err := parseNum(input, off); err != nil {
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

func parseNum(input []rune, offset int) (Token, int, error) {
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
        return Token{}, offset, errors.New("parseNum失敗")
    }

    str := string(input[offset:a])
    token := Token {
        kind: TokenNum,
        str: str,
        pos: a,
    }
    if result, err := strconv.Atoi(str); err != nil {
        return Token{}, offset, errors.New("parseNum失敗")
    } else {
        token.val = result
    }

    return token, a, nil
}

func consumeLeftBracket(tokens []Token, offset *int) bool {
    if *offset >= len(tokens) {
        return false
    }
    token := tokens[*offset]
    if token.kind == TokenLeftBracket {
        (*offset)++
        return true
    }
    return false
}

func consumeRightBracket(tokens []Token, offset *int) bool {
    if *offset >= len(tokens) {
        return false
    }
    token := tokens[*offset]
    if token.kind == TokenRightBracket {
        (*offset)++
        return true
    }
    return false
}

func consumeOp(tokens []Token, offset *int, op string) bool {
    if *offset >= len(tokens) {
        return false
    }
    token := tokens[*offset]
    if token.kind == TokenReserved && token.str == op {
        (*offset)++
        return true
    }
    return false
}

func consumeNum(tokens []Token, offset *int) (int, bool) {
    if *offset >= len(tokens) {
        return 0, false
    }
    token := tokens[*offset]
    if token.kind == TokenNum {
        (*offset)++
        return token.val, true
    }
    return 0, false
}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
    node := Node { kind, lhs, rhs, 0 }
    return &node
}

func newNodeNum(val int) *Node {
    p := newNode(NodeNum, nil, nil)
    (*p).val = val
    return p
}

func primary(tokens []Token, offset *int) (*Node, error) {
    if v, consumed := consumeNum(tokens, offset); consumed {
        return newNodeNum(v), nil
    }

    if consumeLeftBracket(tokens, offset) {
        if n, err := expr(tokens, offset); err != nil {
            return nil, err
        } else {
            if consumeRightBracket(tokens, offset) {
                return n, nil
            } else {
                return nil, errors.New("右括弧が不足しています")
            }
        }
    }

    return nil, errors.New("primaryのパースに失敗しました。")
}

func unary(tokens []Token, offset *int) (*Node, error) {
    if consumeOp(tokens, offset, "+") {
        return primary(tokens, offset)
    } else if consumeOp(tokens, offset, "-") {
        if a, err := primary(tokens, offset); err != nil {
            return nil, err
        } else {
            return newNode(NodeSub, newNodeNum(0), a), nil
        }
    }

    return primary(tokens, offset)
}

func mul(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := unary(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if consumeOp(tokens, offset, "*") {
            if rhs, err := unary(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeMul, node, rhs)
            }
        } else if consumeOp(tokens, offset, "/") {
            if rhs, err := unary(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeDiv, node, rhs)
            }
        } else {
            break
        }
    }
    return node, nil
}

func expr(tokens []Token, offset *int) (*Node, error) {
    return equality(tokens, offset)
}

func add(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := mul(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if consumeOp(tokens, offset, "+") {
            if rhs, err := mul(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeAdd, node, rhs)
            }
        } else if consumeOp(tokens, offset, "-") {
            if rhs, err := mul(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeSub, node, rhs)
            }
        } else {
            break
        }
    }
    return node, nil
}

func equality(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := relational(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if consumeOp(tokens, offset, "==") {
            if rhs, err := relational(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeEq, node, rhs)
            }
        } else if consumeOp(tokens, offset, "!=") {
            if rhs, err := relational(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeNeq, node, rhs)
            }
        } else {
            break
        }
    }
    return node, nil
}

func relational(tokens []Token, offset *int) (*Node, error) {
    var node *Node
    if lhs, err := add(tokens, offset); err != nil {
        return nil, err
    } else {
        node = lhs
    }

    for {
        if consumeOp(tokens, offset, "<") {
            if rhs, err := add(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeLt, node, rhs)
            }
        } else if consumeOp(tokens, offset, "<=") {
            if rhs, err := add(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeLe, node, rhs)
            }
        } else if consumeOp(tokens, offset, ">") {
            if rhs, err := add(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeGt, node, rhs)
            }
        } else if consumeOp(tokens, offset, ">=") {
            if rhs, err := add(tokens, offset); err != nil {
                return nil, err
            } else {
                node = newNode(NodeGe, node, rhs)
            }
        } else {
            break
        }
    }
    return node, nil
}

func eval(node *Node) (int, error) {
    if node == nil {
        return 0, errors.New("evalにnilが渡されました。")
    }

    switch(node.kind) {
    case NodeAdd:
        var l, r int
        if lhs, err := eval((*node).lhs); err != nil {
            return 0, err
        } else {
            l = lhs
        }
        if rhs, err := eval((*node).rhs); err != nil {
            return 0, err
        } else {
            r = rhs
        }
        return l + r, nil
    case NodeSub:
        var l, r int
        if lhs, err := eval((*node).lhs); err != nil {
            return 0, err
        } else {
            l = lhs
        }
        if rhs, err := eval((*node).rhs); err != nil {
            return 0, err
        } else {
            r = rhs
        }
        return l - r, nil
    case NodeMul:
        var l, r int
        if lhs, err := eval((*node).lhs); err != nil {
            return 0, err
        } else {
            l = lhs
        }
        if rhs, err := eval((*node).rhs); err != nil {
            return 0, err
        } else {
            r = rhs
        }
        return l * r, nil
    case NodeDiv:
        var l, r int
        if lhs, err := eval((*node).lhs); err != nil {
            return 0, err
        } else {
            l = lhs
        }
        if rhs, err := eval((*node).rhs); err != nil {
            return 0, err
        } else {
            r = rhs
        }
        return l / r, nil
    case NodeNum:
        return node.val, nil
    }

    return 0, errors.New("不明なノードカインド")
}

func gen(node *Node) {
    if ((*node).kind == NodeNum) {
        fmt.Printf("  push %d\n", (*node).val)
        return
    }

    gen((*node).lhs)
    gen((*node).rhs)

    fmt.Printf("  pop rdi\n")
    fmt.Printf("  pop rax\n")

    switch((*node).kind) {
    case NodeAdd:
        fmt.Printf("  add rax, rdi\n")
    case NodeSub:
        fmt.Printf("  sub rax, rdi\n")
    case NodeMul:
        fmt.Printf("  imul rax, rdi\n")
    case NodeDiv:
        fmt.Printf("  cqo\n")
        fmt.Printf("  idiv rdi\n")
    }

    fmt.Printf("  push rax\n")
}

func printErrorAt(input string, pos int, err string) {
    fmt.Fprintf(os.Stderr, "%s\n", input)
    format := fmt.Sprintf("%%%ds", pos)
    fmt.Fprintf(os.Stderr, format, "")
    fmt.Fprintf(os.Stderr, "^ %s\n", err)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "引数の個数が正しくありません\n")
        os.Exit(1)
    }

    // トークナイズする
    input := []rune(os.Args[1])
    var tokens []Token
    if tokenized, err := tokenize(input); err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
        os.Exit(1)
    } else {
        tokens = tokenized
    }

    // exprパース
    offset := 0
    if node, err := expr(tokens, &offset); err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
        os.Exit(1)
    } else {
        // アセンブラ生成
        fmt.Printf(".intel_syntax noprefix\n")
        fmt.Printf(".globl main\n")
        fmt.Printf("main:\n")
        gen(node);
        fmt.Printf("  pop rax\n")
        fmt.Printf("  ret\n")
    }
}
