package main

import (
    "fmt"
)

func Gen(node *Node) {
    if ((*node).Kind == NodeNum) {
        fmt.Printf("  push %d\n", (*node).Val)
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
