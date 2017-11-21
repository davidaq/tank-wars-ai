package tactics

import (
	f "framework"
)

type Nearest struct {
	round int
}

func NewNearest() *Nearest {
	return &Nearest { round: 0 }
}

func (self *Nearest) Init(state *f.GameState) {
}

func (self *Nearest) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	self.round++
	// if self.round < 5 {
	// 	return
	// }
	tankloop: for _, tank := range state.MyTank {
		fireRadar := radar.Fire[tank.Id]
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if fire != nil && fire.Sin < 0.5 && fire.Faith > 0.2 && tank.Bullet == "" {
				objective[tank.Id] = f.Objective {
					Action: fire.Action,
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
				break
			}
		}
		if ttank != nil {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				// Target: ttank.Pos,
				Target: f.Position { X: 0, Y: 0 },
			}
		}
	}
}

func (self *Nearest) End(state *f.GameState) {
}

func abs (val int) int {
	if val < 0 {
		return -val
	} else {
		return val
	}
}
