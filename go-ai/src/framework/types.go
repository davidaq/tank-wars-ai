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
	ActionTravel = iota    // 仅用作策略的返回，不可作为最终行动类型
)

// 策略输出，单个坦克下一步行动或者移动目标
type Objective struct {
	Action int					// 策略决定直接执行的操作
	Target Position			// 只有Action为ActionTravel的时候才生效
}

type SuggestionItem struct {
	Action int				// 建议接下来坦克的行为
	Urgent int				// 越小越紧急，寻路代表还有多远，
										// 躲避代表如果不被采纳还有多少步（子弹）就被命中，
										// 火力系统代表射击目标离自己多远（子弹距离）
}

// 策略系统协议，必须实现计划、决定两种行为
type Tactics interface {
	Init(state *GameState)																			// 根据初始state，初始化
	Plan(state *GameState, objective *map[string]Objective)			// 根据state，填充objective，设定每个坦克的战略目的地
}
