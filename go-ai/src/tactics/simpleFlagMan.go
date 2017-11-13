package tactics

import (
	f "framework"
)

type SimpleFlagMan struct {
    superior *Simple
    tanks    map[string]f.Tank
    policy   *SimplePolicy
    obs      *Observation
}

// 初始化
func (s *SimpleFlagMan) Init(tanks []f.Tank) {
    s.tanks = make(map[string]f.Tank)
    for _, tank := range tanks {
        s.tanks[tank.Id] = tank
    }
}

func (s *SimpleFlagMan) Plan(state *f.GameState, objs map[string]f.Objective) {
    // 本次可支配的坦克
    var tanks []f.Tank
    for _, tank := range state.MyTank {
        if s.tanks[tank.Id] != (f.Tank{}) {
            if objs[tank.Id] == (f.Objective{}) {
    			tanks = append(tanks, tank)
    		}
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
        s.policy.Dispatch(tanks, s.obs.Kps[0:len(tanks)], objs)
        // s.policy.Defend(s.obs.Flag.Pos, tanks, s.obs.CurState.Params.TankSpeed, objs)
    }
}

// 追击旗点坦克
// 开火判断已经在上一层做了，这里只要追击旗点附近的坦克即可
// func (s *SimpleFlagMan) huntFlagEmy(tanks []f.Tank, objs map[string]f.Objective) {
// 	for _, tank := range tanks {
// 		s.policy.Hunt(tank, s.obs.Flag.EmyTank.Pos, objs)
// 	}
// }
