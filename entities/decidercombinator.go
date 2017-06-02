package entities

import (
	"factorio-assembly/entities/controls"
)

type DeciderCombinator struct {
	*Entity
	Behavior deciderBehavior `json:"control_behavior"`
}

type deciderBehavior struct {
	Conditions deciderConditions `json:"decider_conditions"`
}

type deciderConditions struct {
	FirstSignal        controls.Signal     `json:"first_signal"`
	SecondSignal       *controls.Signal    `json:"second_signal,omitempty"`
	Constant           int                 `json:"constant"`
	Comparator         controls.Comparator `json:"comparator"`
	OutputSignal       controls.Signal     `json:"output_signal"`
	CopyCountFromInput bool                `json:"copy_count_from_input"`
}

func (comb *DeciderCombinator) Init() {
	comb.Entity = &Entity{Name: "decider-combinator", EnityNumber: nextID()}
	comb.Entity.initEntity()
}

func NewDeciderCombinator(firstSignal controls.Signal, secondSignal controls.Signal, comparator controls.Comparator, outputSignal controls.Signal, copyCountFromInput bool) DeciderCombinator {
	return newDeciderCombinator(firstSignal, &secondSignal, 0, comparator, outputSignal, copyCountFromInput)
}

func NewDeciderCombinatorWithConstant(firstSignal controls.Signal, constant int, comparator controls.Comparator, outputSignal controls.Signal, copyCountFromInput bool) DeciderCombinator {
	return newDeciderCombinator(firstSignal, nil, constant, comparator, outputSignal, copyCountFromInput)
}

func newDeciderCombinator(firstSignal controls.Signal, secondSignal *controls.Signal, constant int, comparator controls.Comparator, outputSignal controls.Signal, copyCountFromInput bool) DeciderCombinator {
	dc := DeciderCombinator{
		Behavior: deciderBehavior{
			deciderConditions{
				FirstSignal:        firstSignal,
				SecondSignal:       secondSignal,
				Constant:           constant,
				Comparator:         comparator,
				OutputSignal:       outputSignal,
				CopyCountFromInput: copyCountFromInput,
			},
		},
	}
	dc.Init()

	return dc
}
