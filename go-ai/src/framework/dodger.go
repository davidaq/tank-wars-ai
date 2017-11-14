/**
 * 高性能躲避行动子系统
 * author: linxingchen
 */
package framework

import (
	"math"
    "sort"
)

type Dodger struct {
}

const MAX = 10000

/**
 * 躲避系统
 *
 * 火线为威胁判断
 * 其他象限的威胁为躲避判断参考
 * 同时参考墙、草、队友阻挡调度	队友在一格范围内，同步协调
 */
func (self *Radar) dodge(state *GameState, bulletApproach bool, bullets *map[string][]BulletThreat, enemyApproach bool, enemys *map[string][]EnemyThreat) (map[string]RadarDodge) {
	radarDodge := make(map[string]RadarDodge)

    // 各种操作的紧急度
    moveUrgent := make(map[string]map[int]int)

    // 预存坦克数据
    tankData := make(map[string]Tank)

    // 预存坦克威胁度
    threat := make(map[string]int)

	for _, tank := range state.MyTank {
        tankData[tank.Id] = tank
        // STEP0 初始化各种操作的威胁程度
        tmpMoveUrgent := make(map[int]int)
        tmpMoveUrgent[ActionStay] = MAX
        tmpMoveUrgent[ActionMove] = MAX
        tmpMoveUrgent[ActionLeft] = MAX
        tmpMoveUrgent[ActionBack] = MAX
        tmpMoveUrgent[ActionRight] = MAX
        // 初始化紧急程度
        tmpUrgent := MAX

        // 计算
		if bulletApproach == true && len((*bullets)[tank.Id]) > 0 {
			for _, b := range (*bullets)[tank.Id] {
				if b.Quadrant == QUADRANT_U || b.Quadrant == QUADRANT_L || b.Quadrant == QUADRANT_D || b.Quadrant == QUADRANT_R {
                    // 设置最终紧急度 如果不躲肯定挂掉
                    tmpUrgent = 1

                    // 影响火线上的操作
                    if b.Quadrant == QUADRANT_U || b.Quadrant == QUADRANT_D {
                        // 影响直行、后退、停止
                        if tmpMoveUrgent[ActionMove] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionMove] = b.Distances[b.Quadrant]
                        }
                        if tmpMoveUrgent[ActionBack] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionBack] = b.Distances[b.Quadrant] - state.Params.BulletSpeed   //后退需要加一步
                        }
                        if tmpMoveUrgent[ActionStay] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionStay] = b.Distances[b.Quadrant] - state.Params.BulletSpeed   // 停下视作加一步
                        }
                    } else {
                        // 影响左转、右转、停止
                        if tmpMoveUrgent[ActionLeft] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionLeft] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionRight] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionRight] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionStay] > b.Distances[b.Quadrant] {
                            tmpMoveUrgent[ActionStay] = b.Distances[b.Quadrant] - state.Params.BulletSpeed
                        }
                    }
				} else {
                    // 非火线的情况
                    // 这里计算 如果目前朝向 正好被击中 则为1 如不会被击中 则参与下面计算 计算只记录面向的象限
                    if b.Quadrant == QUADRANT_L_U && b.BulletPosition.Direction == DirectionRight {
                        bulletMove := math.Ceil(float64(b.Distances[QUADRANT_L]) / float64(state.Params.TankSpeed)) * float64(state.Params.BulletSpeed)
                        if bulletMove - float64(state.Params.BulletSpeed / 2) <= float64(b.Distances[QUADRANT_U]) && bulletMove + float64(state.Params.BulletSpeed / 2) > float64(b.Distances[QUADRANT_U]) {
                            // 如果直行，则被侧面的子弹干掉
                            tmpUrgent = 1
                        }
                    }

                    if b.Quadrant == QUADRANT_R_U && b.BulletPosition.Direction == DirectionLeft {
                        bulletMove := math.Ceil(float64(b.Distances[QUADRANT_R]) / float64(state.Params.TankSpeed)) * float64(state.Params.BulletSpeed)
                        if bulletMove - float64(state.Params.BulletSpeed / 2) <= float64(b.Distances[QUADRANT_U]) && bulletMove + float64(state.Params.BulletSpeed / 2) > float64(b.Distances[QUADRANT_U]) {
                            // 如果直行，则被侧面的子弹干掉
                            tmpUrgent = 1
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
                        }
                    }

                    if b.Direction == QUADRANT_L {
                        // 影响左转
                        if tmpMoveUrgent[ActionLeft] > distance {
                            tmpMoveUrgent[ActionLeft] = distance
                        }
                    }

                    if b.Direction == QUADRANT_D {
                        // 影响后退
                        if tmpMoveUrgent[ActionBack] > distance {
                            tmpMoveUrgent[ActionBack] = distance
                        }
                    }

                    if b.Direction == QUADRANT_R {
                        // 影响右转
                        if tmpMoveUrgent[ActionRight] > distance {
                            tmpMoveUrgent[ActionRight] = distance
                        }
                    }
                }
			}
		}

        // 处理敌军情况
		if enemyApproach == true && len((*enemys)[tank.Id]) > 0 {
			for _, e := range ((*enemys)[tank.Id]) {
				if e.Quadrant == QUADRANT_U || e.Quadrant == QUADRANT_L || e.Quadrant == QUADRANT_D || e.Quadrant == QUADRANT_R {
                    // 敌军在火线上的处理
                    // 两回合炮弹距离则为紧急
                    if e.Distances[e.Quadrant] <= state.Params.BulletSpeed * 2 {
                        tmpUrgent = 1
                    }

                    // 影响直行、后退、停止
                    if e.Quadrant == QUADRANT_U || e.Quadrant == QUADRANT_D {
                        if tmpMoveUrgent[ActionMove] > e.Distances[e.Quadrant] {
                            tmpMoveUrgent[ActionMove] = e.Distances[e.Quadrant]
                        }
                        if tmpMoveUrgent[ActionBack] > e.Distances[e.Quadrant] {
                            tmpMoveUrgent[ActionBack] = e.Distances[e.Quadrant] - state.Params.BulletSpeed //后退需要加一步
                        }
                        if tmpMoveUrgent[ActionStay] > e.Distances[e.Quadrant] {
                            tmpMoveUrgent[ActionStay] = e.Distances[e.Quadrant] - state.Params.BulletSpeed
                        }
                    } else {
                        // 影响左转、右转、停止
                        if tmpMoveUrgent[ActionLeft] > e.Distances[e.Quadrant] {
                            tmpMoveUrgent[ActionLeft] = e.Distances[e.Quadrant] - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionRight] > e.Distances[e.Quadrant] {
                            tmpMoveUrgent[ActionRight] = e.Distances[e.Quadrant] - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionStay] > e.Distances[e.Quadrant] {
                            tmpMoveUrgent[ActionStay] = e.Distances[e.Quadrant] - state.Params.BulletSpeed
                        }
                    }
				} else {
                    // 敌军在其他象限的处理
                    if e.Quadrant == QUADRANT_L_U {
                        bulletMove := math.Ceil(float64(e.Distances[QUADRANT_L]) / float64(state.Params.TankSpeed)) * float64(state.Params.BulletSpeed)
                        if bulletMove - float64(state.Params.BulletSpeed / 2) <= float64(e.Distances[QUADRANT_U]) && bulletMove + float64(state.Params.BulletSpeed / 2) > float64(e.Distances[QUADRANT_U]) {
                            // 如果直行，则被侧面的子弹干掉
                            tmpUrgent = 1
                        }
                    }

                    if e.Quadrant == QUADRANT_R_U {
                        bulletMove := math.Ceil(float64(e.Distances[QUADRANT_R]) / float64(state.Params.TankSpeed)) * float64(state.Params.BulletSpeed)
                        if bulletMove - float64(state.Params.BulletSpeed / 2) <= float64(e.Distances[QUADRANT_U]) && bulletMove + float64(state.Params.BulletSpeed / 2) > float64(e.Distances[QUADRANT_U]) {
                            // 如果直行，则被侧面的子弹干掉
                            tmpUrgent = 1
                        }
                    }

                    // 计算作死的可能性
                    distance := 0
                    for _, v := range e.Distances {
                        distance += int(math.Pow(float64(v), 2))
                    }
                    distance = int(math.Sqrt(float64(distance)))

                    if e.Quadrant == QUADRANT_R_U {
                        // 影响直行和右转
                        if tmpMoveUrgent[ActionMove] > distance {
                            tmpMoveUrgent[ActionMove] = distance
                        }
                        if tmpMoveUrgent[ActionRight] > distance {
                            tmpMoveUrgent[ActionRight] = distance
                        }
                    }

                    if e.Quadrant == QUADRANT_L_U {
                        // 影响直行和左转
                        if tmpMoveUrgent[ActionMove] > distance {
                            tmpMoveUrgent[ActionMove] = distance
                        }
                        if tmpMoveUrgent[ActionLeft] > distance {
                            tmpMoveUrgent[ActionLeft] = distance
                        }
                    }

                    if e.Quadrant == QUADRANT_L_D {
                        // 影响左转和后退
                        if tmpMoveUrgent[ActionLeft] > distance {
                            tmpMoveUrgent[ActionLeft] = distance
                        }
                        if tmpMoveUrgent[ActionBack] > distance {
                            tmpMoveUrgent[ActionBack] = distance
                        }
                    }

                    if e.Quadrant == QUADRANT_R_D {
                        // 影响右转和后退
                        if tmpMoveUrgent[ActionRight] > distance {
                            tmpMoveUrgent[ActionRight] = distance
                        }
                        if tmpMoveUrgent[ActionBack] > distance {
                            tmpMoveUrgent[ActionBack] = distance
                        }
                    }
				}
			}
		}
        threat[tank.Id] = tmpUrgent
        moveUrgent[tank.Id] = tmpMoveUrgent
	}

	// STEP4 综合场上局势进行协同调度
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
        } else {
            //for _, action := range actionSequence {
            //    if urgent[action] == MAX {
            //        continue
            //    }
            //    if urgent[ActionMove] == urgent[action] {
            //        // 如果相对最高，直接直行
            //        tmpActionSequence := actionSequence[0]
            //        for k, v := range actionSequence {
            //            if v == ActionMove {
            //                actionSequence[0] = ActionMove
            //                actionSequence[k] = tmpActionSequence
            //                fin = true
            //                break
            //            }
            //        }
            //    }
            //    break
            //}
        }

        // 尽可能去行动，而非停止
        if fin == false {
        //    // 没有被命中，继续执行
        //    // 最大的一堆中，如果行动也是MAX，则优先行动
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
        if threat[tankId] != 1 {
            // 如果不是1，则需要计算
            for i := 0; i < len(actionSequence); i++ {
                if threat[tankId] > urgent[actionSequence[i]] {
                    threat[tankId] = urgent[actionSequence[i]]
                }
            }
            if threat[tankId] <= 1 {
                threat[tankId] = 1
            }
        }


        // 先直接选最高的策略
        // 最高策略的下一步行进位置
        // 计算威胁度
        finThreat := 1 / float64(threat[tankId])
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
    if state.Terain.Get(positionRet.X, positionRet.Y) == 1{
        return false, Position{}
    }
    return true, positionRet
}


