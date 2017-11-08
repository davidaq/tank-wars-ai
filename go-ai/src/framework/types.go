package framework

const (
	DirectionNone = iota
	DirectionUp = iota
	DirectionLeft = iota
	DirectionDown = iota
	DirectionRight = iota
)

const (
	ActionNone = iota
	ActionStay = iota
	ActionMove = iota
	ActionLeft = iota
	ActionRight = iota
	ActionBack = iota
	ActionFireUp = iota
	ActionFireLeft = iota
	ActionFireDown = iota
	ActionFireRight = iota
	ActionTravel = iota    					// 仅用作策略的返回
	ActionTravelWithDodge = iota		 // 仅用作策略的返回
)


// 雷达输出
type RadarResult struct {
	Dodge map[string]RadarDodge
	Fire map[string]RadarFire
}

type RadarDodge struct {
	Threat float64		// 受威胁程度，0到1，1就是如果不采纳肯定会命中
	SafePos Position	// 建议躲避位置，可以直接设定为坦克当前位置表示原地不动（前进方向受威胁）
}

type RadarFire struct {
	Faith float64			// 命中信仰，0到1，1就是如果采纳肯定会命中
	Action int				// ActionFireUp, ActionFireLeft, ActionFireDown, ActionFireRight
	Cost int					// 如果没命中，需要多少回合才能恢复弹药
}

// 策略系统协议，必须实现计划、决定两种行为
type Tactics interface {
	Init(state *GameState)																													// 根据初始state，初始化
	Plan(state *GameState, radar *RadarResult, objective map[string]Objective)			// 填充objective，设定每个坦克的战略目的地
}

// 策略输出，单个坦克下一步行动或者移动目标
type Objective struct {
	Action int					// 策略决定直接执行的操作
	Target Position			// Action为ActionTravel或ActionTravelWithDodge时传入，都是战略目的地，只是ActionTravelWithDodge会遵守雷达躲避建议
}
