package tactics

import (
	f "framework"
)

type Nearest struct {
}

func NewNearest() *Nearest {
	return &Nearest {}
}

func (self *Nearest) Init(state *f.GameState) {
}

func (self *Nearest) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	tankloop: for _, tank := range state.MyTank {
		fireRadar := radar.Fire[tank.Id]
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if fire != nil && fire.Sin < 0.5 && fire.Faith > 0.5 {
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
			}
		}
		if ttank != nil {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravelWithDodge,
				Target: ttank.Pos,
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
