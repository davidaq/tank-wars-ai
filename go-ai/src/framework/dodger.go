/**
 * 高性能躲避行动子系统
 *
 * 说明
 子弹威胁算法
 假设子弹在
 *
 * author: linxingchen
 */
package framework

import (
	"math"
    "sort"
    "fmt"
)

type Dodger struct {
}

const MAX = 10000

/**
 * 获取隔墙情况数量
 */
func (self *Radar) getWallSum(state *GameState, my Tank, enemy EnemyThreat) int {
    wallSum := 0
    // 注意方向和象限系问题
    if enemy.Enemy.Y < my.Pos.Y {
        for y := enemy.Enemy.Y; y < my.Pos.Y; y++ {
            if state.Terain.Get(my.Pos.X, y) == TerainObstacle {
                wallSum++
            }
        }
    }

    if enemy.Enemy.X < my.Pos.X {
        for x := enemy.Enemy.X; x < my.Pos.X; x++ {
            if state.Terain.Get(x, my.Pos.Y) == TerainObstacle {
                wallSum++
            }
        }
    }

    if enemy.Enemy.Y > my.Pos.Y {
        for y := enemy.Enemy.Y; y > my.Pos.Y; y-- {
            if state.Terain.Get(my.Pos.X, y) == TerainObstacle {
                wallSum++
            }
        }
    }

    if enemy.Enemy.X > my.Pos.X {
        for x := enemy.Enemy.X; x > my.Pos.X; x-- {
            if state.Terain.Get(x, my.Pos.Y) == TerainObstacle {
                wallSum++
            }
        }
    }
    return wallSum
}

/**
 * 躲避系统
 *
 * 火线为威胁判断
 * 其他象限的威胁为躲避判断参考
 * 同时参考墙、草、队友阻挡调度	队友在一格范围内，同步协调
 */
func (self *Radar) dodge(state *GameState, bulletApproach bool, bullets *map[string][]BulletThreat, enemyApproach bool, enemys *map[string][]EnemyThreat) (map[string]RadarDodge, map[string]RadarDodge, map[string]RadarDodge, map[string][]ExtDangerSrc) {
    // 各种操作的紧急度
    moveUrgent := make(map[string]map[int]int)
    moveUrgentBullet := make(map[string]map[int]int)
    moveUrgentEnemy  := make(map[string]map[int]int)

    // 预存坦克数据
    tankData := make(map[string]Tank)

    // 预存坦克威胁度
    threat := make(map[string]int)
    threatBullet := make(map[string]int)
    threatEnemy  := make(map[string]int)

    // 各个躲不掉威胁的来源 ID - threat
    extDangerSrc := make(map[string][]ExtDangerSrc)

	for _, tank := range state.MyTank {
        tankData[tank.Id] = tank
        extDangerSrc[tank.Id] = []ExtDangerSrc{}
        // STEP0 初始化各种操作的威胁程度
        tmpMoveUrgent := make(map[int]int)
        tmpMoveUrgent[ActionStay] = MAX
        tmpMoveUrgent[ActionMove] = MAX
        tmpMoveUrgent[ActionLeft] = MAX
        tmpMoveUrgent[ActionBack] = MAX
        tmpMoveUrgent[ActionRight] = MAX
        // 初始化紧急程度
        tmpUrgent := MAX

        // 炮弹情况
        tmpMoveUrgentBullet := make(map[int]int)
        tmpMoveUrgentBullet[ActionStay] = MAX
        tmpMoveUrgentBullet[ActionMove] = MAX
        tmpMoveUrgentBullet[ActionLeft] = MAX
        tmpMoveUrgentBullet[ActionBack] = MAX
        tmpMoveUrgentBullet[ActionRight] = MAX
        tmpBulletUrgent := MAX

        // 敌军情况
        tmpMoveUrgentEnemy := make(map[int]int)
        tmpMoveUrgentEnemy[ActionStay] = MAX
        tmpMoveUrgentEnemy[ActionMove] = MAX
        tmpMoveUrgentEnemy[ActionLeft] = MAX
        tmpMoveUrgentEnemy[ActionBack] = MAX
        tmpMoveUrgentEnemy[ActionRight] = MAX
        tmpEnemyUrgent := MAX

        // 计算
        if bulletApproach == true && len((*bullets)[tank.Id]) > 0 {
            for _, b := range (*bullets)[tank.Id] {
                if b.Quadrant == QUADRANT_U || b.Quadrant == QUADRANT_L || b.Quadrant == QUADRANT_D || b.Quadrant == QUADRANT_R {
                    // 设置最终紧急度 如果不躲 两回合肯定挂掉
                    if b.Distances[b.Quadrant] <= state.Params.BulletSpeed * 4 {
                        tmpUrgent = 1
                        tmpBulletUrgent = 1
                    }

                    // 影响火线上的操作
                    if b.Quadrant == QUADRANT_U || b.Quadrant == QUADRANT_D {
                        // 影响直行、后退、停止
                        if tmpMoveUrgent[ActionMove] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionMove] = b.Distances[b.Quadrant]
                            tmpMoveUrgentBullet[ActionMove] = b.Distances[b.Quadrant]
                        }
                        if tmpMoveUrgent[ActionBack] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionBack] = b.Distances[b.Quadrant] - state.Params.BulletSpeed   //后退需要加一步
                            tmpMoveUrgentBullet[ActionBack] = b.Distances[b.Quadrant] - state.Params.BulletSpeed   //后退需要加一步
                        }
                        if tmpMoveUrgent[ActionStay] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionStay] = b.Distances[b.Quadrant] - state.Params.BulletSpeed   // 停下视作加一步
                            tmpMoveUrgentBullet[ActionStay] = b.Distances[b.Quadrant] - state.Params.BulletSpeed   // 停下视作加一步
                        }
                    } else {
                        // 影响左转、右转、停止
                        if tmpMoveUrgent[ActionLeft] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionLeft] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                            tmpMoveUrgentBullet[ActionLeft] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionRight] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionRight] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                            tmpMoveUrgentBullet[ActionRight] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionStay] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionStay] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                            tmpMoveUrgentBullet[ActionStay] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                        }
                    }
                    // 火线上躲不掉的情况，因为发射时子弹出现在前面位置，所以加1
                    // 两种情况 顺着行进方向（-1， -3），则需要两子弹速度 + 1。如旁边方向（-2， -4），则需要一子弹速度 + 1
                    if (b.Quadrant == QUADRANT_U || b.Quadrant == QUADRANT_D) && b.Distances[b.Quadrant] <= state.Params.BulletSpeed * 2 {
                        tmpUrgent = -1
                        tmpBulletUrgent = -1
                    }
                    if (b.Quadrant == QUADRANT_L || b.Quadrant == QUADRANT_R) && b.Distances[b.Quadrant] <= state.Params.BulletSpeed {
                        tmpUrgent = -1
                        tmpBulletUrgent = -1
                    }

                } else {
                    // 非火线的情况
                    // 这里计算 如果目前朝向 正好被击中 则为1 如不会被击中 则参与下面计算 计算只记录面向的象限
                    if b.Quadrant == QUADRANT_L_U && b.BulletPosition.Direction == DirectionRight {
                        tankMove := math.Ceil(float64(b.Distances[QUADRANT_U]) / float64(state.Params.BulletSpeed)) * float64(state.Params.TankSpeed)
                        if tankMove - float64(state.Params.TankSpeed) < float64(b.Distances[QUADRANT_L]) && tankMove >= float64(b.Distances[QUADRANT_L]) {
                            // 如果直行，则被侧面的子弹干掉
                            tmpMoveUrgent[ActionMove] = b.Distances[QUADRANT_U]
                            tmpMoveUrgentBullet[ActionMove] = b.Distances[QUADRANT_U]
                            tmpUrgent = 1
                            tmpBulletUrgent = 1
                            continue
                        }
                    }

                    if b.Quadrant == QUADRANT_R_U && b.BulletPosition.Direction == DirectionLeft {
                        tankMove := math.Ceil(float64(b.Distances[QUADRANT_U]) / float64(state.Params.BulletSpeed)) * float64(state.Params.TankSpeed)
                        if tankMove - float64(state.Params.TankSpeed) < float64(b.Distances[QUADRANT_R]) && tankMove >= float64(b.Distances[QUADRANT_R]) {
                            // 如果直行，则被侧面的子弹干掉
                            tmpMoveUrgent[ActionMove] = b.Distances[QUADRANT_U]
                            tmpMoveUrgentBullet[ActionMove] = b.Distances[QUADRANT_U]
                            tmpUrgent = 1
                            tmpBulletUrgent = 1
                            continue
                        }
                    }

                    // 计算作死的可能性
                    distance := 0
                    for _, v := range b.Distances {
                        // 因为是单方向子弹，所以直接计算
                        distance += int(math.Pow(float64(v), 2))
                    }
                    distance = int(math.Sqrt(float64(distance)))

                    if b.Direction == QUADRANT_U {
                        // 影响直行
                        if tmpMoveUrgent[ActionMove] > distance {
                            tmpMoveUrgent[ActionMove] = distance
                            tmpMoveUrgentBullet[ActionMove] = distance
                        }
                    }

                    if b.Direction == QUADRANT_L {
                        // 影响左转
                        if tmpMoveUrgent[ActionLeft] > distance {
                            tmpMoveUrgent[ActionLeft] = distance
                            tmpMoveUrgentBullet[ActionLeft] = distance
                        }
                    }

                    if b.Direction == QUADRANT_D {
                        // 影响后退
                        if tmpMoveUrgent[ActionBack] > distance {
                            tmpMoveUrgent[ActionBack] = distance
                            tmpMoveUrgentBullet[ActionBack] = distance
                        }
                    }

                    if b.Direction == QUADRANT_R {
                        // 影响右转
                        if tmpMoveUrgent[ActionRight] > distance {
                            tmpMoveUrgent[ActionRight] = distance
                            tmpMoveUrgentBullet[ActionRight] = distance
                        }
                    }
                }
                if tmpBulletUrgent == 1 || tmpBulletUrgent == -1 {
                    extDangerSrc[tank.Id] = append(extDangerSrc[tank.Id], ExtDangerSrc{
                        Source: b.BulletId,
                        Type:   BULLET_THREAT,
                        Urgent: tmpBulletUrgent,
                        Distance: b.Distances[b.Quadrant],
                    })
                }
            }
        }

        // 处理敌军情况 如果自己在草里，不用考虑敌军 威胁
        if enemyApproach == true && len((*enemys)[tank.Id]) > 0 && state.Terain.Get(tank.Pos.X, tank.Pos.Y) != TerainForest {
            var tmpThreat []EnemyThreat
            tmpLong := state.Params.BulletSpeed * 2 + 1
            tmpShort := state.Params.TankSpeed
            for _, e := range ((*enemys)[tank.Id]) {
                // 只考虑己方一步，不考虑自己的斜方向
                // 考虑敌军4条火线和子弹情况，距离为2个子弹射程。如果敌军tankbullet不为空，则坦克没有威胁
                // 因此考虑红十字形状的敌军范围，假设向前行进，则左右两侧（1， 2象限）和前方（-1象限）的坦克具有威胁
                for _, enemyBullet := range state.EnemyBullet {
                    if e.EnemyId == enemyBullet.From {
                        // 如果已经射出则没有威胁
                        continue
                    }
                }

                tmpAbsX := int(math.Abs(float64(e.Enemy.X - tank.Pos.X)))
                tmpAbsY := int(math.Abs(float64(e.Enemy.Y - tank.Pos.Y)))
                if e.OriginQuadrant == QUADRANT_U || e.OriginQuadrant == QUADRANT_L_U || e.OriginQuadrant == QUADRANT_R_U {
                    if (tmpAbsX <= tmpShort && tmpAbsY <= tmpLong) {
                        tmpThreat = append(tmpThreat, e)
                    }
                }

                if e.OriginQuadrant == QUADRANT_L || e.OriginQuadrant == QUADRANT_L_U || e.OriginQuadrant == QUADRANT_L_D {
                    if (tmpAbsX <= tmpLong && tmpAbsY <= tmpShort) {
                        tmpThreat = append(tmpThreat, e)
                    }
                }

                if e.OriginQuadrant == QUADRANT_D || e.OriginQuadrant == QUADRANT_L_D || e.OriginQuadrant == QUADRANT_R_D {
                    if (tmpAbsX <= tmpShort && tmpAbsY <= tmpLong) {
                        tmpThreat = append(tmpThreat, e)
                    }
                }

                if e.OriginQuadrant == QUADRANT_R || e.OriginQuadrant == QUADRANT_R_U || e.OriginQuadrant == QUADRANT_R_D {
                    if (tmpAbsX <= tmpLong && tmpAbsY <= tmpShort) {
                        tmpThreat = append(tmpThreat, e)
                    }
                }
            }

            // 对已经筛选出来的进行分析，这样对debug友好 TODO 计算墙
            for _, enemy := range tmpThreat {
                if enemy.Quadrant == QUADRANT_U || enemy.Quadrant == QUADRANT_L || enemy.Quadrant == QUADRANT_D || enemy.Quadrant == QUADRANT_R {
                    tmpUrgent = 1
                    tmpEnemyUrgent = 1
                    if enemy.Distances[enemy.Quadrant] < state.Params.BulletSpeed {
                        tmpUrgent = -1
                        tmpEnemyUrgent = -1
                    }
                } else {
                    // 其他的威胁度直接距离来做，最后除速度向上取整
                    distance := 0
                    for _, v := range enemy.Distances {
                        distance += int(math.Pow(float64(v), 2))
                    }
                    distance = int(math.Ceil(math.Sqrt(float64(distance)) / float64(state.Params.BulletSpeed)))

                    // 两侧永远不要小于等于1，否则容易误判紧急事件
                    if tmpUrgent > distance {
                        tmpUrgent = distance
                        tmpEnemyUrgent = distance
                    }


                }
                // 假设，要向前一步 那么有威胁的是相对象限的1，2，距离减1坦克速度的-1，和距离加1坦克速度的-3
                if enemy.Quadrant == QUADRANT_U || enemy.Quadrant == QUADRANT_R_U || enemy.Quadrant == QUADRANT_L_U || enemy.Quadrant == QUADRANT_D {
                    if enemy.Quadrant == QUADRANT_U || enemy.Quadrant == QUADRANT_D {
                        if enemy.Quadrant == QUADRANT_U {
                            if tmpMoveUrgent[ActionMove] > enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed {
                                tmpMoveUrgent[ActionMove] = enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed
                            }
                            if tmpMoveUrgentEnemy[ActionMove] > enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed {
                                tmpMoveUrgentEnemy[ActionMove] = enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed
                            }
                        }
                        if enemy.Quadrant == QUADRANT_D {
                            if tmpMoveUrgent[ActionMove] > enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed {
                                tmpMoveUrgent[ActionMove] = enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed
                            }
                            if tmpMoveUrgentEnemy[ActionMove] > enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed {
                                tmpMoveUrgentEnemy[ActionMove] = enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed
                            }
                        }
                    } else {
                        if enemy.Distances[QUADRANT_L] <= tmpShort || enemy.Distances[QUADRANT_R] <= tmpShort {
                            // 如果到U有墙，则跳过
                            if tmpMoveUrgent[ActionMove] > enemy.Distances[QUADRANT_U] {
                                tmpMoveUrgent[ActionMove] = enemy.Distances[QUADRANT_U]
                            }
                            if tmpMoveUrgentEnemy[ActionMove] > enemy.Distances[QUADRANT_U] {
                                tmpMoveUrgentEnemy[ActionMove] = enemy.Distances[QUADRANT_U]
                            }
                        }
                    }
                }

                // 假设，向左一步 那么有威胁的是相对象限的2,3象限距离减1坦克速度的-2和加1坦克速度的-4
                if enemy.Quadrant == QUADRANT_L || enemy.Quadrant == QUADRANT_L_U || enemy.Quadrant == QUADRANT_L_D || enemy.Quadrant == QUADRANT_R {
                    if enemy.Quadrant == QUADRANT_L || enemy.Quadrant == QUADRANT_R {
                        if enemy.Quadrant == QUADRANT_L {
                            if tmpMoveUrgent[ActionLeft] > enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed {
                                tmpMoveUrgent[ActionLeft] = enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed
                            }
                            if tmpMoveUrgentEnemy[ActionLeft] > enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed {
                                tmpMoveUrgentEnemy[ActionLeft] = enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed
                            }
                        }
                        if enemy.Quadrant == QUADRANT_R {
                            if tmpMoveUrgent[ActionLeft] > enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed {
                                tmpMoveUrgent[ActionLeft] = enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed
                            }
                            if tmpMoveUrgentEnemy[ActionLeft] > enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed {
                                tmpMoveUrgentEnemy[ActionLeft] = enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed
                            }
                        }
                    } else {
                        if enemy.Distances[QUADRANT_U] <= tmpShort || enemy.Distances[QUADRANT_D] <= tmpShort {
                            if tmpMoveUrgent[ActionLeft] > enemy.Distances[QUADRANT_L] {
                                tmpMoveUrgent[ActionLeft] = enemy.Distances[QUADRANT_L]
                            }
                            if tmpMoveUrgentEnemy[ActionLeft] > enemy.Distances[QUADRANT_L] {
                                tmpMoveUrgentEnemy[ActionLeft] = enemy.Distances[QUADRANT_L]
                            }
                        }
                    }
                }

                // 假设向下一步 那么有威胁的是相对象限的-3、-1 3,4
                if enemy.Quadrant == QUADRANT_D || enemy.Quadrant == QUADRANT_L_D || enemy.Quadrant == QUADRANT_R_D || enemy.Quadrant == QUADRANT_U {
                    if enemy.Quadrant == QUADRANT_D || enemy.Quadrant == QUADRANT_U {
                        if enemy.Quadrant == QUADRANT_D {
                            if tmpMoveUrgent[ActionBack] > enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed {
                                tmpMoveUrgent[ActionBack] = enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed
                            }
                            if tmpMoveUrgentEnemy[ActionBack] > enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed {
                                tmpMoveUrgentEnemy[ActionBack] = enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed
                            }
                        }
                        if enemy.Quadrant == QUADRANT_U {
                            if tmpMoveUrgent[ActionBack] > enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed {
                                tmpMoveUrgent[ActionBack] = enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed
                            }
                            if tmpMoveUrgentEnemy[ActionBack] > enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed {
                                tmpMoveUrgentEnemy[ActionBack] = enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed
                            }
                        }
                    } else {
                        if enemy.Distances[QUADRANT_L] <= tmpShort || enemy.Distances[QUADRANT_R] <= tmpShort {
                            if tmpMoveUrgent[ActionBack] > enemy.Distances[QUADRANT_D] {
                                tmpMoveUrgent[ActionBack] = enemy.Distances[QUADRANT_D]
                            }
                            if tmpMoveUrgentEnemy[ActionBack] > enemy.Distances[QUADRANT_D] {
                                tmpMoveUrgentEnemy[ActionBack] = enemy.Distances[QUADRANT_D]
                            }
                        }
                    }
                }

                // 假设向右一步，那么有威胁的是-2 -4 1 4
                if enemy.Quadrant == QUADRANT_R || enemy.Quadrant == QUADRANT_R_U || enemy.Quadrant == QUADRANT_R_D || enemy.Quadrant == QUADRANT_L {
                    if enemy.Quadrant == QUADRANT_R || enemy.Quadrant == QUADRANT_L {
                        if enemy.Quadrant == QUADRANT_R {
                            if tmpMoveUrgent[ActionRight] > enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed {
                                tmpMoveUrgent[ActionRight] = enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed
                            }
                            if tmpMoveUrgentEnemy[ActionRight] > enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed {
                                tmpMoveUrgentEnemy[ActionRight] = enemy.Distances[enemy.Quadrant] - state.Params.TankSpeed
                            }
                        }
                        if enemy.Quadrant == QUADRANT_L {
                            if tmpMoveUrgent[ActionRight] > enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed {
                                tmpMoveUrgent[ActionRight] = enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed
                            }
                            if tmpMoveUrgentEnemy[ActionRight] > enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed {
                                tmpMoveUrgentEnemy[ActionRight] = enemy.Distances[enemy.Quadrant] + state.Params.TankSpeed
                            }
                        }
                    } else {
                        if enemy.Distances[QUADRANT_U] <= tmpShort || enemy.Distances[QUADRANT_D] <= tmpShort {
                            if tmpMoveUrgent[ActionRight] > enemy.Distances[QUADRANT_R] {
                                tmpMoveUrgent[ActionRight] = enemy.Distances[QUADRANT_R]
                            }
                            if tmpMoveUrgentEnemy[ActionRight] > enemy.Distances[QUADRANT_R] {
                                tmpMoveUrgentEnemy[ActionRight] = enemy.Distances[QUADRANT_R]
                            }
                        }
                    }
                }

                // 假设停止则四个火线
                if enemy.Quadrant == QUADRANT_U || enemy.Quadrant == QUADRANT_L || enemy.Quadrant == QUADRANT_D || enemy.Quadrant == QUADRANT_R {
                    if tmpMoveUrgent[ActionStay] > enemy.Distances[enemy.Quadrant] {
                        tmpMoveUrgent[ActionStay] = enemy.Distances[enemy.Quadrant]
                    }
                    if tmpMoveUrgentEnemy[ActionStay] > enemy.Distances[enemy.Quadrant] {
                        tmpMoveUrgentEnemy[ActionStay] = enemy.Distances[enemy.Quadrant]
                    }
                }
                if tmpEnemyUrgent == 1 || tmpEnemyUrgent == -1 {
                    extDangerSrc[tank.Id] = append(extDangerSrc[tank.Id], ExtDangerSrc{
                        Source: enemy.EnemyId,
                        Type:   ENEMY_THREAT,
                        Urgent: tmpEnemyUrgent,
                        Distance: enemy.Distances[enemy.Quadrant],
                    })
                }
            }
        }

        threat[tank.Id] = tmpUrgent
        threatBullet[tank.Id] = tmpBulletUrgent
        threatEnemy[tank.Id] = tmpEnemyUrgent
        moveUrgent[tank.Id] = tmpMoveUrgent
        moveUrgentBullet[tank.Id] = tmpMoveUrgentBullet
        moveUrgentEnemy[tank.Id] = tmpMoveUrgentEnemy
    }

	// STEP4 综合场上局势进行协同调度
    radarDodge := make(map[string]RadarDodge)
    radarDodgeBullet := make(map[string]RadarDodge)
    radarDodgeEnemy  := make(map[string]RadarDodge)
    // 综合威胁
    radarDodge = self.calcDodge(moveUrgent, threat, state, tankData)
    // 子弹威胁
    radarDodgeBullet = self.calcDodge(moveUrgentBullet, threatBullet, state, tankData)
    // 敌军威胁
    radarDodgeEnemy = self.calcDodge(moveUrgentEnemy, threatEnemy, state, tankData)

    fmt.Println("###")
    fmt.Println(radarDodge)
    fmt.Println(radarDodgeEnemy)
    fmt.Println(radarDodgeBullet)
    fmt.Println("###")

	return radarDodge, radarDodgeBullet, radarDodgeEnemy, extDangerSrc
}

func (self *Radar) calcDodge(moveUrgent map[string]map[int]int, threat map[string]int, state *GameState, tankData map[string]Tank) map[string]RadarDodge {
    radarDodge := make(map[string]RadarDodge)
    for tankId, urgent := range moveUrgent {
        // STEP4.1 行动威胁排序，保存行动从大到小的key
        urgentV := []int{}  // 距离value列表
        urgentA := []int{}  // 行动key列表
        for _, vu := range urgent {
            urgentV = append(urgentV, vu)
        }
        sort.Sort(sort.Reverse(sort.IntSlice(urgentV)))
        // 根据行动去找类型
        tmpUniqueCheck := make(map[int]int)
        for _, sortUrgentV := range urgentV {
            for ku, vu := range urgent {
                if vu == sortUrgentV && tmpUniqueCheck[ku] != 1 {
                    urgentA = append(urgentA, ku)
                    tmpUniqueCheck[ku] = 1
                }
            }
        }

        // 如果都是MAX则无需推荐
        if urgentV[4] == MAX {
            radarDodge[tankId] = RadarDodge{
                Threat: 0,
                SafePos: Position{},
            }
            continue
        }

        // STEP4.2 环境变量分析
        // 尽可能往没有草的地方躲
        // 不能碰墙壁2回合判断，防止被逼到墙角
        actionWall := make(map[int]bool)
        nextPos    := make(map[int]Position)

        for _, action := range urgentA {
            canmove, tmpNextPos := self.convertActionToPosition(state, tankData[tankId], action, 1)
            if canmove == false {
                actionWall[action] = true
                continue
            }
            nextPos[action] = tmpNextPos
        }

        // 排除墙列表
        actionSequence := []int{}
        for _, action := range urgentA {
            if actionWall[action] != true {
                actionSequence = append(actionSequence, action)
            }
        }

        // 如果全都撞墙，则不去推荐
        if len(actionSequence) == 0 {
            radarDodge[tankId] = RadarDodge{
                Threat:  -1,
                SafePos: Position{},
            }
            continue
        }

        // STEP4.3 对行动类型列表进行分析，后面会对附近坦克进行计算
        // 优先度直接扔到第一位，方便后面处理
        // 前进优先，如果前进是满距离，而且相对最高（重复值一样）
        fin := false
        if urgent[ActionMove] == MAX {
            // 直接直行
            // 找key
            tmpActionSequence := actionSequence[0]
            for k, v := range actionSequence {
                if v == ActionMove {
                    actionSequence[0] = ActionMove
                    actionSequence[k] = tmpActionSequence
                    fin = true
                    break
                }
            }
        }



        // 尽可能去行动，而非停止
        if fin == false {
            // 没有被命中，继续执行
            // 最大的一堆中，如果行动也是MAX，则优先行动
            maxset := []int{}
            tmp := MAX
            for _, action := range actionSequence {
                if tmp == MAX {
                    tmp = urgent[action]
                }
                if tmp == urgent[action] {
                    maxset = append(maxset, action)
                } else {
                    break
                }
            }
            // 去计算已有的maxset中是否有行动
            tmpActionSequence := actionSequence[0]
            for k, a := range maxset {
                if a == ActionMove || a == ActionRight || a == ActionBack || a == ActionLeft {
                    actionSequence[0] = a
                    actionSequence[k] = tmpActionSequence
                    fin = true
                    break
                }
            }

        }

        // 计算最小的威胁作为闪避威胁度
        if threat[tankId] != 1 && threat[tankId] != -1{
            // 如果不是1或-1，则需要计算
            for i := 0; i < len(actionSequence); i++ {
                if threat[tankId] > urgent[actionSequence[i]] {
                    threat[tankId] = urgent[actionSequence[i]]
                }
            }
            if threat[tankId] <= 1 {
                threat[tankId] = 1
            }
        }

        // 如果为max则停止执行
        if threat[tankId] == MAX {
            radarDodge[tankId] = RadarDodge{
                Threat:  0,
                SafePos: Position{},
            }
            continue
        }

        // 先直接选最高的策略
        // 最高策略的下一步行进位置
        // 计算威胁度
        finThreat := 1 / (float64(threat[tankId]))
        radarDodge[tankId] = RadarDodge{
            Threat:  finThreat,
            SafePos: nextPos[actionSequence[0]],
        }
        // 草丛躲避
        // 查找是否有阻挡的己方坦克
    }
    return radarDodge
}

/**
 * 返回行动后坐标和能否行动
 */
func (self *Radar) convertActionToPosition(state *GameState, tank Tank, action int, step int) (bool, Position){
    distance := step * state.Params.TankSpeed
    positionRet := Position{}
    if action == ActionMove && tank.Pos.Direction == DirectionUp || action == ActionLeft && tank.Pos.Direction == DirectionRight || action == ActionBack && tank.Pos.Direction == DirectionDown || action == ActionRight && tank.Pos.Direction == DirectionLeft{
        positionRet.X = tank.Pos.X
        positionRet.Y = tank.Pos.Y - distance
    }
    if action == ActionLeft && tank.Pos.Direction == DirectionUp || action == ActionBack && tank.Pos.Direction == DirectionRight || action == ActionRight && tank.Pos.Direction == DirectionDown || action == ActionMove && tank.Pos.Direction == DirectionLeft{
        positionRet.X = tank.Pos.X - distance
        positionRet.Y = tank.Pos.Y
    }
    if action == ActionBack && tank.Pos.Direction == DirectionUp || action == ActionRight && tank.Pos.Direction == DirectionRight || action == ActionMove && tank.Pos.Direction == DirectionDown || action == ActionLeft && tank.Pos.Direction == DirectionLeft{
        positionRet.X = tank.Pos.X
        positionRet.Y = tank.Pos.Y + distance
    }
    if action == ActionRight && tank.Pos.Direction == DirectionUp || action == ActionMove && tank.Pos.Direction == DirectionRight || action == ActionLeft && tank.Pos.Direction == DirectionDown || action == ActionBack && tank.Pos.Direction == DirectionLeft{
        positionRet.X = tank.Pos.X + distance
        positionRet.Y = tank.Pos.Y
    }
    if action == ActionStay {
        positionRet.X = tank.Pos.X
        positionRet.Y = tank.Pos.Y
    }

    // 撞墙、超出地图边界判断
    // 如果能走但撞墙，则允许执行
    if positionRet.X == tank.Pos.X {
        if positionRet.Y >= tank.Pos.Y {
            for y := tank.Pos.Y; y <= positionRet.Y; y++ {
                if state.Terain.Get(positionRet.X, y) == TerainObstacle {
                    // 如果直行，则看能否挪动位置
                    if action == ActionMove && y - tank.Pos.Y > 1 {
                        return true, Position{X:tank.Pos.X, Y: y - 1}
                    }
                    return false, Position{}
                }
            }
        } else {
            for y := tank.Pos.Y; y >= positionRet.Y; y-- {
                if state.Terain.Get(positionRet.X, y) == TerainObstacle {
                    if action == ActionMove && tank.Pos.Y - y > 1 {
                        return true, Position{X: tank.Pos.X, Y: y + 1}
                    }
                    return false, Position{}
                }
            }
        }
    }
    if positionRet.Y == tank.Pos.Y {
        if positionRet.X >= tank.Pos.X {
            for x := tank.Pos.X; x <= positionRet.X; x++ {
                if state.Terain.Get(x, positionRet.Y) == TerainObstacle {
                    if action == ActionMove && x - tank.Pos.X > 1 {
                        return true, Position{X: x - 1, Y: tank.Pos.Y}
                    }
                    return false, Position{}
                }
            }
        } else {
            for x := tank.Pos.X; x >= positionRet.X; x-- {
                if state.Terain.Get(x, positionRet.Y) == TerainObstacle {
                    if action == ActionMove && tank.Pos.X - x > 1 {
                        return true, Position{X: x + 1, Y: tank.Pos.Y}
                    }
                    return false, Position{}
                }
            }
        }
    }

    // 检查碰撞己方和敌方坦克
    for _, s := range state.MyTank {
        // 排除自己
        if s.Id == tank.Id {
            continue
        }
        // distance可能是多个格
        if positionRet.X == tank.Pos.X {
            if positionRet.Y >= tank.Pos.Y {
                for y := tank.Pos.Y; y <= positionRet.Y; y++ {
                    if positionRet.X == s.Pos.X && y == s.Pos.Y {
                        if action == ActionMove && y - tank.Pos.Y > 1 {
                            return true, Position{X: tank.Pos.X, Y: y}
                        }
                        return false, Position{}
                    }
                }
            } else {
                for y := tank.Pos.Y; y >= positionRet.Y; y-- {
                    if positionRet.X == s.Pos.X && y == s.Pos.Y {
                        if action == ActionMove && tank.Pos.Y - y > 1 {
                            return true, Position{X: tank.Pos.X, Y: y}
                        }
                        return false, Position{}
                    }
                }
            }
        }
        if positionRet.Y == tank.Pos.Y {
            if positionRet.X >= tank.Pos.X {
                for x := tank.Pos.X; x <= positionRet.X; x++ {
                    if x == s.Pos.X && positionRet.Y == s.Pos.Y {
                        if action == ActionMove && x - tank.Pos.X > 1 {
                            return true, Position{X: x, Y: tank.Pos.Y}
                        }
                        return false, Position{}
                    }
                }
            } else {
                for x := tank.Pos.X; x >= positionRet.X; x-- {
                    if x == s.Pos.X && positionRet.Y == s.Pos.Y {
                        if action == ActionMove && tank.Pos.X - x > 1 {
                            return true, Position{X: x, Y: tank.Pos.Y}
                        }
                        return false, Position{}
                    }
                }
            }
        }
    }
    for _, e := range state.EnemyTank {
        // distance可能是多个格
        if positionRet.X == tank.Pos.X {
            if positionRet.Y >= tank.Pos.Y {
                for y := tank.Pos.Y; y <= positionRet.Y; y++ {
                    if positionRet.X == e.Pos.X && y == e.Pos.Y {
                        return false, Position{}
                    }
                }
            } else {
                for y := tank.Pos.Y; y >= positionRet.Y; y-- {
                    if positionRet.X == e.Pos.X && y == e.Pos.Y {
                        return false, Position{}
                    }
                }
            }
        }
        if positionRet.Y == tank.Pos.Y {
            if positionRet.X >= tank.Pos.X {
                for x := tank.Pos.X; x <= positionRet.X; x++ {
                    if x == e.Pos.X && positionRet.Y == e.Pos.Y {
                        return false, Position{}
                    }
                }
            } else {
                for x := tank.Pos.X; x >= positionRet.X; x-- {
                    if x == e.Pos.X && positionRet.Y == e.Pos.Y {
                        return false, Position{}
                    }
                }
            }
        }
    }

    return true, positionRet
}