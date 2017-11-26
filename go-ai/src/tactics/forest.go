package tactics

import (
	f "framework"
	// "math/rand"
	// "fmt"
)

type Forest struct {
}

func NewForest() *Forest {
	inst := &Forest {}
	return inst
}

func (self *Forest) Init(state *f.GameState) {
}

func (self *Forest) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	
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

		// least := 99999
		// var ttank *f.Tank
		// for _, etank := range state.EnemyTank {
		// 	dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
		// 	if dist < least {
		// 		ttank = &etank
		// 		least = dist
		// 	}
		// }
		// if ttank != nil {
		// 	travel := f.ActionTravel
		// 	if radar.Dodge[tank.Id].Threat > 0.9 && rand.Int() % 3 > 0 {
		// 		travel = f.ActionTravelWithDodge
		// 	}
		// 	p := ttank.Pos
		// 	switch p.Direction {
		// 	case f.DirectionUp:
		// 		p.Y -= state.Params.TankSpeed
		// 	case f.DirectionDown:
		// 		p.Y += state.Params.TankSpeed
		// 	case f.DirectionLeft:
		// 		p.X -= state.Params.TankSpeed
		// 	case f.DirectionRight:
		// 		p.X += state.Params.TankSpeed
		// 	}
		// 	objective[tank.Id] = f.Objective {
		// 		Action: travel,
		// 		Target: p,
		// 		// Target: f.Position { X: 0, Y: 0 },
		// 	}
		// }
		p := f.Position {
			X: state.Terain.Width/2,
			Y: state.Terain.Height/2,
			Direction: f.DirectionUp,
		}
		objective[tank.Id] = f.Objective {
			Action: f.ActionTravel,
			Target: p,
		}

		// 草丛巡逻
		if state.Terain.Data[tank.Pos.Y][tank.Pos.X] == 2 {
			pos := forestPartol(tank.Pos, *state.Terain, state.Params.TankSpeed)
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: pos,
			}
		}
	}

}

func (self *Forest) End(state *f.GameState) {
}
