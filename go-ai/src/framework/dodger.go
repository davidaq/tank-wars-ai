/**
 * 子弹躲避行动子系统
 * author: linxingchen
 */
package framework

import (
	"math"
)

type Dodger struct {
}

type EnemyBullet struct {
	BulletPosition Position
	Distance 	int
	Quadrant	int
}

type EnemyThreat struct {
	Enemy 		Position
	Distance	int
	Quadrant	int
}

func NewDodger() *Dodger {
	inst := &Dodger {
	}
	return inst
}

func (self *Dodger) avoidBullet(tank *Tank, state *GameState) (bulletApproach bool, enemyBullet []EnemyBullet) {
	radius := 8	//子弹雷达半径

	// 计算子弹在自己的什么方位 子弹方向必须面向火线 从自己方位向四个方位探索墙
	var tmpBullet []EnemyBullet
	for _, enemyBullet := range state.EnemyBullet {
		// 检查是否需要关注
		if math.Pow(float64(enemyBullet.Pos.X - tank.Pos.X), 2) + math.Pow(float64(enemyBullet.Pos.Y - tank.Pos.Y), 2) <= float64(radius * radius) {
			// 判断方向 并 看方向上有没有墙
			// 象限
			tmpQuadrant := 0
			tmpDistance := 0	// 注意要除2
			if enemyBullet.Pos.X < tank.Pos.X {
				if enemyBullet.Pos.Y < tank.Pos.Y && (enemyBullet.Pos.Direction == DirectionRight || enemyBullet.Pos.Direction == DirectionDown) {
					tmpQuadrant = 2
					if enemyBullet.Pos.Direction == DirectionRight {
						for w := enemyBullet.Pos.X; w <= tank.Pos.X; w++ {
							if 1 == state.Terain.Get(w, tank.Pos.Y) {
								tmpQuadrant = 0
								break
							}
							tmpDistance++
						}
					} else {
						for w := enemyBullet.Pos.Y; w <= tank.Pos.Y; w++ {
							if 1 == state.Terain.Get(tank.Pos.X, w) {
								tmpQuadrant = 0
								break
							}
							tmpDistance++
						}
					}
				}
				if enemyBullet.Pos.Y > tank.Pos.Y && (enemyBullet.Pos.Direction == DirectionRight || enemyBullet.Pos.Direction == DirectionUp) {
					tmpQuadrant = 3
					if enemyBullet.Pos.Direction == DirectionRight {
						for w := enemyBullet.Pos.X; w <= tank.Pos.X; w++ {
							if 1 == state.Terain.Get(w, tank.Pos.Y) {
								tmpQuadrant = 0
								break
							}
							tmpDistance++
						}
					} else {
						for w := enemyBullet.Pos.Y; w >= tank.Pos.Y; w-- {
							if 1 == state.Terain.Get(tank.Pos.X, w) {
								tmpQuadrant = 0
								break
							}
							tmpDistance++
						}
					}
				}
			}
			if enemyBullet.Pos.X > tank.Pos.X {
				if enemyBullet.Pos.Y < tank.Pos.Y && (enemyBullet.Pos.Direction == DirectionLeft || enemyBullet.Pos.Direction == DirectionDown) {
					tmpQuadrant = 1
					if enemyBullet.Pos.Direction == DirectionLeft {
						for w := enemyBullet.Pos.X; w >= tank.Pos.X; w-- {
							if 1 == state.Terain.Get(w, tank.Pos.Y) {
								tmpQuadrant = 0
								break
							}
							tmpDistance++
						}
					} else {
						for w := enemyBullet.Pos.Y; w <= tank.Pos.Y; w++ {
							if 1 == state.Terain.Get(tank.Pos.X, w) {
								tmpQuadrant = 0
								break
							}
							tmpDistance++
						}
					}
				}
				if enemyBullet.Pos.Y > tank.Pos.Y && (enemyBullet.Pos.Direction == DirectionLeft || enemyBullet.Pos.Direction == DirectionUp) {
					tmpQuadrant = 4
					if enemyBullet.Pos.Direction == DirectionLeft {
						for w := enemyBullet.Pos.X; w >= tank.Pos.X; w-- {
							if 1 == state.Terain.Get(w, tank.Pos.Y) {
								tmpQuadrant = 0
								break
							}
							tmpDistance++
						}
					} else {
						for w := enemyBullet.Pos.Y; w <= tank.Pos.Y; w++ {
							if 1 == state.Terain.Get(tank.Pos.X, w) {
								tmpQuadrant = 0
								break
							}
							tmpDistance++
						}
					}
				}
			}
			if enemyBullet.Pos.X == tank.Pos.X {
				//在Y火线上
				if enemyBullet.Pos.Y < tank.Pos.Y && enemyBullet.Pos.Direction == DirectionDown {
					tmpQuadrant = -1
					for w := enemyBullet.Pos.Y; w <= tank.Pos.Y; w++ {
						if 1 == state.Terain.Get(tank.Pos.X, w) {
							tmpQuadrant = 0
							break
						}
						tmpDistance++
					}
				}
				if enemyBullet.Pos.Y > tank.Pos.Y && enemyBullet.Pos.Direction == DirectionUp {
					tmpQuadrant = -3
					for w := enemyBullet.Pos.Y; w >= tank.Pos.Y; w-- {
						if 1 == state.Terain.Get(tank.Pos.X, w) {
							tmpQuadrant = 0
							break
						}
						tmpDistance++
					}
				}
			}
			if enemyBullet.Pos.Y == tank.Pos.Y {
				//在X火线上
				if enemyBullet.Pos.X < tank.Pos.X && enemyBullet.Pos.Direction == DirectionRight {
					tmpQuadrant = -2
					for w := enemyBullet.Pos.X; w <= tank.Pos.X; w++ {
						if 1 == state.Terain.Get(w, tank.Pos.Y) {
							tmpQuadrant = 0
							break
						}
						tmpDistance++
					}
				}
				if enemyBullet.Pos.X > tank.Pos.X && enemyBullet.Pos.Direction == DirectionLeft {
					tmpQuadrant = -4
					for w := enemyBullet.Pos.X; w >= tank.Pos.X; w-- {
						if 1 == state.Terain.Get(w, tank.Pos.Y) {
							tmpQuadrant = 0
							break
						}
						tmpDistance++
					}
				}
			}

			if tmpQuadrant == 0 || tmpDistance == 0 {
				continue
			}

			tmpBullet = append(tmpBullet, EnemyBullet{
				BulletPosition: enemyBullet.Pos,
				Distance: int(tmpDistance / 2),
				Quadrant: tmpQuadrant,
			})
		}
	}

	if 0 == len(tmpBullet) {
		return false, tmpBullet
	}
	return true, tmpBullet
}


func (self *Dodger) oldAvoidBullet(tank *Tank, state *GameState) (bulletApproach bool, enemyBullet []EnemyBullet) {
	// 注意射过来的炮弹可能多个
	// 没有子弹射过来的情况
	enemyBulletList := []EnemyBullet{}
	for _, bullet := range state.EnemyBullet {
		if tank.Pos.X == bullet.Pos.X && int(math.Abs(float64(bullet.Pos.Y - tank.Pos.Y))) <= 8 || tank.Pos.Y == bullet.Pos.Y && int(math.Abs(float64(bullet.Pos.X - tank.Pos.X))) <= 8 {
			// 子弹方向判断
			bulletDirection := DirectionNone
			if tank.Pos.X == bullet.Pos.X {
				if tank.Pos.Y > bullet.Pos.Y && bullet.Pos.Direction == DirectionDown {
					// 墙判断
					for wallY := bullet.Pos.Y; wallY < tank.Pos.Y; wallY++ {
						if 1 == state.Terain.Get(tank.Pos.X, wallY) {
							continue
						}
					}
					bulletDirection = DirectionDown
				} else if tank.Pos.Y < bullet.Pos.Y && bullet.Pos.Direction == DirectionUp {
					for wallY := tank.Pos.Y; wallY < bullet.Pos.Y; wallY++ {
						if 1 == state.Terain.Get(tank.Pos.X, wallY) {
							continue
						}
					}
					bulletDirection = DirectionUp
				}
			}
			if tank.Pos.Y == bullet.Pos.Y {
				if tank.Pos.X > bullet.Pos.X && bullet.Pos.Direction == DirectionRight {
					for wallX := bullet.Pos.X; wallX < tank.Pos.X; wallX++ {
						if 1 == state.Terain.Get(wallX, tank.Pos.Y) {
							continue
						}
					}
					bulletDirection = DirectionRight
				} else if tank.Pos.Y < bullet.Pos.Y && bullet.Pos.Direction == DirectionLeft {
					for wallX := tank.Pos.Y; wallX < bullet.Pos.Y; wallX++ {
						if 1 == state.Terain.Get(wallX, tank.Pos.Y) {
							continue
						}
					}
					bulletDirection = DirectionLeft
				}
			}

			if bulletDirection == DirectionNone {
				continue
			}
			tmpBulletPosition := Position{
				X: bullet.Pos.X,
				Y: bullet.Pos.Y,
				Direction: bulletDirection,
			}
			enemyBulletList = append(enemyBulletList, EnemyBullet{
				Distance: int(math.Abs(float64(bullet.Pos.Y - tank.Pos.Y)) + math.Abs(float64(bullet.Pos.X - tank.Pos.X))),
				BulletPosition: tmpBulletPosition,
			})
		}
	}

	if 0 == len(enemyBulletList) {
		return false, enemyBulletList
	}

	return true, enemyBulletList
}

func (self *Dodger) threat(tank *Tank, state *GameState) (threat bool, enemyThreat []EnemyThreat) {
	radius := 5	//雷达半径

	// 计算敌军在自己的什么方位 墙无视，防止LYB苟在墙里
	var tmpEnemyThreat []EnemyThreat
	for _, enemyTank := range state.EnemyTank {
		// 检查是否需要关注 圆形雷达
		if math.Pow(float64(enemyTank.Pos.X - tank.Pos.X), 2) + math.Pow(float64(enemyTank.Pos.Y - tank.Pos.Y), 2) <= float64(radius * radius) {
			tmpEnemyThreat = append(tmpEnemyThreat, EnemyThreat{
				Enemy: enemyTank.Pos,
				//Distance: int(math.Sqrt(math.Pow(float64(enemyTank.Pos.X - tank.Pos.X), 2) + math.Pow(float64(enemyTank.Pos.Y - tank.Pos.Y), 2))),
			})
		}
	}

	if 0 == len(tmpEnemyThreat) {
		return false, tmpEnemyThreat
	}

	// 计算方位和象限
	for k, enemy := range tmpEnemyThreat {
		if enemy.Enemy.Y < tank.Pos.Y {
			if enemy.Enemy.X > tank.Pos.X {
				// 计算象限
				tmpEnemyThreat[k].Quadrant = 1
				// 计算distance 非面向火线则distance+1
				if enemy.Enemy.Direction == DirectionLeft {
					tmpEnemyThreat[k].Distance = enemy.Enemy.X - tank.Pos.X
				} else if enemy.Enemy.Direction == DirectionDown {
					tmpEnemyThreat[k].Distance = tank.Pos.Y - enemy.Enemy.Y
				} else {
					tmpEnemyThreat[k].Distance = int(math.Max(float64(enemy.Enemy.X - tank.Pos.X), float64(tank.Pos.Y - enemy.Enemy.Y))) + 1
				}
			}
			if enemy.Enemy.X < tank.Pos.X {
				tmpEnemyThreat[k].Quadrant = 2
				// 计算distance
				if enemy.Enemy.Direction == DirectionRight {
					tmpEnemyThreat[k].Distance = tank.Pos.X - enemy.Enemy.X
				} else if enemy.Enemy.Direction == DirectionDown {
					tmpEnemyThreat[k].Distance = tank.Pos.Y - enemy.Enemy.Y
				} else {
					tmpEnemyThreat[k].Distance = int(math.Max(float64(tank.Pos.X - enemy.Enemy.X), float64(tank.Pos.Y - enemy.Enemy.Y))) + 1
				}
			}
		}

		if enemy.Enemy.Y > tank.Pos.Y {
			if enemy.Enemy.X < tank.Pos.X {
				tmpEnemyThreat[k].Quadrant = 3
				// 计算distance
				if enemy.Enemy.Direction == DirectionUp {
					tmpEnemyThreat[k].Distance = enemy.Enemy.Y - tank.Pos.Y
				} else if enemy.Enemy.Direction == DirectionRight {
					tmpEnemyThreat[k].Distance = tank.Pos.X - enemy.Enemy.X
				} else {
					tmpEnemyThreat[k].Distance = int(math.Max(float64(enemy.Enemy.Y - tank.Pos.Y), float64(tank.Pos.X - enemy.Enemy.X))) + 1
				}
			}
			if enemy.Enemy.X > tank.Pos.X {
				tmpEnemyThreat[k].Quadrant = 4
				// 计算distance
				if enemy.Enemy.Direction == DirectionUp {
					tmpEnemyThreat[k].Distance = enemy.Enemy.Y - tank.Pos.Y
				} else if enemy.Enemy.Direction == DirectionLeft {
					tmpEnemyThreat[k].Distance = enemy.Enemy.X - tank.Pos.X
				} else {
					tmpEnemyThreat[k].Distance = int(math.Max(float64(enemy.Enemy.Y - tank.Pos.Y), float64(enemy.Enemy.X - tank.Pos.X))) + 1
				}
			}
		}

		if enemy.Enemy.X == tank.Pos.X {
			if enemy.Enemy.Y < tank.Pos.Y {
				tmpEnemyThreat[k].Quadrant = -1
				tmpEnemyThreat[k].Distance = tank.Pos.Y - enemy.Enemy.Y
				if enemy.Enemy.Direction == DirectionLeft || enemy.Enemy.Direction == DirectionRight {
					tmpEnemyThreat[k].Distance++
				}
				if enemy.Enemy.Direction == DirectionUp {
					// 背对着，加2
					tmpEnemyThreat[k].Distance += 2
				}
			}
			if enemy.Enemy.Y > tank.Pos.Y {
				tmpEnemyThreat[k].Quadrant = -3
				tmpEnemyThreat[k].Distance = enemy.Enemy.Y - tank.Pos.Y
				if enemy.Enemy.Direction == DirectionLeft || enemy.Enemy.Direction == DirectionRight {
					tmpEnemyThreat[k].Distance++
				}
				if enemy.Enemy.Direction == DirectionDown {
					tmpEnemyThreat[k].Distance += 2
				}
			}
		}

		if enemy.Enemy.Y == tank.Pos.Y {
			if enemy.Enemy.X < tank.Pos.X {
				tmpEnemyThreat[k].Quadrant = -2
				tmpEnemyThreat[k].Distance = tank.Pos.X - enemy.Enemy.X
				if enemy.Enemy.Direction == DirectionUp || enemy.Enemy.Direction == DirectionDown {
					tmpEnemyThreat[k].Distance++
				}
				if enemy.Enemy.Direction == DirectionLeft {
					tmpEnemyThreat[k].Distance += 2
				}
			}
			if enemy.Enemy.X > tank.Pos.X {
				tmpEnemyThreat[k].Quadrant = -4
				tmpEnemyThreat[k].Distance = enemy.Enemy.X - tank.Pos.X
				if enemy.Enemy.Direction == DirectionUp || enemy.Enemy.Direction == DirectionDown {
					tmpEnemyThreat[k].Distance++
				}
				if enemy.Enemy.Direction == DirectionRight {
					tmpEnemyThreat[k].Distance += 2
				}
			}
		}
	}

	return true, tmpEnemyThreat
}


func (self *Dodger) determine(tank *Tank, state *GameState, bulletApproach bool, enemyBullets *[]EnemyBullet, enemyApproach bool, enemyThreat *[]EnemyThreat) (action int, urgent int) {
	// 注意方向撞墙情况，如果没有则确认方向
	// 象限旋转（假定为上，然后旋转）
	quadrant := make(map[int]int)
	switch tank.Pos.Direction {
	case DirectionUp:
		quadrant[0] = 0
		quadrant[1] = 1
		quadrant[2] = 2
		quadrant[3] = 3
		quadrant[4] = 4
		quadrant[-1] = -1
		quadrant[-2] = -2
		quadrant[-3] = -3
		quadrant[-4] = -4

	case DirectionLeft:
		quadrant[0] = 0
		quadrant[2] = 1
		quadrant[3] = 2
		quadrant[4] = 3
		quadrant[1] = 4
		quadrant[-4] = -1
		quadrant[-1] = -2
		quadrant[-2] = -3
		quadrant[-3] = -4

	case DirectionDown:
		quadrant[0] = 0
		quadrant[3] = 1
		quadrant[4] = 2
		quadrant[1] = 3
		quadrant[2] = 4
		quadrant[-3] = -1
		quadrant[-4] = -2
		quadrant[-1] = -3
		quadrant[-2] = -4

	case DirectionRight:
		quadrant[0] = 0
		quadrant[4] = 1
		quadrant[1] = 2
		quadrant[2] = 3
		quadrant[3] = 4
		quadrant[-4] = -1
		quadrant[-1] = -2
		quadrant[-2] = -3
		quadrant[-3] = -4
	}

	// 最后收敛到几个方向上，直接在方向上标出最小的紧急度，最后走方向中紧急度排名第一，但是方向中不紧急的
	BulletMoveUrgent := [6]int{}
	BulletMoveUrgent[ActionMove] = math.MaxInt32
	BulletMoveUrgent[ActionBack] = math.MaxInt32
	BulletMoveUrgent[ActionLeft] = math.MaxInt32
	BulletMoveUrgent[ActionRight] = math.MaxInt32

	BulletMoveUrgent[ActionStay] = math.MaxInt32


	// 炮弹为第一优先级
	if bulletApproach == true {
		// 炮弹的一条线都走不了
		for _, bullet := range *enemyBullets {
			switch quadrant[bullet.Quadrant] {
			case 1:
				// 影响直行和右转
				if BulletMoveUrgent[ActionMove] > bullet.Distance {
					BulletMoveUrgent[ActionMove] = bullet.Distance
				}
				if BulletMoveUrgent[ActionRight] > bullet.Distance {
					BulletMoveUrgent[ActionRight] = bullet.Distance
				}
			case 2:
				// 影响直行和左转
				if BulletMoveUrgent[ActionMove] > bullet.Distance {
					BulletMoveUrgent[ActionMove] = bullet.Distance
				}
				if BulletMoveUrgent[ActionLeft] > bullet.Distance {
					BulletMoveUrgent[ActionLeft] = bullet.Distance
				}
			case 3:
				// 影响左转
				if BulletMoveUrgent[ActionLeft] > bullet.Distance {
					BulletMoveUrgent[ActionLeft] = bullet.Distance
				}
			case 4:
				// 影响右转
				if BulletMoveUrgent[ActionRight] > bullet.Distance {
					BulletMoveUrgent[ActionRight] = bullet.Distance
				}
			case -1:
				// 影响直行、后退、开火和停止
				if BulletMoveUrgent[ActionMove] > bullet.Distance {
					BulletMoveUrgent[ActionMove] = bullet.Distance
				}
				if BulletMoveUrgent[ActionBack] > bullet.Distance {
					BulletMoveUrgent[ActionBack] = bullet.Distance
				}
				if BulletMoveUrgent[ActionStay] > bullet.Distance {
					BulletMoveUrgent[ActionStay] = bullet.Distance
				}
			case -2:
				// 影响左转、右转、开火和停止
				if BulletMoveUrgent[ActionLeft] > bullet.Distance {
					BulletMoveUrgent[ActionLeft] = bullet.Distance
				}
				if BulletMoveUrgent[ActionRight] > bullet.Distance {
					BulletMoveUrgent[ActionRight] = bullet.Distance
				}
				if BulletMoveUrgent[ActionStay] > bullet.Distance {
					BulletMoveUrgent[ActionStay] = bullet.Distance
				}
			case -3:
				// 影响直行、后退、开火和停止
				if BulletMoveUrgent[ActionMove] > bullet.Distance {
					BulletMoveUrgent[ActionMove] = bullet.Distance
				}
				if BulletMoveUrgent[ActionBack] > bullet.Distance {
					BulletMoveUrgent[ActionBack] = bullet.Distance
				}
				if BulletMoveUrgent[ActionStay] > bullet.Distance {
					BulletMoveUrgent[ActionStay] = bullet.Distance
				}
			case -4:
				// 影响左转、右转、开火和停止
				if BulletMoveUrgent[ActionLeft] > bullet.Distance {
					BulletMoveUrgent[ActionLeft] = bullet.Distance
				}
				if BulletMoveUrgent[ActionRight] > bullet.Distance {
					BulletMoveUrgent[ActionRight] = bullet.Distance
				}
				if BulletMoveUrgent[ActionStay] > bullet.Distance {
					BulletMoveUrgent[ActionStay] = bullet.Distance
				}
			}
		}
	}


	// 躲避为第二优先级
	//EnemyMoveUrgent := make(map[int]int)
	//EnemyMoveUrgent[ActionMove] = math.MaxInt32
	//EnemyMoveUrgent[ActionBack] = math.MaxInt32
	//EnemyMoveUrgent[ActionLeft] = math.MaxInt32
	//EnemyMoveUrgent[ActionRight] = math.MaxInt32
	//
	//EnemyMoveUrgent[ActionStay] = math.MaxInt32
	//EnemyMoveUrgent[ActionNone] = math.MaxInt32

	if enemyApproach == true {
		// 不能走敌军的象限
		for _, enemy := range *enemyThreat {
			switch quadrant[enemy.Quadrant] {
			case 1:
				// 影响直行和右转
				if BulletMoveUrgent[ActionMove] > enemy.Distance {
					BulletMoveUrgent[ActionMove] = enemy.Distance
				}
				if BulletMoveUrgent[ActionRight] > enemy.Distance {
					BulletMoveUrgent[ActionRight] = enemy.Distance
				}
			case 2:
				// 影响直行和左转
				if BulletMoveUrgent[ActionMove] > enemy.Distance {
					BulletMoveUrgent[ActionMove] = enemy.Distance
				}
				if BulletMoveUrgent[ActionLeft] > enemy.Distance {
					BulletMoveUrgent[ActionLeft] = enemy.Distance
				}
			case 3:
				// 影响左转和后退
				if BulletMoveUrgent[ActionLeft] > enemy.Distance {
					BulletMoveUrgent[ActionLeft] = enemy.Distance
				}
				if BulletMoveUrgent[ActionBack] > enemy.Distance {
					BulletMoveUrgent[ActionBack] = enemy.Distance
				}
			case 4:
				// 影响右转和后退
				if BulletMoveUrgent[ActionRight] > enemy.Distance {
					BulletMoveUrgent[ActionRight] = enemy.Distance
				}
				if BulletMoveUrgent[ActionBack] > enemy.Distance {
					BulletMoveUrgent[ActionBack] = enemy.Distance
				}
			case -1:
				// 枪口对准了，影响直行、后退、停止
				if BulletMoveUrgent[ActionMove] > enemy.Distance {
					BulletMoveUrgent[ActionMove] = enemy.Distance
				}
				if BulletMoveUrgent[ActionBack] > enemy.Distance {
					BulletMoveUrgent[ActionBack] = enemy.Distance
				}
				if BulletMoveUrgent[ActionStay] > enemy.Distance {
					BulletMoveUrgent[ActionStay] = enemy.Distance
				}
			case -2:
				// 影响左转、右转、停止
				if BulletMoveUrgent[ActionLeft] > enemy.Distance {
					BulletMoveUrgent[ActionLeft] = enemy.Distance
				}
				if BulletMoveUrgent[ActionRight] > enemy.Distance {
					BulletMoveUrgent[ActionRight] = enemy.Distance
				}
				if BulletMoveUrgent[ActionStay] > enemy.Distance {
					BulletMoveUrgent[ActionStay] = enemy.Distance
				}
			case -3:
				// 影响直行、后退、停止
				if BulletMoveUrgent[ActionMove] > enemy.Distance {
					BulletMoveUrgent[ActionMove] = enemy.Distance
				}
				if BulletMoveUrgent[ActionBack] > enemy.Distance {
					BulletMoveUrgent[ActionBack] = enemy.Distance
				}
				if BulletMoveUrgent[ActionStay] > enemy.Distance {
					BulletMoveUrgent[ActionStay] = enemy.Distance
				}
			case -4:
				// 影响左转、右转、停止
				if BulletMoveUrgent[ActionLeft] > enemy.Distance {
					BulletMoveUrgent[ActionLeft] = enemy.Distance
				}
				if BulletMoveUrgent[ActionRight] > enemy.Distance {
					BulletMoveUrgent[ActionRight] = enemy.Distance
				}
				if BulletMoveUrgent[ActionStay] > enemy.Distance {
					BulletMoveUrgent[ActionStay] = enemy.Distance
				}
			}
		}
	}

	// 对子弹威胁和敌军威胁进行找最大的，然后按照子弹威胁优先
	bulletMaxAction	 := -1
	bulletMaxUrgent  := -1
	bulletMinUrgent  := math.MaxInt32
	for i := 1; i < len(BulletMoveUrgent); i++ {
		if bulletMinUrgent > BulletMoveUrgent[i] {
			bulletMinUrgent = BulletMoveUrgent[i]
		}
		// 行动优先
		if bulletMaxUrgent < BulletMoveUrgent[i] && (i == ActionMove || i == ActionRight || i == ActionLeft || i == ActionBack){
			bulletMaxAction = i
			bulletMaxUrgent = BulletMoveUrgent[i]
		}
	}

	// 遵从本来的方向(2号行动)，如果原来的方向不为MAX，则顺序去找第一个大的行动。
	if BulletMoveUrgent[ActionMove] == math.MaxInt32 {
		// 继续行进
		return ActionMove, bulletMinUrgent
	}

	// 已经有了行动推荐
	if bulletMaxAction != -1 {
		return bulletMaxAction, bulletMinUrgent
	} else {
		// 行动推荐失败，除非停留绝对安全，否则进行其他判断
		if BulletMoveUrgent[ActionStay] == math.MaxInt32 {
			return ActionStay, bulletMinUrgent
		} else {
			// 此处光荣弹策略需要添加
			return bulletMaxAction, bulletMinUrgent
		}
	}

	return bulletMaxAction, bulletMinUrgent
}

func (self *Dodger) Suggest(tank *Tank, state *GameState, objective *Objective) SuggestionItem {
	// 如果采纳，计算几步被干掉
	//self.calcObjectiveUrgent(tank, state, objective)

	// 躲避炮弹，炮弹还有几步打到
	bulletApproach, enemyBullets := self.avoidBullet(tank, state)

	// 敌军的威胁，注意敌军到火线距离是可能炮弹的1/2
	enemyApproach, enemyThreat := self.threat(tank, state)

	// 最终决定往哪里躲
	if bulletApproach == true || enemyApproach == true{
		action, urgent := self.determine(tank, state, bulletApproach, &enemyBullets, enemyApproach, &enemyThreat)
		//fmt.Println("########")
		//fmt.Println(tank)
		//fmt.Println("--------")
		//fmt.Println(enemyBullets)
		//fmt.Println(enemyThreat)
		//fmt.Println("--------")
		//fmt.Println(action)
		//fmt.Println(urgent)
		//fmt.Println("########")
		return SuggestionItem {
			Action: action,
			Urgent: urgent,
		}
	}

	ret := SuggestionItem {
		Action: ActionNone,
		Urgent: 0,
	}
	return ret
}
