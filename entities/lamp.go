package entities

import (
	"factorio-assembly/entities/controls"
)

type Lamp struct {
	*Entity
	Behavior lampBehavior `json:"control_behavior"`
}

type lampBehavior struct {
	Conditions lampConditions `json:"circuit_condition"`
}

type lampConditions struct {
	FirstSignal  controls.Signal     `json:"first_signal"`
	SecondSignal *controls.Signal    `json:"second_signal,omitempty"`
	Constant     int                 `json:"constant"`
	Comparator   controls.Comparator `json:"comparator"`
}

func (lamp *Lamp) Init() {
	lamp.Entity = &Entity{Name: "small-lamp", EnityNumber: nextID()}
	lamp.Entity.initEntity()
}

func NewLamp(firstSignal controls.Signal, secondSignal controls.Signal, comparator controls.Comparator) Lamp {
	return newLamp(firstSignal, &secondSignal, 0, comparator)
}

func NewLampWithConstant(firstSignal controls.Signal, constant int, comparator controls.Comparator) Lamp {
	return newLamp(firstSignal, nil, constant, comparator)
}

func newLamp(firstSignal controls.Signal, secondSignal *controls.Signal, constant int, comparator controls.Comparator) Lamp {
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
