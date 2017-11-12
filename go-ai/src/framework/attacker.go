/**
 * 开火攻击行动子系统
 * author: linxingchen
 */
package framework

import (
	"math";
)

type Attacker struct {
}

func NewAttacker() *Attacker {
	inst := &Attacker {
	}
	return inst
}

func calcSin (theTank Tank, tanks []Tank, fireDirection int, bulletSpeed int) float64 {
	for _, otherTank := range tanks {
		if theTank.Id == otherTank.Id {
			continue
		}

		fx := theTank.Pos.X
		fy := theTank.Pos.Y
		ox := otherTank.Pos.X
		oy := otherTank.Pos.Y

		switch fireDirection {
		case QUADRANT_U: 
			if ox == fx && fy > oy &&  fy - oy <= bulletSpeed {
				return float64(1)
			}
		case QUADRANT_L:
			if oy == fy && fx > ox && fx - ox <= bulletSpeed {
				return float64(1)
			}
		case QUADRANT_D:
			if ox == fx && oy > fy && oy - fy <= bulletSpeed {
				return float64(1)
			}
		case QUADRANT_R:
			if oy == fy && ox > fx && ox - fx <= bulletSpeed {
				return float64(1)
			}
		}
	}
	return float64(0)
}

func calcFaith (distance, bulletSpeed int) float64 {
	if distance <= bulletSpeed {
		return float64(1)
	}
	return float64(0)
}

func calcCost (tank Tank, fireDirection int, bulletSpeed int, terain Terain) int {
	cost := 1
	switch fireDirection {
	case QUADRANT_U:
		for i := tank.Pos.Y - 1; i >= 0; i-- {
			if terain.Get(tank.Pos.X, i) == 1 {
				return int(math.Ceil(float64(cost) / float64(bulletSpeed)))
			}
			cost += 1
		}
	case QUADRANT_L:
		for i := tank.Pos.X - 1; i >= 0; i-- {
			if terain.Get(i, tank.Pos.Y) == 1 {
				return int(math.Ceil(float64(cost) / float64(bulletSpeed)))
			}
			cost += 1
		}
	case QUADRANT_D:
		for i := tank.Pos.Y + 1; i < terain.Height; i++ {
			if terain.Get(tank.Pos.X, i) == 1 {
				return int(math.Ceil(float64(cost) / float64(bulletSpeed)))
			}
			cost += 1
		}
	case QUADRANT_R:
		for i := tank.Pos.X + 1; i < terain.Width; i++ {
			if terain.Get(i, tank.Pos.Y) == 1 {
				return int(math.Ceil(float64(cost) / float64(bulletSpeed)))
			}
			cost += 1
		}
	}
	return 0
}

func (self *Radar) attack(state *GameState, enemyThreats *map[string][]EnemyThreat) (map[string]*RadarFireAll) {
	radarFireAlls := make(map[string]*RadarFireAll)

	for _, tank := range state.MyTank {
		for _, enemyThreat := range (*enemyThreats)[tank.Id] {
			faith := float64(0)
			sin := float64(0)
			cost := 0
			if len(enemyThreat.Distances) == 1 {
				for fireDirection, dist := range enemyThreat.Distances {
					faith = calcFaith(dist, state.Params.BulletSpeed)
					sin = calcSin(tank, state.MyTank, fireDirection, state.Params.BulletSpeed)
					cost = calcCost(tank, fireDirection, state.Params.BulletSpeed, state.Terain)
					switch fireDirection {
					case QUADRANT_U:
						radarFireAlls[tank.Id].Up = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireUp}
					case QUADRANT_L:
						radarFireAlls[tank.Id].Left = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireLeft}
					case QUADRANT_D:
						radarFireAlls[tank.Id].Down = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireDown}
					case QUADRANT_R:
						radarFireAlls[tank.Id].Right = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireRight}
					}
				}
			}
		}
	}

	return radarFireAlls
}
/**
 * 计算火线位置
 */
// func (self *Attacker) calcFireLine(tank *Tank, state *GameState, fireline *[]Position) {

// 	switch tank.Pos.Direction {
// 		case DirectionUp:
// 			for y := tank.Pos.Y - 1; y >= 0; y-- {
// 				if 1 == state.Terain.Get(tank.Pos.X, y) {
// 					break
// 				}
// 				*fireline = append(*fireline, Position{X: tank.Pos.X, Y: y, Direction: tank.Pos.Direction})
// 			}
// 		case DirectionDown:
// 			for y := tank.Pos.Y + 1; y < state.Terain.Height; y++ {
// 				if 1 == state.Terain.Get(tank.Pos.X, y) {
// 					break
// 				}
// 				*fireline = append(*fireline, Position{X: tank.Pos.X, Y: y, Direction: tank.Pos.Direction})
// 			}
// 		case DirectionLeft:
// 			for x := tank.Pos.X - 1; x >= 0; x-- {
// 				if 1 == state.Terain.Get(x, tank.Pos.Y) {
// 					break
// 				}
// 				*fireline = append(*fireline, Position{X: x, Y: tank.Pos.Y, Direction: tank.Pos.Direction})
// 			}
// 		case DirectionRight:
// 			for x := tank.Pos.X + 1; x < state.Terain.Width; x++ {
// 				if 1 == state.Terain.Get(x, tank.Pos.Y) {
// 					break
// 				}
// 				*fireline = append(*fireline, Position{X: x, Y: tank.Pos.Y, Direction: tank.Pos.Direction})
// 			}
// 	}
// }

// func (self *Attacker) firelineEmeny(tank *Tank, state *GameState, fireline *[]Position) (fire bool, urgent int){
// 	// 初始化
// 	fire = false
// 	urgent = 0

// 	// 如果火线上，障碍物以前有地方坦克，直接射击
// 	for _, enemytank := range state.EnemyTank{
// 		if (enemytank.Pos.X == tank.Pos.X || enemytank.Pos.Y == tank.Pos.Y) {
// 			// 如果根本不在一条直线上，无需考虑
// 			if 0 == len(*fireline) {
// 				self.calcFireLine(tank, state, fireline)
// 			}
// 			// 如果敌方坦克正好在fireline上则发射 并计算还有多少步
// 			for _, firelinev := range (*fireline) {
// 				// 这个地方之后再优化
// 				if (firelinev.X == enemytank.Pos.X && firelinev.Y == enemytank.Pos.Y) {
// 					return true, int(math.Abs(float64(enemytank.Pos.X - tank.Pos.X)) + math.Abs(float64(enemytank.Pos.Y - tank.Pos.Y)))
// 				}
// 			}
// 		}
// 	}

// 	return fire, urgent
// }

// func (self *Attacker) prefire(tank *Tank, state *GameState, fireline *[]Position) (fire bool, urgent int){
// 	// 初始化
// 	fire = false
// 	urgent = 0

// 	// 如果火线两侧的位置，有朝向火线的坦克，则射击
// 	for _, enemytank := range state.EnemyTank{
// 		switch tank.Pos.Direction {
// 		case DirectionUp:
// 		case DirectionDown:
// 			if (enemytank.Pos.Direction != DirectionLeft && enemytank.Pos.Direction != DirectionRight) {
// 				continue
// 			}
// 		case DirectionLeft:
// 		case DirectionRight:
// 			if (enemytank.Pos.Direction != DirectionDown && enemytank.Pos.Direction != DirectionUp) {
// 				continue
// 			}
// 		}

// 		// 预测6格即可，即6格火线
// 		for i := 0; i < 6; i++ {
// 			if len(*fireline) > i {
// 				// 监测火线两侧敌军
// 				if (*fireline)[i].X == enemytank.Pos.X || (*fireline)[i].Y == enemytank.Pos.Y {
// 					// 敌军到火线的距离
// 					tmpFEDis := int(math.Abs(float64(enemytank.Pos.X - (*fireline)[i].X)) + math.Abs(float64(enemytank.Pos.Y - (*fireline)[i].Y)))
// 					if tmpFEDis == int(math.Ceil(float64(tmpFEDis) / 2)) {
// 						return true, i
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return fire, urgent
// }

// func (self *Attacker) Suggest(tank *Tank, state *GameState, objective *Objective) {
// 	// 定义默认值
// 	fire	:= false
// 	action 	:= ActionNone
// 	urgent 	:= 0

// 	// 火线
// 	var fireline []Position

// 	// 当前火线方向是否有敌军
// 	fire, fireUrgent := self.firelineEmeny(tank, state, &fireline)
// 	if fire == true {
// 		action = ActionFireUp
// 		urgent = int(math.Ceil(float64(fireUrgent / 2)))
// 	}

// 	// 当前火线方向附近是否有敌军
// 	if fire == false {
// 		fire, fireUrgent := self.prefire(tank, state, &fireline)
// 		if fire == true {
// 			action = ActionFireUp
// 			urgent = int(math.Ceil(float64(fireUrgent) / 2))
// 		}
// 	}

// 	_, _ = action, urgent
// 	// ret := SuggestionItem {
// 	// 	Action: action,
// 	// 	Urgent: urgent,
// 	// }
// 	// return ret
// }
