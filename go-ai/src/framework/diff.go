// 检查当前state和之前state，判断草丛威胁
package framework

type Diff struct {
	prevState *GameState
}

func NewDiff() &Diff {
	return &Diff {
		prevState: nil
	}
}

func (self *Diff) compare(newState *GameState) DiffResult {
	ret := DiffResult {
		ForestThreat: make(map[Position]float64),
	}
	if self.prevState != nil {
		// TODO
		ret.ForestThreat[Position { X: 0, Y: 0 }] = 1.
	}
	self.prevState = newState
	return ret
}
