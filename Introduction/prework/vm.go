package vm

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
func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {

		// FETCH
		instruction_binary_format := memory[registers[0]]
		var address byte = 0
		var instruction string = ""
		var register_one_index byte = 0
		var register_two_index byte = 0
		var constants byte = 0

		// DECODE
		switch instruction_binary_format {
		case Load:
			instruction = "Load"
			register_one_index = memory[registers[0]+1]
			address = memory[registers[0]+2]
		case Store:
			instruction = "Store"
			register_one_index = memory[registers[0]+1]
			address = memory[registers[0]+2]
		case Add:
			instruction = "Add"
			register_one_index = memory[registers[0]+1]
			register_two_index = memory[registers[0]+2]
		case Sub:
			instruction = "Sub"
			register_one_index = memory[registers[0]+1]
			register_two_index = memory[registers[0]+2]
		case Addi:
			instruction = "Addi"
			register_one_index = memory[registers[0]+1]
			constants = memory[registers[0]+2]
		case Subi:
			instruction = "Subi"
			register_one_index = memory[registers[0]+1]
			constants = memory[registers[0]+2]
		case Jump:
			instruction = "Jump"
			constants = memory[registers[0]+1]
		case Beqz:
			instruction = "Beqz"
			register_one_index = memory[registers[0]+1]
			constants = memory[registers[0]+2]
		case Halt:
			instruction = "Halt"
		default:
			panic("Invalid operation: " + string(instruction_binary_format))
		}

		// EXECUTE
		switch instruction {
		case "Load":
			registers[register_one_index] = memory[address]
			registers[0] += 3
		case "Store":
			memory[address] = registers[register_one_index]
			registers[0] += 3
		case "Add":
			registers[register_one_index] = registers[register_one_index] + registers[register_two_index]
			registers[0] += 3
		case "Sub":
			registers[register_one_index] = registers[register_one_index] - registers[register_two_index]
			registers[0] += 3
		case "Addi":
			registers[register_one_index] = registers[register_one_index] + constants
			registers[0] += 3
		case "Subi":
			registers[register_one_index] = registers[register_one_index] - constants
			registers[0] += 3
		case "Jump":
			registers[0] = constants
		case "Beqz":
			registers[0] += 3
			if registers[register_one_index] == 0 {
				registers[0] += constants
			}
		case "Halt":
			return
		default:
			panic("Invalid opcode: " + instruction)
		}
	}
}
