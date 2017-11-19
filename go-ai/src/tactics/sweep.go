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
	tankGroupB map[string]*TankGroup
	gameCount int
}

func NewSweep() *Sweep {
	return &Sweep {
		tankGroup: make(map[string]*TankGroup),
		tankGroupB: make(map[string]*TankGroup),
		gameCount: 0,
	}
}

func (self *Sweep) Init(state *f.GameState) {
	i:=0
	for _, tank := range state.MyTank {
		if i<3 {
			temp := TankGroup {
				fireStatus: false,
				target: f.Position {X:4,Y:8+i},
				leave: f.Position {X:4,Y:6+i},
			}
			self.tankGroup[tank.Id] = &temp
		} else {
			temp := TankGroup {
				target: f.Position {X:4,Y:8+i},
			}
			self.tankGroupB[tank.Id] = &temp
		}
		i++
	}
}

func (self *Sweep) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	i := 0
	stage := 0
	for _, tank := range state.MyTank {
		if _,ok:=self.tankGroup[tank.Id];ok {
			if tank.Pos.X == self.tankGroup[tank.Id].leave.X && tank.Pos.Y == self.tankGroup[tank.Id].leave.Y {
				self.tankGroup[tank.Id].fireStatus = true
			}
			if !self.tankGroup[tank.Id].fireStatus && self.gameCount>20 {
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
		if _,ok:=self.tankGroupB[tank.Id];ok {
			if tank.Pos.X == self.tankGroupB[tank.Id].target.X && tank.Pos.Y == self.tankGroupB[tank.Id].target.Y {
				if stage == 0 {
					objective[tank.Id] = f.Objective {
						Action: f.ActionFireRight,
					}
					self.tankGroupB[tank.Id].target = f.Position {X:9+i,Y:12,Direction:f.DirectionUp}
					i++
					stage = 1
				}
				if stage == 1 {
					objective[tank.Id] = f.Objective {
						Action: f.ActionFireUp,
					}
				}
			}
		}
	}
	
}

func (self *Sweep) End(state *f.GameState) {
}