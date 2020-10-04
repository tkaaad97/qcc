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
    TokenNum
    TokenEof
)

type Token struct {
    kind TokenKind
    val int
    str string
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

        if (input[off] == '+' || input[off] == '-') {
            token := Token {
                kind: TokenReserved,
                str: string([]rune{input[off]}),
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
    }
    if result, err := strconv.Atoi(str); err != nil {
        return Token{}, offset, errors.New("parseNum失敗")
    } else {
        token.val = result
    }

    return token, a, nil
}

func consumeOp(tokens []Token, offset *int, op string) bool {
    token := tokens[*offset]
    if token.kind == TokenReserved && token.str == op {
        (*offset)++
        return true
    }
    return false
}

func consumeNum(tokens []Token, offset *int) (int, bool) {
    token := tokens[*offset]
    if token.kind == TokenNum {
        (*offset)++
        return token.val, true
    }
    return 0, false
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

    fmt.Printf(".intel_syntax noprefix\n")
    fmt.Printf(".globl main\n")
    fmt.Printf("main:\n")

    // 最初の数
    tl := len(tokens)
    offset := 0
    if a0, consumed := consumeNum(tokens, &offset); consumed  {
        fmt.Printf("  mov rax, %d\n", a0)
    } else {
        fmt.Fprintf(os.Stderr, "最初のトークンが数ではありません。")
        os.Exit(1)
    }

    for {
        if (offset >= tl) {
            break
        }

        if consumeOp(tokens, &offset, "+") {
            if a, consumed := consumeNum(tokens, &offset); consumed {
                fmt.Printf("  add rax, %d\n", a)
            } else {
                fmt.Fprintf(os.Stderr, "+の後のトークンが数ではありません。 str: %s", tokens[offset].str)
                os.Exit(1)
            }
        } else if consumeOp(tokens, &offset, "-") {
            if a, consumed := consumeNum(tokens, &offset); consumed {
                fmt.Printf("  sub rax, %d\n", a)
            } else {
                fmt.Fprintf(os.Stderr, "-の後のトークンが数ではありません。 str: %s", tokens[offset].str)
                os.Exit(1)
            }
        } else {
            fmt.Fprintf(os.Stderr, "演算子があるべきところで別トークン str: %s", tokens[offset].str)
            os.Exit(1)
        }
    }

    fmt.Printf("  ret\n")
}
