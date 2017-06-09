package compiler

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"text/tabwriter"
)

func TestMOV(t *testing.T) {
	source := `
	MOV ax, bx
	MOV cx, 101
	MOV dx, [102]
	MOV [103], ex
	MOV [105], 106
	MOV [ax], bx
	MOV [cx], 104
	MOV dx, [ex]
	`

	expected := []Instruction{
		Instruction{Enc: 1, C: 10, A: 11},
		Instruction{Enc: 1, C: 12, Cmux: 1, Imm: 101},
		Instruction{Enc: 1, C: 13, Amux: 1, Mar: 2, Addr: 102, Rd: 1},
		Instruction{Mar: 2, Addr: 103, Wr: 1, A: 14, Mbr: 1},
		Instruction{Mar: 2, Addr: 105, Wr: 1, Mbr: 1, Cmux: 1, Imm: 106},
		Instruction{B: 10, Wr: 1, Mar: 1, A: 11, Mbr: 1},
		Instruction{B: 12, Wr: 1, Mar: 1, Cmux: 1, Imm: 104, Mbr: 1},
		Instruction{Enc: 1, C: 13, B: 14, Rd: 1, Mar: 1, Amux: 1},
	}

	testBuild(t, source, expected)
}

func TestALU(t *testing.T) {
	source := `
	ADD ax, bx, cx
	SUB ax, bx, cx
	MUL ax, bx, cx
	DIV ax, bx, cx
	MOD ax, bx, cx
	AND ax, bx, cx
	OR ax, bx, cx
	XOR ax, bx, cx
	SAL ax, bx, cx
	SAR ax, bx, cx
	
	ADD ax, bx, 123
	MUL cx, [124], dx
	DIV [125], ex, fx
	MOD [126], ax, 127
	AND bx, [128], 129
	`

	expected := []Instruction{
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 1},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 2},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 3},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 4},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 5},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 6},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 7},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 8},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 9},
		Instruction{Enc: 1, A: 11, B: 12, C: 10, Alu: 10},
		Instruction{Enc: 1, A: 11, Bmux: 1, C: 10, Alu: 1, Imm: 123},
		Instruction{Enc: 1, B: 13, C: 12, Alu: 3, Addr: 124, Rd: 1, Mar: 2, Amux: 1},
		Instruction{Addr: 125, Wr: 1, Mar: 2, A: 14, B: 15, Alu: 4, Mbr: 1},
		Instruction{Addr: 126, Wr: 1, Mar: 2, A: 10, Imm: 127, Bmux: 1, Alu: 5, Mbr: 1},
		Instruction{Enc: 1, C: 11, Alu: 6, Addr: 128, Rd: 1, Mar: 2, Amux: 1, Bmux: 1, Imm: 129},
	}

	testBuild(t, source, expected)
}

func TestVAR(t *testing.T) {
	source := `
	VAR foo, 123
	VAR bar, ax
	MOV ax, [$foo]
	`

	expected := []Instruction{
		Instruction{Mbr: 1, Mar: 2, Wr: 1, Addr: 16, Cmux: 1, Imm: 123},
		Instruction{Mbr: 1, Mar: 2, Wr: 1, Addr: 17, A: 10},
		Instruction{Enc: 1, C: 10, Amux: 1, Mar: 2, Addr: 16, Rd: 1},
	}

	testBuild(t, source, expected)
}

func TestJMP(t *testing.T) {
	source := `
	foo:MOV ax, 0
	
	MOV bx, 0
	bar:MOV cx, 0
	
	JMP foo
	JMP bar
	`

	expected := []Instruction{
		Instruction{Enc: 1, C: 10, Cmux: 1},
		Instruction{Enc: 1, C: 11, Cmux: 1},
		Instruction{Enc: 1, C: 12, Cmux: 1},
		Instruction{Cond: 7, Bran: 0},
		Instruction{Cond: 7, Bran: 2},
	}

	testBuild(t, source, expected)
}

func TestCondJMP(t *testing.T) {
	source := `
	label0:MOV ax, 0
	label1:MOV ax, 0
	label2:MOV ax, 0
	label3:MOV ax, 0
	label4:MOV ax, 0
	label5:MOV ax, 0
	
	JE label0, ax, bx
	JNE label1, ax, bx
	JL label2, ax, bx
	JG label3, ax, bx
	JLE label4, ax, bx
	JGE label5, ax, bx
	
	JE label0, ax, 123
	JE label1, [111], bx
	JE label2, [222], 456
	`

	expected := []Instruction{
		Instruction{Enc: 1, C: 10, Cmux: 1},
		Instruction{Enc: 1, C: 10, Cmux: 1},
		Instruction{Enc: 1, C: 10, Cmux: 1},
		Instruction{Enc: 1, C: 10, Cmux: 1},
		Instruction{Enc: 1, C: 10, Cmux: 1},
		Instruction{Enc: 1, C: 10, Cmux: 1},

		Instruction{A: 10, B: 11, Alu: 2, Bran: 0, Cond: 1},
		Instruction{A: 10, B: 11, Alu: 2, Bran: 1, Cond: 2},
		Instruction{A: 10, B: 11, Alu: 2, Bran: 2, Cond: 3},
		Instruction{A: 10, B: 11, Alu: 2, Bran: 3, Cond: 4},
		Instruction{A: 10, B: 11, Alu: 2, Bran: 4, Cond: 5},
		Instruction{A: 10, B: 11, Alu: 2, Bran: 5, Cond: 6},

		Instruction{A: 10, Alu: 2, Cond: 1, Bran: 0, Bmux: 1, Imm: 123},
		Instruction{Amux: 1, Mar: 2, Addr: 111, Rd: 1, Alu: 2, Cond: 1, Bran: 1, B: 11},
		Instruction{Amux: 1, Mar: 2, Addr: 222, Rd: 1, Alu: 2, Cond: 1, Bran: 2, Bmux: 1, Imm: 456},
	}

	testBuild(t, source, expected)
}

func TestIO(t *testing.T) {
	source := `
	IN ax, p_1
	OUT p_2, bx
	OUT p_3, 120
	`

	expected := []Instruction{
		Instruction{Enc: 1, C: 10, Amux: 1, Mar: 2, Addr: 1, Rd: 1},
		Instruction{Mar: 2, Addr: 2, Wr: 1, A: 11, Mbr: 1},
		Instruction{Mar: 2, Addr: 3, Wr: 1, Mbr: 1, Cmux: 1, Imm: 120},
	}

	testBuild(t, source, expected)
}

func TestConst(t *testing.T) {
	source := `
	#CONST foo, 200
	#CONST bar, 201
	#CONST fun, 202
	MOV ax, $foo
	VAR test, $bar
	ADD bx, [$test], $fun
	`
	expected := []Instruction{
		Instruction{Enc: 1, C: 10, Cmux: 1, Imm: 200},
		Instruction{Mbr: 1, Mar: 2, Wr: 1, Addr: 16, Cmux: 1, Imm: 201},
		Instruction{Enc: 1, C: 11, Alu: 1, Addr: 16, Rd: 1, Mar: 2, Amux: 1, Bmux: 1, Imm: 202},
	}

	testBuild(t, source, expected)
}

func testBuild(t *testing.T, source string, expected []Instruction) {
	reader := bytes.NewReader([]byte(source))

	result, err := BuildSource(reader)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(result, expected) {
		aligned := alignResultExpected(result, expected)
		t.Errorf("Build results do not match expected\n%s", aligned)
	}
}

func alignResultExpected(result interface{}, expected interface{}) string {

	str := fmt.Sprintf("Result:\t%+v\nExpected:\t%+v", result, expected)
	str = strings.Replace(str, " ", "\t", -1)

	buf := bytes.Buffer{}

	aligner := tabwriter.NewWriter(&buf, 0, 1, 1, ' ', 0)
	aligner.Write([]byte(str))
	aligner.Flush()

	return buf.String()
}
