package tactics

import (
	f "framework"
)

type SimpleKiller struct {
    superior *Simple
    tanks    []f.Tank
    policy   *SimplePolicy
	obs      *Observation
}

// 初始化
func (s *SimpleKiller) Init() {
}

// 制定执行计划
func (s *SimpleKiller) Plan(state *f.GameState, objs map[string]f.Objective){
	// 本次可支配的坦克
	var tanks []f.Tank
	for _, tank := range state.MyTank {
		if objs[tank.Id] == (f.Objective{}) {
			tanks = append(tanks, tank)
		}
	}
	if len(tanks) > 0 {
		s.huntEmyTank(tanks, objs)
	}
}

func (s *SimpleKiller) huntEmyTank(tanks []f.Tank, objs map[string]f.Objective) {
	s.policy.Patrol(s.obs.Flag.Pos, tanks, s.obs.CurState.EnemyTank, objs)
}
