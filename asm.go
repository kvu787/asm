package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var help string = `usage: asm [-h]

Asm implements a simple x86-like processor.

Instructions are read from standard input.

The instructions behave similarly to their x86 counterparts.
The following instructions are supported:

	mov
	push
	pop

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
	ret`

// global vars
var (
	ip           int64
	instructions []string
	regs         [8]int64
	mem          [1000]int64
	zf           bool
	sf           bool
	labels       map[string]int64
)

// operand formats
var (
	immre  = regexp.MustCompile(`^[$]\d+$`)
	regre  = regexp.MustCompile(`^%\w+$`)
	mem1re = regexp.MustCompile(`^\d+$`)
	mem2re = regexp.MustCompile(`^[(]%\w+[)]$`)
	mem3re = regexp.MustCompile(`^([-]?\d+)[(](%\w+)[)]$`)
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-h" {
		fmt.Println(help)
		os.Exit(0)
	}

	// init global vars
	instructions = make([]string, 0)
	labels = make(map[string]int64)

	// init stack and frame pointers
	*getPtr("%sp") = int64(len(mem) - 1)
	*getPtr("%fp") = int64(len(mem) - 1)

	// load instructions from standard input
	i := int64(0)
	for s := bufio.NewScanner(os.Stdin); s.Scan(); {
		text := s.Text()
		if text[0] == '.' {
			labels[text] = i
		} else {
			instructions = append(instructions, s.Text())
			i++
		}
	}

	// execute
	for ip < int64(len(instructions)) {
		instruction := instructions[ip]
		ip++
		Exec(instruction)
	}

	os.Exit(0)
}

// Exec runs an assembly instruction. It may change registers, change memory,
// and/or set flags.
func Exec(s string) {
	operands := strings.Split(s, " ")
	switch operands[0] {
	case "p":
		var nextInstruction string
		if ip >= int64(len(instructions)) {
			nextInstruction = "no more instructions"
		} else {
			nextInstruction = instructions[ip]
		}
		fmt.Printf(
			"a: %d\n"+
				"b: %d\n"+
				"c: %d\n"+
				"d: %d\n"+
				"e: %d\n"+
				"f: %d\n"+
				"sp: %d\n"+
				"fp: %d\n"+
				"ip: %d\n"+
				"next instruction: %v\n",
			regs[0], regs[1], regs[2], regs[3], regs[4], regs[5], regs[6], regs[7], ip, nextInstruction)
	case "cmp":
		n := getVal(operands[1])
		m := getVal(operands[2])
		result := m - n
		zf = result == 0
		sf = result < 0
	case "jmp", "je", "jne", "jz", "jnz", "jg", "jge", "jl", "jle":
		label := operands[1]
		addr := labels[label]
		switch operands[0] {
		case "jmp":
			ip = addr
		case "je", "jz":
			if zf {
				ip = addr
			}
		case "jne", "jnz":
			if !zf {
				ip = addr
			}
		case "jg":
			if !zf && !sf {
				ip = addr
			}
		case "jge":
			if zf || !sf {
				ip = addr
			}
		case "jl":
			if !zf && sf {
				ip = addr
			}
		case "jle":
			if zf || sf {
				ip = addr
			}
		}
	case "add", "sub", "mul", "mov":
		src := getVal(operands[1])
		pDst := getPtr(operands[2])
		switch operands[0] {
		case "add":
			*pDst = *pDst + src
		case "sub":
			*pDst = *pDst - src
		case "mul":
			*pDst = *pDst * src
		case "mov":
			*pDst = src
		}
	case "pop":
		pDst := getPtr(operands[1])
		*pDst = mem[getVal("%sp")]
		*getPtr("%sp") = getVal("%sp") + 1
	case "push":
		*getPtr("%sp") = getVal("%sp") - 1
		src := getVal(operands[1])
		mem[getVal("%sp")] = src
	case "call":
		Exec("push %ip")
		label := operands[1]
		Exec(fmt.Sprintf("jmp %s", label))
	case "leave":
		Exec("mov %fp %sp")
		Exec("pop %fp")
	case "ret":
		Exec("pop %ip")
	default:
		panic("Exec: unrecognized instruction")
	}
}

// getVal returns the value of the operand.
// Operand is an immediate, memory, or register name.
func getVal(operand string) int64 {
	if operand[0] == '$' {
		i, _ := strconv.ParseInt(operand[1:], 10, 64)
		return i
	} else {
		return *getPtr(operand)
	}
}

// getPtr returns a pointer to the memory location of the operand
// Operand is a memory or register name.
func getPtr(operand string) *int64 {
	switch true {
	case regre.MatchString(operand):
		switch operand[1:] {
		case "a":
			return &regs[0]
		case "b":
			return &regs[1]
		case "c":
			return &regs[2]
		case "d":
			return &regs[3]
		case "e":
			return &regs[4]
		case "f":
			return &regs[5]
		case "sp":
			return &regs[6]
		case "fp":
			return &regs[7]
		case "ip":
			return &ip
		default:
			panic("getPtr: invalid register name: " + operand[1:])
		}
	case mem1re.MatchString(operand):
		i, _ := strconv.ParseInt(operand, 10, 64)
		return &mem[i]
	case mem2re.MatchString(operand):
		reg := operand[1 : len(operand)-1]
		regval := *getPtr(reg)
		return &mem[regval]
	case mem3re.MatchString(operand):
		matches := mem3re.FindStringSubmatch(operand)
		reg := matches[2]
		regval := *getPtr(reg)
		intval, _ := strconv.ParseInt(matches[1], 10, 64)
		return &mem[regval+intval]
	default:
		panic("getPtr: unrecognized operand")
	}
}
