package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ip int = 0
var instructions []string = []string{}
var regs [4]int
var mem [1000]int
var zf bool
var sf bool

func main() {
	// load instructions
	for s := bufio.NewScanner(os.Stdin); s.Scan(); {
		instructions = append(instructions, s.Text())
	}

	// execute
	for ip < len(instructions) {
		instruction := instructions[ip]
		ip++
		exec(instruction)
	}
	os.Exit(0)
}

func exec(s string) {
	operands := strings.Split(s, " ")
	switch operands[0] {
	case "p":
		var nextInstruction string
		if ip >= len(instructions) {
			nextInstruction = "no more instructions"
		} else {
			nextInstruction = instructions[ip]
		}
		fmt.Printf("a: %d\n"+
			"b: %d\n"+
			"c: %d\n"+
			"d: %d\n"+
			"ip: %d\n"+
			"next instruction: %v\n",
			regs[0], regs[1], regs[2], regs[3], ip, nextInstruction)
	case "cmp":
		m := getSrc(operands[1])
		n := getSrc(operands[2])
		result := m - n
		zf = result == 0
		sf = result < 0
	case "jmp",
		"je", "jne", "jz", "jnz",
		"jg", "jge", "jl", "jle":
		addr := getSrc(operands[1])
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
	default:
		src := getSrc(operands[1])
		pDst := getDst(operands[2])
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
	}
}

func getSrc(name string) int {
	if name[0] == '%' {
		return *getDst(name)
	} else if name[0] == '$' {
		i, _ := strconv.Atoi(name[1:])
		return i
	} else {
		i, _ := strconv.Atoi(name)
		return mem[i]
	}
}

func getDst(name string) *int {
	if name[0] == '%' {
		switch name[1:] {
		case "a":
			return &regs[0]
		case "b":
			return &regs[1]
		case "c":
			return &regs[2]
		case "d":
			return &regs[3]
		default:
			panic("invalid register name")
		}
	} else {
		i, _ := strconv.Atoi(name)
		return &mem[i]
	}
}
