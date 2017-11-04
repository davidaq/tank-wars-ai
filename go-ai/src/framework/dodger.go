// 子弹躲避行动子系统
package framework

type Dodger struct {
}

func NewDodger() *Dodger {
	inst := &Dodger {
	}
	return inst
}

func (self *Dodger) Suggest(tank *Tank, state *GameState, objective *Objective) SuggestionItem {
	ret := SuggestionItem {
		Action: ActionMove,
		Urgent: 1,
	}
	return ret
}
