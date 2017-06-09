package entities

type ConstantCombinator struct {
	*Entity
	Behavior constantBehavior `json:"control_behavior"`
}

type constantBehavior struct {
	Filters FilterGroup `json:"filters"`
}

func (con *ConstantCombinator) Init() {
	con.Entity = &Entity{Name: "constant-combinator", EnityNumber: nextID()}
	con.Entity.initEntity()
}

func NewConstantCombinator(signalCountPairs ...interface{}) ConstantCombinator {
	return NewConstantCombinatorWithGroup(NewFilterGroup(signalCountPairs...))
}

func NewConstantCombinatorWithGroup(filterGroup FilterGroup) ConstantCombinator {
	cc := ConstantCombinator{
		Behavior: constantBehavior{
			Filters: filterGroup,
		},
	}
	cc.Init()

	return cc
}
