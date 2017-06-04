package compiler

type Instruction struct {
	Amux byte
	Bmux byte
	Cmux byte
	Cond byte
	Alu  byte
	Mbr  byte
	Mar  byte
	Rd   byte
	Wr   byte
	Enc  byte
	A    byte
	B    byte
	C    byte
	Addr int
	Bran int
	Imm  int
}

type instructionSignatureGroup [][]argumentType

func joinInstructions(instructions ...Instruction) (out Instruction) {
	if len(instructions) == 0 {
		return
	}

	out = instructions[0]

	for _, instruction := range instructions[1:] {
		setInstructionByteValueOnce(&out.Amux, instruction.Amux)
		setInstructionByteValueOnce(&out.Bmux, instruction.Bmux)
		setInstructionByteValueOnce(&out.Cmux, instruction.Cmux)
		setInstructionByteValueOnce(&out.Cond, instruction.Cond)
		setInstructionByteValueOnce(&out.Alu, instruction.Alu)
		setInstructionByteValueOnce(&out.Mbr, instruction.Mbr)
		setInstructionByteValueOnce(&out.Mar, instruction.Mar)
		setInstructionByteValueOnce(&out.Rd, instruction.Rd)
		setInstructionByteValueOnce(&out.Wr, instruction.Wr)
		setInstructionByteValueOnce(&out.Enc, instruction.Enc)
		setInstructionByteValueOnce(&out.A, instruction.A)
		setInstructionByteValueOnce(&out.B, instruction.B)
		setInstructionByteValueOnce(&out.C, instruction.C)
		setInstructionIntValueOnce(&out.Addr, instruction.Addr)
		setInstructionIntValueOnce(&out.Bran, instruction.Bran)
		setInstructionIntValueOnce(&out.Imm, instruction.Imm)
	}

	return
}

func setInstructionByteValueOnce(dest *byte, value byte) {
	if value == 0 {
		return
	}

	if *dest != 0 {
		panic("Instruction field is given a nonzero value more than once")
	}

	*dest = value
}

func setInstructionIntValueOnce(dest *int, value int) {
	if value == 0 {
		return
	}

	if *dest != 0 {
		panic("Instruction field is given a nonzero value more than once")
	}

	*dest = value
}

func createSignatureGroup(signatureStringsGroups ...[]string) instructionSignatureGroup {
	group := make(instructionSignatureGroup, len(signatureStringsGroups))
	for i, signatureStrings := range signatureStringsGroups {
		group[i] = make([]argumentType, len(signatureStrings))

		for j, typeStr := range signatureStrings {
			var typ argumentType
			switch typeStr {
			case "xx":
				typ = argumentTypeRegister
			case "[xx]":
				typ = argumentTypeRegisterPointer
			case "123":
				typ = argumentTypeImmediate
			case "[123]":
				typ = argumentTypeImmediatePointer
			case "pp":
				typ = argumentTypePort
			case "label":
				typ = argumentTypeLabel
			case "undefined":
				typ = argumentTypeUndefined
			default:
				panic("Failed to parse argument type")
			}
			group[i][j] = typ
		}
	}

	return group
}

func (sigGroup instructionSignatureGroup) matches(args []argument) bool {
	for _, signature := range sigGroup {
		if len(args) != len(signature) {
			continue
		}

		match := true
		for i := range args {
			if args[i].argType != signature[i] {
				match = false
				break
			}
		}
		if match == true {
			return true
		}
	}

	return false
}
