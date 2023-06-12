package vm

import (
	"fmt"
)

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	for {
		pc := registers[0]
		op := memory[pc]

		if op == Halt {
			return
		}

		// fetch opargs: most instructions have 2
		arg1, arg2 := memory[pc+1], memory[pc+2]

		// execute
		switch op {
		case Load:
			// load a value into the specified register
			registers[arg1] = memory[arg2]
			registers[0] += 3
		case Store:
			// store a value from a register
			memory[arg2] = registers[arg1]
			registers[0] += 3
		case Add:
			// perform addition and replace arg1 with the result
			registers[arg1] += registers[arg2]
			registers[0] += 3
		case Sub:
			// perofrm substraction and replace arg1 with the result
			registers[arg1] -= registers[arg2]
			registers[0] += 3
		case Addi:
			registers[arg1] += arg2
			registers[0] += 3
		case Subi:
			registers[arg1] -= arg2
			registers[0] += 3
		case Jump:
			registers[0] = arg1
		case Beqz:
			registers[0] += 3
			if registers[arg1] == 0 {
				registers[0] += arg2
			}
		default:
			panic(fmt.Sprintf("Invalid opcode: 0x%02x", op))
		}
	}
}
