// 子弹躲避行动子系统
package framework

import (
	"math"
)

type Dodger struct {
}

func NewDodger() *Dodger {
	inst := &Dodger {
	}
	return inst
}


func (self *Dodger) avoidBullet(tank *Tank, state *GameState) (avoid bool, urgent int) {
	// 注意射过来的炮弹可能多个
	// 没有子弹射过来的情况
	enemyBulletList := []Position{}
	for _, bullet := range state.EnemyBullet {
		if tank.Pos.X == bullet.Pos.X && int(math.Abs(float64(bullet.Pos.Y - tank.Pos.Y))) <= 7 || tank.Pos.Y == bullet.Pos.Y && int(math.Abs(float64(bullet.Pos.X - tank.Pos.X))) <= 7 {
			enemyBulletList = append(enemyBulletList, Position{
				X: bullet.Pos.X,
				Y: bullet.Pos.Y,
			})
		}
	}

	if 0 == len(enemyBulletList) {
		return false, 0
	}

	// 判断方向，然后准备跑
	

	// 从哪个方向
	return false, 1
}

func (self *Dodger) Suggest(tank *Tank, state *GameState, objective *Objective) SuggestionItem {
	// 如果采纳，计算几步被干掉
	//self.calcObjectiveUrgent(tank, state, objective)
	// 躲避炮弹，炮弹还有几步打到
	self.avoidBullet(tank, state)
	// 敌军的威胁，注意敌军到火线距离是可能炮弹的1/2

	ret := SuggestionItem {
		Action: ActionMove,
		Urgent: 1,
	}
	return ret
}
