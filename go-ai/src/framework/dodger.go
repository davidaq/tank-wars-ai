/**
 * 高性能躲避行动子系统
 * author: linxingchen
 */
package framework

import (
	"fmt"
	"math"
    "sort"
)

type Dodger struct {
}

/**
 * 躲避系统
 *
 * 火线为威胁判断
 * 其他象限的威胁为躲避判断参考
 * 同时参考墙、草、队友阻挡调度	队友在一格范围内，同步协调
 */
func (self *Radar) dodge(state *GameState, bulletApproach bool, bullets *map[string][]BulletThreat, enemyApproach bool, enemys *map[string][]EnemyThreat) (map[string]RadarDodge) {
	radarDodge := make(map[string]RadarDodge)

	// 约定紧急程度为击中回合数倒数
	firelineThreat := make(map[string]map[int]int)
	quadrantThreat := make(map[string]map[int]int)
	for _, tank := range state.MyTank {
		// STEP1 计算火线上的威胁
		// 先算最紧急的火线威胁步数
		tmpFirelineThreat := make(map[int]int)
		// STEP2 观察除火线外四个象限
		// 再看非火线象限上的威胁步数
		tmpQuadrantThreat := make(map[int]int)
		if bulletApproach == true && len((*bullets)[tank.Id]) > 0 {
			for _, b := range (*bullets)[tank.Id] {
				if b.Quadrant == QUADRANT_U || b.Quadrant == QUADRANT_L || b.Quadrant == QUADRANT_D || b.Quadrant == QUADRANT_R {
					if tmpFirelineThreat[b.Quadrant] == 0 || tmpFirelineThreat[b.Quadrant] > b.Distance {
						tmpFirelineThreat[b.Quadrant] = b.Distance
					}
				} else {
					if tmpQuadrantThreat[b.Quadrant] == 0 || tmpQuadrantThreat[b.Quadrant] > b.Distance {
						tmpQuadrantThreat[b.Quadrant] = b.Distance
					}
				}
			}
		}

		if enemyApproach == true && len((*enemys)[tank.Id]) > 0 {
			for _, e := range ((*enemys)[tank.Id]) {
				if e.Quadrant == QUADRANT_U || e.Quadrant == QUADRANT_L || e.Quadrant == QUADRANT_D || e.Quadrant == QUADRANT_R {
					for _, distance := range e.Distances {
						if tmpFirelineThreat[e.Quadrant] == 0 || tmpFirelineThreat[e.Quadrant] > distance {
							tmpFirelineThreat[e.Quadrant] = distance
						}
					}
				} else {
					for quadrant, distance := range e.Distances {
						if tmpQuadrantThreat[quadrant] == 0 || tmpQuadrantThreat[quadrant] > distance {
							tmpQuadrantThreat[quadrant] = distance
						}
					}
				}
			}
		}

		firelineThreat[tank.Id] = tmpFirelineThreat
		quadrantThreat[tank.Id] = tmpQuadrantThreat
	}

    // 各种操作的紧急度
    moveUrgent := make(map[string]map[int]int)

	// STEP3 对威胁程度进行分析
	if bulletApproach == true || enemyApproach == true {
		for _, tank := range state.MyTank {
            // STEP3.0 初始化
            tmpMoveUrgent := make(map[int]int)
            tmpMoveUrgent[ActionStay] = math.MaxInt32
            tmpMoveUrgent[ActionMove] = math.MaxInt32
            tmpMoveUrgent[ActionLeft] = math.MaxInt32
            tmpMoveUrgent[ActionBack] = math.MaxInt32
            tmpMoveUrgent[ActionRight] = math.MaxInt32

			// STEP3.1 火线上的必须躲
			if len(firelineThreat[tank.Id]) >= 0 {
                radarDodge[tank.Id] = RadarDodge{
                    Threat: 1, // 火线上不进行行动肯定会被击毁
                }
				// 坦克在某个火线上遭遇袭击
				for quadrant, distance := range firelineThreat[tank.Id] {
                    if quadrant == QUADRANT_U || quadrant == QUADRANT_D {
                        // 影响直行、后退、停止
                        if tmpMoveUrgent[ActionMove] > distance {
                            tmpMoveUrgent[ActionMove] = distance
                        }
                        if tmpMoveUrgent[ActionBack] > distance {
                            tmpMoveUrgent[ActionBack] = distance - state.Params.BulletSpeed //后退需要加一步
                        }
                        if tmpMoveUrgent[ActionStay] > distance {
                            tmpMoveUrgent[ActionStay] = distance
                        }
                    } else {
                        // 影响左转、右转、停止
                        if tmpMoveUrgent[ActionLeft] > distance {
                            tmpMoveUrgent[ActionLeft] = distance - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionRight] > distance {
                            tmpMoveUrgent[ActionRight] = distance - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionStay] > distance {
                            tmpMoveUrgent[ActionStay] = distance
                        }
                    }
				}
			}

            // STEP3.2 其他象限情况
            if len(quadrantThreat[tank.Id]) >= 0 {
                for quadrant, distance := range quadrantThreat[tank.Id] {
                    if quadrant == QUADRANT_R_U {
                        // 影响直行和右转
                        if tmpMoveUrgent[ActionMove] > distance {
                            tmpMoveUrgent[ActionMove] = distance
                        }
                        if tmpMoveUrgent[ActionRight] > distance {
                            tmpMoveUrgent[ActionRight] = distance  - state.Params.BulletSpeed
                        }
                    }
                    if quadrant == QUADRANT_L_U {
                        // 影响直行和左转
                        if tmpMoveUrgent[ActionMove] > distance {
                            tmpMoveUrgent[ActionMove] = distance
                        }
                        if tmpMoveUrgent[ActionLeft] > distance {
                            tmpMoveUrgent[ActionLeft] = distance - state.Params.BulletSpeed
                        }
                    }
                    if quadrant == QUADRANT_L_D {
                        // 影响左转和后退
                        if tmpMoveUrgent[ActionLeft] > distance {
                            tmpMoveUrgent[ActionLeft] = distance - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionBack] > distance {
                            tmpMoveUrgent[ActionBack] = distance - state.Params.BulletSpeed
                        }
                    }
                    if quadrant == QUADRANT_R_D {
                        // 影响右转和后退
                        if tmpMoveUrgent[ActionRight] > distance {
                            tmpMoveUrgent[ActionRight] = distance - state.Params.BulletSpeed
                        }
                        if tmpMoveUrgent[ActionBack] > distance {
                            tmpMoveUrgent[ActionBack] = distance - state.Params.BulletSpeed
                        }
                    }
                }
            }
            moveUrgent[tank.Id] = tmpMoveUrgent
		}
	} else {
		for _, tank := range state.MyTank {
			radarDodge[tank.Id] = RadarDodge{
				Threat: 0,
				SafePos: Position{},
			}
		}
        return radarDodge
	}

	// STEP4 综合场上局势进行协同调度
    for tankId, urgent := range moveUrgent {
        // STEP4.1 行动威胁排序，保存行动从大到小的key
        urgentV := []int{}
        urgentA := []int{}
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
        
        // STEP4.2 环境变量分析
        // 尽可能往没有草的地方躲
        // 不能碰墙壁2回合判断，防止被逼到墙角
        actionWall := make(map[int]bool)
        nextPos    := make(map[int]Position)
        // 获取坦克数据
        tankData := Tank{}
        for _, tank := range state.MyTank {
            if tank.Id == tankId {
                tankData = tank
            }
        }
        for _, action := range urgentA {
            canmove, tmpNextPos := self.convertActionToPosition(state, tankData, action, 1)
            if canmove == false {
                actionWall[action] = true
                continue
            }
            nextPos[action] = tmpNextPos
        }

        // 排除列表
        toAction := []int{}
        for _, action := range urgentA {
            if actionWall[action] != true {
                toAction = append(toAction, action)
            }
        }

        // STEP4.3 对行动类型列表进行分析，后面会对附近坦克进行计算
        // 前进优先

        // 其他选最高的

        // 草丛躲避

        fmt.Println("######")
        fmt.Println(urgentA)
        fmt.Println(actionWall)
        fmt.Println(toAction)
        fmt.Println("######")
        //for i := 1; i < len(urgent); i++ {
        //    if radarDodge[tankId].Threat > urgent[i] {
        //       radarDodge[tankId] = RadarDodge{
        //           Threat: urgent[i],
        //       }
        //    }
        //}



        // 查找是否有阻挡的己方坦克
    }

    // STEP5 动作转坐标点

	return radarDodge
    //
	//// 对子弹威胁和敌军威胁进行找最大的，然后按照子弹威胁优先
	//bulletMaxAction	 := -1
	//bulletMaxUrgent  := -1
	//bulletMinUrgent  := math.MaxInt32
	//for i := 1; i < len(BulletMoveUrgent); i++ {
	//	if bulletMinUrgent > BulletMoveUrgent[i] {
	//		bulletMinUrgent = BulletMoveUrgent[i]
	//	}
	//	// 行动优先
	//	if bulletMaxUrgent < BulletMoveUrgent[i] && (i == ActionMove || i == ActionRight || i == ActionLeft || i == ActionBack){
	//		bulletMaxAction = i
	//		bulletMaxUrgent = BulletMoveUrgent[i]
	//	}
	//}
    //
	//// 遵从本来的方向(2号行动)，如果原来的方向不为MAX，则顺序去找第一个大的行动。
	//if BulletMoveUrgent[ActionMove] == math.MaxInt32 {
	//	// 继续行进
	//	return ActionMove, bulletMinUrgent
	//}
    //
	//// 已经有了行动推荐
	//if bulletMaxAction != -1 {
	//	return bulletMaxAction, bulletMinUrgent
	//} else {
	//	// 行动推荐失败，除非停留绝对安全，否则进行其他判断
	//	if BulletMoveUrgent[ActionStay] == math.MaxInt32 {
	//		return ActionStay, bulletMinUrgent
	//	} else {
	//		// 此处光荣弹策略需要添加
	//		return bulletMaxAction, bulletMinUrgent
	//	}
	//}
    //
	//return bulletMaxAction, bulletMinUrgent
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

    // 撞墙、超出地图边界判断
    if state.Terain.Get(positionRet.X, positionRet.Y) == 1{
        return false, Position{}
    }

    return true, positionRet
}


