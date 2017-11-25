package tactics

import (
	f "framework"
	// "math"
)

type Terminator struct {
	obs *Observation
	Roles map[string]*CattyRole
	tankGroupA map[string] f.Tank
	tankGroupB map[string] f.Tank
}


func NewTerminator() *Terminator {
	return &Terminator {
		Roles: make(map[string]*CattyRole),
		tankGroupA: make(map[string]f.Tank),
		tankGroupB: make(map[string]f.Tank),
	}
}

func (self *Terminator) Init(state *f.GameState) {
	i:=0
	self.obs = NewObservation(state)
	for _, tank := range state.MyTank {
		self.Roles[tank.Id] = &CattyRole { obs: self.obs, Tank: tank}
		if i<2 {
			self.tankGroupA[tank.Id] = tank
		} else {
			self.tankGroupB[tank.Id] = tank
		}
		i++
	}
}

func (self *Terminator) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	n := 0
	// checker := false

	for tankid := range objective {
		delete(objective, tankid)
	}

	self.obs.makeObservation(state, radar, objective)

	self.updateRole()
	
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
		// if radar.DodgeBullet[tank.Id].Threat == -1 && tempRadarFire != nil && tempRadarFire.Sin < 0.5 {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: tempRadarFire.Action,
		// 	}
		// 	continue tankloop
		// }

		// 子弹躲避
		// if radar.DodgeBullet[tank.Id].Threat > 0.2 {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: f.ActionTravel,
		// 		Target: radar.DodgeBullet[tank.Id].SafePos,
		// 	}
		// 	continue tankloop
		// }

		// 无子弹躲避
		// if tank.Bullet != "" {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: f.ActionTravelWithDodge,
		// 		Target: radar.Dodge[tank.Id].SafePos,
		// 	}
		// 	continue tankloop
		// }

		// 开火
		// if tempRadarFire != nil && tempRadarFire.Sin < 0.5 && tempRadarFire.Faith > 0.2 && tank.Bullet == "" {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: tempRadarFire.Action,
		// 	}
		// 	continue tankloop
		// }

		// 寻路
		// least := 99999
		// furthest := -99999
		// var ttank *f.Tank
		distance := state.Terain.Width/6
		patrolPos := []f.Position{
			{ X: state.Terain.Width/2-distance, Y: state.Terain.Height/2 },
			{ X: state.Terain.Width/2, Y: state.Terain.Height/2-distance },
			{ X: state.Terain.Width/2+distance, Y: state.Terain.Height/2 },
			{ X: state.Terain.Width/2, Y: state.Terain.Height/2+distance },
		}
		// 战斗A组
		if _, ok := self.tankGroupA[tank.Id]; ok {
			// faith排序
			// fireRadar := radar.Fire[tank.Id]
			// var tempRadarFire *f.RadarFire
			// for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			// 	if tempRadarFire == nil {
			// 		tempRadarFire = fire
			// 	}
			// 	if fire!=nil {
			// 		if fire.Faith > tempRadarFire.Faith {
			// 			tempRadarFire = fire
			// 		}
			// 	}
			// }

			// 光荣弹开火
			if radar.DodgeBullet[tank.Id].Threat == -1 && tempRadarFire != nil && tempRadarFire.Sin < 0.5 {
				objective[tank.Id] = f.Objective {
					Action: tempRadarFire.Action,
				}
				continue tankloop
			}

			// 子弹躲避
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
					Target: tank.Pos,
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

			// flagPartol
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: patrolPos[(n-1)%4],
			}

			// Stalker
			// for _, etank := range state.EnemyTank {
			// 	dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
			// 	if dist < least {
			// 		ttank = &etank
			// 		least = dist
			// 	}
			// }
			// if ttank != nil {
			// 	resPos := ttank.Pos
			// 	mid := state.Terain.Width/2
			// 	dis := state.Params.BulletSpeed * 3+1
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
			role := self.Roles[tank.Id]
			if self.obs.Flag.Exist && self.obs.Flag.Next <= 5 {
				role.occupyFlag()
				continue tankloop
			}
			role.hunt()
                        role.act()
			//if role.Tank.Bullet != "" || !role.checkDone() {
			//	role.move()
			//} else {
			//	role.act()
			//}
		}

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

func (self *Terminator) End(state *f.GameState) {
}

func (self *Terminator) updateRole() {
	for id, role := range self.Roles {
		if self.obs.MyTank[id] != (f.Tank{}) {
			role.Tank  = self.obs.MyTank[id]
			role.Dodge = self.obs.Radar.DodgeBullet[id]
            role.Fire  = self.obs.Radar.Fire[id]
		} else {
			delete(self.Roles, id)
		}

		if role.Target != (CattyTarget{}) && role.Target.Tank != (f.Tank{}) {
			if role.obs.EmyTank[role.Target.Tank.Id] == (f.Tank{}) {
				role.Target = CattyTarget{}
			} else {
				role.Target.Tank = role.obs.EmyTank[role.Target.Tank.Id]
			}
		}
	}
}
