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

	inc
	dec
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
	ret

	p: prints processor state`

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
	setVal("%sp", int64(len(mem)-1))
	setVal("%fp", int64(len(mem)-1))

	// load instructions from standard input
	i := int64(0)
	for s := bufio.NewScanner(os.Stdin); s.Scan(); {
		text := s.Text()
		if len(text) == 0 {
			continue
		} else if text[0] == '.' {
			labels[text] = i
		} else {
			instructions = append(instructions, s.Text())
			i++
		}
	}

	// execute instructions
	for ip < int64(len(instructions)) {
		instruction := instructions[ip]
		ip++
		err := Exec(instruction)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}

	os.Exit(0)
}

// Exec runs an assembly instruction. It may change registers, change memory,
// and/or set flags.
func Exec(s string) error {
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
		return nil
	case "cmp":
		n, err := getVal(operands[1])
		if err != nil {
			return err
		}
		m, err := getVal(operands[2])
		if err != nil {
			return err
		}
		result := m - n
		zf = result == 0
		sf = result < 0
		return nil
	case "jmp", "je", "jne", "jz", "jnz", "jg", "jge", "jl", "jle":
		label := operands[1]
		addr, valid := labels[label]
		if !valid {
			return fmt.Errorf("invalid label: %s", label)
		}
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
		return nil
	case "add", "sub", "mul", "mov":
		src, err := getVal(operands[1])
		if err != nil {
			return err
		}
		pDst, err := getPtr(operands[2])
		if err != nil {
			return err
		}
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
		return nil
	case "inc":
		err := Exec("add $1 " + operands[1])
		return err
	case "dec":
		err := Exec("sub $1 " + operands[1])
		return err
	case "pop":
		err := Exec("mov (%sp) " + operands[1])
		if err != nil {
			return err
		}
		Exec("inc %sp")
		return nil
	case "push":
		val, err := getVal(operands[1])
		if err != nil {
			return err
		}
		Exec("dec %sp")
		return Exec(fmt.Sprintf("mov $%d (%%sp)", val))
	case "call":
		Exec("push %ip")
		err := Exec(fmt.Sprintf("jmp %s", operands[1]))
		if err != nil {
			return err
		}
		return nil
	case "leave":
		Exec("mov %fp %sp")
		Exec("pop %fp")
		return nil
	case "ret":
		Exec("pop %ip")
		return nil
	default:
		return fmt.Errorf("invalid instruction: %s", operands[0])
	}
}

// getVal returns the value of the operand.
// Operand is an immediate, memory, or register name.
func getVal(operand string) (int64, error) {
	if operand[0] == '$' {
		i, err := strconv.ParseInt(operand[1:], 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	} else {
		ptr, err := getPtr(operand)
		if err != nil {
			return 0, err
		}
		return *ptr, nil
	}
}

// setVal sets the value of a storage location of the operand.
// Operand is a memory or register name.
func setVal(operand string, val int64) error {
	ptr, err := getPtr(operand)
	if err != nil {
		return err
	}
	*ptr = val
	return nil
}

// getPtr returns a pointer to the storage location of the operand.
// Operand is a memory or register name.
func getPtr(operand string) (*int64, error) {
	err := fmt.Errorf("invalid operand: %s", operand)
	switch true {
	case regre.MatchString(operand):
		switch operand[1:] {
		case "a":
			return &regs[0], nil
		case "b":
			return &regs[1], nil
		case "c":
			return &regs[2], nil
		case "d":
			return &regs[3], nil
		case "e":
			return &regs[4], nil
		case "f":
			return &regs[5], nil
		case "sp":
			return &regs[6], nil
		case "fp":
			return &regs[7], nil
		case "ip":
			return &ip, nil
		default:
			return nil, err
		}
	case mem1re.MatchString(operand):
		i, err := strconv.ParseInt(operand, 10, 64)
		if err != nil {
			return nil, err
		}
		return &mem[i], nil
	case mem2re.MatchString(operand):
		reg := operand[1 : len(operand)-1]
		regPtr, err := getPtr(reg)
		if err != nil {
			return nil, err
		}
		return &mem[*regPtr], nil
	case mem3re.MatchString(operand):
		matches := mem3re.FindStringSubmatch(operand)
		reg := matches[2]
		retPtr, err := getPtr(reg)
		if err != nil {
			return nil, err
		}
		intval, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, err
		}
		return &mem[*retPtr+intval], nil
	default:
		return nil, err
	}
}
