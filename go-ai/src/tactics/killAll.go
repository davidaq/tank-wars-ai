package tactics

import (
	"fmt"
	f "framework";
)

type KillAll struct {
}

func NewKillAll() *KillAll {
	return &KillAll {}
}

func (self *KillAll) Init(state *f.GameState) {
}

var count = 0
func (self *KillAll) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	count++
	// count := 0
	// targetX := 6
	// targetY := 8
	tankloop: for _, tank := range state.MyTank {
		delete(objective, tank.Id)
		// if radar.Dodge[tank.Id].Threat >= 0.4 && radar.Dodge[tank.Id].Threat < 1 && (tank.Pos.X != targetX && tank.Pos.Y != targetY + count) {
		fmt.Println("-----------------------", count, tank.Id, radar.Dodge[tank.Id].Threat)
		if tank.Bullet != "" {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravelWithDodge,
			}
			continue tankloop
		}
		if radar.Dodge[tank.Id].Threat >= 0.4 && radar.Dodge[tank.Id].Threat <= 1 {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: radar.Dodge[tank.Id].SafePos,
			}
			continue tankloop
		} 
		// else if radar.Dodge[tank.Id].Threat == -1 {
		// 	for _, bullet := range radar.ExtDangerSrc[tank.Id] {
		// 		if bullet.Urgent == -1 {
		// 			tankInVis := true
		// 			for _, e := range state.EnemyTank {
		// 				if bullet.Source == e.Id {
		// 					tankInVis = false
		// 					break
		// 				}
		// 			}
		// 			if tankInVis {
		// 				objective[tank.Id] = f.Objective {
		// 					Action: f.ActionTravel,
		// 					Target: radar.Dodge[tank.Id].SafePos,
		// 				}
		// 				continue tankloop
		// 			}
		// 		}
		// 	}
		// }

		fireRadar := radar.Fire[tank.Id]
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if fire != nil && fire.Sin < 0.5 && fire.Faith > 0 {
				objective[tank.Id] = f.Objective {
					Action: fire.Action,
				}
				continue tankloop
			}
		}	

		var ttank *f.Tank
		least := 99999
		dist := 0
		for _, etank := range state.EnemyTank {
			dist = abs(tank.Pos.X - etank.Pos.X) + abs(tank.Pos.Y - etank.Pos.Y)
			if dist < least {
				ttank = &etank
				least = dist
			}
		}

		// if ttank != nil && dist < 10 {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: f.ActionTravelWithDodge,
		// 		Target: f.Position{X: ttank.Pos.X - 8, Y: ttank.Pos.Y},
		// 	}
		// 	continue tankloop
		// }
		if ttank != nil {
			objective[tank.Id] = f.Objective {
				Action: f.ActionTravel,
				Target: ttank.Pos,
			}
			continue tankloop
		}			
		// if tank.Pos.X == targetX && tank.Pos.Y == targetY + count {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: f.ActionFireRight,
		// 	}
		// } else {
		// 	objective[tank.Id] = f.Objective {
		// 		Action: f.ActionTravel,
		// 		Target: f.Position { X: targetX, Y:targetY + count, Direction: f.DirectionLeft },
		// 	}
		// }
		// count += 1
	}
	// bottomTank := state.MyTank[0]
	// maxY := 0	
	// count = 0
	// for _, tank := range state.MyTank {		
	// 	if tank.Pos.X == targetX && tank.Pos.Y == targetY + count {
	// 		if tank.Pos.Y > maxY {
	// 			maxY = tank.Pos.Y
	// 			bottomTank = tank
	// 		}
	// 	}
	// 	count += 1
	// }
	// if bottomTank.Pos.X == targetX && bottomTank.Pos.Y == maxY {	
	// 	if _, has := objective[bottomTank.Id]; !has {
	// 		objective[bottomTank.Id] = f.Objective {
	// 			Action: f.ActionFireDown,
	// 		}
	// 	}
	// }
}

func (self *KillAll) End(state *f.GameState) {
}