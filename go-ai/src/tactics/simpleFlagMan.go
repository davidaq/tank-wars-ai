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
	// 占领旗点
    if (s.obs.Flag.Exist) {
        s.policy.Occupy(s.obs.Flag.Pos, tanks[0], objs)
    	tanks = tanks[1:]
    }
    // 其余回到苟点
    if len(tanks) > 0 {
        ftanks := s.policy.Dispatch(tanks, s.obs.FlagKps[0:len(tanks)], objs)
        // 已抵达目的地坦克根据友伤选择是否开火
        var radarFire f.RadarFireAll
        for _, ftank := range ftanks {
            radarFire = s.obs.Radar.Fire[ftank.Id]
            // 无友伤则开火
            if s.obs.Side == "red" {
                if radarFire == (f.RadarFireAll{}) || radarFire.Down.Sin <= 0.3 {
                    objs[ftank.Id] = f.Objective{ Action: f.ActionFireDown}
                } else {
                    objs[ftank.Id] = f.Objective{ Action: f.ActionStay}
                }
            } else {
                if radarFire == (f.RadarFireAll{}) || radarFire.Up.Sin <= 0.3 {
                    objs[ftank.Id] = f.Objective{ Action: f.ActionFireUp}
                } else {
                    objs[ftank.Id] = f.Objective{ Action: f.ActionStay}
                }
            }
        }
    }
}
