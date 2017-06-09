package entities

type DeciderCombinator struct {
	*Entity
	Behavior deciderBehavior `json:"control_behavior"`
}

type deciderBehavior struct {
	Conditions deciderConditions `json:"decider_conditions"`
}

type deciderConditions struct {
	FirstSignal        Signal     `json:"first_signal"`
	SecondSignal       *Signal    `json:"second_signal,omitempty"`
	Constant           int        `json:"constant"`
	Comparator         Comparator `json:"comparator"`
	OutputSignal       Signal     `json:"output_signal"`
	CopyCountFromInput bool       `json:"copy_count_from_input"`
}

func (comb *DeciderCombinator) Init() {
	comb.Entity = &Entity{Name: "decider-combinator", EnityNumber: nextID()}
	comb.Entity.initEntity()
}

func NewDeciderCombinator(firstSignal Signal, secondSignal Signal, comparator Comparator, outputSignal Signal, copyCountFromInput bool) DeciderCombinator {
	return newDeciderCombinator(firstSignal, &secondSignal, 0, comparator, outputSignal, copyCountFromInput)
}

func NewDeciderCombinatorWithConstant(firstSignal Signal, constant int, comparator Comparator, outputSignal Signal, copyCountFromInput bool) DeciderCombinator {
	return newDeciderCombinator(firstSignal, nil, constant, comparator, outputSignal, copyCountFromInput)
}

func newDeciderCombinator(firstSignal Signal, secondSignal *Signal, constant int, comparator Comparator, outputSignal Signal, copyCountFromInput bool) DeciderCombinator {
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
