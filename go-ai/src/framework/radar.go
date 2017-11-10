// 雷达系统，用于侦测威胁以及可以开火的目标
package framework;

type Radar struct {
}

func NewRadar() *Radar {
	inst := &Radar {
	}
	return inst
}

// 检查自己当前的位置以及面朝方向前进的位置周围是否安全
// 检查每个坦克四周开火命中率与代价
func (self *Radar) Scan(state *GameState) *RadarResult {
	ret := &RadarResult {
		Dodge: make(map[string]RadarDodge),
		Fire: make(map[string]RadarFireAll),
	}
	for _, tank := range state.MyTank {
		ret.Fire[tank.Id] = RadarFireAll {
			Up: &RadarFire {
				Faith: 1.,
				Action: ActionFireUp,
				Cost: 10,
			},
		}
	}
	return ret
}
