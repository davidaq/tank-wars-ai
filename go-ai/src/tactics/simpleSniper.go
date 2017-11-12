package tactics

import (
	f "framework"
)

type SimpleSniper struct {
    superior *Simple
    tanks    []f.Tank
    policy   *SimplePolicy
	obs      *Observation
}

// 初始化
func (s *SimpleSniper) Init() {
}

// 制定执行计划
func (s *SimpleSniper) Plan(state *f.GameState, objs map[string]f.Objective){
	// 本次可支配的坦克
	var tanks []f.Tank
	for _, tank := range state.MyTank {
		if objs[tank.Id] == (f.Objective{}) {
			tanks = append(tanks, tank)
		}
	}
	if len(tanks) > 0 {
		s.Hide(tanks, objs)
	}
}

func (s *SimpleSniper) Hide(tanks []f.Tank, objs map[string]f.Objective) {
	arrPos := SortByPos(s.obs.Flag.Pos, s.obs.Kps)
	s.policy.Dispatch(tanks, arrPos[0:len(tanks)], objs)
}
