package tactics

import (
	f "framework"
)

type Fox struct {
	tankGroupA map[string]f.Tank
	tankGroupB map[string]f.Tank
}

func NewFox() *Fox {
	return &Fox {
		tankGroupA: make(map[string]f.Tank),
		tankGroupB: make(map[string]f.Tank),
	}
}

func (self *Fox) Init(state *f.GameState) {
	i:=0
	for _, tank := range state.MyTank {
		if i<3 {
			self.tankGroupA[tank.Id] = tank
		} else {
			self.tankGroupB[tank.Id] = tank
		}
		i++
	}
}

func (self *Fox) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	tankloop: for _, tank := range state.MyTank {
		// 开火
		fireRadar := radar.Fire[tank.Id]
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if fire != nil && fire.Sin < 0.7 && fire.Faith >= 0.3 && tank.Bullet == "" {
				objective[tank.Id] = f.Objective {
					Action: fire.Action,
				}
				continue tankloop
			}
		}

		// 躲避
		if radar.Dodge[tank.Id].Threat > 0.7 {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: radar.Dodge[tank.Id].SafePos,
			}
			continue tankloop
		}

		// 寻路
		least := 99999
		furthest := -99999
		var ttank *f.Tank
		// 战斗A组
		if _, ok := self.tankGroupA[tank.Id]; ok {
			for _, etank := range state.EnemyTank {
				dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
				if dist < least {
					ttank = &etank
					least = dist
				}
			}
			if ttank != nil {
				objective[tank.Id] = f.Objective {
					Action: f.ActionTravel,
					Target: ttank.Pos,
				}
			}
		}
		// 战斗B组
		if _, ok := self.tankGroupB[tank.Id]; ok {
			for _, etank := range state.EnemyTank {
				dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
				if dist > furthest {
					ttank = &etank
					furthest = dist
				}
			}
			if ttank != nil {
				objective[tank.Id] = f.Objective {
					Action: f.ActionTravel,
					Target: ttank.Pos,
				}
			}
		}

		// 夺旗
		if state.FlagWait == 0 {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravelWithDodge,
				Target: f.Position { X: state.Terain.Width/2, Y: state.Terain.Height/2 },
			}
		}

	}
}

func (self *Fox) End(state *f.GameState) {
}
