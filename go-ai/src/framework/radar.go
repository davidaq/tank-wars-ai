/**
 * 高性能最怂雷达系统，用于侦测威胁以及可以开火的目标
 * author: linxingchen
 */
package framework;

import (
	"math"
)

type Radar struct {
}

func NewRadar() *Radar {
	inst := &Radar {
	}
	return inst
}

// 8象限系统
const (
	QUADRANT_NONE = 0
	QUADRANT_U = -1		// 上方
	QUADRANT_L = -2		// 左侧
	QUADRANT_D = -3		// 下方
	QUADRANT_R = -4		// 右侧
	QUADRANT_R_U = 1	// 第一象限 右上角
	QUADRANT_L_U = 2	// 第二象限 左上角
	QUADRANT_L_D = 3	// 第三象限 左下角
	QUADRANT_R_D = 4	// 第四象限 右下角
)

type BulletThreat struct {
	BulletPosition Position
	Distance 	int	// 距离四方向火线的威胁度
	Quadrant	int	// 相对于坦克的第几象限
}

type EnemyThreat struct {
	Enemy 		Position
    Quadrant    int // 敌军坦克所在的象限
    Distances   map[int]int // 坦克火线象限 - 垂直坦克火线的距离，如果敌军在坦克火线上，则为水平距离
}

// 侦测几回合的威胁
const RADAR_BULLET_STEP = 4
const RADAR_ENEMY_STEP	= 3

func (self *Radar) avoidBullet(state *GameState) (bulletApproach bool, bulletThreat map[string][]BulletThreat){
	// 雷达半径由步数实际算出
	radius := state.Params.BulletSpeed * RADAR_BULLET_STEP

	// 双方子弹merge
	bulletMerge := []Bullet{}
	for _, tmpEBullet := range state.EnemyBullet {
		bulletMerge = append(bulletMerge, tmpEBullet)
	}
	for _, tmpMBullet := range state.MyBullet {
		bulletMerge = append(bulletMerge, tmpMBullet)
	}

	// 计算子弹在自己的什么位置
	bulletRadar := make(map[string][]BulletThreat)
	// 循环计算各个坦克
	for _, tank := range state.MyTank {
		var tmpBullet []BulletThreat
		for _, bullet := range bulletMerge {
			// 检查是否需要关注
			if math.Pow(float64(bullet.Pos.X - tank.Pos.X), 2) + math.Pow(float64(bullet.Pos.Y - tank.Pos.Y), 2) <= float64(radius * radius) {
				// 象限
				tmpQuadrant := QUADRANT_NONE
				tmpDistance := 0	// 注意要除子弹速度，这里是**距离火线位置**。威胁度还需要另外处理
				if bullet.Pos.X < tank.Pos.X {
					if bullet.Pos.Y < tank.Pos.Y && (bullet.Pos.Direction == DirectionRight || bullet.Pos.Direction == DirectionDown) {
						tmpQuadrant = QUADRANT_L_U
						if bullet.Pos.Direction == DirectionRight {
							for w := bullet.Pos.X; w <= tank.Pos.X; w++ {
								if 1 == state.Terain.Get(w, tank.Pos.Y) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
								tmpDistance++
							}
						} else {
							for w := bullet.Pos.Y; w <= tank.Pos.Y; w++ {
								if 1 == state.Terain.Get(tank.Pos.X, w) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
								tmpDistance++
							}
						}
					}
					if bullet.Pos.Y > tank.Pos.Y && (bullet.Pos.Direction == DirectionRight || bullet.Pos.Direction == DirectionUp) {
						tmpQuadrant = QUADRANT_L_D
						if bullet.Pos.Direction == DirectionRight {
							for w := bullet.Pos.X; w <= tank.Pos.X; w++ {
								if 1 == state.Terain.Get(w, tank.Pos.Y) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
								tmpDistance++
							}
						} else {
							for w := bullet.Pos.Y; w >= tank.Pos.Y; w-- {
								if 1 == state.Terain.Get(tank.Pos.X, w) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
								tmpDistance++
							}
						}
					}
				}
				if bullet.Pos.X > tank.Pos.X {
					if bullet.Pos.Y < tank.Pos.Y && (bullet.Pos.Direction == DirectionLeft || bullet.Pos.Direction == DirectionDown) {
						tmpQuadrant = QUADRANT_R_U
						if bullet.Pos.Direction == DirectionLeft {
							for w := bullet.Pos.X; w >= tank.Pos.X; w-- {
								if 1 == state.Terain.Get(w, tank.Pos.Y) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
								tmpDistance++
							}
						} else {
							for w := bullet.Pos.Y; w <= tank.Pos.Y; w++ {
								if 1 == state.Terain.Get(tank.Pos.X, w) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
								tmpDistance++
							}
						}
					}
					if bullet.Pos.Y > tank.Pos.Y && (bullet.Pos.Direction == DirectionLeft || bullet.Pos.Direction == DirectionUp) {
						tmpQuadrant = QUADRANT_R_D
						if bullet.Pos.Direction == DirectionLeft {
							for w := bullet.Pos.X; w >= tank.Pos.X; w-- {
								if 1 == state.Terain.Get(w, tank.Pos.Y) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
								tmpDistance++
							}
						} else {
							for w := bullet.Pos.Y; w <= tank.Pos.Y; w++ {
								if 1 == state.Terain.Get(tank.Pos.X, w) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
								tmpDistance++
							}
						}
					}
				}
				if bullet.Pos.X == tank.Pos.X {
					//在Y火线上
					if bullet.Pos.Y < tank.Pos.Y && bullet.Pos.Direction == DirectionDown {
						tmpQuadrant = QUADRANT_U
						for w := bullet.Pos.Y; w <= tank.Pos.Y; w++ {
							if 1 == state.Terain.Get(tank.Pos.X, w) {
								tmpQuadrant = QUADRANT_NONE
								break
							}
							tmpDistance++
						}
					}
					if bullet.Pos.Y > tank.Pos.Y && bullet.Pos.Direction == DirectionUp {
						tmpQuadrant = QUADRANT_D
						for w := bullet.Pos.Y; w >= tank.Pos.Y; w-- {
							if 1 == state.Terain.Get(tank.Pos.X, w) {
								tmpQuadrant = QUADRANT_NONE
								break
							}
							tmpDistance++
						}
					}
				}
				if bullet.Pos.Y == tank.Pos.Y {
					//在X火线上
					if bullet.Pos.X < tank.Pos.X && bullet.Pos.Direction == DirectionRight {
						tmpQuadrant = QUADRANT_L
						for w := bullet.Pos.X; w <= tank.Pos.X; w++ {
							if 1 == state.Terain.Get(w, tank.Pos.Y) {
								tmpQuadrant = QUADRANT_NONE
								break
							}
							tmpDistance++
						}
					}
					if bullet.Pos.X > tank.Pos.X && bullet.Pos.Direction == DirectionLeft {
						tmpQuadrant = QUADRANT_R
						for w := bullet.Pos.X; w >= tank.Pos.X; w-- {
							if 1 == state.Terain.Get(w, tank.Pos.Y) {
								tmpQuadrant = QUADRANT_NONE
								break
							}
							tmpDistance++
						}
					}
				}

				if tmpQuadrant == QUADRANT_NONE || tmpDistance == 0 {
					continue
				}

				tmpBullet = append(tmpBullet, BulletThreat{
					BulletPosition: bullet.Pos,
					Distance: tmpDistance,	// 这里不做处理，后面再综合除
					Quadrant: tmpQuadrant,
				})
			}
		}
		bulletRadar[tank.Id] = tmpBullet
	}

	for _, tmp := range bulletRadar {
		if len(tmp) > 0 {
			return true, bulletRadar
		}
	}
	return false, bulletRadar
}

func (self *Radar) threat(state *GameState) (threat bool, enemyThreat map[string][]EnemyThreat) {
	// 雷达半径由步数实际算出
	radius := state.Params.TankSpeed * RADAR_ENEMY_STEP

	// 计算敌军在自己的什么方位 无视墙，防止LYB苟墙角
	enemyRadar := make(map[string][]EnemyThreat)

	// 循环计算多个坦克
	for _, tank := range state.MyTank {
		var tmpEnemyThreat []EnemyThreat
		for _, enemyTank := range state.EnemyTank {
			// 检查是否需要关注 圆形雷达
			if math.Pow(float64(enemyTank.Pos.X - tank.Pos.X), 2) + math.Pow(float64(enemyTank.Pos.Y - tank.Pos.Y), 2) <= float64(radius * radius) {
				tmpEnemyThreat = append(tmpEnemyThreat, EnemyThreat{
					Enemy: enemyTank.Pos,
				})
			}
		}

		if 0 == len(tmpEnemyThreat) {
			enemyRadar[tank.Id] = tmpEnemyThreat
			continue
		}

		// 计算方位和象限
		for k, enemy := range tmpEnemyThreat {
            tmpEnemyThreat[k].Distances = make(map[int]int)
			if enemy.Enemy.Y < tank.Pos.Y {
				if enemy.Enemy.X > tank.Pos.X {
					// 计算象限
                    tmpEnemyThreat[k].Quadrant = QUADRANT_R_U

                    // 计算相对火线的distances
                    tmpEnemyThreat[k].Distances[QUADRANT_U] = enemy.Enemy.X - tank.Pos.X
                    tmpEnemyThreat[k].Distances[QUADRANT_R] = tank.Pos.Y - enemy.Enemy.Y
				}
				if enemy.Enemy.X < tank.Pos.X {
					tmpEnemyThreat[k].Quadrant = QUADRANT_L_U

                    // 计算相对火线的distances
                    tmpEnemyThreat[k].Distances[QUADRANT_U] = tank.Pos.X - enemy.Enemy.X
                    tmpEnemyThreat[k].Distances[QUADRANT_L] = tank.Pos.Y - enemy.Enemy.Y
				}
			}

			if enemy.Enemy.Y > tank.Pos.Y {
				if enemy.Enemy.X < tank.Pos.X {
					tmpEnemyThreat[k].Quadrant = QUADRANT_L_D

                    // 计算相对火线的distances
                    tmpEnemyThreat[k].Distances[QUADRANT_L] = enemy.Enemy.Y - tank.Pos.Y
                    tmpEnemyThreat[k].Distances[QUADRANT_D] = tank.Pos.X - enemy.Enemy.X
				}
				if enemy.Enemy.X > tank.Pos.X {
					tmpEnemyThreat[k].Quadrant = QUADRANT_R_D

                    // 计算相对火线的distances
                    tmpEnemyThreat[k].Distances[QUADRANT_R] = enemy.Enemy.Y - tank.Pos.Y
                    tmpEnemyThreat[k].Distances[QUADRANT_D] = enemy.Enemy.X - tank.Pos.X
				}
			}

			if enemy.Enemy.X == tank.Pos.X {
				if enemy.Enemy.Y < tank.Pos.Y {
					tmpEnemyThreat[k].Quadrant = QUADRANT_U
                    tmpEnemyThreat[k].Distances[QUADRANT_U] = tank.Pos.Y - enemy.Enemy.Y // 当在火线上时由垂直距离改为水平距离
				}
				if enemy.Enemy.Y > tank.Pos.Y {
					tmpEnemyThreat[k].Quadrant = QUADRANT_D
                    tmpEnemyThreat[k].Distances[QUADRANT_D] = enemy.Enemy.Y - tank.Pos.Y
				}
			}

			if enemy.Enemy.Y == tank.Pos.Y {
				if enemy.Enemy.X < tank.Pos.X {
					tmpEnemyThreat[k].Quadrant = QUADRANT_L
                    tmpEnemyThreat[k].Distances[QUADRANT_L] = tank.Pos.X - enemy.Enemy.X
				}
				if enemy.Enemy.X > tank.Pos.X {
					tmpEnemyThreat[k].Quadrant = QUADRANT_R
                    tmpEnemyThreat[k].Distances[QUADRANT_R] = enemy.Enemy.X - tank.Pos.X
				}
			}
		}
		enemyRadar[tank.Id] = tmpEnemyThreat
	}

	for _, tmp := range enemyRadar {
		if len(tmp) != 0 {
			return true, enemyRadar
		}
	}
	return false, enemyRadar
}

/**
 * 象限转置 从假定向上的象限系统转换为相对于实际坦克方向的象限系统
 * **禁止随意改动**
 */
func (self *Radar) convertQuadrant(state *GameState, bulletApproach bool, bullets *map[string][]BulletThreat, enemyApproach bool, enemy *map[string][]EnemyThreat) (bool){
	if bulletApproach == false && enemyApproach == false {
		return false
	}

	for _, tank := range state.MyTank {

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
			quadrant[-2] = -1
			quadrant[-3] = -2
			quadrant[-4] = -3
			quadrant[-1] = -4

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

		if bulletApproach == true {
			if len((*bullets)[tank.Id]) > 0 {
				for k := range (*bullets)[tank.Id] {
					(*bullets)[tank.Id][k].Quadrant = quadrant[(*bullets)[tank.Id][k].Quadrant]
				}
			}
		}

		if enemyApproach == true {
			if len((*enemy)[tank.Id]) > 0 {
				for k := range (*enemy)[tank.Id] {
					(*enemy)[tank.Id][k].Quadrant = quadrant[(*enemy)[tank.Id][k].Quadrant]
                    for kds, vds := range (*enemy)[tank.Id][k].Distances {
                        tmpQuadrantKey := quadrant[kds]
                        delete((*enemy)[tank.Id][k].Distances, kds)
                        (*enemy)[tank.Id][k].Distances[tmpQuadrantKey] = vds
                    }
				}
			}
		}
	}
	return true
}

// 检查自己当前的位置以及面朝方向前进的位置周围是否安全
// 检查每个坦克四周开火命中率与代价
func (self *Radar) Scan(state *GameState) *RadarResult {
	// 躲避炮弹，炮弹还有几步打到
	bulletApproach, bullets := self.avoidBullet(state)

	// 敌军威胁
	enemyApproach, enemy := self.threat(state)

	// 转换象限
	self.convertQuadrant(state, bulletApproach, &bullets, enemyApproach, &enemy)

	// 准备躲避返回字段
	radarDodge := make(map[string]RadarDodge)
	if bulletApproach != true && enemyApproach != true {
		for _, tank := range state.MyTank {
			radarDodge[tank.Id] = RadarDodge{}
		}
	} else {
		// 躲避系统（撞墙、友军、草丛警戒）
		radarDodge = self.dodge(state, bulletApproach, &bullets, enemyApproach, &enemy)
	}

	// 开火系统


	// 返回
	ret := &RadarResult {
		Dodge: make(map[string]RadarDodge),
		Fire: make(map[string]RadarFireAll),
	}
	for _, tank := range state.MyTank {
		ret.Fire[tank.Id] = RadarFireAll {
			Up: &RadarFire {
				Faith: 1.,
				Action: ActionFireUp,
				Cost: 10,
			},
		}
		ret.Dodge[tank.Id] = radarDodge[tank.Id]
	}
	return ret
}
