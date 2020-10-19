package main

import (
    "fmt"
    "os"
)

type GenState struct {
    LabelCounter int
}

func GenProgram(nodes []*Node, localsLen int) {
    fmt.Printf(".intel_syntax noprefix\n")
    fmt.Printf(".globl main\n")
    fmt.Printf("main:\n")

    // プロローグ
    fmt.Printf("  push rbp\n")
    fmt.Printf("  mov rbp, rsp\n")
    fmt.Printf("  sub rsp, %d\n", localsLen * 8)

    state := GenState { 1 }
    for _, node := range(nodes) {
        Gen(node, &state)
    }

    // エピローグ
    fmt.Printf("  mov rsp, rbp\n")
    fmt.Printf("  pop rbp\n")
    fmt.Printf("  ret\n")
}

func GenLVarAddress(node *Node) {
    if (*node).Kind != NodeLVar {
        fmt.Fprintf(os.Stderr, "代入の左辺値が変数ではありません。\n")
        os.Exit(1)
    }

    fmt.Printf("  mov rax, rbp\n")
    fmt.Printf("  sub rax, %d\n", (*node).Offset)
}

func Gen(node *Node, state *GenState) {
    switch ((*node).Kind) {
    case NodeReturn:
        Gen((*node).Lhs, state)
        fmt.Printf("  mov rsp, rbp\n")
        fmt.Printf("  pop rbp\n")
        fmt.Printf("  ret\n")
        return
    case NodeNum:
        fmt.Printf("  mov rax, %d\n", (*node).Val)
        return
    case NodeLVar:
        GenLVarAddress(node)
        fmt.Printf("  mov rax, [rax]\n")
        return
    case NodeAssign:
        GenLVarAddress((*node).Lhs)
        fmt.Printf("  push rax\n")
        Gen((*node).Rhs, state)
        fmt.Printf("  mov rdi, rax\n")
        fmt.Printf("  pop rax\n")
        fmt.Printf("  mov [rax], rdi\n")
        fmt.Printf("  mov rax, rdi\n")
        return
    case NodeIf:
        Gen((*node).Lhs, state)
        fmt.Printf("  cmp rax, 0\n")
        rhs := (*node).Rhs
        label := (*state).LabelCounter
        if (*rhs).Kind == NodeEither {
            fmt.Printf("  je .Lelse%d\n", label)
            Gen((*rhs).Lhs, state)
            fmt.Printf("  jmp .Lend%d\n", label)
            fmt.Printf(".Lelse%d:\n", label)
            Gen((*rhs).Rhs, state)
            fmt.Printf(".Lend%d:\n", label)
        } else {
            fmt.Printf("  je .Lend%d\n", label)
            Gen((*node).Rhs, state)
            fmt.Printf(".Lend%d:\n", label)
        }
        (*state).LabelCounter++
        return
    case NodeWhile:
        label := (*state).LabelCounter
        fmt.Printf(".Lstart%d:\n", label)
        Gen((*node).Lhs, state)
        fmt.Printf("  cmp rax, 0\n")
        fmt.Printf("  je .Lend%d\n", label)
        Gen((*node).Rhs, state)
        fmt.Printf("  jmp .Lstart%d\n", label)
        fmt.Printf(".Lend%d:\n", label)
        (*state).LabelCounter++
        return
    }

    Gen((*node).Lhs, state)
    fmt.Printf("  push rax\n")
    Gen((*node).Rhs, state)

    fmt.Printf("  mov rdi, rax\n")
    fmt.Printf("  pop rax\n")

    switch((*node).Kind) {
    case NodeAdd:
        fmt.Printf("  add rax, rdi\n")
    case NodeSub:
        fmt.Printf("  sub rax, rdi\n")
    case NodeMul:
        fmt.Printf("  imul rax, rdi\n")
    case NodeDiv:
        fmt.Printf("  cqo\n")
        fmt.Printf("  idiv rdi\n")
    case NodeEq:
        fmt.Printf("cmp rax, rdi\n")
        fmt.Printf("sete al\n")
        fmt.Printf("movzb rax, al\n")
    case NodeNeq:
        fmt.Printf("cmp rax, rdi\n")
        fmt.Printf("setne al\n")
        fmt.Printf("movzb rax, al\n")
    case NodeLt:
        fmt.Printf("cmp rax, rdi\n")
        fmt.Printf("setl al\n")
        fmt.Printf("movzb rax, al\n")
    case NodeLe:
        fmt.Printf("cmp rax, rdi\n")
        fmt.Printf("setle al\n")
        fmt.Printf("movzb rax, al\n")
    case NodeGt:
        fmt.Printf("cmp rdi, rax\n")
        fmt.Printf("setl al\n")
        fmt.Printf("movzb rax, al\n")
    case NodeGe:
        fmt.Printf("cmp rdi, rax\n")
        fmt.Printf("setle al\n")
        fmt.Printf("movzb rax, al\n")
    }

    return
}
