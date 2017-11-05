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
	Urgent int				// 越小越紧急，寻路代表还有多远，
							// 躲避代表如果不被采纳还有多少步（子弹）就被命中，
							// 火力系统代表射击目标离自己多远（子弹距离）
}
type Suggestion struct {
	Dodge 	SuggestionItem		// 躲避系统建议
	Attack	SuggestionItem		// 火力系统建议
	Travel	SuggestionItem		// 寻路系统建议
}

// 行动子系统协议
type reactor interface {
	Suggest(tank *Tank, state *GameState, objective *Objective) SuggestionItem
}

// 策略系统协议，必须实现计划、决定两种行为
type Tactics interface {
	Init(state *GameState)													// 根据初始state，初始化
	Plan(state *GameState, objective *map[string]Objective)					// 根据state，填充objective，设定每个坦克的战略目的地
	Decide(tank *Tank, state *GameState, suggestion Suggestion) int			// 3个行动系统得出建议后，最终决定采取哪个行动
}
