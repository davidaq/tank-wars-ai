package tactics

import (
	f "framework"
	// "fmt"
)

type Simplest struct {
	tankGroupA map[string]f.Tank
	tankGroupB map[string]f.Tank
}

func NewSimplest() *Simplest {
	return &Simplest {
		tankGroupA: make(map[string]f.Tank),
		tankGroupB: make(map[string]f.Tank),
	}
}

func (self *Simplest) Init(state *f.GameState) {
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

func (self *Simplest) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	n := 0
	// checker := false

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
		if tempRadarFire != nil && tempRadarFire.Sin < 0.5 && tempRadarFire.Faith >= 0.6 && tank.Bullet == "" {
			objective[tank.Id] = f.Objective {
				Action: tempRadarFire.Action,
			}
			continue tankloop
		}

		// 寻路
		least := 99999
		var ttank *f.Tank

		distance := state.Params.BulletSpeed +1
		patrolPos := []f.Position{
			{ X: state.Terain.Width/2-distance, Y: state.Terain.Height/2 },
			{ X: state.Terain.Width/2, Y: state.Terain.Height/2-distance },
			{ X: state.Terain.Width/2+distance, Y: state.Terain.Height/2 },
			{ X: state.Terain.Width/2, Y: state.Terain.Height/2+distance },
		}
		// 战斗A组
		if _, ok := self.tankGroupA[tank.Id]; ok {
			// flagPartol
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: patrolPos[(n-1)%4],
			}
		}
		// 战斗B组
		if _, ok := self.tankGroupB[tank.Id]; ok {
			// nearest
			for _, etank := range state.EnemyTank {
				dist := abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
				if dist < least {
					ttank = &etank
					least = dist
				}
			}
			if ttank != nil {
				p := ttank.Pos
				// var tbullet f.Bullet
				// for _, bullet := range state.EnemyBullet {
				// 	if bullet.From == ttank.Id {
				// 		tbullet = bullet
				// 	}
				// }
				// ebcost := caculateEnemyCost(tbullet, state.Terain, state.Params.BulletSpeed)
				// fmt.Println("bulletCost:",ebcost)
				switch p.Direction {
				case f.DirectionUp:
					p.Y -= 3 * state.Params.TankSpeed + 1
					if ttank.Bullet != "" {
						p.Y -= 2 * state.Params.TankSpeed
					}
				case f.DirectionDown:
					p.Y += 3 * state.Params.TankSpeed + 1
					if ttank.Bullet != "" {
						p.Y += 2 * state.Params.TankSpeed
					}
				case f.DirectionLeft:
					p.X -= 3 * state.Params.TankSpeed + 1
					if ttank.Bullet != "" {
						p.Y -= 2 * state.Params.TankSpeed
					}
				case f.DirectionRight:
					p.X += 3 * state.Params.TankSpeed + 1
					if ttank.Bullet != "" {
						p.Y += 2 * state.Params.TankSpeed
					}
				}
				objective[tank.Id] = f.Objective {
					Action: f.ActionTravel,
					Target: p,
				}
			}
		}

		// 草丛巡逻
		// if state.Terain.Data[tank.Pos.Y][tank.Pos.X] == 2 {
		// 	pos := forestPartol(tank.Pos, *state.Terain, state.Params.TankSpeed)
		// 	objective[tank.Id] = f.Objective {
		// 		Action: f.ActionTravel,
		// 		Target: pos,
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

func (self *Simplest) End(state *f.GameState) {
}

