// 行走寻路行动子系统
package framework

type Traveller struct {
}

func NewTraveller() *Traveller {
	inst := &Traveller {
	}
	return inst
}

func (self *Traveller) Init(state *GameState, tankid string) {
}

func (self *Traveller) Suggest(state *GameState, objective *Objective) SuggestionItem {
	ret := SuggestionItem {
		Action: ActionMove,
		Urgent: 1,
	}
	return ret
}
