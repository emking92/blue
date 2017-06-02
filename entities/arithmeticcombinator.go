package entities

import (
	"blue-asm/entities/controls"
)

type ArithmaticCombinator struct {
	*Entity
	Behavior arithmaticBehavior `json:"control_behavior"`
}

type arithmaticBehavior struct {
	Conditions arithmaticConditions `json:"arithmetic_conditions"`
}

type arithmaticConditions struct {
	FirstSignal        controls.Signal    `json:"first_signal"`
	SecondSignal       *controls.Signal   `json:"second_signal,omitempty"`
	Constant           int                `json:"constant"`
	Operation          controls.Operation `json:"operation"`
	OutputSignal       controls.Signal    `json:"output_signal"`
	CopyCountFromInput int                `json:"copy_count_from_input"`
}

func (comb *ArithmaticCombinator) Init() {
	comb.Entity = &Entity{Name: "arithmetic-combinator", EnityNumber: nextID()}
	comb.initEntity()
}

func NewArithmaticCombinator(firstSignal controls.Signal, secondSignal controls.Signal, operation controls.Operation, outputSignal controls.Signal) ArithmaticCombinator {
	return newArithmaticCombinator(firstSignal, &secondSignal, 0, operation, outputSignal)
}

func NewArithmaticCombinatorWithConstant(firstSignal controls.Signal, constant int, operation controls.Operation, outputSignal controls.Signal) ArithmaticCombinator {
	return newArithmaticCombinator(firstSignal, nil, constant, operation, outputSignal)
}

func newArithmaticCombinator(firstSignal controls.Signal, secondSignal *controls.Signal, constant int, operation controls.Operation, outputSignal controls.Signal) ArithmaticCombinator {
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
