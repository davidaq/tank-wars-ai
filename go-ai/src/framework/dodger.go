/**
 * 高性能躲避行动子系统
 * author: linxingchen
 */
package framework

import (
	//"math"
)

type Dodger struct {
}

//func (self *Radar) dodge(state *GameState, bulletApproach bool, bullets *[]BulletThreat, enemyApproach bool, enemys *[]EnemyThreat) (map[string]RadarDodge) {
    //
	//// 最后收敛到几个方向上，直接在方向上标出最小的紧急度，最后走方向中紧急度排名第一，但是方向中不紧急的
	//BulletMoveUrgent := [6]int{}
	//BulletMoveUrgent[ActionMove] = math.MaxInt32
	//BulletMoveUrgent[ActionBack] = math.MaxInt32
	//BulletMoveUrgent[ActionLeft] = math.MaxInt32
	//BulletMoveUrgent[ActionRight] = math.MaxInt32
    //
	//BulletMoveUrgent[ActionStay] = math.MaxInt32
    //
    //
	//// 炮弹为第一优先级
	//if bulletApproach == true {
	//	// 炮弹的一条线都走不了
	//	for _, bullet := range *enemyBullets {
	//		switch quadrant[bullet.Quadrant] {
	//		case 1:
	//			// 影响直行和右转
	//			if BulletMoveUrgent[ActionMove] > bullet.Distance {
	//				BulletMoveUrgent[ActionMove] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionRight] > bullet.Distance {
	//				BulletMoveUrgent[ActionRight] = bullet.Distance
	//			}
	//		case 2:
	//			// 影响直行和左转
	//			if BulletMoveUrgent[ActionMove] > bullet.Distance {
	//				BulletMoveUrgent[ActionMove] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionLeft] > bullet.Distance {
	//				BulletMoveUrgent[ActionLeft] = bullet.Distance
	//			}
	//		case 3:
	//			// 影响左转
	//			if BulletMoveUrgent[ActionLeft] > bullet.Distance {
	//				BulletMoveUrgent[ActionLeft] = bullet.Distance
	//			}
	//		case 4:
	//			// 影响右转
	//			if BulletMoveUrgent[ActionRight] > bullet.Distance {
	//				BulletMoveUrgent[ActionRight] = bullet.Distance
	//			}
	//		case -1:
	//			// 影响直行、后退、开火和停止
	//			if BulletMoveUrgent[ActionMove] > bullet.Distance {
	//				BulletMoveUrgent[ActionMove] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionBack] > bullet.Distance {
	//				BulletMoveUrgent[ActionBack] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionStay] > bullet.Distance {
	//				BulletMoveUrgent[ActionStay] = bullet.Distance
	//			}
	//		case -2:
	//			// 影响左转、右转、开火和停止
	//			if BulletMoveUrgent[ActionLeft] > bullet.Distance {
	//				BulletMoveUrgent[ActionLeft] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionRight] > bullet.Distance {
	//				BulletMoveUrgent[ActionRight] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionStay] > bullet.Distance {
	//				BulletMoveUrgent[ActionStay] = bullet.Distance
	//			}
	//		case -3:
	//			// 影响直行、后退、开火和停止
	//			if BulletMoveUrgent[ActionMove] > bullet.Distance {
	//				BulletMoveUrgent[ActionMove] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionBack] > bullet.Distance {
	//				BulletMoveUrgent[ActionBack] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionStay] > bullet.Distance {
	//				BulletMoveUrgent[ActionStay] = bullet.Distance
	//			}
	//		case -4:
	//			// 影响左转、右转、开火和停止
	//			if BulletMoveUrgent[ActionLeft] > bullet.Distance {
	//				BulletMoveUrgent[ActionLeft] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionRight] > bullet.Distance {
	//				BulletMoveUrgent[ActionRight] = bullet.Distance
	//			}
	//			if BulletMoveUrgent[ActionStay] > bullet.Distance {
	//				BulletMoveUrgent[ActionStay] = bullet.Distance
	//			}
	//		}
	//	}
	//}
    //
    //
	//// 躲避为第二优先级
	////EnemyMoveUrgent := make(map[int]int)
	////EnemyMoveUrgent[ActionMove] = math.MaxInt32
	////EnemyMoveUrgent[ActionBack] = math.MaxInt32
	////EnemyMoveUrgent[ActionLeft] = math.MaxInt32
	////EnemyMoveUrgent[ActionRight] = math.MaxInt32
	////
	////EnemyMoveUrgent[ActionStay] = math.MaxInt32
	////EnemyMoveUrgent[ActionNone] = math.MaxInt32
    //
	//if enemyApproach == true {
	//	// 不能走敌军的象限
	//	for _, enemy := range *enemyThreat {
	//		switch quadrant[enemy.Quadrant] {
	//		case 1:
	//			// 影响直行和右转
	//			if BulletMoveUrgent[ActionMove] > enemy.Distance {
	//				BulletMoveUrgent[ActionMove] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionRight] > enemy.Distance {
	//				BulletMoveUrgent[ActionRight] = enemy.Distance
	//			}
	//		case 2:
	//			// 影响直行和左转
	//			if BulletMoveUrgent[ActionMove] > enemy.Distance {
	//				BulletMoveUrgent[ActionMove] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionLeft] > enemy.Distance {
	//				BulletMoveUrgent[ActionLeft] = enemy.Distance
	//			}
	//		case 3:
	//			// 影响左转和后退
	//			if BulletMoveUrgent[ActionLeft] > enemy.Distance {
	//				BulletMoveUrgent[ActionLeft] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionBack] > enemy.Distance {
	//				BulletMoveUrgent[ActionBack] = enemy.Distance
	//			}
	//		case 4:
	//			// 影响右转和后退
	//			if BulletMoveUrgent[ActionRight] > enemy.Distance {
	//				BulletMoveUrgent[ActionRight] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionBack] > enemy.Distance {
	//				BulletMoveUrgent[ActionBack] = enemy.Distance
	//			}
	//		case -1:
	//			// 枪口对准了，影响直行、后退、停止
	//			if BulletMoveUrgent[ActionMove] > enemy.Distance {
	//				BulletMoveUrgent[ActionMove] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionBack] > enemy.Distance {
	//				BulletMoveUrgent[ActionBack] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionStay] > enemy.Distance {
	//				BulletMoveUrgent[ActionStay] = enemy.Distance
	//			}
	//		case -2:
	//			// 影响左转、右转、停止
	//			if BulletMoveUrgent[ActionLeft] > enemy.Distance {
	//				BulletMoveUrgent[ActionLeft] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionRight] > enemy.Distance {
	//				BulletMoveUrgent[ActionRight] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionStay] > enemy.Distance {
	//				BulletMoveUrgent[ActionStay] = enemy.Distance
	//			}
	//		case -3:
	//			// 影响直行、后退、停止
	//			if BulletMoveUrgent[ActionMove] > enemy.Distance {
	//				BulletMoveUrgent[ActionMove] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionBack] > enemy.Distance {
	//				BulletMoveUrgent[ActionBack] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionStay] > enemy.Distance {
	//				BulletMoveUrgent[ActionStay] = enemy.Distance
	//			}
	//		case -4:
	//			// 影响左转、右转、停止
	//			if BulletMoveUrgent[ActionLeft] > enemy.Distance {
	//				BulletMoveUrgent[ActionLeft] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionRight] > enemy.Distance {
	//				BulletMoveUrgent[ActionRight] = enemy.Distance
	//			}
	//			if BulletMoveUrgent[ActionStay] > enemy.Distance {
	//				BulletMoveUrgent[ActionStay] = enemy.Distance
	//			}
	//		}
	//	}
	//}
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
//}


