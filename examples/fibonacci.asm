jmp .main

.fib
push %fp
mov %sp %fp
push %d
mov 2(%fp) %d
cmp %d $0      # check if n = 0
je .fib_0
cmp %d $1      # check if n = 1
je .fib_1
jmp .fib_else  # handle n >= 2
.fib_0
mov $0 %a
jmp .fib_ret
.fib_1
mov $1 %a
jmp .fib_ret
.fib_else
push %d
dec (%sp)
call .fib      # first recursive call
mov %a %d
dec (%sp)
call .fib      # second recursive call
add %d %a
.fib_ret
mov -1(%fp) %d
leave
ret

# push argument and call .fib
.main
push $7
call .fib
p
