// 开火攻击行动子系统
package framework

import (
	"math"
)

type Attacker struct {
}

func NewAttacker() *Attacker {
	inst := &Attacker {
	}
	return inst
}

/**
 * 计算火线位置
 */
func (self *Attacker) calcFireLine(tank *Tank, state *GameState, fireline *[]Position) {

	switch tank.Pos.Direction {
		case DirectionUp:
			for y := tank.Pos.Y - 1; y >= 0; y-- {
				if 1 == state.Terain.Get(tank.Pos.X, y) {
					break
				}
				*fireline = append(*fireline, Position{X: tank.Pos.X, Y: y, Direction: tank.Pos.Direction})
			}
		case DirectionDown:
			for y := tank.Pos.Y + 1; y < state.Terain.Height; y++ {
				if 1 == state.Terain.Get(tank.Pos.X, y) {
					break
				}
				*fireline = append(*fireline, Position{X: tank.Pos.X, Y: y, Direction: tank.Pos.Direction})
			}
		case DirectionLeft:
			for x := tank.Pos.X - 1; x >= 0; x-- {
				if 1 == state.Terain.Get(x, tank.Pos.Y) {
					break
				}
				*fireline = append(*fireline, Position{X: x, Y: tank.Pos.Y, Direction: tank.Pos.Direction})
			}
		case DirectionRight:
			for x := tank.Pos.X + 1; x < state.Terain.Width; x++ {
				if 1 == state.Terain.Get(x, tank.Pos.Y) {
					break
				}
				*fireline = append(*fireline, Position{X: x, Y: tank.Pos.Y, Direction: tank.Pos.Direction})
			}
	}
}

func (self *Attacker) firelineEmeny(tank *Tank, state *GameState, fireline *[]Position) (fire bool, urgent int){
	// 初始化
	fire = false
	urgent = 0

	// 如果火线上，障碍物以前有地方坦克，直接射击
	for _, enemytank := range state.EnemyTank{
		if (enemytank.Pos.X == tank.Pos.X || enemytank.Pos.Y == tank.Pos.Y) {
			// 如果根本不在一条直线上，无需考虑
			if 0 == len(*fireline) {
				self.calcFireLine(tank, state, fireline)
			}
			// 如果敌方坦克正好在fireline上则发射 并计算还有多少步
			for _, firelinev := range (*fireline) {
				// 这个地方之后再优化
				if (firelinev.X == enemytank.Pos.X && firelinev.Y == enemytank.Pos.Y) {
					return true, int(math.Abs(float64(enemytank.Pos.X - tank.Pos.X)) + math.Abs(float64(enemytank.Pos.Y - tank.Pos.Y)))
				}
			}
		}
	}

	return fire, urgent
}

func (self *Attacker) prefire(tank *Tank, state *GameState, fireline *[]Position) (fire bool, urgent int){
	// 初始化
	fire = false
	urgent = 0

	// 如果火线两侧的位置，有朝向火线的坦克，则射击
	for _, enemytank := range state.EnemyTank{
		switch tank.Pos.Direction {
		case DirectionUp:
		case DirectionDown:
			if (enemytank.Pos.Direction != DirectionLeft && enemytank.Pos.Direction != DirectionRight) {
				continue
			}
		case DirectionLeft:
		case DirectionRight:
			if (enemytank.Pos.Direction != DirectionDown && enemytank.Pos.Direction != DirectionUp) {
				continue
			}
		}

		// 预测6格即可，即6格火线
		for i := 0; i < 6; i++ {
			if len(*fireline) > i {
				// 监测火线两侧敌军
				if (*fireline)[i].X == enemytank.Pos.X || (*fireline)[i].Y == enemytank.Pos.Y {
					// 敌军到火线的距离
					tmpFEDis := int(math.Abs(float64(enemytank.Pos.X - (*fireline)[i].X)) + math.Abs(float64(enemytank.Pos.Y - (*fireline)[i].Y)))
					if tmpFEDis == int(math.Ceil(float64(tmpFEDis) / 2)) {
						return true, i
					}
				}
			}
		}
	}

	return fire, urgent
}

func (self *Attacker) Suggest(tank *Tank, state *GameState, objective *Objective) SuggestionItem {
	// 定义默认值
	fire	:= false
	action 	:= ActionNone
	urgent 	:= 0

	// 火线
	var fireline []Position

	// 当前火线方向是否有敌军
	fire, fireUrgent := self.firelineEmeny(tank, state, &fireline)
	if fire == true {
		action = ActionFire
		urgent = int(math.Ceil(float64(fireUrgent / 2)))
	}

	// 当前火线方向附近是否有敌军
	if fire == false {
		fire, fireUrgent := self.prefire(tank, state, &fireline)
		if fire == true {
			action = ActionFire
			urgent = int(math.Ceil(float64(fireUrgent) / 2))
		}
	}

	ret := SuggestionItem {
		Action: action,
		Urgent: urgent,
	}
	return ret
}
