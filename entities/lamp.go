package entities

type Lamp struct {
	*Entity
	Behavior lampBehavior `json:"control_behavior"`
}

type lampBehavior struct {
	Conditions lampConditions `json:"circuit_condition"`
}

type lampConditions struct {
	FirstSignal  Signal     `json:"first_signal"`
	SecondSignal *Signal    `json:"second_signal,omitempty"`
	Constant     int        `json:"constant"`
	Comparator   Comparator `json:"comparator"`
}

func (lamp *Lamp) Init() {
	lamp.Entity = &Entity{Name: "small-lamp", EnityNumber: nextID()}
	lamp.Entity.initEntity()
}

func NewLamp(firstSignal Signal, secondSignal Signal, comparator Comparator) Lamp {
	return newLamp(firstSignal, &secondSignal, 0, comparator)
}

func NewLampWithConstant(firstSignal Signal, constant int, comparator Comparator) Lamp {
	return newLamp(firstSignal, nil, constant, comparator)
}

func newLamp(firstSignal Signal, secondSignal *Signal, constant int, comparator Comparator) Lamp {
	lamp := Lamp{
		Behavior: lampBehavior{
			lampConditions{
				FirstSignal:  firstSignal,
				SecondSignal: secondSignal,
				Constant:     constant,
				Comparator:   comparator,
			},
		},
	}

	lamp.Init()

	return lamp
}
