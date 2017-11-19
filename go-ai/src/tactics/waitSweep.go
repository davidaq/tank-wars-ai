package tactics

import (
	f "framework"
	"math"
)

type TankTarget struct {
	target f.Position
	stage int
}

type WaitSweep struct {
	tankGroupA map[string]*TankTarget
	tankGroupB map[string]*TankTarget
}

func NewWaitSweep() *WaitSweep {
	return &WaitSweep {
		tankGroupA: make(map[string]*TankTarget),
		tankGroupB: make(map[string]*TankTarget),
	}
}

func (self *WaitSweep) Init(state *f.GameState) {
	i:=0
	for _, tank := range state.MyTank {
		if i<3 {
			temp := TankTarget {
				target: f.Position {X:5,Y:8+i},
			}
			self.tankGroupA[tank.Id] = &temp
		} else {
			temp := TankTarget {
				target: f.Position {X:5,Y:8+i},
				stage: 0,
			}
			self.tankGroupB[tank.Id] = &temp
		}
		i++
	}
}

func (self *WaitSweep) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	i := 0
	for _, tank := range state.MyTank {
		_, flag := self.tankGroupA[tank.Id]
		if flag {
			if tank.Pos.X == self.tankGroupA[tank.Id].target.X && tank.Pos.Y == self.tankGroupA[tank.Id].target.Y {
				objective[tank.Id] = f.Objective {
					Action: f.ActionFireRight,
				}
			} else {
				objective[tank.Id] = f.Objective {
					Action: f.ActionTravelWithDodge,
					Target: f.Position { X: self.tankGroupA[tank.Id].target.X, Y: self.tankGroupA[tank.Id].target.Y, Direction: f.DirectionUp },
				}
			}
		} else {
			if self.tankGroupB[tank.Id].stage == 1 {
				if tank.Pos.X == self.tankGroupB[tank.Id].target.X && tank.Pos.Y == self.tankGroupB[tank.Id].target.Y {
					objective[tank.Id] = f.Objective {
						Action: f.ActionFireUp,
					}
				} else {
					objective[tank.Id] = f.Objective {
						Action: f.ActionTravelWithDodge,
						Target: f.Position { X: self.tankGroupB[tank.Id].target.X, Y: self.tankGroupB[tank.Id].target.Y, Direction: f.DirectionUp },
					}
				}
			} else {
				if tank.Pos.X == self.tankGroupB[tank.Id].target.X && tank.Pos.Y == self.tankGroupB[tank.Id].target.Y {
					objective[tank.Id] = f.Objective {
						Action: f.ActionFireRight,
					}
					self.tankGroupB[tank.Id].target = f.Position {X:9+i,Y:12,Direction:f.DirectionUp}
					self.tankGroupB[tank.Id].stage = 1
				} else {
					objective[tank.Id] = f.Objective {
						Action: f.ActionTravelWithDodge,
						Target: f.Position { X: self.tankGroupB[tank.Id].target.X, Y: self.tankGroupB[tank.Id].target.Y, Direction: f.DirectionRight },
					}
				}
			}
			if state.FlagWait == 0 {
				if len(self.tankGroupB) == 2 {
					if tank.Pos.X == 9 {
						objective[tank.Id] = f.Objective {
							Action: f.ActionTravel,
							Target: f.Position { X: 9, Y: 9 },
						}
						self.tankGroupB[tank.Id].stage = 0
					}
				} else {
					objective[tank.Id] = f.Objective {
						Action: f.ActionTravelWithDodge,
						Target: f.Position { X: 9, Y: 9 },
					}
					self.tankGroupB[tank.Id].stage = 0
				}
			}
			if len(radar.ForestThreat) > 0 {
				for k := range(radar.ForestThreat) {
					if math.Sqrt(float64((k.X-tank.Pos.X)*(k.X-tank.Pos.X)+(k.Y-tank.Pos.Y)*(k.Y-tank.Pos.Y))) <=1 {
						if k.X == tank.Pos.X {
							if k.Y < tank.Pos.Y {
								objective[tank.Id] = f.Objective {
									Action: f.ActionFireUp,
								}
							} else {
								objective[tank.Id] = f.Objective {
									Action: f.ActionFireDown,
								}
							}
						}
						if k.Y == tank.Pos.Y {
							if k.X < tank.Pos.X {
								objective[tank.Id] = f.Objective {
									Action: f.ActionFireLeft,
								}
							} else {
								objective[tank.Id] = f.Objective {
									Action: f.ActionFireRight,
								}
							}
						}
					}
				}
			} 
			i++
		}
		
	}
}

func (self *WaitSweep) End(state *f.GameState) {
}
