package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

var ip int = 0
var instructions []string = []string{}
var regs [4]int
var mem [1000]int

func exec(s string) {
	if s == "p" {
		fmt.Println(regs)
	} else {
		operands := strings.Split(s, " ")
		src := getSrc(operands[1])
		pDst := getDst(operands[2])
		switch operands[0] {
		case "add":
			*pDst = *pDst + src
		case "sub":
			*pDst = *pDst - src
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
