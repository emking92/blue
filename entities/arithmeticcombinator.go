package entities

type ArithmaticCombinator struct {
	*Entity
	Behavior arithmaticBehavior `json:"control_behavior"`
}

type arithmaticBehavior struct {
	Conditions arithmaticConditions `json:"arithmetic_conditions"`
}

type arithmaticConditions struct {
	FirstSignal        Signal    `json:"first_signal"`
	SecondSignal       *Signal   `json:"second_signal,omitempty"`
	Constant           int       `json:"constant"`
	Operation          Operation `json:"operation"`
	OutputSignal       Signal    `json:"output_signal"`
	CopyCountFromInput int       `json:"copy_count_from_input"`
}

func (comb *ArithmaticCombinator) Init() {
	comb.Entity = &Entity{Name: "arithmetic-combinator", EnityNumber: nextID()}
	comb.initEntity()
}

func NewArithmaticCombinator(firstSignal Signal, secondSignal Signal, operation Operation, outputSignal Signal) ArithmaticCombinator {
	return newArithmaticCombinator(firstSignal, &secondSignal, 0, operation, outputSignal)
}

func NewArithmaticCombinatorWithConstant(firstSignal Signal, constant int, operation Operation, outputSignal Signal) ArithmaticCombinator {
	return newArithmaticCombinator(firstSignal, nil, constant, operation, outputSignal)
}

func newArithmaticCombinator(firstSignal Signal, secondSignal *Signal, constant int, operation Operation, outputSignal Signal) ArithmaticCombinator {
	ac := ArithmaticCombinator{
		Behavior: arithmaticBehavior{
			arithmaticConditions{
				FirstSignal:  firstSignal,
				SecondSignal: secondSignal,
				Constant:     constant,
				Operation:    operation,
				OutputSignal: outputSignal,
			},
		},
	}
	ac.Init()
	return ac
}
