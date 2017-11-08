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
	for _, tank := range state.MyTank {
		(*objective)[tank.Id] = f.Objective {
			Action: rand.Int() % 6,
		}
	}
}
