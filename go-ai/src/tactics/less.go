package tactics

import (
	f "framework"
	"math/rand"
	"fmt"
)

type Less struct {
	prevTarget map[string]*f.Position
	justFired map[string]int
	justFiredY map[string]bool
	enemyY map[int]int
	round int
}

func NewLess() *Less {
	inst := &Less {
		prevTarget: make(map[string]*f.Position),
		justFired: make(map[string]int),
		justFiredY: make(map[string]bool),
		enemyY: make(map[int]int),
		round: 0,
	}
	return inst
}

func (self *Less) Init(state *f.GameState) {
}

func (self *Less) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	self.round++
	fireForest := make(map[f.Position]FireForest)
	for _, bullet := range state.EnemyBullet {
		if bullet.Pos.Y >= 8 && bullet.Pos.Y <= 12 && bullet.Pos.Direction == f.DirectionLeft {
			self.enemyY[bullet.Pos.Y] = self.round
		}
	}
	isFirst := true
	tankloop: for _, tank := range state.MyTank {
		preferDodge := false
		fireRadar := radar.Fire[tank.Id]
		if fired, ok := self.justFired[tank.Id]; !ok || self.round - fired > 1 {
			fireForest[f.Position { X: tank.Pos.X - 1, Y: tank.Pos.Y }] = FireForest { tank.Id, f.ActionFireLeft }
			fireForest[f.Position { X: tank.Pos.X + 1, Y: tank.Pos.Y }] = FireForest { tank.Id, f.ActionFireRight }
			fireForest[f.Position { X: tank.Pos.X, Y: tank.Pos.Y - 1 }] = FireForest { tank.Id, f.ActionFireUp }
			fireForest[f.Position { X: tank.Pos.X, Y: tank.Pos.Y + 1}] = FireForest { tank.Id, f.ActionFireDown }
			for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
				if fire != nil && fire.Sin < 0.2 && fire.Faith > 0 {
					objective[tank.Id] = f.Objective {
						Action: fire.Action,
					}
					self.justFired[tank.Id] = self.round
					if fire.Action == f.ActionFireUp || fire.Action == f.ActionFireDown {
						self.justFiredY[tank.Id] = true
					}
					continue tankloop
				}
			}
		}
		if isFirst {
			isFirst = false
			if state.FlagWait < 3 {
				objective[tank.Id] = f.Objective {
					Target: f.Position {
						X: state.Terain.Width / 2,
						Y: state.Terain.Height / 2,
					},
					Action: f.ActionTravelWithDodge,
				}
				continue tankloop
			}
		}

		least := 99999
		var ttank *f.Tank
		if self.round > 20 {
			for _, etank := range state.EnemyTank {
				dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
				if etank.Pos.Y >= 6 && etank.Pos.X <= 12 {
					if dist < least {
						ttank = &etank
						least = dist
					}
				}
			}
		}
		pos := f.Position {}
		if fired, ok := self.justFired[tank.Id]; ok && self.round - fired < 12 / state.Params.BulletSpeed {
			if self.justFiredY[tank.Id] {
				preferDodge = true
				pos.Y = tank.Pos.Y
				if tank.Pos.X > 5 {
					pos.X = tank.Pos.X - 1
				} else {
					pos.X = tank.Pos.X + 1
				}
			} else {
				pos.X = tank.Pos.X
				if tank.Pos.Y > 3 && tank.Pos.Y < 12 {
					if tank.Pos.Y == 9 {
						if tank.Pos.Direction == f.DirectionUp {
							pos.Y = tank.Pos.Y - 1
						} else {
							pos.Y = tank.Pos.Y + 1
						}
					} else if tank.Pos.Y >= 10 {
						pos.Y = tank.Pos.Y + 1
					} else {
						pos.Y = tank.Pos.Y - 1
					}
				} else {
					pos.Y = tank.Pos.Y + 1
				}
			}
		} else if ttank != nil {
			pos = ttank.Pos
			if tank.Pos.X <= 8 && pos.Y == tank.Pos.Y {
				if fireRadar.Right != nil && fireRadar.Right.Sin < 0.1 {
					objective[tank.Id] = f.Objective {
						Action: f.ActionFireRight,
					}
					self.justFired[tank.Id] = self.round
					continue tankloop
				}
			} else {
				// 5,8 - 10,12
				if pos.X < 5 {
					pos.X = 5
				}
				if pos.X > 10 {
					pos.X = 10
				}
				if pos.Y < 8 {
					pos.Y = 8
				}
				if pos.Y > 12 {
					pos.Y = 12
				}	
			}
		} else {
			if r, ok := self.enemyY[tank.Pos.Y]; ok && self.round - r < 12 / state.Params.BulletSpeed {
				if fireRadar.Right != nil && fireRadar.Right.Sin < 0.1 {
					self.justFired[tank.Id] = self.round
					objective[tank.Id] = f.Objective {
						Action: f.ActionFireRight,
					}
					continue tankloop
				}
			}
			target := self.prevTarget[tank.Id]
			if target == nil || target.SDist(tank.Pos) < state.Params.TankSpeed || rand.Int() % 5 == 0 {
				y := rand.Int() % (12 - 8) + 8
				if self.round < 20 {
					fmt.Println("hurry")
					if y > 8 {
						y = 8
					}
				}
				target = &f.Position {
					X: 5,
					Y: y,
					Direction: f.DirectionRight,
				}
				self.prevTarget[tank.Id] = target
			}
			pos = *target
			if tank.Pos.X > 4 && tank.Pos.X < 7 && state.Terain.Get(tank.Pos.X, tank.Pos.Y) == 2 && rand.Int() % 3 == 0 {
				if fireRadar.Right != nil && fireRadar.Right.Sin < 0.1 {
					objective[tank.Id] = f.Objective {
						Action: f.ActionFireRight,
					}
					self.justFired[tank.Id] = self.round
					continue tankloop
				}
			}
			fmt.Println(pos)
		}
		action := f.ActionTravel
		_ = preferDodge
		// if preferDodge {
			action = f.ActionTravelWithDodge
		// }
		objective[tank.Id] = f.Objective {
			Action: action,
			Target: pos,
		}
	}
	for position, posibility := range radar.ForestThreat {
		if posibility > 0.9 {
			if fire, ok := fireForest[f.Position { X: position.X, Y: position.Y }]; ok {
				objective[fire.tankId] = f.Objective {
					Action: fire.action,
				}
			}
		}
	}
	if state.FlagWait < 3 {
		for _, tank := range state.MyTank {
			obj := objective[tank.Id]
			if obj.Action == f.ActionFireRight && tank.Pos.Y == state.Terain.Height / 2 {
				objective[tank.Id] = f.Objective {
					Action: f.ActionStay,
				}
			}
		}
	}
}

func (self *Less) End(state *f.GameState) {
}
