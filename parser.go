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

func IsType(token Token) bool {
    switch (token.Kind) {
    case TokenChar:
        return true
    case TokenInt:
        return true
    }
    return false
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
            kind := TokenIdent
            switch (string(ident)) {
            case "char":
                kind = TokenChar
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
            case "sizeof":
                kind = TokenSizeOf
            }
            token := Token {
                Kind: kind,
                Str: string(ident),
                Pos: off,
            }
            tokens = append(tokens, token)
            continue
        }

        if (input[off] == '"') {
            off++
            start := off
            lit := ""
            for {
                if input[off] == '"' {
                    lit = string(input[start:off])
                    off++
                    break
                }
                off++
            }
            token := Token {
                Kind: TokenStringLiteral,
                Str: lit,
                Pos: start,
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

        if input[off] == '[' {
            token := Token {
                Kind: TokenLeftBracket,
                Str: "[",
                Pos: off,
            }
            tokens = append(tokens, token)
            off++
            continue
        }

        if input[off] == ']' {
            token := Token {
                Kind: TokenRightBracket,
                Str: "]",
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

func ConsumeLeftBracket(state *ParserState) bool {
    return ConsumeTokenKind(state, TokenLeftBracket)
}

func ConsumeRightBracket(state *ParserState) bool {
    return ConsumeTokenKind(state, TokenRightBracket)
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
    if token.Kind == TokenIdent && IsIdent(token.Str) {
        (*state).Offset++
        return token.Str, true
    }
    return "", false
}

func ConsumeStringLiteral(state *ParserState) (string, bool) {
    if (*state).Offset >= len((*state).Tokens) {
        return "", false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenStringLiteral {
        (*state).Offset++
        return token.Str, true
    }
    return "", false
}

func ConsumeSizeOf(state *ParserState) bool {
    if (*state).Offset >= len((*state).Tokens) {
        return false
    }
    token := (*state).Tokens[(*state).Offset]
    if token.Kind == TokenSizeOf {
        (*state).Offset++
        return true
    }
    return false
}

func ConsumeType(state *ParserState) (*CType, bool) {
    if (*state).Offset >= len((*state).Tokens) {
        return nil, false
    }
    token := (*state).Tokens[(*state).Offset]
    switch (token.Kind) {
    case TokenChar:
        (*state).Offset++
        return Char(), true
    case TokenInt:
        (*state).Offset++
        return Int(), true
    }
    return nil, false
}

func SatisfyOp(state *ParserState, op string) bool {
    if (*state).Offset >= len((*state).Tokens) {
        return false
    }
    token := (*state).Tokens[(*state).Offset]
    return token.Kind == TokenReserved && token.Str == op
}

func SatisfyTokenKind(state *ParserState, kind TokenKind) bool {
    if (*state).Offset >= len((*state).Tokens) {
        return false
    }
    token := (*state).Tokens[(*state).Offset]
    return token.Kind == kind
}

func SatisfyType(state *ParserState) bool {
    if (*state).Offset >= len((*state).Tokens) {
        return false
    }
    token := (*state).Tokens[(*state).Offset]
    return IsType(token)
}

func NewNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
    node := Node { kind, lhs, rhs, 0, 0, "", nil }
    return &node
}

func NewNodeNum(val int) *Node {
    p := NewNode(NodeNum, nil, nil)
    (*p).Val = val
    p.Type = Int();
    return p
}

func NewNodeAdd(lhs *Node, rhs *Node) *Node {
    node := Node { NodeAdd, lhs, rhs, 0, 0, "", (*lhs).Type }
    return &node
}

func NewNodeSub(lhs *Node, rhs *Node) *Node {
    t := (*lhs).Type
    lt := (*lhs).Type
    rt := (*rhs).Type
    if lt != nil && rt != nil && (*lt).Kind == CTypePointer && (*rt).Kind == CTypePointer {
        t = Int()
    }
    node := Node { NodeSub, lhs, rhs, 0, 0, "", t }
    return &node
}

func NewNodeDecl(state *ParserState, name string, t *CType) (*Node, error) {
    locals := &((*state).Locals)
    if _, exists := (*locals)[name]; exists {
        return nil, fmt.Errorf("変数名はすでに使われています。name: %s", name)
    }
    (*state).LocalOffset += Lcm(SizeOf(t), 8)
    o := (*state).LocalOffset
    p := NewNode(NodeLVar, nil, nil)
    (*p).Offset = o
    (*p).Type = t
    (*locals)[name] = p
    return p, nil
}

func NewNodeGVar(state *ParserState, name string, t *CType) (*Node, error) {
    globals := &((*state).Globals)
    if _, exists := (*globals)[name]; exists {
        return nil, fmt.Errorf("グローバル変数名はすでに使われています。name: %s", name)
    }
    p := NewNode(NodeGVar, nil, nil)
    (*p).Type = t
    (*p).Ident = name
    (*globals)[name] = p
    return p, nil
}

func NewNodeStringLiteral(state *ParserState, lit string) (*Node, error) {
    literals := &(state.StringLiterals)
    off := len(*literals)
    p := NewNode(NodeStringLiteral, nil, nil)
    p.Type = PointerTo(Char())
    p.Offset = off
    *literals = append(*literals, lit)
    return p, nil
}

func RefNode(state *ParserState, name string) (*Node, error) {
    locals := &((*state).Locals)
    if n, ok := (*locals)[name]; ok {
        return n, nil
    }
    globals := &((*state).Globals)
    if n, ok := (*globals)[name]; ok {
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

func NewNodeFuncCall(state* ParserState, name string, args []*Node) *Node {
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
    funcs := (*state).Funcs
    if t, exists := funcs[name]; exists {
        // TODO extern宣言で外部の関数の型を拾えるようにする
        (*node).Type = t
    }
    return node
}

func NewNodeFuncDef(state* ParserState, funcName string, params []*Node, block *Node, returnType *CType) (*Node, error) {
    funcs := &((*state).Funcs)
    if _, exists := (*funcs)[funcName]; exists {
        return nil, fmt.Errorf("関数名が重複しています name: %s", funcName)
    }
    node := NewNode(NodeFuncDef, nil, block)
    (*node).Ident = funcName
    (*node).Type = returnType
    current := node
    for _, param := range(params) {
        (*current).Lhs = param
        current = param
    }
    (*funcs)[funcName] = returnType
    return node, nil
}

func NewNodeAddr(a *Node) *Node {
    node := NewNode(NodeAddr, a, nil)
    (*node).Type = PointerTo((*a).Type)
    return node
}

func NewNodeDeref(a *Node) *Node {
    node := NewNode(NodeDeref, a, nil)
    if t, ok := DerefType(a.Type); ok {
        node.Type = t
    }
    return node
}

func ArrayToPointer(node *Node) *Node {
    if node != nil && ((*node).Kind == NodeLVar || (*node).Kind == NodeGVar) {
        t := (*node).Type
        if t != nil && (*t).Kind == CTypeArray {
            pt := CType { CTypePointer, (*t).PointerTo, 0, nil, nil }
            q := NewNode(NodeAddr, node, nil)
            (*q).Type = &pt
            return q
        }
    }
    return node
}

func Program(state *ParserState) ([]*Node, []NodeAndLocalSize, error) {
    globals := []*Node{}
    defs := []NodeAndLocalSize{}
    for {
        if (*state).Offset >= len((*state).Tokens) {
            break
        }

        if node, err := FuncDefOrDecl(state); err != nil {
            return nil, nil, err
        } else {
            t := (*node).Type
            if (*node).Kind == NodeGVar {
                globals = append(globals, node)
            } else if (*t).Kind != CTypeFunction {
                def := NodeAndLocalSize { node, (*state).LocalOffset }
                defs = append(defs, def)
            }
            (*state).Locals = make(map[string]*Node)
            (*state).LocalOffset = 0
        }
    }

    return globals, defs, nil
}

func FuncDefOrDecl(state *ParserState) (*Node, error) {
    var baseType *CType
    if t, consumed := ConsumeType(state); !consumed {
        return nil, errors.New("パース失敗。型がありません")
    } else {
        baseType = t
    }

    for {
        if ConsumeOp(state, "*") {
            baseType = PointerTo(baseType)
        } else {
            break
        }
    }

    if (*state).Offset >= len((*state).Tokens) {
        return nil, errors.New("パース失敗")
    }

    if ident, cont, err := Declarator(state); err != nil {
        return nil, err
    } else {
        if SatisfyTokenKind(state, TokenLeftBrace) {
            t := cont(baseType)
            paramNodes := []*Node{}
            if (*t).Kind == CTypeFunction {
                params := (*t).Parameters
                for _, param := range(params) {
                    if paramNode, err := NewNodeDecl(state, param.Name, param.Type); err != nil {
                        return nil, err
                    } else {
                        paramNodes = append(paramNodes, paramNode)
                    }
                }
            }

            if block, err := Block(state); err != nil {
                return nil, err
            } else {
                return NewNodeFuncDef(state, ident, paramNodes, block, (*t).ReturnType)
            }
        } else if ConsumeOp(state, ";") {
            return NewNodeGVar(state, ident, cont(baseType));
        } else {
            return nil, errors.New("パース失敗")
        }
    }
}

func Declarator(state *ParserState) (string, func (*CType) *CType, error) {
    if (*state).Offset >= len((*state).Tokens) {
        return "", nil, errors.New("Declaratorパース失敗")
    }

    if ConsumeOp(state, "*") {
        if i, cont, err := DirectDeclarator(state); err != nil {
            return "", nil, err
        } else {
            f := func (baseType *CType) *CType {
                return PointerTo(cont(baseType))
            }
            return i, f, nil
        }
    }

    return DirectDeclarator(state)
}

func DirectDeclarator(state *ParserState) (string, func (*CType) *CType, error) {
    if (*state).Offset >= len((*state).Tokens) {
        return "", nil, errors.New("DirectDeclaratorパース失敗")
    }

    token := (*state).Tokens[(*state).Offset]
    ident := ""
    cont := func (a *CType) *CType { return a }
    if ConsumeLeftParenthesis(state) {
        if i, f, err := Declarator(state); err != nil {
            return "", nil, errors.New("DirectDeclaratorパース失敗")
        } else {
            ident = i
            cont = f
        }
        if !ConsumeRightParenthesis(state) {
            return "", nil, errors.New("DirectDeclaratorパース失敗")
        }
    } else if i, consumed := ConsumeIdent(state); consumed {
        ident = i
    } else {
        return "", nil, errors.New("DirectDeclaratorパース失敗")
    }

    if (*state).Offset >= len((*state).Tokens) {
        return  ident, cont, nil
    }

    token = (*state).Tokens[(*state).Offset]
    if token.Kind == TokenLeftBracket {
        if sizes, err := ArrayQualifiers(state); err != nil {
            return "", nil, err
        } else {
            cont0 := cont
            cont = func (baseType *CType) *CType {
                currentType := baseType
                for i := len(sizes) - 1; i >= 0; i-- {
                    currentType = Array(currentType, sizes[i])
                }
                return cont0(currentType)
            }
        }
    } else if token.Kind == TokenLeftParenthesis {
        if params, err := FuncParameters(state); err != nil {
            return "", nil, err
        } else {
            cont0 := cont
            cont = func (baseType *CType) *CType {
                return cont0(Function(baseType, params))
            }
        }
    }
    return ident, cont, nil
}

func ArrayQualifiers(state *ParserState) ([]int, error) {
    if (*state).Offset >= len((*state).Tokens) {
        return nil, errors.New("配列修飾部分パース失敗")
    }

    sizes := []int{}
    for {
        if ConsumeLeftBracket(state) {
            if size, consumed := ConsumeNum(state); consumed {
                if ConsumeRightBracket(state) {
                    sizes = append(sizes, size)
                } else {
                    return nil, errors.New("配列修飾部分パース失敗")
                }
            } else {
                return nil, errors.New("配列修飾部分パース失敗")
            }
        } else {
            break
        }
    }

    return sizes, nil
}

func FuncParameters(state *ParserState) ([]Parameter, error) {
    if !ConsumeLeftParenthesis(state) {
        return nil, errors.New("関数定義パース失敗")
    }

    params := []Parameter{}
    if !ConsumeRightParenthesis(state) {
        for {
            if paramName, t, err := Declaration(state); err != nil {
                return nil, errors.New("関数定義仮引数パース失敗")
            } else {
                params = append(params, Parameter{ paramName, t })
                if ConsumeRightParenthesis(state) {
                    break
                } else if !ConsumeComma(state) {
                    return nil, errors.New("関数定義仮引数の後にカンマがありません")
                }
            }
        }
    }
    return params, nil
}

func Declaration(state *ParserState) (string, *CType, error) {
    var baseType *CType
    if t, consumed := ConsumeType(state); !consumed {
        return "", nil, errors.New("Declarationパース失敗。型がありません")
    } else {
        baseType = t
    }

    for {
        if ConsumeOp(state, "*") {
            baseType = PointerTo(baseType)
        } else {
            break
        }
    }

    if ident, cont, err := Declarator(state); err != nil {
        return "", nil, err
    } else {
        return ident, cont(baseType), nil
    }
}

func DeclArray(state *ParserState, baseType *CType) (*CType, error) {
    if (*state).Offset >= len((*state).Tokens) {
        return nil, errors.New("DeclArrayパース失敗")
    }

    if ConsumeLeftBracket(state) {
        if size, consumed := ConsumeNum(state); consumed {
            if ConsumeRightBracket(state) {
                return Array(baseType, size), nil
            } else {
                return nil, errors.New("DeclArrayパース失敗。\"]\"がありません。")
            }
        } else {
            return nil, errors.New("DeclArrayパース失敗。配列サイズが数値ではありません。")
        }
    } else {
        return nil, errors.New("DeclArrayパース失敗。\"[\"がありません。")
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

    if SatisfyType(state) {
        if ident, t, err := Declaration(state); err != nil {
            return nil, err
        } else {
            if !ConsumeOp(state, ";") {
                return nil, errors.New("Stmtパース失敗。\";\"が不足しています。")
            }
            return NewNodeDecl(state, ident, t)
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
        return nil, err
    } else {
        if !ConsumeOp(state, ";") {
            return nil, errors.New("Returnパース失敗")
        }
        return NewNode(NodeReturn, e, nil), nil
    }
}

func Postfix(state *ParserState) (*Node, error) {
    if (*state).Offset >= len((*state).Tokens) {
        return nil, errors.New("Postfixパース失敗")
    }

    if prim, err0 := Primary(state); err0 != nil {
        return nil, err0
    } else {
        node := prim
        if ConsumeLeftBracket(state) {
            if expr, err1 := Expr(state); err1 != nil {
                return nil, err1
            } else {
                if !ConsumeRightBracket(state) {
                    return nil, errors.New("\"]\"が不足しています。")
                }
                node = NewNodeDeref(NewNodeAdd(ArrayToPointer(prim), ArrayToPointer(expr)))
            }
        }
        return node, nil
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
                return NewNodeFuncCall(state, ident, args), nil
            }
            for {
                if expr, err := Expr(state); err != nil {
                    return nil, err
                } else {
                    args = append(args, expr)
                    if ConsumeRightParenthesis(state) {
                        return NewNodeFuncCall(state, ident, args), nil
                    } else if !ConsumeComma(state) {
                        return nil, errors.New("関数呼び出し引数の後にカンマがありません")
                    }
                }
            }
        }
        return RefNode(state, ident)
    }

    if lit, consumed := ConsumeStringLiteral(state); consumed {
        return NewNodeStringLiteral(state, lit)
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
            n := NewNode(NodeSub, NewNodeNum(0), a)
            n.Type = Int()
            return n, nil
        }
    } else if ConsumeOp(state, "&") {
        if a, err := Primary(state); err != nil {
            return nil, err
        } else {
            return NewNodeAddr(a), nil
        }
    } else if ConsumeOp(state, "*") {
        if a, err := Primary(state); err != nil {
            return nil, err
        } else {
            a = ArrayToPointer(a)
            return NewNodeDeref(a), nil
        }
    } else if ConsumeSizeOf(state) {
        if a, err := Unary(state); err != nil {
            return nil, err
        } else {
            t := (*a).Type
            if t == nil {
                return nil, errors.New("sizeof引数の型が不明です")
            } else {
                return NewNodeNum(SizeOf(t)), nil
            }
        }
    }
    return Postfix(state)
}

func Mul(state *ParserState) (*Node, error) {
    var node *Node
    var t *CType
    if lhs, err := Unary(state); err != nil {
        return nil, err
    } else {
        node = lhs
        t = (*lhs).Type
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
    (*node).Type = t
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
                node = NewNodeAdd(ArrayToPointer(node), ArrayToPointer(rhs))
            }
        } else if ConsumeOp(state, "-") {
            if rhs, err := Mul(state); err != nil {
                return nil, err
            } else {
                node = NewNodeSub(ArrayToPointer(node), ArrayToPointer(rhs))
            }
        } else {
            break
        }
    }
    return node, nil
}

func Assign(state *ParserState) (*Node, error) {
    var node *Node
    var t *CType
    if lhs, err := Equality(state); err != nil {
        return nil, err
    } else {
        node = lhs
        t = (*lhs).Type
    }

    if ConsumeOp(state, "=") {
        if rhs, err := Assign(state); err != nil {
            return nil, err
        } else {
            node = NewNode(NodeAssign, node, ArrayToPointer(rhs))
        }
    }
    node.Type = t
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
                node = NewNode(NodeEq, ArrayToPointer(node), ArrayToPointer(rhs))
                // TODO boolある?
                (*node).Type = Int()
            }
        } else if ConsumeOp(state, "!=") {
            if rhs, err := Relational(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeNeq, ArrayToPointer(node), ArrayToPointer(rhs))
                // TODO boolある?
                (*node).Type = Int()
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
                // TODO boolある?
                (*node).Type = Int()
            }
        } else if ConsumeOp(state, "<=") {
            if rhs, err := Add(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeLe, node, rhs)
                // TODO boolある?
                (*node).Type = Int()
            }
        } else if ConsumeOp(state, ">") {
            if rhs, err := Add(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeGt, node, rhs)
                // TODO boolある?
                (*node).Type = Int()
            }
        } else if ConsumeOp(state, ">=") {
            if rhs, err := Add(state); err != nil {
                return nil, err
            } else {
                node = NewNode(NodeGe, node, rhs)
                // TODO boolある?
                (*node).Type = Int()
            }
        } else {
            break
        }
    }
    return node, nil
}
