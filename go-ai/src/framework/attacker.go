/**
 * 开火攻击行动子系统
 * author: linxingchen
 */
package framework

import (
	"fmt";
	"math";
)

type Attacker struct {
}

func NewAttacker() *Attacker {
	inst := &Attacker {
	}
	return inst
}

func calcSin (theTank Tank, tanks []Tank, enemyPos Position, fireDirection int, bulletSpeed int, tankSpeed int) float64 {
	for _, otherTank := range tanks {
		if theTank.Id == otherTank.Id {
			continue
		}

		fx := theTank.Pos.X
		fy := theTank.Pos.Y
		ox := otherTank.Pos.X
		oy := otherTank.Pos.Y

		switch fireDirection {
		case DirectionUp: 
			if ox == fx && fy > oy {
				if enemyPos.X != -1 && enemyPos.X == ox && enemyPos.Y > oy && enemyPos.Y < fy {
					return float64(0)
				} 
				return float64(1)
			} else if fy > oy {
				if ox - fx == tankSpeed && otherTank.Pos.Direction == DirectionLeft {
					return float64(0.5)
				}
				if fx - ox == tankSpeed && otherTank.Pos.Direction == DirectionRight {
					return float64(0.5)
				}
			}
		case DirectionLeft:
			if oy == fy && fx > ox {
				if enemyPos.X != -1 && enemyPos.Y == oy && enemyPos.X > ox && enemyPos.X < fx {
					return float64(0)
				}
				return float64(1)
			} else if fx > ox {
				if oy - fy == tankSpeed && otherTank.Pos.Direction == DirectionUp {
					return float64(0.5)
				}
				if fy - oy == tankSpeed && otherTank.Pos.Direction == DirectionDown {
					return float64(0.5)
				}
			}
		case DirectionDown:
			if ox == fx && oy > fy {
				if enemyPos.X != -1 &&  enemyPos.X == ox && enemyPos.Y > fy && enemyPos.Y < oy {
					return float64(0)
				}
				return float64(1)
			} else if oy > fy {
				if ox - fx == tankSpeed && otherTank.Pos.Direction == DirectionLeft {
					return float64(0.5)
				}
				if fx - ox == tankSpeed && otherTank.Pos.Direction == DirectionRight {
					return float64(0.5)
				}
			}
		case DirectionRight:
			if oy == fy && ox > fx {
				if enemyPos.X != -1 && enemyPos.Y == oy && enemyPos.X > fx && enemyPos.X < ox {
					return float64(0)
				}
				return float64(1)
			} else if ox > fx {
				if oy - fy == tankSpeed && otherTank.Pos.Direction == DirectionUp {
					return float64(0.5)
				}
				if fy - oy == tankSpeed && otherTank.Pos.Direction == DirectionDown {
					return float64(0.5)
				}
			}
		}
	}
	return float64(0)
}

// faith计算方法：

// 首先根据与敌方垂直距离给出faith基准值:
// 1(1个子弹速度以内), 0.8（2个）, 0.6（3个）, 0.4（4个）, 0.3（10个以内），其他情况是0

// 然后判断是否在火线上，若在火线上：
// 如果我方与敌方坦克垂直距离 == 1，直接返回最终faith = 1
// 如果敌方与开火方向是相同的或者相反的，返回faith基准值
// 如果敌方与开火方向垂直，返回faith基准值 * 0.9

// 若敌方坦克不在火线上：
// 若敌方坦克下回合朝向火线，且正好有可能走到火线上，返回(faith基准值 / 2) - 0.15
// 若敌方坦克与火线距离 > 1个坦克速度，返回faith = 0
// 若敌方坦克与火线距离在1个坦克速度以内，但不朝向火线，返回(faith基准值 / 2) - 0.3

func calcFaith (verticalDistance, bulletSpeed int, tankSpeed int, fireLine bool, fireDirection int, enemyPos Position, tankPos Position, terain *Terain) float64 {
	faith := float64(0)

	if verticalDistance <= bulletSpeed + 1 {
		faith = float64(1)
	} else if verticalDistance <= bulletSpeed * 2 + 1 {
		faith = 0.8
	} else if verticalDistance <= bulletSpeed * 3 + 1 {
		faith = 0.6
	} else if verticalDistance <= bulletSpeed * 4 + 1 {
		faith = 0.4
	} else if verticalDistance <= bulletSpeed * 10 + 1 {
		faith = 0.3
	} else {
		return faith
	}

	// fmt.Println("1 ------- faith enemy direction", enemyPos.Direction, "count", count)

	if fireLine {
		// 与敌方坦克相邻，不管敌方朝向都必中
		if verticalDistance <= 1 {
			return 1.
		}

		// 敌方朝向和开火方向相同或相反，且在火线上		
		if enemyPos.Direction == fireDirection || enemyPos.Direction == fireDirection + 2 || enemyPos.Direction == fireDirection - 2  {
			return faith
		}

		// 敌方朝向和开火方向垂直，且在火线上
		return faith * 0.9
	} else {
		// 敌方不在火线，开火方向是上或下
		faith = faith / float64(2)
		if fireDirection == DirectionUp || fireDirection == DirectionDown {

			// 坦克下回合走不到火线上
			if int(math.Abs(float64(tankPos.X - enemyPos.X))) > tankSpeed {
				return float64(0)
			}

			// 敌方坦克在火线左侧，朝向火线
			if tankPos.X > enemyPos.X && enemyPos.Direction == DirectionRight && terain.Get(tankPos.X, enemyPos.Y) != 1 {
				return faith - 0.15
			}
			// 敌方坦克在火线右侧，朝向火线
			if tankPos.X < enemyPos.X && enemyPos.Direction == DirectionLeft && terain.Get(tankPos.X, enemyPos.Y) != 1{
				return faith - 0.15
			}
			// 敌方坦克朝向与火线相反或不朝向火线
			return faith - 0.3
		}

		// 开火方向是左或右
		if fireDirection == DirectionLeft || fireDirection == DirectionRight {

			// 坦克下回合走不到火线上
			if int(math.Abs(float64(tankPos.Y - enemyPos.Y))) > tankSpeed {
				return float64(0)
			}

			// 敌方坦克在火线上面，朝向火线
			if tankPos.Y > enemyPos.Y && enemyPos.Direction == DirectionDown && terain.Get(enemyPos.X, tankPos.Y) != 1 {
				return faith - 0.15
			}
			// 敌方坦克在火线下面，朝向火线
			if tankPos.Y < enemyPos.Y && enemyPos.Direction == DirectionUp && terain.Get(enemyPos.X, tankPos.Y) != 1 {
				return faith - 0.15
			}
			// 敌方坦克朝向与火线相反或不朝向火线
			return faith - 0.3
		}
	}

	return faith
}

func calcCost (tank Tank, fireDirection int, bulletSpeed int, terain *Terain) int {
	cost := 0
	switch fireDirection {
	case DirectionUp:
		// fmt.Println("2 ----- UP", "tank", tank.Id, "count", count)
		for i := tank.Pos.Y - 1; i >= 0; i-- {
			if terain.Get(tank.Pos.X, i) == 1 {
				return cost
			}
			cost += 1
		}
	case DirectionLeft:
		// fmt.Println("2 ----- LEFT", "tank", tank.Id, "count", count)		
		for i := tank.Pos.X - 1; i >= 0; i-- {
			if terain.Get(i, tank.Pos.Y) == 1 {
				return cost
			}
			cost += 1
		}
	case DirectionDown:
		// fmt.Println("2 ----- DOWN", "tank", tank.Id, "count", count)		
		for i := tank.Pos.Y + 1; i < terain.Height; i++ {
			if terain.Get(tank.Pos.X, i) == 1 {
				return cost
			}
			cost += 1
		}
	case DirectionRight:
		// fmt.Println("2 ----- RIGHT", "tank", tank.Id, "count", count)		
		for i := tank.Pos.X + 1; i < terain.Width; i++ {
			if terain.Get(i, tank.Pos.Y) == 1 {
				return cost
			}
			cost += 1
		}
	}
	return cost
}

func directionConvert(fireDirection int, tank Tank) int {
	realDirection := 0
	switch fireDirection {
	case QUADRANT_U:
		realDirection = DirectionUp
	case QUADRANT_L:
		realDirection = DirectionLeft
	case QUADRANT_D:
		realDirection = DirectionDown
	case QUADRANT_R:
		realDirection = DirectionRight
	}
	return DirectionUp + ((realDirection - DirectionUp) + (tank.Pos.Direction - DirectionUp) + 4) % 4
}

var count = 0

func (self *Radar) Attack(state *GameState, enemyThreats *map[string][]EnemyThreat) (map[string]*RadarFireAll) {
	radarFireAlls := make(map[string]*RadarFireAll)
	
	count += 1

	for _, tank := range state.MyTank {
		if tank.Bullet != "" {
			continue
		}

		radarFireAlls[tank.Id] = &RadarFireAll {}
		
		var noEnemy Position
		noEnemy.X = -1
		noEnemy.Y = -1

		sin := calcSin(tank, state.MyTank, noEnemy, DirectionUp, state.Params.BulletSpeed, state.Params.TankSpeed)
		cost := calcCost(tank, DirectionUp, state.Params.BulletSpeed, state.Terain)
		
		radarFireAlls[tank.Id].Up = &RadarFire {Faith: 0, Cost: cost, Sin: sin, Action: ActionFireUp}

		sin = calcSin(tank, state.MyTank, noEnemy, DirectionLeft, state.Params.BulletSpeed, state.Params.TankSpeed)
		cost = calcCost(tank, DirectionLeft, state.Params.BulletSpeed, state.Terain)
		radarFireAlls[tank.Id].Left = &RadarFire {Faith: 0, Cost: cost, Sin: sin, Action: ActionFireLeft}

		sin = calcSin(tank, state.MyTank, noEnemy, DirectionDown, state.Params.BulletSpeed, state.Params.TankSpeed)
		cost = calcCost(tank, DirectionDown, state.Params.BulletSpeed, state.Terain)		
		radarFireAlls[tank.Id].Down = &RadarFire {Faith: 0, Cost: cost, Sin: sin, Action: ActionFireDown}

		sin = calcSin(tank, state.MyTank, noEnemy, DirectionRight, state.Params.BulletSpeed, state.Params.TankSpeed)
		cost = calcCost(tank, DirectionRight, state.Params.BulletSpeed, state.Terain)		
		radarFireAlls[tank.Id].Right = &RadarFire {Faith: 0, Cost: cost, Sin: sin, Action: ActionFireRight}

		// radarFireAlls[tank.Id].Up = &RadarFire {Faith: 0, Cost: 0, Sin: 0, Action: ActionFireUp}
		// radarFireAlls[tank.Id].Left = &RadarFire {Faith: 0, Cost: 0, Sin: 0, Action: ActionFireLeft}
		// radarFireAlls[tank.Id].Down = &RadarFire {Faith: 0, Cost: 0, Sin: 0, Action: ActionFireDown}
		// radarFireAlls[tank.Id].Right = &RadarFire {Faith: 0, Cost: 0, Sin: 0, Action: ActionFireRight}		

		for _, enemyThreat := range (*enemyThreats)[tank.Id] {
			faith := float64(0)
			sin := float64(0)
			cost := 0

			// 敌方不在火线，但在火线两侧
			if len(enemyThreat.Distances) == 2 {
				verticalDist := 0
				fireDirection := 0
				needToCalc := false
				horizontalDist := 0

				var dirA, dirB int
				var distA, distB int
				first := true
				for direction, dist := range enemyThreat.Distances {
					if first {
						first = false
						dirA, distA = direction, dist
					} else {
						dirB, distB = direction, dist
					}
				}
				if distA < distB {
					dirA, dirB = dirB, dirA
					distA, distB = distB, distA
				}

				verticalDist = distA
				horizontalDist = distB
				fireDirection = dirA


				if horizontalDist >= 1 && horizontalDist <= 2 * state.Params.TankSpeed {
					needToCalc = true
				}

				realDirection := directionConvert(fireDirection, tank)
				// fmt.Println("0 ------- DIRECTION", "real", realDirection, "fire", fireDirection)
				
				if needToCalc {
					faith = calcFaith(verticalDist, state.Params.BulletSpeed, state.Params.TankSpeed, false, realDirection, enemyThreat.Enemy, tank.Pos, state.Terain)										
				} else {
					faith = 0
				}
				sin = calcSin(tank, state.MyTank, enemyThreat.Enemy, realDirection, state.Params.BulletSpeed, state.Params.TankSpeed)
				cost = calcCost(tank, realDirection, state.Params.BulletSpeed, state.Terain)
				
				if cost < verticalDist {
					faith = float64(0)
					sin = float64(0)
				}

				fmt.Println("FIRE INFO -------- NO FIRELINE", "count", count, "faith", faith, "sin", sin, "cost", cost, "dist", verticalDist, "tank", tank.Id, "enemyX", enemyThreat.Enemy.X, "enemyY", enemyThreat.Enemy.Y)
				
				cost = int(math.Ceil(float64(cost) / float64(state.Params.BulletSpeed)))


				switch realDirection {
				case DirectionUp:
					if faith < radarFireAlls[tank.Id].Up.Faith {
						continue
					}
					radarFireAlls[tank.Id].Up = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireUp}
				case DirectionLeft:
					if faith < radarFireAlls[tank.Id].Left.Faith {
						continue
					}
					radarFireAlls[tank.Id].Left = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireLeft}
				case DirectionDown:
					if faith < radarFireAlls[tank.Id].Down.Faith {
						continue
					}
					radarFireAlls[tank.Id].Down = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireDown}
				case DirectionRight:
					if faith < radarFireAlls[tank.Id].Right.Faith {
						continue
					}
					radarFireAlls[tank.Id].Right = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireRight}
				}
			}

			// 敌方在火线
			if len(enemyThreat.Distances) == 1 {
				for fireDirection, dist := range enemyThreat.Distances {
					realDirection := directionConvert(fireDirection, tank)
					// fmt.Println("0 ------- DIRECTION", "real", realDirection, "fire", fireDirection)

					faith = calcFaith(dist, state.Params.BulletSpeed, state.Params.TankSpeed, true, realDirection, enemyThreat.Enemy, tank.Pos, state.Terain)
					sin = calcSin(tank, state.MyTank, enemyThreat.Enemy, realDirection, state.Params.BulletSpeed, state.Params.TankSpeed)
					cost = calcCost(tank, realDirection, state.Params.BulletSpeed, state.Terain)

					if cost < dist {
						faith = float64(0)
						sin = float64(0)
					}
					fmt.Println("FIRE INFO -------- FIRELINE", "count", count, "faith", faith, "sin", sin, "cost", cost, "dist", dist, "tank", tank.Id, "enemyX", enemyThreat.Enemy.X, "enemyY", enemyThreat.Enemy.Y)
					
					cost = int(math.Ceil(float64(cost) / float64(state.Params.BulletSpeed)))					
					switch realDirection {
					case DirectionUp:
						radarFireAlls[tank.Id].Up = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireUp}
					case DirectionLeft:
						radarFireAlls[tank.Id].Left = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireLeft}
					case DirectionDown:
						radarFireAlls[tank.Id].Down = &RadarFire {Faith: faith, Cost: cost, Sin: sin, Action: ActionFireDown}
					case DirectionRight:
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
