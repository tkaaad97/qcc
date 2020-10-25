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
            if IsAlpha(c) {
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
            kind := TokenReserved
            switch (string(ident)) {
            case "int":
                kind = TokenInt
            case "return":
                kind = TokenReturn
            case "if":
                kind = TokenIf
            case "else":
                kind = TokenElse
            case "for":
                kind = TokenFor
            case "while":
                kind = TokenWhile
            }
            token := Token {
                Kind: kind,
                Str: string(ident),
                Pos: off,
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

        if (input[off] == ',') {
            token := Token {
                Kind: TokenComma,
                Str: ",",
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if (input[off] == '&') {
            token := Token {
                Kind: TokenReserved,
                Str: "&",
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
                Kind: TokenLeftParenthesis,
                Str: "(",
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if input[off] == ')' {
            token := Token {
                Kind: TokenRightParenthesis,
                Str: ")",
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if input[off] == '{' {
            token := Token {
                Kind: TokenLeftBrace,
                Str: "{",
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if input[off] == '}' {
            token := Token {
                Kind: TokenRightBrace,
                Str: "}",
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

func ConsumeTokenKind(state *ParserState, kind TokenKind) bool {
    if (*state).Offset >= len((*state).Tokens) {
        return false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == kind {
        (*state).Offset++
        return true
    }
    return false
}

func ConsumeLeftParenthesis(state *ParserState) bool {
    return ConsumeTokenKind(state, TokenLeftParenthesis)
}

func ConsumeRightParenthesis(state *ParserState) bool {
    return ConsumeTokenKind(state, TokenRightParenthesis)
}

func ConsumeLeftBrace(state *ParserState) bool {
    return ConsumeTokenKind(state, TokenLeftBrace)
}

func ConsumeRightBrace(state *ParserState) bool {
    return ConsumeTokenKind(state, TokenRightBrace)
}

func ConsumeElse(state *ParserState) bool {
    return ConsumeTokenKind(state, TokenElse)
}

func ConsumeComma(state *ParserState) bool {
    return ConsumeTokenKind(state, TokenComma)
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

func ConsumeType(state *ParserState) (string, bool) {
    if (*state).Offset >= len((*state).Tokens) {
        return "", false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenInt {
        (*state).Offset++
        return token.Str, true
    }
    return "", false
}

func NewNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
    node := Node { kind, lhs, rhs, 0, 0, "" }
    return &node
}

func NewNodeNum(val int) *Node {
    p := NewNode(NodeNum, nil, nil)
    (*p).Val = val
    return p
}

func NewNodeDecl(state *ParserState, name string) (*Node, error) {
    locals := &((*state).Locals)
    if _, exists := (*locals)[name]; exists {
        return nil, fmt.Errorf("変数名はすでに使われています。name: %s", name)
    }
    o := (len(*locals) + 1) * 8
    p := NewNode(NodeLVar, nil, nil)
    (*p).Offset = o
    (*locals)[name] = p
    return p, nil
}

func RefNodeLVar(state *ParserState, name string) (*Node, error) {
    locals := &((*state).Locals)
    if n, ok := (*locals)[name]; ok {
        return n, nil
    }
    return nil, fmt.Errorf("変数が見つかりません。name: %s", name)
}

func NewNodeBlock(nodes []*Node) *Node {
    node := NewNode(NodeBlock, nil, nil)
    current := node
    for i, a := range(nodes) {
        if i == 0 {
            (*current).Lhs = a
        } else {
            (*current).Rhs = NewNode(NodeBlockChild, a, nil)
            current = (*current).Rhs
        }
    }
    return node
}

func NewNodeFuncCall(name string, args []*Node) *Node {
    node := NewNode(NodeFuncCall, nil, nil)
    (*node).Ident = name
    current := node
    if len(args) > 6 {
        args = args[:6]
    }
    for i, arg := range(args) {
        if i == 0 {
            (*current).Lhs = NewNode(NodeFuncArg, arg, nil)
            current = (*current).Lhs
        } else {
            (*current).Rhs = NewNode(NodeFuncArg, arg, nil)
            current = (*current).Rhs
        }
    }
    return node
}

func NewNodeFuncDef(funcName string, params []*Node, block *Node) *Node {
    node := NewNode(NodeFuncDef, nil, block)
    (*node).Ident = funcName
    current := node
    for _, param := range(params) {
        (*current).Lhs = param
        current = param
    }
    return node
}

func Program(state *ParserState) ([]NodeAndLocals, error) {
    defs := []NodeAndLocals{}
    for {
        if (*state).Offset >= len((*state).Tokens) {
            break
        }

        if node, err := FuncDef(state); err != nil {
            return []NodeAndLocals{}, err
        } else {
            def := NodeAndLocals { node, (*state).Locals }
            defs = append(defs, def)
            (*state).Locals = make(map[string]*Node)
        }
    }

    return defs, nil
}

func FuncDef(state *ParserState) (*Node, error) {
    if _, consumed := ConsumeType(state); !consumed {
        return nil, errors.New("関数定義パース失敗。型がありません")
    }

    if funcName, consumed := ConsumeIdent(state); consumed {
        if !ConsumeLeftParenthesis(state) {
            return nil, errors.New("関数定義パース失敗")
        }
        paramNodes := []*Node{}
        if !ConsumeRightParenthesis(state) {
            for {
                if _, consumed := ConsumeType(state); !consumed {
                    return nil, errors.New("関数定義パース失敗。型がありません")
                }

                if paramName, consumed := ConsumeIdent(state); consumed {
                    // 引数はローカル変数と同じように扱う
                    if paramNode, err := NewNodeDecl(state, paramName); err != nil {
                        return nil, err
                    } else {
                        paramNodes = append(paramNodes, paramNode)
                    }
                    if ConsumeRightParenthesis(state) {
                        break
                    } else if !ConsumeComma(state) {
                        return nil, errors.New("関数定義仮引数の後にカンマがありません")
                    }
                } else {
                    return nil, errors.New("関数定義仮引数パース失敗")
                }
            }
        }

        if block, err := Block(state); err != nil {
            return nil, err
        } else {
            return NewNodeFuncDef(funcName, paramNodes, block), nil
        }
    } else {
        return nil, errors.New("関数定義パース失敗")
    }
}

func Block(state *ParserState) (*Node, error) {
    if ConsumeLeftBrace(state) {
        nodes := []*Node{}
        for {
            if !ConsumeRightBrace(state) {
                if stmt, err := Stmt(state); err != nil {
                    return nil, err
                } else {
                    nodes = append(nodes, stmt)
                }
            } else {
                break
            }
        }
        return NewNodeBlock(nodes), nil
    } else {
        return nil, errors.New("ブロックパース失敗。\"{\"がありません。")
    }
}

func Stmt(state *ParserState) (*Node, error) {
    if (*state).Offset >= len((*state).Tokens) {
        return nil, fmt.Errorf("Stmtパース失敗 %v", *state)
    }

    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenLeftBrace {
        return Block(state)
    }

    if token.Kind == TokenReturn {
        return Return(state)
    }

    if token.Kind == TokenIf {
        (*state).Offset++
        if !ConsumeLeftParenthesis(state) {
            return nil, fmt.Errorf("Stmtパース失敗 %v", *state)
        }

        var cond *Node
        if expr, err := Expr(state); err != nil {
            return nil, err
        } else {
            cond = expr
        }

        if !ConsumeRightParenthesis(state) {
            return nil, errors.New("Stmtパース失敗。\")\"が不足しています。")
        }

        var rhs *Node
        if stmt, err := Stmt(state); err != nil {
            return nil, err
        } else {
            rhs = stmt
        }

        if ConsumeElse(state) {
            if stmt, err := Stmt(state); err != nil {
                return nil, err
            } else {
                rhs = NewNode(NodeEither, rhs, stmt)
            }
        }

        return NewNode(NodeIf, cond, rhs), nil
    }

    if token.Kind == TokenFor {
        (*state).Offset++
        if !ConsumeLeftParenthesis(state) {
            return nil, fmt.Errorf("Stmtパース失敗 %v", *state)
        }

        var pre *Node
        if !ConsumeOp(state, ";") {
            if expr, err := Expr(state); err != nil {
                return nil, err
            } else {
                pre = expr
            }
            if !ConsumeOp(state, ";") {
                return nil, errors.New("for文パース失敗。\";\"が不足しています。")
            }
        } else {
            pre = NewNodeNum(1)
        }

        var cond *Node
        if !ConsumeOp(state, ";") {
            if expr, err := Expr(state); err != nil {
                return nil, err
            } else {
                cond = expr
            }
            if !ConsumeOp(state, ";") {
                return nil, errors.New("for文パース失敗。\";\"が不足しています。")
            }
        } else {
            cond = NewNodeNum(1)
        }

        var post *Node
        if !ConsumeRightParenthesis(state) {
            if expr, err := Expr(state); err != nil {
                return nil, err
            } else {
                post = expr
            }

            if !ConsumeRightParenthesis(state) {
                return nil, errors.New("Stmtパース失敗。\")\"が不足しています。")
            }
        } else {
            post = NewNodeNum(1)
        }

        var action *Node
        if stmt, err := Stmt(state); err != nil {
            return nil, err
        } else {
            action = stmt
        }

        return NewNode(NodeFor, NewNode(NodeForFirst, pre, cond), NewNode(NodeForSecond, post, action)), nil
    }

    if token.Kind == TokenWhile {
        (*state).Offset++
        if !ConsumeLeftParenthesis(state) {
            return nil, fmt.Errorf("Stmtパース失敗 %v", *state)
        }

        var cond *Node
        if expr, err := Expr(state); err != nil {
            return nil, err
        } else {
            cond = expr
        }

        if !ConsumeRightParenthesis(state) {
            return nil, errors.New("Stmtパース失敗。\")\"が不足しています。")
        }

        var rhs *Node
        if stmt, err := Stmt(state); err != nil {
            return nil, err
        } else {
            rhs = stmt
        }

        return NewNode(NodeWhile, cond, rhs), nil
    }

    if _, consumed0 := ConsumeType(state); consumed0 {
        if ident, consumed := ConsumeIdent(state); consumed {
            if !ConsumeOp(state, ";") {
                return nil, errors.New("Stmtパース失敗。\";\"が不足しています。")
            }
            return NewNodeDecl(state, ident)
        } else {
            return nil, errors.New("Stmtパース失敗。変数名がありません。")
        }
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
        if ConsumeLeftParenthesis(state) {
            args := []*Node{}
            if ConsumeRightParenthesis(state) {
                return NewNodeFuncCall(ident, args), nil
            }
            for {
                if expr, err := Expr(state); err != nil {
                    return nil, err
                } else {
                    args = append(args, expr)
                    if ConsumeRightParenthesis(state) {
                        return NewNodeFuncCall(ident, args), nil
                    } else if !ConsumeComma(state) {
                        return nil, errors.New("関数呼び出し引数の後にカンマがありません")
                    }
                }
            }
        }
        return RefNodeLVar(state, ident)
    }

    if ConsumeLeftParenthesis(state) {
        if n, err := Expr(state); err != nil {
            return nil, err
        } else {
            if ConsumeRightParenthesis(state) {
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
    } else if ConsumeOp(state, "&") {
        if a, err := Primary(state); err != nil {
            return nil, err
        } else {
            return NewNode(NodeAddr, a, nil), nil
        }
    } else if ConsumeOp(state, "*") {
        if a, err := Primary(state); err != nil {
            return nil, err
        } else {
            return NewNode(NodeDeref, a, nil), nil
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
