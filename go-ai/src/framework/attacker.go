// 开火攻击行动子系统
package framework

type Attacker struct {
}

func NewAttacker() *Attacker {
	inst := &Attacker {
	}
	return inst
}

func (self *Attacker) Init(state *GameState, tankid string) {
}

func (self *Attacker) Suggest(state *GameState, objective *Objective) SuggestionItem {
	ret := SuggestionItem {
		Action: ActionFire,
		Urgent: 0.5,
	}
	return ret
}
