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
	Imm  int
}

type instructionBuildStrategy struct {
	typeSignatures instructionSignatureGroup
	build          instructionBuilder
}

type instructionBuilder func(pgm *programBuilder, op string, args []argument) Instruction

type instructionSignatureGroup [][]argumentType

var opBuildStrategies = map[string]instructionBuildStrategy{
	"mov": instructionBuildStrategy{
		typeSignatures: buildSignatureGroup(
			[]string{"xx", "xx"},
			[]string{"xx", "[xx]"},
			[]string{"xx", "123"},
			[]string{"xx", "[123]"},
			[]string{"[123]", "xx"},
			[]string{"[xx]", "xx"},
		),
		build: instructionBuilderMOV,
	},
	"add": aluInstructionBuildStrategy,
	"sub": aluInstructionBuildStrategy,
	"mul": aluInstructionBuildStrategy,
	"div": aluInstructionBuildStrategy,
	"mod": aluInstructionBuildStrategy,
	"and": aluInstructionBuildStrategy,
	"or":  aluInstructionBuildStrategy,
	"xor": aluInstructionBuildStrategy,
	"sal": aluInstructionBuildStrategy,
	"sar": aluInstructionBuildStrategy,
}

var (
	aluInstructionBuildStrategy = instructionBuildStrategy{
		typeSignatures: buildSignatureGroup(
			[]string{"xx", "xx", "xx"},
			[]string{"xx", "xx", "[123]"},
			[]string{"xx", "xx", "123"},
			[]string{"[123]", "xx", "xx"},
		),
		build: instructionBuilderALU,
	}

	aluOperationValues = map[string]byte{
		"add": 1,
		"sub": 2,
		"mul": 3,
		"div": 4,
		"mod": 5,
		"and": 6,
		"or":  7,
		"xor": 8,
		"sal": 9,
		"sar": 10,
	}
)

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

func instructionBuilderMOV(pgm *programBuilder, op string, args []argument) Instruction {
	ins0 := args[0].build(pgm)
	ins1 := args[1].build(pgm)

	return joinInstructions(ins0, ins1)
}

func instructionBuilderALU(pgm *programBuilder, op string, args []argument) Instruction {
	aluByte, _ := aluOperationValues[op]
	insAlu := Instruction{Alu: aluByte}

	ins0 := args[0].build(pgm)
	ins1 := args[1].build(pgm)
	ins2 := args[2].build(pgm)

	return joinInstructions(insAlu, ins0, ins1, ins2)
}

func buildArg(pgm *programBuilder, arg argument) {

}

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

func buildSignatureGroup(signatureStringsGroups ...[]string) instructionSignatureGroup {
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
			case "unknown":
				typ = argumentTypeUnknown
			default:
				panic("Failed to parse argument type")
			}
			group[i][j] = typ
		}
	}

	return group
}
