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
	ActionFire = iota
)

// 策略输出，单个坦克的战略目标，作为参考传输给行动系统
type Objective struct {
	Target Position
}

// 3个行动系统给出的行动建议
type SuggestionItem struct {
	Action int				// 建议接下来坦克的行为
	Urgent float32		// 该建议的迫切程度：[0-1]
}
type Suggestion struct {
	dodge 	SuggestionItem		// 躲避系统建议
	attack	SuggestionItem		// 火力系统建议
	travel	SuggestionItem		// 寻路系统建议
}

// 行动子系统协议
type reactor interface {
	Init(state *GameState, tankid string)
	Suggest(state *GameState, objective *Objective) SuggestionItem
}

// 策略系统协议，必须实现计划、决定两种行为
type Tactics interface {
	Init(state *GameState)																					// 根据初始state，初始化
	Plan(state *GameState, objective *map[string]Objective)					// 根据state，填充objective，设定每个坦克的战略目的地
	Decide(tankid string, suggestion Suggestion) int								// 3个行动系统得出建议后，最终决定采取哪个行动
}
