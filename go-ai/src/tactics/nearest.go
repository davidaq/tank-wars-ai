package tactics

import (
	f "framework"
	"math/rand"
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
		faith := 0.
		var pfire *f.RadarFire
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if fire != nil && fire.Sin < 0.5 && fire.Faith > faith {
				pfire = fire
			}
		}
		if pfire != nil {
			objective[tank.Id] = f.Objective {
				Action: pfire.Action,
			}
			continue tankloop
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
			travel := f.ActionTravel
			_ = rand.Int()
			// if radar.Dodge[tank.Id].Threat > 0.9 && rand.Int() % 3 > 0 {
			// 	travel = f.ActionTravelWithDodge
			// }
			p := ttank.Pos
			switch p.Direction {
			case f.DirectionUp:
				p.Y -= state.Params.TankSpeed
			case f.DirectionDown:
				p.Y += state.Params.TankSpeed
			case f.DirectionLeft:
				p.X -= state.Params.TankSpeed
			case f.DirectionRight:
				p.X += state.Params.TankSpeed
			}
			objective[tank.Id] = f.Objective {
				Action: travel,
				Target: p,
				// Target: f.Position { X: 15, Y: 15 },
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
