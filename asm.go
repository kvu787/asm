package main

import (
	"bufio"
	"fmt"
	"strconv"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		exec(scanner.Text())
	}
}

var registers [4]int
var memory [1000]int

func exec(s string) {
	if s == "p" {
		fmt.Println(registers)
	} else {
		operands := strings.Split(s, " ")
		src, _ := strconv.Atoi(operands[1])
		pDst := nameToRegister(operands[2])
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

func nameToRegister(name string) *int {
	switch name {
	case "a":
		return &registers[0]
	case "b":
		return &registers[1]
	case "c":
		return &registers[2]
	case "d":
		return &registers[3]
	default:
		panic("invalid register name")
	}
}

