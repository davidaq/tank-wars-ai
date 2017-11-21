package tactics

import (
	f "framework"
)

type Fox struct {
}

func NewFox() *Fox {
	return &Fox {}
}

func (self *Fox) Init(state *f.GameState) {
}

func (self *Fox) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	tankloop: for _, tank := range state.MyTank {
		fireRadar := radar.Fire[tank.Id]
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if fire != nil && fire.Sin < 0.4 && fire.Faith >= 0.3 && tank.Bullet == "" {
				objective[tank.Id] = f.Objective {
					Action: fire.Action,
				}
				if radar.Dodge[tank.Id].Threat == 1 {
					objective[tank.Id] = f.Objective {
						Action: f.ActionTravel,
						Target: radar.Dodge[tank.Id].SafePos,
					}
				}
				continue tankloop
			}
		}

		least := 99999
		var ttank *f.Tank
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
				// Target: f.Position { X: 14, Y: 9 },
			}
		}
	}
}

func (self *Fox) End(state *f.GameState) {
}
