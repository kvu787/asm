usage: asm [-h]

Asm implements a simple x86-like processor.

Instructions are read from standard input and executed.

INSTRUCTIONS

  The instructions are similar to their x86 counterparts.
  The following instructions are supported:

    mov
    push
    pop

    inc
    dec
    add
    sub
    mul

    cmp
    jmp
    je, jz
    jne, jnz
    jg
    jge
    jl
    jle

    call
    leave
    ret

    p (prints processor state)
    exit (terminates this program)

REGISTERS

  6 general purpose registers are available:

    %a, %b, %c, %d, %e, %f

  Special registers:

    %ip: stores index of next instruction
    %fp: frame pointer
    %sp: stack pointer

CALLING PROCEDURE

  Stack usage is identical to x86.

  Return register:        %a
  Caller-saved registers: %a, %b, %c
  Callee-saved registers: %d, %e, %f

MEMORY ADDRESSING MODES

  imm      = mem[imm]
  (reg)    = mem[reg]
  num(reg) = mem[num + reg]