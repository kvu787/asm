jmp .main
.factorial
push %fp
mov %sp %fp
mov $1 %a
mov $0 %b
.factorial_loop
cmp 2(%fp) %b
je .factorial_return
add $1 %b
mul %b %a
jmp .factorial_loop
.factorial_return
leave
ret
.main
push $6
call .factorial
p
