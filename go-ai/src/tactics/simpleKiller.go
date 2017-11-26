package tactics

import (
	f "framework"
)

type BattleGroup struct {
    parent *Simple
    roles  map[string]*SimpleRole
}

// 初始化
func (b *BattleGroup) Init(tanks []f.Tank) {
    b.roles = make(map[string]*SimpleRole)

    if len(tanks) == 0 { return }

    for _, tank := range tanks {
        b.roles[tank.Id] = &SimpleRole { Tank: tank }
    }
}

func (b *BattleGroup) Plan(){

}
