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

const (
	BULLET_THREAT = 1
	ENEMY_THREAT  = 2
)


type BulletThreat struct {
	BulletPosition Position
	BulletId	string // 子弹id
	Quadrant	int	// 相对于坦克的第几象限
	Direction 	int // 朝向哪个象限（能判断出转换后的方向）
	Distances 	map[int]int	// 距离四方向火线的威胁度
}

type EnemyThreat struct {
	Enemy 		Position
	EnemyId		string 	// 敌军id
	OriginQuadrant int  // 原始象限
    Quadrant    int 	// 敌军坦克所在的象限
    Distances   map[int]int // 坦克火线象限 - 垂直坦克火线的距离，如果敌军在坦克火线上，则为水平距离
}

type ExtDangerSrc struct {
	Source		string 	// 威胁来源
	Type		int     // 威胁来源种类 BULLET_THREAT = 1 ENEMY_THREAT = 2
	Urgent 		int		// 威胁度
	Distance 	int		// 距离
}

// 侦测几回合的威胁
const RADAR_BULLET_STEP = 5
const RADAR_ENEMY_STEP	= 5

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
			if math.Abs(float64(bullet.Pos.X - tank.Pos.X)) <= float64(radius) && math.Abs(float64(bullet.Pos.Y - tank.Pos.Y)) <= float64(radius) {
				// 象限
				tmpQuadrant 	:= QUADRANT_NONE
				tmpDirection 	:= QUADRANT_NONE
				tmpBulletQuadrantThreat := make(map[int]int)
				if bullet.Pos.X < tank.Pos.X {
					if bullet.Pos.Y < tank.Pos.Y && (bullet.Pos.Direction == DirectionRight || bullet.Pos.Direction == DirectionDown) {
						tmpQuadrant = QUADRANT_L_U
						if bullet.Pos.Direction == DirectionRight {
							tmpDirection = QUADRANT_U
							for w := bullet.Pos.X; w <= tank.Pos.X; w++ {
								if 1 == state.Terain.Get(w, tank.Pos.Y) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
							}
						} else {
							tmpDirection = QUADRANT_L
							for w := bullet.Pos.Y; w <= tank.Pos.Y; w++ {
								if 1 == state.Terain.Get(tank.Pos.X, w) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
							}
						}
						// 计算相对火线的distances
						tmpBulletQuadrantThreat[QUADRANT_U] = tank.Pos.X - bullet.Pos.X
						tmpBulletQuadrantThreat[QUADRANT_L] = tank.Pos.Y - bullet.Pos.Y
					}
					if bullet.Pos.Y > tank.Pos.Y && (bullet.Pos.Direction == DirectionRight || bullet.Pos.Direction == DirectionUp) {
						tmpQuadrant = QUADRANT_L_D
						if bullet.Pos.Direction == DirectionRight {
							tmpDirection = QUADRANT_D
							for w := bullet.Pos.X; w <= tank.Pos.X; w++ {
								if 1 == state.Terain.Get(w, tank.Pos.Y) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
							}
						} else {
							tmpDirection = QUADRANT_L
							for w := bullet.Pos.Y; w >= tank.Pos.Y; w-- {
								if 1 == state.Terain.Get(tank.Pos.X, w) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
							}
						}
						tmpBulletQuadrantThreat[QUADRANT_L] = bullet.Pos.Y - tank.Pos.Y
						tmpBulletQuadrantThreat[QUADRANT_D] = tank.Pos.X - bullet.Pos.X
					}
				}
				if bullet.Pos.X > tank.Pos.X {
					if bullet.Pos.Y < tank.Pos.Y && (bullet.Pos.Direction == DirectionLeft || bullet.Pos.Direction == DirectionDown) {
						tmpQuadrant = QUADRANT_R_U
						if bullet.Pos.Direction == DirectionLeft {
							tmpDirection = QUADRANT_U
							for w := bullet.Pos.X; w >= tank.Pos.X; w-- {
								if 1 == state.Terain.Get(w, tank.Pos.Y) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
							}
						} else {
							tmpDirection = QUADRANT_R
							for w := bullet.Pos.Y; w <= tank.Pos.Y; w++ {
								if 1 == state.Terain.Get(tank.Pos.X, w) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
							}
						}
						tmpBulletQuadrantThreat[QUADRANT_U] = bullet.Pos.X - tank.Pos.X
						tmpBulletQuadrantThreat[QUADRANT_R] = tank.Pos.Y - bullet.Pos.Y
					}
					if bullet.Pos.Y > tank.Pos.Y && (bullet.Pos.Direction == DirectionLeft || bullet.Pos.Direction == DirectionUp) {
						tmpQuadrant = QUADRANT_R_D
						if bullet.Pos.Direction == DirectionLeft {
							tmpDirection = QUADRANT_D
							for w := bullet.Pos.X; w >= tank.Pos.X; w-- {
								if 1 == state.Terain.Get(w, tank.Pos.Y) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
							}
						} else {
							tmpDirection = QUADRANT_R
							for w := bullet.Pos.Y; w <= tank.Pos.Y; w++ {
								if 1 == state.Terain.Get(tank.Pos.X, w) {
									tmpQuadrant = QUADRANT_NONE
									break
								}
							}
						}
						tmpBulletQuadrantThreat[QUADRANT_D] = bullet.Pos.X - tank.Pos.X
						tmpBulletQuadrantThreat[QUADRANT_R] = bullet.Pos.Y - tank.Pos.Y
					}
				}
				if bullet.Pos.X == tank.Pos.X {
					//在Y火线上
					if bullet.Pos.Y < tank.Pos.Y && bullet.Pos.Direction == DirectionDown {
						tmpQuadrant = QUADRANT_U
						tmpDirection = QUADRANT_U
						for w := bullet.Pos.Y; w <= tank.Pos.Y; w++ {
							if 1 == state.Terain.Get(tank.Pos.X, w) {
								tmpQuadrant = QUADRANT_NONE
								break
							}
						}
						tmpBulletQuadrantThreat[QUADRANT_U] = tank.Pos.Y - bullet.Pos.Y
					}
					if bullet.Pos.Y > tank.Pos.Y && bullet.Pos.Direction == DirectionUp {
						tmpQuadrant = QUADRANT_D
						tmpDirection = QUADRANT_D
						for w := bullet.Pos.Y; w >= tank.Pos.Y; w-- {
							if 1 == state.Terain.Get(tank.Pos.X, w) {
								tmpQuadrant = QUADRANT_NONE
								break
							}
						}
						tmpBulletQuadrantThreat[QUADRANT_D] = bullet.Pos.Y - tank.Pos.Y
					}
				}
				if bullet.Pos.Y == tank.Pos.Y {
					//在X火线上
					if bullet.Pos.X < tank.Pos.X && bullet.Pos.Direction == DirectionRight {
						tmpQuadrant = QUADRANT_L
						tmpDirection = QUADRANT_L
						for w := bullet.Pos.X; w <= tank.Pos.X; w++ {
							if 1 == state.Terain.Get(w, tank.Pos.Y) {
								tmpQuadrant = QUADRANT_NONE
								break
							}
						}
						tmpBulletQuadrantThreat[QUADRANT_L] = tank.Pos.X - bullet.Pos.X
					}
					if bullet.Pos.X > tank.Pos.X && bullet.Pos.Direction == DirectionLeft {
						tmpQuadrant = QUADRANT_R
						tmpDirection = QUADRANT_R
						for w := bullet.Pos.X; w >= tank.Pos.X; w-- {
							if 1 == state.Terain.Get(w, tank.Pos.Y) {
								tmpQuadrant = QUADRANT_NONE
								break
							}
						}
						tmpBulletQuadrantThreat[QUADRANT_R] = bullet.Pos.X - tank.Pos.X
					}
				}

				if tmpQuadrant == QUADRANT_NONE {
					continue
				}

				tmpBullet = append(tmpBullet, BulletThreat{
					BulletPosition: bullet.Pos,
					BulletId: bullet.From,
					Quadrant: tmpQuadrant,
					Direction: tmpDirection,
					Distances: tmpBulletQuadrantThreat,
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
	radius := state.Params.BulletSpeed * RADAR_ENEMY_STEP

	// 计算敌军在自己的什么方位 无视墙，防止LYB苟墙角
	enemyRadar := make(map[string][]EnemyThreat)

	// 循环计算多个坦克
	for _, tank := range state.MyTank {
		var tmpEnemyThreat []EnemyThreat
		for _, enemyTank := range state.EnemyTank {
			// 检查是否需要关注 方型雷达
			if math.Abs(float64(enemyTank.Pos.X - tank.Pos.X)) <= float64(radius) && math.Abs(float64(enemyTank.Pos.Y - tank.Pos.Y)) <= float64(radius) {
				tmpEnemyThreat = append(tmpEnemyThreat, EnemyThreat{
					Enemy: enemyTank.Pos,
					EnemyId: enemyTank.Id,
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
					tmpEnemyThreat[k].OriginQuadrant = QUADRANT_R_U
                    tmpEnemyThreat[k].Quadrant = QUADRANT_R_U

                    // 计算相对火线的distances
                    tmpEnemyThreat[k].Distances[QUADRANT_U] = enemy.Enemy.X - tank.Pos.X
                    tmpEnemyThreat[k].Distances[QUADRANT_R] = tank.Pos.Y - enemy.Enemy.Y
				}
				if enemy.Enemy.X < tank.Pos.X {
					tmpEnemyThreat[k].OriginQuadrant = QUADRANT_L_U
					tmpEnemyThreat[k].Quadrant = QUADRANT_L_U

                    // 计算相对火线的distances
                    tmpEnemyThreat[k].Distances[QUADRANT_U] = tank.Pos.X - enemy.Enemy.X
                    tmpEnemyThreat[k].Distances[QUADRANT_L] = tank.Pos.Y - enemy.Enemy.Y
				}
			}

			if enemy.Enemy.Y > tank.Pos.Y {
				if enemy.Enemy.X < tank.Pos.X {
					tmpEnemyThreat[k].OriginQuadrant = QUADRANT_L_D
					tmpEnemyThreat[k].Quadrant = QUADRANT_L_D

                    // 计算相对火线的distances
                    tmpEnemyThreat[k].Distances[QUADRANT_L] = enemy.Enemy.Y - tank.Pos.Y
                    tmpEnemyThreat[k].Distances[QUADRANT_D] = tank.Pos.X - enemy.Enemy.X
				}
				if enemy.Enemy.X > tank.Pos.X {
					tmpEnemyThreat[k].OriginQuadrant = QUADRANT_R_D
					tmpEnemyThreat[k].Quadrant = QUADRANT_R_D

                    // 计算相对火线的distances
                    tmpEnemyThreat[k].Distances[QUADRANT_R] = enemy.Enemy.Y - tank.Pos.Y
                    tmpEnemyThreat[k].Distances[QUADRANT_D] = enemy.Enemy.X - tank.Pos.X
				}
			}

			if enemy.Enemy.X == tank.Pos.X {
				if enemy.Enemy.Y < tank.Pos.Y {
					tmpEnemyThreat[k].OriginQuadrant = QUADRANT_U
					tmpEnemyThreat[k].Quadrant = QUADRANT_U
                    tmpEnemyThreat[k].Distances[QUADRANT_U] = tank.Pos.Y - enemy.Enemy.Y // 当在火线上时由垂直距离改为水平距离
				}
				if enemy.Enemy.Y > tank.Pos.Y {
					tmpEnemyThreat[k].OriginQuadrant = QUADRANT_D
					tmpEnemyThreat[k].Quadrant = QUADRANT_D
                    tmpEnemyThreat[k].Distances[QUADRANT_D] = enemy.Enemy.Y - tank.Pos.Y
				}
			}

			if enemy.Enemy.Y == tank.Pos.Y {
				if enemy.Enemy.X < tank.Pos.X {
					tmpEnemyThreat[k].OriginQuadrant = QUADRANT_L
					tmpEnemyThreat[k].Quadrant = QUADRANT_L
                    tmpEnemyThreat[k].Distances[QUADRANT_L] = tank.Pos.X - enemy.Enemy.X
				}
				if enemy.Enemy.X > tank.Pos.X {
					tmpEnemyThreat[k].OriginQuadrant = QUADRANT_R
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
					(*bullets)[tank.Id][k].Direction = quadrant[(*bullets)[tank.Id][k].Direction]
					tmpConverted := make(map[int]int)
					for kds, vds := range (*bullets)[tank.Id][k].Distances {
						tmpConverted[quadrant[kds]] = vds
					}
					(*bullets)[tank.Id][k].Distances = tmpConverted
				}
			}
		}

		if enemyApproach == true {
			if len((*enemy)[tank.Id]) > 0 {
				for k := range (*enemy)[tank.Id] {
					(*enemy)[tank.Id][k].Quadrant = quadrant[(*enemy)[tank.Id][k].Quadrant]
					tmpConverted := make(map[int]int)
                    for kds, vds := range (*enemy)[tank.Id][k].Distances {
						tmpConverted[quadrant[kds]] = vds
                    }
					(*enemy)[tank.Id][k].Distances = tmpConverted
				}
			}
		}
	}
	return true
}

// 检查自己当前的位置以及面朝方向前进的位置周围是否安全
// 检查每个坦克四周开火命中率与代价
func (self *Radar) Scan(state *GameState, diff *DiffResult) *RadarResult {
	// 首先进行全图扫描
	fullmapThreat := self.fullMapThreat(state)

	// 躲避炮弹，炮弹还有几步打到
	bulletApproach, bullets := self.avoidBullet(state)

	// 敌军威胁
	enemyApproach, enemy := self.threat(state)

	// 转换象限
	self.convertQuadrant(state, bulletApproach, &bullets, enemyApproach, &enemy)

	// 准备躲避返回字段
	radarDodge := make(map[string]RadarDodge)
	radarDodgeBullet := make(map[string]RadarDodge)
	radarDodgeEnemy  := make(map[string]RadarDodge)
	extDangerSrc := make(map[string][]ExtDangerSrc)
	if bulletApproach != true && enemyApproach != true {
		for _, tank := range state.MyTank {
			radarDodge[tank.Id] = RadarDodge{}
			radarDodgeBullet[tank.Id] = RadarDodge{}
			radarDodgeEnemy[tank.Id] = RadarDodge{}
			extDangerSrc[tank.Id] = []ExtDangerSrc{}
		}
	} else {
		// 躲避系统（撞墙、友军、草丛警戒）
		radarDodge, radarDodgeBullet, radarDodgeEnemy, extDangerSrc = self.dodge(state, bulletApproach, &bullets, enemyApproach, &enemy)
	}

	// 开火系统
	attack := self.Attack(state, &enemy)

	// 返回
	ret := &RadarResult {
		Dodge: make(map[string]RadarDodge),
		DodgeBullet: make(map[string]RadarDodge),
		DodgeEnemy: make(map[string]RadarDodge),
		Fire: make(map[string]RadarFireAll),
		Bullet: make(map[string][]BulletThreat),
		Enemy: make(map[string][]EnemyThreat),
		ExtDangerSrc: make(map[string][]ExtDangerSrc),
		FullMapThreat: make(map[Position]float64),
	}
	for _, tank := range state.MyTank {
		if atk, ok := attack[tank.Id]; ok {
			ret.Fire[tank.Id] = *atk
		}
		ret.Dodge[tank.Id] = radarDodge[tank.Id]
		ret.DodgeBullet[tank.Id] = radarDodgeBullet[tank.Id]
		ret.DodgeEnemy[tank.Id] = radarDodgeEnemy[tank.Id]
		ret.Bullet[tank.Id] = bullets[tank.Id]
		ret.Enemy[tank.Id] = enemy[tank.Id]
		ret.ExtDangerSrc[tank.Id] = extDangerSrc[tank.Id]
	}
	ret.FullMapThreat = fullmapThreat
	return ret
}
