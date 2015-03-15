jmp .main
.factorial
push %fp
mov %sp %fp
mov $1 %a
mov $0 %b
.factorial-loop
cmp 2(%fp) %b
je .factorial-return
add $1 %b
mul %b %a
jmp .factorial-loop
.factorial-return
leave
ret
.main
push $6
call .factorial
p
