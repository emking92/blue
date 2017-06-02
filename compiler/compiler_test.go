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
	MOV cx, 222
	`

	expected := []Instruction{
		Instruction{Enc: 1, C: 10, A: 11},
		Instruction{Enc: 1, C: 12, Cmux: 1, Imm: 222},
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
	}

	testBuild(t, source, expected)
}

func TestVAR(t *testing.T) {
	source := `
	VAR foo, 123
	VAR bar, ax
	`

	expected := []Instruction{
		Instruction{Mbr: 1, Mar: 2, Wr: 1, Addr: 16, Cmux: 1, Imm: 123},
		Instruction{Mbr: 1, Mar: 2, Wr: 1, Addr: 17, A: 10},
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
