package tactics

import (
	f "framework"
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
	tankloop: for _, tank := range state.MyTank {
		n++
		// 躲避
		if radar.DodgeBullet[tank.Id].Threat > 0.7 {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravelWithDodge,
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
		if tempRadarFire != nil && tempRadarFire.Sin < 0.5 && tempRadarFire.Faith > 0.2 && tank.Bullet == "" {
			objective[tank.Id] = f.Objective {
				Action: tempRadarFire.Action,
			}
			continue tankloop
		}

		// 躲避
		// if radar.DodgeEnemy[tank.Id].Threat > 0.9 {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: f.ActionTravel,
		// 		Target: radar.DodgeBullet[tank.Id].SafePos,
		// 	}
		// 	continue tankloop
		// }

		// 寻路
		// least := 99999
		// furthest := -99999
		// var ttank *f.Tank
		patrolPos := []f.Position{
			{ X: state.Terain.Width/2-5, Y: state.Terain.Height/2 },
			{ X: state.Terain.Width/2, Y: state.Terain.Height/2-5 },
			{ X: state.Terain.Width/2+5, Y: state.Terain.Height/2 },
			{ X: state.Terain.Width/2, Y: state.Terain.Height/2+5 },
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
		}

		// 夺旗
		if len(self.tankGroupA) > 0 {
			if state.FlagWait <= 8 {
				if _, ok := self.tankGroupA[tank.Id]; ok {
					objective[tank.Id] = f.Objective {
						Action: f.ActionTravel,
						Target: f.Position { X: state.Terain.Width/2, Y: state.Terain.Height/2 },
					}
				}
			}
		} else {
			if state.FlagWait <= 8 {
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
