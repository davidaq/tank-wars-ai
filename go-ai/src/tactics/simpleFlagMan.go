package tactics

import (
	f "framework"
    // "fmt"
)

type SimpleFlagMan struct {
    superior *Simple
    Tanks    map[string]f.Tank
    policy   *SimplePolicy
    obs      *Observation
}

// 初始化
func (s *SimpleFlagMan) Init(tanks []f.Tank) {
    s.Tanks = make(map[string]f.Tank)
    if len(tanks) == 0 { return }
    for _, tank := range tanks {
        s.Tanks[tank.Id] = tank
    }
}

func (s *SimpleFlagMan) Plan(state *f.GameState, objs map[string]f.Objective) {
    // 本次可支配的坦克
    var tanks []f.Tank
    for id, tank := range s.Tanks {
        if objs[id] == (f.Objective{}) {
            tanks = append(tanks, tank)
        }
    }
	if len(tanks) > 0 {
		s.occupyFlag(tanks, objs)
	}
}

// 占据旗点
func (s *SimpleFlagMan) occupyFlag(tanks []f.Tank, objs map[string]f.Objective) {
	// 前往旗点
    if (s.obs.Flag.Exist) {
        s.policy.MoveTo(s.obs.Flag.Pos, tanks[0], objs)
    	tanks = tanks[1:]
    }
    // 其它坦克
    if len(tanks) > 0 {
        ftanks := make(map[string]f.Tank)
        if len(s.obs.FlagKps) < len(tanks) {
            ftanks = s.policy.Dispatch(tanks[0:len(s.obs.FlagKps)], s.obs.FlagKps, objs)
            for _, tank := range tanks[len(s.obs.FlagKps):]{
                ftanks[tank.Id] = tank
            }
        } else {
            ftanks = s.policy.Dispatch(tanks, s.obs.FlagKps, objs)
        }
        // 自由坦克判断能否开火
        if len(ftanks) > 0 {
            s.policy.FireToFlag(ftanks, s.obs, objs)
        }
    }
}
