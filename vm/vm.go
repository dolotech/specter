package vm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

// The VM, with its program and memory abstractions.
type VM struct {
	p *program
	m *memory
}

// Not safe for multi-threads
var b *bufio.Writer = bufio.NewWriter(os.Stdout)

// Create a new VM.
func New() *VM {
	return &VM{newProgram(), newMemory()}
}

// Run executes the vm bytecode read by the reader.
func Run(vm *VM, r io.Reader) {
	var i int32

	// Parse the content to execute.
	parse(vm, r)

	// Execution loop.
	defer b.Flush()
	for i = vm.p.start; vm.p.instrs[i] != _OP_END; i++ {
		runInstruction(vm, &i)
	}
}

// Run a single instruction.
func runInstruction(vm *VM, instrIndex *int32) {
	a0, a1 := vm.p.args[*instrIndex][0], vm.p.args[*instrIndex][1]

	//printInstr("before", *instrIndex, opcode(vm.p.instrs.sl[*instrIndex]), a0, a1)

	switch vm.p.instrs[*instrIndex] {
	case _OP_NOP:
		// Nothing
	case _OP_INT:
		// Not implemented
	case _OP_MOV:
		*a0 = *a1
	case _OP_PUSH:
		pushStack(vm.m, *a0)
	case _OP_POP:
		popStack(vm.m, a0)
	case _OP_PUSHF:
		pushStack(vm.m, vm.m.FLAGS)
	case _OP_POPF:
		popStack(vm.m, a0)
	case _OP_INC:
		(*a0)++
	case _OP_DEC:
		(*a0)--
	case _OP_ADD:
		*a0 += *a1
	case _OP_SUB:
		*a0 -= *a1
	case _OP_MUL:
		*a0 *= *a1
	case _OP_DIV:
		*a0 /= *a1
	case _OP_MOD:
		vm.m.remainder = *a0 % *a1
	case _OP_REM:
		*a0 = vm.m.remainder
	case _OP_NOT:
		*a0 = ^(*a0)
	case _OP_XOR:
		*a0 ^= *a1
	case _OP_OR:
		*a0 |= *a1
	case _OP_AND:
		*a0 &= *a1
	case _OP_SHL:
		// cannot shift on signed int32
		if *a1 > 0 {
			*a0 <<= uint(*a1)
		}
	case _OP_SHR:
		// cannot shift on signed int32
		if *a1 > 0 {
			*a0 >>= uint(*a1)
		}
	case _OP_CMP:
		if *a0 == *a1 {
			vm.m.FLAGS = 0x1
		} else if *a0 > *a1 {
			vm.m.FLAGS = 0x2
		} else {
			vm.m.FLAGS = 0x0
		}
	case _OP_CALL:
		pushStack(vm.m, *instrIndex)
		fallthrough
	case _OP_JMP:
		*instrIndex = *a0 - 1
	case _OP_RET:
		popStack(vm.m, instrIndex)
	case _OP_JE:
		if vm.m.FLAGS&0x1 != 0 {
			*instrIndex = *a0 - 1
		}
	case _OP_JNE:
		if vm.m.FLAGS&0x1 == 0 {
			*instrIndex = *a0 - 1
		}
	case _OP_JG:
		if vm.m.FLAGS&0x2 != 0 {
			*instrIndex = *a0 - 1
		}
	case _OP_JGE:
		if vm.m.FLAGS&0x3 != 0 {
			*instrIndex = *a0 - 1
		}
	case _OP_JL:
		if vm.m.FLAGS&0x3 == 0 {
			*instrIndex = *a0 - 1
		}
	case _OP_JLE:
		if vm.m.FLAGS&0x2 == 0 {
			*instrIndex = *a0 - 1
		}
	case _OP_PRN:
		//fmt.Printf("%d\n", *a0)
		b.WriteString(strconv.FormatInt(int64(*a0), 10))
		// WriteRune calls WriteByte, so save a call
		b.WriteByte('\n')
	}
	/*
		if *instrIndex >= 0 {
			printInstr("after", *instrIndex, opcode(vm.p.instrs.sl[*instrIndex]), a0, a1)
		} else {
			printInstr("after", *instrIndex, opcode(vm.p.instrs.sl[*instrIndex+1]), a0, a1)
		}
	*/
}

func printInstr(prefix string, idx int32, op opcode, a0, a1 *int32) {
	switch {
	case a0 == nil && a1 == nil:
		fmt.Printf("[%s] instr=%d: %d (%s) a0=nil, a1=nil\n", prefix, idx, op, op)
	case a1 == nil:
		fmt.Printf("[%s] instr=%d: %d (%s) a0=%d, a1=nil\n", prefix, idx, op, op, *a0)
	case a0 == nil:
		fmt.Printf("[%s] instr=%d: %d (%s) a0=nil, a1=%d\n", prefix, idx, op, op, *a1)
	default:
		fmt.Printf("[%s] instr=%d: %d (%s) a0=%d, a1=%d\n", prefix, idx, op, op, *a0, *a1)
	}
}
