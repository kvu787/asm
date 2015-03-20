usage: asm [-h]

Asm implements a simple x86-like processor.

Instructions are read from standard input.

INSTRUCTIONS

The instructions are similar to their x86 counterparts.
The following instructions are supported:

  mov  src dst 
  push src
  pop  dst

  inc dst
  dec dst
  add src dst
  sub src dst
  mul src dst

  cmp src1 src2
  jmp label
  je, jz label
  jne, jnz label
  jg label
  jge label
  jl label
  jle label

  call label
  leave
  ret

  p: prints processor state

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