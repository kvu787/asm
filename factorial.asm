# not tested
jmp .main
.factorial
push %fp
mov %sp %fp
mov $1 %a
mov $0 %b
.factorial-loop
cmp 1(%fp) %b
je .factorial-return
add $1 %b
mul %b %a
jmp .factorial-loop
.factorial-return
mov %fp %sp
pop %fp
ret
.main
push $3
call .factorial
mov %a %d
mov $4 %sp
call .factorial
add %a %d
mov $5 %sp
call .factorial
add %a %d
p
