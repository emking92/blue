package compiler

import (
	"strconv"
)

var precompilerStrategies = map[string]precompiler{
	"const": precompiler{
		precompile: precompileConstants,
		argCount:   2,
	},
}

type precompiler struct {
	argCount   int
	precompile func(pgm *programBuilder, op string, args []string)
}

func precompileConstants(pgm *programBuilder, op string, args []string) {
	c := args[0]

	_, err := strconv.Atoi(args[1])
	if err == strconv.ErrRange {
		pgm.compileErrorString("const value out of int32 range: %s", args[1])
		return
	} else if err != nil {
		pgm.compileErrorString("const value must be an integer: %s", args[1])
		return
	}

	if pgm.variables.IsVarDefined(c) {
		pgm.compileErrorString(`const variable already defined: "%s"`, c)
		return
	}

	pgm.variables.CreateVariable(c, args[1])
}
