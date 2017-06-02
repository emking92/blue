package entities

import (
	"factorio-assembly/entities/controls"
)

type ConstantCombinator struct {
	*Entity
	Behavior constantBehavior `json:"control_behavior"`
}

type constantBehavior struct {
	Filters controls.FilterGroup `json:"filters"`
}

func (con *ConstantCombinator) Init() {
	con.Entity = &Entity{Name: "constant-combinator", EnityNumber: nextID()}
	con.Entity.initEntity()
}

func NewConstantCombinator(signalCountPairs ...interface{}) ConstantCombinator {
	return NewConstantCombinatorWithGroup(controls.NewFilterGroup(signalCountPairs...))
}

func NewConstantCombinatorWithGroup(filterGroup controls.FilterGroup) ConstantCombinator {
	cc := ConstantCombinator{
		Behavior: constantBehavior{
			Filters: filterGroup,
		},
	}
	cc.Init()

	return cc
}
