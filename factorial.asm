mov $1 %a
mov $1 %b
.loop
mul %b %a
add $1 %b
cmp $10 %b
jne .loop
p
