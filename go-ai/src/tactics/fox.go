package tactics

import (
	f "framework"
	// "fmt"
)

type Fox struct {
	tankGroupA map[string]f.Tank
	tankGroupB map[string]f.Tank
}

func NewFox() *Fox {
	return &Fox {
		tankGroupA: make(map[string]f.Tank),
		tankGroupB: make(map[string]f.Tank),
	}
}

func caculateQuadrant (mid int, pos f.Position) int {
	if pos.X > mid && pos.Y < mid {
		return 1
	} else if pos.X < mid && pos.Y < mid {
		return 2
	} else if pos.X < mid && pos.Y > mid {
		return 3
	} else if pos.X > mid && pos.Y > mid {
		return 4
	}
	return 0
}

func (self *Fox) Init(state *f.GameState) {
	i:=0
	for _, tank := range state.MyTank {
		if i<2 {
			self.tankGroupA[tank.Id] = tank
		} else {
			self.tankGroupB[tank.Id] = tank
		}
		i++
	}
}

func (self *Fox) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	n := 0
	checker := false

	// 分组存活判断
	tempTankGroupA := make(map[string]f.Tank)
	tempTankGroupB := make(map[string]f.Tank)
	for _, tank := range state.MyTank {
		if _, ok := self.tankGroupA[tank.Id]; ok {
			tempTankGroupA[tank.Id] = tank
		}
		if _, ok := self.tankGroupB[tank.Id]; ok {
			tempTankGroupB[tank.Id] = tank
		}
	}
	self.tankGroupA = tempTankGroupA
	self.tankGroupB = tempTankGroupB

	tankloop: for _, tank := range state.MyTank {
		n++

		// 动态分组
		if len(self.tankGroupA) <= 1 && len(self.tankGroupB) <= 1 {
			self.tankGroupB[tank.Id] = tank
			delete(self.tankGroupA, tank.Id)
		}

		// faith排序
		fireRadar := radar.Fire[tank.Id]
		var tempRadarFire *f.RadarFire
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if tempRadarFire == nil {
				tempRadarFire = fire
			}
			if fire!=nil {
				if fire.Faith > tempRadarFire.Faith {
					tempRadarFire = fire
				}
			}
		}

		// 光荣弹开火
		if radar.DodgeBullet[tank.Id].Threat == -1 && tempRadarFire != nil && tempRadarFire.Sin < 0.5 {
			objective[tank.Id] = f.Objective {
				Action: tempRadarFire.Action,
			}
			continue tankloop
		}

		// 子弹躲避
		// if _, ok := self.tankGroupA[tank.Id]; ok {
		// 	if radar.DodgeBullet[tank.Id].Threat > 0.2 {
		// 		objective[tank.Id] = f.Objective {
		// 			Action: f.ActionTravel,
		// 			Target: radar.DodgeBullet[tank.Id].SafePos,
		// 		}
		// 		continue tankloop
		// 	}
		// }

		// if _, ok := self.tankGroupB[tank.Id]; ok {
		// 	if radar.DodgeBullet[tank.Id].Threat > 0.2 {
		// 		objective[tank.Id] = f.Objective {
		// 			Action: f.ActionTravel,
		// 			Target: radar.DodgeBullet[tank.Id].SafePos,
		// 		}
		// 		continue tankloop
		// 	}
		// }
		if radar.DodgeBullet[tank.Id].Threat > 0.2 {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: radar.DodgeBullet[tank.Id].SafePos,
			}
			continue tankloop
		}

		// 无子弹躲避
		if tank.Bullet != "" {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravelWithDodge,
				Target: radar.Dodge[tank.Id].SafePos,
			}
			continue tankloop
		}

		// 开火
		if tempRadarFire != nil && tempRadarFire.Sin < 0.5 && tempRadarFire.Faith > 0.2 && tank.Bullet == "" {
			objective[tank.Id] = f.Objective {
				Action: tempRadarFire.Action,
			}
			continue tankloop
		}

		// 寻路
		least := 99999
		// furthest := -99999
		var ttank *f.Tank
		distance := state.Terain.Width/6
		patrolPos := []f.Position{
			{ X: state.Terain.Width/2-distance, Y: state.Terain.Height/2 },
			{ X: state.Terain.Width/2, Y: state.Terain.Height/2-distance },
			{ X: state.Terain.Width/2+distance, Y: state.Terain.Height/2 },
			{ X: state.Terain.Width/2, Y: state.Terain.Height/2+distance },
		}
		// 战斗A组
		if _, ok := self.tankGroupA[tank.Id]; ok {
			// nearest
			// for _, etank := range state.EnemyTank {
			// 	dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
			// 	if dist < least {
			// 		ttank = &etank
			// 		least = dist
			// 	}
			// }
			// if ttank != nil {
			// 	objective[tank.Id] = f.Objective {
			// 		Action: f.ActionTravel,
			// 		Target: ttank.Pos,
			// 	}
			// }

			// flagPartol
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: patrolPos[(n-1)%4],
			}
			// if radar.DodgeEnemy[tank.Id].Threat > 0.9 {
			// 	objective[tank.Id] = f.Objective {
			// 		Action: f.ActionTravel,
			// 		Target: patrolPos[n%4],
			// 	}
			// }

			// Stalker
			// for _, etank := range state.EnemyTank {
			// 	dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
			// 	if dist > furthest {
			// 		ttank = &etank
			// 		least = dist
			// 	}
			// }
			// if ttank != nil {
			// 	resPos := ttank.Pos
			// 	mid := state.Terain.Width/2
			// 	dis := state.Params.BulletSpeed * 2+1
			// 	targetQuadrant := caculateQuadrant(mid, ttank.Pos)
			// 	switch targetQuadrant {
			// 	case 0:
			// 		break
			// 	case 1:
			// 		if !checker {
			// 			resPos.X -= dis
			// 		} else {
			// 			resPos.Y += dis
			// 		}
			// 		checker = true
			// 	case 2:
			// 		if !checker {
			// 			resPos.X += dis
			// 		} else {
			// 			resPos.Y += dis
			// 		}
			// 		checker = true
			// 	case 3:
			// 		if !checker {
			// 			resPos.X += dis
			// 		} else {
			// 			resPos.Y -= dis
			// 		}
			// 		checker = true
			// 	case 4:
			// 		if !checker {
			// 			resPos.X -= dis
			// 		} else {
			// 			resPos.Y -= dis
			// 		}
			// 		checker = true
			// 	}
			// 	objective[tank.Id] = f.Objective {
			// 		Action: f.ActionTravel,
			// 		Target: resPos,
			// 	}	
			// }
		}
		// 战斗B组
		if _, ok := self.tankGroupB[tank.Id]; ok {
			// furthest
			// for _, etank := range state.EnemyTank {
			// 	dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
			// 	if dist > furthest {
			// 		ttank = &etank
			// 		furthest = dist
			// 	}
			// }
			// if ttank != nil {
			// 	objective[tank.Id] = f.Objective {
			// 		Action: f.ActionTravel,
			// 		Target: ttank.Pos,
			// 	}
			// }

			// flagPartol
			// objective[tank.Id] = f.Objective {
			// 	Action: f.ActionTravel,
			// 	Target: patrolPos[(n-1)%4],
			// }
			// if radar.DodgeEnemy[tank.Id].Threat > 0.9 {
			// 	objective[tank.Id] = f.Objective {
			// 		Action: f.ActionTravel,
			// 		Target: patrolPos[n%4],
			// 	}
			// }

			// Stalker
			for _, etank := range state.EnemyTank {
				dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
				if dist < least {
					ttank = &etank
					least = dist
				}
			}
			if ttank != nil {
				resPos := ttank.Pos
				mid := state.Terain.Width/2
				dis := state.Params.BulletSpeed * 3+1
				targetQuadrant := caculateQuadrant(mid, ttank.Pos)
				switch targetQuadrant {
				case 0:
					break
				case 1:
					if !checker {
						resPos.X -= dis
					} else {
						resPos.Y += dis
					}
					checker = true
				case 2:
					if !checker {
						resPos.X += dis
					} else {
						resPos.Y += dis
					}
					checker = true
				case 3:
					if !checker {
						resPos.X += dis
					} else {
						resPos.Y -= dis
					}
					checker = true
				case 4:
					if !checker {
						resPos.X -= dis
					} else {
						resPos.Y -= dis
					}
					checker = true
				}
				objective[tank.Id] = f.Objective {
					Action: f.ActionTravel,
					Target: resPos,
				}
				if abs(tank.Pos.X - ttank.Pos.X) < dis/3 ||	abs(tank.Pos.Y - ttank.Pos.Y) < dis/3 {
					objective[tank.Id] = f.Objective {
						Action: f.ActionTravel,
						Target: patrolPos[(n-1)%4],
					}
				}
			}
		}

		// 坦克躲避
		// if radar.Dodge[tank.Id].Threat == 1 {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: f.ActionTravel,
		// 		Target: radar.Dodge[tank.Id].SafePos,
		// 	}
		// }

		// 夺旗
		if len(self.tankGroupA) > 0 {
			if state.FlagWait <= 5 {
				if _, ok := self.tankGroupA[tank.Id]; ok {
					objective[tank.Id] = f.Objective {
						Action: f.ActionTravel,
						Target: f.Position { X: state.Terain.Width/2, Y: state.Terain.Height/2 },
					}
				}
			}
		} else {
			if state.FlagWait <= 5 {
				if _, ok := self.tankGroupB[tank.Id]; ok {
					objective[tank.Id] = f.Objective {
						Action: f.ActionTravel,
						Target: f.Position { X: state.Terain.Width/2, Y: state.Terain.Height/2 },
					}
				}
			}
		}
	}
}

func (self *Fox) End(state *f.GameState) {
}
