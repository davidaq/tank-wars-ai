package tactics

import (
	f "framework"
)

type FlagGroup struct {
    parent *Simple
    roles  map[string]*SimpleRole
}

// 初始化
func (f *FlagGroup) Init(tanks []f.Tank) {
    f.roles = make(map[string]*SimpleRole)

    if len(tanks) == 0 { return }

    for _, tank := range tanks {
        f.roles[tank.Id] = &SimpleRole { Tank: tank }
    }
}

func (f *FlagGroup) Plan() {

}
