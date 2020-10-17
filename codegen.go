package main

import (
    "fmt"
    "os"
)

func GenProgram(nodes []*Node) {
    fmt.Printf(".intel_syntax noprefix\n")
    fmt.Printf(".globl main\n")
    fmt.Printf("main:\n")

    // プロローグ
    fmt.Printf("  push rbp\n")
    fmt.Printf("  mov rbp, rsp\n")
    fmt.Printf("  sub rsp, 208\n")

    for _, node := range(nodes) {
        Gen(node)
        fmt.Printf("  pop rax\n")
    }

    // エピローグ
    fmt.Printf("  mov rsp, rbp\n")
    fmt.Printf("  pop rbp\n")
    fmt.Printf("  ret\n")
}

func GenLVar(node *Node) {
    if (*node).Kind != NodeLVar {
        fmt.Fprintf(os.Stderr, "代入の左辺値が変数ではありません。\n")
        os.Exit(1)
    }

    fmt.Printf("  mov rax, rbp\n")
    fmt.Printf("  sub rax, %d\n", (*node).Offset)
    fmt.Printf("  push rax\n")
}

func Gen(node *Node) {
    switch ((*node).Kind) {
    case NodeNum:
        fmt.Printf("  push %d\n", (*node).Val)
        return
    case NodeLVar:
        GenLVar(node)
        fmt.Printf("  pop rax\n")
        fmt.Printf("  mov rax, [rax]\n")
        fmt.Printf("  push rax\n")
        return
    case NodeAssign:
        GenLVar((*node).Lhs)
        Gen((*node).Rhs)

        fmt.Printf("  pop rdi\n")
        fmt.Printf("  pop rax\n")
        fmt.Printf("  mov [rax], rdi\n")
        fmt.Printf("  push rdi\n")
        return
    }

    Gen((*node).Lhs)
    Gen((*node).Rhs)

    fmt.Printf("  pop rdi\n")
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

    fmt.Printf("  push rax\n")
}
