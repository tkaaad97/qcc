package main

import (
    "fmt"
    "os"
)

type GenState struct {
    LabelCounter int
}

func GenProgram(defs []NodeAndLocals) {
    fmt.Printf(".intel_syntax noprefix\n")
    fmt.Printf(".globl main\n")

    state := GenState { 1 }
    for _, def := range(defs) {
        GenDef(def.Node, len(def.Locals), &state)
    }
}

func GenParams(node *Node) {
    registers := []string{ "rdi", "rsi", "rdx", "rcx", "r8", "r9", }
    i := 0
    for {
        if node == nil {
            break
        }
        fmt.Printf("  mov rax, rbp\n")
        fmt.Printf("  sub rax, %d\n", (*node).Offset)
        fmt.Printf("  mov [rax], %s\n", registers[i])
        i++
        node = (*node).Lhs
    }
}

func GenDef(node *Node, localsLen int, state *GenState) {
    if (node == nil || (*node).Kind != NodeFuncDef) {
        fmt.Fprintf(os.Stderr, "関数定義のノードではありません\n")
        os.Exit(1)
    }

    fmt.Printf("%s:\n", (*node).Ident)

    // プロローグ
    fmt.Printf("  push rbp\n")
    fmt.Printf("  mov rbp, rsp\n")
    fmt.Printf("  sub rsp, %d\n", localsLen * 8)

    // 引数をレジスタからスタックに移動
    GenParams((*node).Lhs)

    Gen((*node).Rhs, state)

    // エピローグ
    fmt.Printf("  mov rsp, rbp\n")
    fmt.Printf("  pop rbp\n")
    fmt.Printf("  ret\n")
}

func GenLVarAddress(node *Node, state *GenState) {
    if (*node).Kind == NodeLVar {
        fmt.Printf("  mov rax, rbp\n")
        fmt.Printf("  sub rax, %d\n", (*node).Offset)
    } else if (*node).Kind == NodeDeref {
        Gen((*node).Lhs, state)
    } else {
        fmt.Fprintf(os.Stderr, "代入の左辺値が変数ではありません。\n",)
        os.Exit(1)
    }
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
        GenLVarAddress(node, state)
        fmt.Printf("  mov rax, [rax]\n")
        return
    case NodeAssign:
        GenLVarAddress((*node).Lhs, state)
        fmt.Printf("  push rax\n")
        Gen((*node).Rhs, state)
        fmt.Printf("  mov rdi, rax\n")
        fmt.Printf("  pop rax\n")
        fmt.Printf("  mov [rax], rdi\n")
        fmt.Printf("  mov rax, rdi\n")
        return
    case NodeAddr:
        GenLVarAddress((*node).Lhs, state)
        return
    case NodeDeref:
        GenLVarAddress((*node).Lhs, state)
        fmt.Printf("  mov rax, [rax]\n")
        return
    case NodeBlock:
        if (*node).Lhs == nil {
            return
        }
        Gen((*node).Lhs, state)
        current := (*node).Rhs
        for {
            if current == nil {
                return
            }
            Gen((*current).Lhs, state)
            current = (*current).Rhs
        }
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
    case NodeFor:
        label := (*state).LabelCounter
        first := (*node).Lhs
        second := (*node).Rhs
        pre := (*first).Lhs
        cond := (*first).Rhs
        post := (*second).Lhs
        action := (*second).Rhs
        Gen(pre, state)
        fmt.Printf(".Lstart%d:\n", label)
        Gen(cond, state)
        fmt.Printf("cmp rax, 0\n")
        fmt.Printf("je .Lend%d\n", label)
        Gen(action, state)
        Gen(post, state)
        fmt.Printf("jmp .Lstart%d\n", label)
        fmt.Printf(".Lend%d:\n", label)
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
    case NodeFuncCall:
        funcName := (*node).Ident
        arg := (*node).Lhs
        argNum := 0
        argRegisters := []string{ "rdi", "rsi", "rdx", "rcx", "r8", "r9", }
        for {
            if arg == nil {
                break
            }
            Gen((*arg).Lhs, state)
            fmt.Printf("  mov %s, rax\n", argRegisters[argNum])
            arg = (*arg).Rhs
            argNum++
        }
        fmt.Printf("  call %s\n", funcName)
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
