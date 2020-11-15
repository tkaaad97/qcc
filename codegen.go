package main

import (
    "fmt"
    "os"
)

type GenState struct {
    LabelCounter int
}

func Movs(dst AsmLocation, src AsmValue) {
    if dst.AsmLocationDataType() == src.AsmValueDataType() {
        fmt.Printf("  mov %s, %s\n", dst.ShowAsmLocation(), src.ShowAsmValue())
    } else {
        fmt.Printf("  movsx %s, %s\n", dst.ShowAsmLocation(), src.ShowAsmValue())
    }
}

func GenProgram(globals []*Node, stringLiterals []string, defs []NodeAndLocalSize) {
    fmt.Printf("  .intel_syntax noprefix\n")

    if len(globals) > 0 || len(stringLiterals) > 0 {
        fmt.Printf("  .data\n")
        for _, node := range(globals) {
            GenGVar(node)
        }

        for i, lit := range(stringLiterals) {
            GenStringLiteralData(i, lit)
        }
    }

    fmt.Printf("  .text\n")
    fmt.Printf("  .globl main\n")

    state := GenState { 1 }
    for _, def := range(defs) {
        GenDef(def.Node, def.LocalSize, &state)
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
        if (node.Offset != 0) {
            fmt.Printf("  sub rax, %d\n", node.Offset)
        }
        fmt.Printf("  mov [rax], %s\n", registers[i])
        i++
        node = node.Lhs
    }
}

func GenGVar(node *Node) {
    size := SizeOf(node.Type)
    fmt.Printf("%s:\n", node.Ident)
    fmt.Printf("  .zero %d\n", size)
}

func StringLiteralLabel(off int) string {
    return fmt.Sprintf(".LC%d", off)
}

func GenStringLiteralData(i int, lit string) {
    fmt.Printf("%s:\n", StringLiteralLabel(i))
    fmt.Printf("  .string %q\n", lit)
}

func GenDef(node *Node, localSize int, state *GenState) {
    if (node == nil || node.Kind != NodeFuncDef) {
        fmt.Fprintf(os.Stderr, "関数定義のノードではありません\n")
        os.Exit(1)
    }

    fmt.Printf("%s:\n", node.Ident)

    // プロローグ
    fmt.Printf("  push rbp\n")
    fmt.Printf("  mov rbp, rsp\n")
    if localSize != 0 {
        fmt.Printf("  sub rsp, %d\n", localSize)
    }

    // 引数をレジスタからスタックに移動
    GenParams(node.Lhs)

    Gen(node.Rhs, state)

    // エピローグ
    fmt.Printf("  mov rsp, rbp\n")
    fmt.Printf("  pop rbp\n")
    fmt.Printf("  ret\n")
}

func GenVarAddress(node *Node, state *GenState) {
    if node.Kind == NodeLVar {
        fmt.Printf("  mov rax, rbp\n")
        if node.Offset != 0 {
            fmt.Printf("  sub rax, %d\n", node.Offset)
        }
    } else if node.Kind == NodeGVar {
        fmt.Printf("  lea rax, %s[rip]\n", node.Ident)
    } else if node.Kind == NodeDeref {
        Gen(node.Lhs, state)
    } else {
        fmt.Fprintf(os.Stderr, "代入の左辺値が変数ではありません。\n")
        os.Exit(1)
    }
}

func GenPushVarAddress(node *Node, state *GenState) {
    if node.Kind == NodeLVar {
        fmt.Printf("  mov rax, rbp\n")
        if node.Offset != 0 {
            fmt.Printf("  sub rax, %d\n", node.Offset)
        }
        fmt.Printf("  push rax\n")
    } else if node.Kind == NodeGVar {
        fmt.Printf("  lea rax, %s[rip]\n", node.Ident)
        fmt.Printf("  push rax\n")
    } else if node.Kind == NodeDeref {
        GenExpr(node.Lhs, state)
    } else {
        fmt.Fprintf(os.Stderr, "代入の左辺値が変数ではありません。\n")
        os.Exit(1)
    }
}

func Gen(node *Node, state *GenState) {
    switch (node.Kind) {
    case NodeReturn:
        Gen(node.Lhs, state)
        fmt.Printf("  mov rsp, rbp\n")
        fmt.Printf("  pop rbp\n")
        fmt.Printf("  ret\n")
        return
    case NodeNum:
        fmt.Printf("  mov rax, %d\n", node.Val)
        return
    case NodeLVar:
        GenVarAddress(node, state)
        dst := ResolveDstRegisterByType(int(Rax), node.Type)
        src := AsmDeref{ Rax , CTypeToAsmDataType(node.Type) }
        Movs(dst, src)
        return
    case NodeGVar:
        GenVarAddress(node, state)
        dst := ResolveDstRegisterByType(int(Rax), node.Type)
        src := AsmDeref{ Rax , CTypeToAsmDataType(node.Type) }
        Movs(dst, src)
        return
    case NodeStringLiteral:
        fmt.Printf("  lea rax, %s[rip]\n", StringLiteralLabel(node.Offset))
        return
    case NodeAssign:
        GenPushVarAddress(node.Lhs, state)
        Gen(node.Rhs, state)
        fmt.Printf("  pop rdi\n")
        dst := AsmDeref{ Rdi , CTypeToAsmDataType(node.Type) }
        src := ResolveRegisterByType(0, node.Type).AsmLocationToValue()
        Movs(dst, src)
        return
    case NodeAddr:
        GenVarAddress(node.Lhs, state)
        return
    case NodeDeref:
        Gen(node.Lhs, state)
        dst := ResolveDstRegisterByType(int(Rax), node.Type)
        src := AsmDeref{ Rax , CTypeToAsmDataType(node.Type) }
        Movs(dst, src)
        return
    case NodeBlock:
        if node.Lhs == nil {
            return
        }
        Gen(node.Lhs, state)
        current := node.Rhs
        for {
            if current == nil {
                break
            }
            Gen(current.Lhs, state)
            current = current.Rhs
        }
        return
    case NodeIf:
        Gen(node.Lhs, state)
        fmt.Printf("  cmp rax, 0\n")
        rhs := node.Rhs
        label := state.LabelCounter
        if rhs.Kind == NodeEither {
            fmt.Printf("  je .Lelse%d\n", label)
            Gen(rhs.Lhs, state)
            fmt.Printf("  jmp .Lend%d\n", label)
            fmt.Printf(".Lelse%d:\n", label)
            Gen(rhs.Rhs, state)
            fmt.Printf(".Lend%d:\n", label)
        } else {
            fmt.Printf("  je .Lend%d\n", label)
            Gen(node.Rhs, state)
            fmt.Printf(".Lend%d:\n", label)
        }
        state.LabelCounter++
        return
    case NodeFor:
        label := state.LabelCounter
        first := node.Lhs
        second := node.Rhs
        pre := first.Lhs
        cond := first.Rhs
        post := second.Lhs
        action := second.Rhs
        Gen(pre, state)
        fmt.Printf(".Lstart%d:\n", label)
        Gen(cond, state)
        fmt.Printf("cmp rax, 0\n")
        fmt.Printf("je .Lend%d\n", label)
        Gen(action, state)
        Gen(post, state)
        fmt.Printf("jmp .Lstart%d\n", label)
        fmt.Printf(".Lend%d:\n", label)
        state.LabelCounter++
        return
    case NodeWhile:
        label := state.LabelCounter
        fmt.Printf(".Lstart%d:\n", label)
        Gen(node.Lhs, state)
        fmt.Printf("  cmp rax, 0\n")
        fmt.Printf("  je .Lend%d\n", label)
        Gen(node.Rhs, state)
        fmt.Printf("  jmp .Lstart%d\n", label)
        fmt.Printf(".Lend%d:\n", label)
        state.LabelCounter++
        return
    case NodeFuncCall:
        funcName := node.Ident
        arg := node.Lhs
        argNum := 0
        argRegisters := []string{ "rdi", "rsi", "rdx", "rcx", "r8", "r9", }
        for {
            if arg == nil {
                break
            }
            Gen(arg.Lhs, state)
            fmt.Printf("  push rax\n");
            arg = arg.Rhs
            argNum++
        }
        for i := argNum - 1; i >= 0; i-- {
            fmt.Printf("  pop rax\n");
            fmt.Printf("  mov %s, rax\n", argRegisters[i]);
        }
        fmt.Printf("  mov al, 0\n")
        fmt.Printf("  call %s\n", funcName)
        return
    case NodeCastIntegral:
        Gen(node.Lhs, state)
        dst := ResolveRegisterByType(int(Rax), node.Type)
        src := ResolveRegisterByType(int(Rax), node.Lhs.Type).AsmLocationToValue()
        Movs(dst, src)
        return
    }

    GenBinOp(node, state)
}

func GenExpr(node *Node, state *GenState) {
    if (!IsExpr(node)) {
        fmt.Fprintf(os.Stderr, "式ノードではありません\n")
        os.Exit(1)
    }

    switch (node.Kind) {
    case NodeNum:
        fmt.Printf("  push %d\n", node.Val)
        return
    case NodeLVar:
        Gen(node, state)
        fmt.Printf("  push rax\n")
        return
    case NodeGVar:
        Gen(node, state)
        fmt.Printf("  push rax\n")
        return
    case NodeStringLiteral:
        fmt.Printf("  lea rax, %s[rip]\n", StringLiteralLabel(node.Offset))
        fmt.Printf("  push rax\n")
    case NodeAssign:
        GenPushVarAddress(node.Lhs, state)
        Gen(node.Rhs, state)
        fmt.Printf("  pop rdi\n")
        dst := AsmDeref{ Rdi , CTypeToAsmDataType(node.Type) }
        src := ResolveRegisterByType(0, node.Type).AsmLocationToValue()
        Movs(dst, src)
        if src != Rax {
            Movs(Rax, src)
        }
        return
    case NodeAddr:
        GenPushVarAddress(node.Lhs, state)
        return
    case NodeDeref:
        Gen(node.Lhs, state)
        fmt.Printf("  push [rax]\n")
        return
    case NodeFuncCall:
        Gen(node, state)
        fmt.Printf("  push rax\n")
        return
    case NodeCastIntegral:
        Gen(node, state)
        fmt.Printf("  push rax\n")
        return
    }

    GenBinOp(node, state)
    fmt.Printf("  push rax\n")
}

func GenBinOp(node *Node, state *GenState) {
    lhs := node.Lhs
    rhs := node.Rhs
    GenExpr(lhs, state)
    GenExpr(rhs, state)

    fmt.Printf("  pop rdi\n")
    fmt.Printf("  pop rax\n")

    cmpa := ResolveRegisterByType(int(Rax), lhs.Type).ShowAsmLocation()
    cmpb := ResolveRegisterByType(int(Rdi), lhs.Type).ShowAsmLocation()

    switch(node.Kind) {
    case NodeAdd:
        if (lhs.Type != nil && lhs.Type.Kind == CTypePointer) {
            size := SizeOf(lhs.Type.PointerTo)
            fmt.Printf("  imul rdi, %d\n", size)
        }
        fmt.Printf("  add rax, rdi\n")
        return
    case NodeSub:
        if (lhs.Type != nil && lhs.Type.Kind == CTypePointer) {
            size := SizeOf(lhs.Type.PointerTo)
            fmt.Printf("  imul rdi, %d\n", size)
        }
        fmt.Printf("  sub rax, rdi\n")
        return
    case NodeMul:
        fmt.Printf("  imul rax, rdi\n")
        return
    case NodeDiv:
        fmt.Printf("  cqo\n")
        fmt.Printf("  idiv rdi\n")
        return
    case NodeEq:
        fmt.Printf("cmp %s, %s\n", cmpa, cmpb)
        fmt.Printf("sete al\n")
        fmt.Printf("movzb rax, al\n")
        return
    case NodeNeq:
        fmt.Printf("cmp %s, %s\n", cmpa, cmpb)
        fmt.Printf("setne al\n")
        fmt.Printf("movzb rax, al\n")
        return
    case NodeLt:
        fmt.Printf("cmp %s, %s\n", cmpa, cmpb)
        fmt.Printf("setl al\n")
        fmt.Printf("movzb rax, al\n")
        return
    case NodeLe:
        fmt.Printf("cmp %s, %s\n", cmpa, cmpb)
        fmt.Printf("setle al\n")
        fmt.Printf("movzb rax, al\n")
        return
    case NodeGt:
        fmt.Printf("cmp %s, %s\n", cmpb, cmpa)
        fmt.Printf("setl al\n")
        fmt.Printf("movzb rax, al\n")
        return
    case NodeGe:
        fmt.Printf("cmp %s, %s\n", cmpb, cmpa)
        fmt.Printf("setle al\n")
        fmt.Printf("movzb rax, al\n")
        return
    }

    fmt.Fprintf(os.Stderr, "二項演算ではありません\n")
    os.Exit(1)
}
