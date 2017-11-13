/**
 * 开火攻击行动子系统
 * author: linxingchen
 */
package framework

import (
	// "fmt";
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
		case DirectionUp:
			if ox == fx && fy > oy &&  fy - oy <= bulletSpeed {
				return float64(1)
			}
		case DirectionLeft:
			if oy == fy && fx > ox && fx - ox <= bulletSpeed {
				return float64(1)
			}
		case DirectionDown:
			if ox == fx && oy > fy && oy - fy <= bulletSpeed {
				return float64(1)
			}
		case DirectionRight:
			if oy == fy && ox > fx && ox - fx <= bulletSpeed {
				return float64(1)
			}
		}
	}
	return float64(0)
}

func calcFaith (verticalDistance, bulletSpeed int, tankSpeed int, fireLine bool, fireDirection int, enemyPos Position, tankPos Position) float64 {
	faith := float64(0)

	if verticalDistance <= bulletSpeed {
		faith = float64(1)
	} else if verticalDistance <= bulletSpeed * 2 {
		faith = 0.5
	} else {
		return faith
	}

	if fireLine {

		// 一个子弹距离内且在火线上，不管敌方朝向都必中
		if verticalDistance <= bulletSpeed {
			return float64(1)
		}

		// 敌方朝向和开火方向相同或相反，且在火线上		
		if enemyPos.Direction == fireDirection || enemyPos.Direction == fireDirection + 2 || enemyPos.Direction == fireDirection - 2  {
			return faith
		}

		// 敌方朝向和开火方向垂直，且在火线上
		return faith / 2
	} else {
		// 敌方不在火线，开火方向是上或下
		if fireDirection == DirectionUp || fireDirection == DirectionDown {

			// 坦克下回合走不到火线上
			if int(math.Abs(float64(tankPos.X - enemyPos.X))) != tankSpeed {
				return float64(0)
			}

			// 敌方坦克在火线左侧，朝向火线
			if tankPos.X > enemyPos.X && enemyPos.Direction == DirectionRight {
				return faith - 0.15
			}
			// 敌方坦克在火线右侧，朝向火线
			if tankPos.X < enemyPos.X && enemyPos.Direction == DirectionLeft {
				return faith - 0.15
			}
			// 敌方坦克朝向与火线相反或不朝向火线
			return float64(0)
		}

		// 开火方向是左或右
		if fireDirection == DirectionLeft || fireDirection == DirectionRight {

			// 坦克下回合走不到火线上
			if int(math.Abs(float64(tankPos.Y - enemyPos.Y))) != tankSpeed {
				return float64(0)
			}

			// 敌方坦克在火线上面，朝向火线
			if tankPos.Y > enemyPos.Y && enemyPos.Direction == DirectionDown {
				return faith - 0.15
			}
			// 敌方坦克在火线下面，朝向火线
			if tankPos.Y < enemyPos.Y && enemyPos.Direction == DirectionUp {
				return faith - 0.15
			}
			// 敌方坦克朝向与火线相反或不朝向火线
			return float64(0)
		}
	}

	return faith
}

func calcCost (tank Tank, fireDirection int, bulletSpeed int, terain *Terain) int {
	cost := 0
	switch fireDirection {
	case DirectionUp:
		for i := tank.Pos.Y - 1; i >= 0; i-- {
			if terain.Get(tank.Pos.X, i) == 1 {
				return cost
			}
			cost += 1
		}
	case DirectionLeft:
		for i := tank.Pos.X - 1; i >= 0; i-- {
			if terain.Get(i, tank.Pos.Y) == 1 {
				return cost
			}
			cost += 1
		}
	case DirectionDown:
		for i := tank.Pos.Y + 1; i < terain.Height; i++ {
			if terain.Get(tank.Pos.X, i) == 1 {
				return cost
			}
			cost += 1
		}
	case DirectionRight:
		for i := tank.Pos.X + 1; i < terain.Width; i++ {
			if terain.Get(i, tank.Pos.Y) == 1 {
				return cost
			}
			cost += 1
		}
	}
	return 0
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

func (self *Radar) Attack(state *GameState, enemyThreats *map[string][]EnemyThreat) (map[string]*RadarFireAll) {
	radarFireAlls := make(map[string]*RadarFireAll)

	for _, tank := range state.MyTank {
		if tank.Bullet != "" {
			continue
		}
		radarFireAlls[tank.Id] = &RadarFireAll {}
		for _, enemyThreat := range (*enemyThreats)[tank.Id] {
			faith := float64(0)
			sin := float64(0)
			cost := 0

			// 敌方不在火线，但在火线两侧
			if len(enemyThreat.Distances) == 2 {
				verticalDist := 0
				for fireDirection, dist := range enemyThreat.Distances {
					if dist == 1 {
						realDirection := directionConvert(fireDirection, tank)

						faith = calcFaith(verticalDist, state.Params.BulletSpeed, state.Params.TankSpeed, false, realDirection, enemyThreat.Enemy, tank.Pos)
						sin = calcSin(tank, state.MyTank, realDirection, state.Params.BulletSpeed)
						cost = calcCost(tank, realDirection, state.Params.BulletSpeed, state.Terain)

						if cost < dist {
							faith = float64(0)
							sin = float64(0)
						}
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
					} else {
						verticalDist = dist
					}
				}
			}

			// 敌方在火线
			if len(enemyThreat.Distances) == 1 {
				for fireDirection, dist := range enemyThreat.Distances {
					realDirection := directionConvert(fireDirection, tank)

					faith = calcFaith(dist, state.Params.BulletSpeed, state.Params.TankSpeed, true, realDirection, enemyThreat.Enemy, tank.Pos)
					sin = calcSin(tank, state.MyTank, realDirection, state.Params.BulletSpeed)
					cost = calcCost(tank, realDirection, state.Params.BulletSpeed, state.Terain)

					if cost < dist {
						faith = float64(0)
						sin = float64(0)
					}

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
