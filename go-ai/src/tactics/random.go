package tactics

import (
	f "framework"
	"math/rand"
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

func (self *Random) Decide(tank *f.Tank, state *f.GameState, suggestion f.Suggestion) int {
	return rand.Int() % 6;
}
