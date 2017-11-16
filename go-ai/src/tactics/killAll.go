package tactics

import (
	"fmt";
	f "framework";
)

type KillAll struct {
}

func NewKillAll() *KillAll {
	return &KillAll {}
}

func (self *KillAll) Init(state *f.GameState) {
}

func (self *KillAll) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	tankloop: for _, tank := range state.MyTank {
		least := 99999
		var ttank *f.Tank
		for _, etank := range state.EnemyTank {
			dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
			if dist < least {
				ttank = &etank
				least = dist
			}
		}

		fmt.Println("----------------------------", tank.Id, radar.Dodge[tank.Id].Threat)
		if radar.Dodge[tank.Id].Threat >= 0.2 && radar.Dodge[tank.Id].Threat < 1 {
			ifDodge := true
			for _, oTank := range state.MyTank {
				if oTank.Id == tank.Id {
					continue
				}
				if oTank.Pos.X == radar.Dodge[tank.Id].SafePos.X && oTank.Pos.Y == radar.Dodge[tank.Id].SafePos.Y {
					ifDodge = false
					break
				}
			}
			if ifDodge {
				objective[tank.Id] = f.Objective {
					Action: f.ActionTravel,
					Target: radar.Dodge[tank.Id].SafePos,
				}
				continue tankloop
			}
		}
		fireRadar := radar.Fire[tank.Id]
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if fire != nil && fire.Sin < 0.5 && fire.Faith > 0 {
				objective[tank.Id] = f.Objective {
					Action: fire.Action,
				}
				continue tankloop
			}
		}


		if ttank != nil {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: ttank.Pos,
			}
		}
	}
}

func (self *KillAll) End(state *f.GameState) {
}