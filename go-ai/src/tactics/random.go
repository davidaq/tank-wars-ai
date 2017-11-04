package tactics

import (
	f "framework"
)

type Random struct {
}

func NewRandom() *Random {
	inst := &Random {}
	return inst
}

func (self *Random) Init(state *f.GameState) {
}

func (self *Random) Plan(state *f.GameState, objective *map[string]f.Objective) {
}

func (self *Random) Decide(tankid string, suggestion f.Suggestion) int {
	return f.ActionLeft
}
