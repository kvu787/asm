## MAIN METHOD
.main
  # push argument and call .fib
  push $7
  call .fib
  p
  exit

## RECURSIVE FIBONACCI SUBROUTINE
.fib
  # standard x86 calling procedure
  push %fp       # save old frame pointer
  mov %sp %fp    # set current frame pointer
  
  # read arguments
  push %d        # save 'callee-saved' register
  mov 2(%fp) %d  # move argument to %d
  
  # if else jump structure
  cmp %d $0      # if n = 0
  je .fib_0
  cmp %d $1      # else if n = 1
  je .fib_1
  jmp .fib_else  # else (if n >= 2)

# if n = 0
.fib_0
  mov $0 %a      # set return value
  jmp .fib_ret

# if n = 1
.fib_1
  mov $1 %a      # set return value
  jmp .fib_ret

# if n >= 2
.fib_else
  push %d
  dec (%sp)
  call .fib      # first recursive call
  mov %a %d
  dec (%sp)
  call .fib      # second recursive call
  add %d %a

.fib_ret
  # standard x86 return procedure
  mov -1(%fp) %d # restore 'callee-saved' register
  leave          # restore stack and frame pointers
  ret            # return to caller's next instruction
