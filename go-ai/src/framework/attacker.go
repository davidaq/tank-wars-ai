// 开火攻击行动子系统
package framework

type Attacker struct {
}

func NewAttacker() *Attacker {
	inst := &Attacker {
	}
	return inst
}

func (self *Attacker) Suggest(tank *Tank, state *GameState, objective *Objective) SuggestionItem {
	ret := SuggestionItem {
		Action: ActionFire,
		Urgent: 1,
	}
	return ret
}
