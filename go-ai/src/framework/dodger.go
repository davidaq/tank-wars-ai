// 子弹躲避行动子系统
package framework

type Dodger struct {
}

func NewDodger() *Dodger {
	inst := &Dodger {
	}
	return inst
}

func (self *Dodger) Init(state *GameState, tankid string) {
}

func (self *Dodger) Suggest(state *GameState, objective *Objective) SuggestionItem {
	ret := SuggestionItem {
		Action: ActionMove,
		Urgent: 1,
	}
	return ret
}
