usage: asm [-h]

Asm implements a simple x86-like processor.

Instructions are read from standard input.

The instructions behave similarly to their x86 counterparts.
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

  p: prints processor state