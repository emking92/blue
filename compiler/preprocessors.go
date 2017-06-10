package compiler

var precompilerStrategies = map[string]precompiler{
	"define": precompiler{
		precompile: precompileMacros,
		argCount:   2,
	},
}

type precompiler struct {
	argCount   int
	precompile func(pgm *programBuilder, op string, args []string)
}

func precompileMacros(pgm *programBuilder, op string, args []string) {
	c := args[0]

	if pgm.variables.IsVarDefined(c) {
		pgm.compileErrorString(`macro already defined: "%s"`, c)
		return
	}

	pgm.variables.CreateVariable(c, args[1])
}
