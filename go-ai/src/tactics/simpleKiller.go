package tactics

import (
	f "framework"
)

type SimpleKiller struct {
    superior *Simple
    Tanks    map[string]f.Tank
    policy   *SimplePolicy
	obs      *Observation
}

// 初始化
func (s *SimpleKiller) Init(tanks []f.Tank) {
    s.Tanks = make(map[string]f.Tank)
    if len(tanks) == 0 { return }
    for _, tank := range tanks {
        s.Tanks[tank.Id] = tank
    }
}

// 制定执行计划
func (s *SimpleKiller) Plan(state *f.GameState, objs map[string]f.Objective){
    // 本次可支配的坦克
    var tanks []f.Tank
    for id, tank := range s.Tanks {
        if objs[id] == (f.Objective{}) {
            tanks = append(tanks, tank)
        }
    }
	if len(tanks) > 0 {
		s.huntEmyTank(tanks, objs)
	}
}

// 待优化
func (s *SimpleKiller) huntEmyTank(tanks []f.Tank, objs map[string]f.Objective) {
	s.policy.Patrol(s.obs.Flag.Pos, tanks, s.obs.CurState.EnemyTank, objs)
}
