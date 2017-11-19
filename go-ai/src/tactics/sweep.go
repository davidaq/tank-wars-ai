package tactics

import (
	// "fmt"
	f "framework"
	// "math/rand"
)

type TankGroup struct {
	fireStatus bool
	target f.Position
	leave f.Position
}

type Sweep struct {
	tankGroup map[string]*TankGroup
	gameCount int
	// fireStatus []bool
	// target []f.Position
}

func NewSweep() *Sweep {
	return &Sweep {
		// fireStatus: []bool{true,true,true},
		// target: []f.Position{{ X: 4, Y: 10 },{ X: 4, Y: 9 },{ X: 4, Y: 8 }},
		tankGroup: make(map[string]*TankGroup),
		gameCount: 0,
	}
}

func (self *Sweep) Init(state *f.GameState) {
	i:=0
	for _, tank := range state.MyTank {
		temp := TankGroup {
			fireStatus: false,
			target: f.Position {X:4,Y:10-i},
			leave: f.Position {X:4,Y:8-i},
		}
		self.tankGroup[tank.Id] = &temp
		i++
	}
}

func (self *Sweep) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	for _, tank := range state.MyTank {
		if tank.Pos.X == self.tankGroup[tank.Id].leave.X && tank.Pos.Y == self.tankGroup[tank.Id].leave.Y {
			self.tankGroup[tank.Id].fireStatus = true
		}
		if !self.tankGroup[tank.Id].fireStatus && self.gameCount>10 {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravelWithDodge,
				Target: f.Position { X: self.tankGroup[tank.Id].leave.X, Y: self.tankGroup[tank.Id].leave.Y },
			}
		} else {
			if tank.Pos.X == self.tankGroup[tank.Id].target.X && tank.Pos.Y == self.tankGroup[tank.Id].target.Y {
				objective[tank.Id] = f.Objective {
					Action: f.ActionFireRight,
				}
				self.tankGroup[tank.Id].fireStatus = false
				self.gameCount++
			} else {
				objective[tank.Id] = f.Objective {
					Action: f.ActionTravelWithDodge,
					Target: f.Position { X: self.tankGroup[tank.Id].target.X, Y: self.tankGroup[tank.Id].target.Y, Direction: f.DirectionUp },
				}
			}
		}
	}
	
}

func (self *Sweep) End(state *f.GameState) {
}