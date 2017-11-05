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


func (self *Dodger) avoidBullet(tank *Tank, state *GameState, availableDirection *[]int) (avoid bool, urgent int) {
	// 注意射过来的炮弹可能多个
	// 没有子弹射过来的情况
	// TODO:注意墙
	enemyBulletList := []Position{}
	for _, bullet := range state.EnemyBullet {
		if tank.Pos.X == bullet.Pos.X && int(math.Abs(float64(bullet.Pos.Y - tank.Pos.Y))) <= 7 || tank.Pos.Y == bullet.Pos.Y && int(math.Abs(float64(bullet.Pos.X - tank.Pos.X))) <= 7 {
			// 子弹方向判断
			bulletDirection := DirectionNone
			if tank.Pos.X == bullet.Pos.X {
				if tank.Pos.Y > bullet.Pos.Y && bullet.Pos.Direction == DirectionDown {
					bulletDirection = DirectionDown
				} else if tank.Pos.Y < bullet.Pos.Y && bullet.Pos.Direction == DirectionUp {
					bulletDirection = DirectionUp
				}
			}
			if tank.Pos.Y == bullet.Pos.Y {
				if tank.Pos.X > bullet.Pos.X && bullet.Pos.Direction == DirectionRight {
					bulletDirection = DirectionRight
				} else if tank.Pos.Y < bullet.Pos.Y && bullet.Pos.Direction == DirectionLeft {
					bulletDirection = DirectionLeft
				}
			}

			if bulletDirection == DirectionNone {
				continue
			}
			enemyBulletList = append(enemyBulletList, Position{
				X: bullet.Pos.X,
				Y: bullet.Pos.Y,
				Direction: bulletDirection,
			})
		}
	}

	if 0 == len(enemyBulletList) {
		return false, 0
	}

	// 判断方向，然后准备跑
	minUrgent := math.MaxInt64
	for _, dangerBullet := range enemyBulletList {
		// 把子弹射过来的方向去掉
		for dKey, dData := range *availableDirection {
			if dData == dangerBullet.Direction {
				*availableDirection = append((*availableDirection)[:dKey], (*availableDirection)[dKey+1:]...)
				bulletDistance := int(math.Abs(float64(dangerBullet.X - tank.Pos.X)) + math.Abs(float64(dangerBullet.Y - tank.Pos.Y)))
				if minUrgent > bulletDistance {
					minUrgent = bulletDistance
				}
			}
		}
	}
	return true, minUrgent
}

func (self *Dodger) threat(tank *Tank, state *GameState, availableDirection *[]int) (avoid bool, urgent int) {
	return false, 0
}

func (self *Dodger) determine(tank *Tank, state *GameState, availableDirection *[]int) (action int) {
	// 注意方向撞墙情况，如果没有则确认方向
	direction := tank.Pos.Direction
	for i := len(*availableDirection) - 1; i >= 0; i-- {
		switch (*availableDirection)[i] {
		case DirectionRight:
			// 后期再改为越开阔越好
			if 0 == state.Terain.Get(tank.Pos.X + 1, tank.Pos.Y) {
				direction = DirectionRight
			}
		case DirectionLeft:
			if 0 == state.Terain.Get(tank.Pos.X - 1, tank.Pos.Y) {
				direction = DirectionLeft
			}
		case DirectionUp:
			if 0 == state.Terain.Get(tank.Pos.X, tank.Pos.Y - 1) {
				direction = DirectionUp
			}

		case DirectionDown:
			if 0 == state.Terain.Get(tank.Pos.X, tank.Pos.Y + 1) {
				direction = DirectionDown
			}
		}
	}

	// 如果实在没地方躲，例如死胡同则强行冲锋
	if direction == DirectionUp{
		switch tank.Pos.Direction {
		case DirectionUp:
			return ActionNone
		case DirectionLeft:
			return ActionRight
		case DirectionRight:
			return ActionLeft
		case DirectionDown:
			return ActionRight
		}
	}

	if direction == DirectionRight {
		switch tank.Pos.Direction {
		case DirectionUp:
			return ActionRight
		case DirectionRight:
			return ActionMove
		case DirectionDown:
			return ActionLeft
		case DirectionLeft:
			return ActionRight
		}
	}

	if direction == DirectionDown {
		switch tank.Pos.Direction {
		case DirectionUp:
			return ActionRight
		case DirectionRight:
			return ActionRight
		case DirectionDown:
			return ActionMove
		case DirectionLeft:
			return ActionLeft
		}
	}

	if direction == DirectionLeft {
		switch tank.Pos.Direction {
		case DirectionUp:
			return ActionLeft
		case DirectionRight:
			return ActionRight
		case DirectionDown:
			return ActionRight
		case DirectionLeft:
			return ActionMove
		}
	}

	return tank.Pos.Direction
}

func (self *Dodger) Suggest(tank *Tank, state *GameState, objective *Objective) SuggestionItem {
	// 可用的方向，用于剔除
	availableDirection := []int{
		DirectionNone,
		DirectionUp,
		DirectionDown,
		DirectionLeft,
		DirectionRight,
	}

	// 如果采纳，计算几步被干掉
	//self.calcObjectiveUrgent(tank, state, objective)
	// 躲避炮弹，炮弹还有几步打到
	avoid, minUrgent := self.avoidBullet(tank, state, &availableDirection)

	// 敌军的威胁，注意敌军到火线距离是可能炮弹的1/2
	self.threat(tank, state, &availableDirection)

	// 最终决定往哪里躲
	if avoid == true {
		action := self.determine(tank, state, &availableDirection)
		return SuggestionItem {
			Action: action,
			Urgent: minUrgent,
		}
	}

	ret := SuggestionItem {
		Action: ActionNone,
		Urgent: 0,
	}
	return ret
}
