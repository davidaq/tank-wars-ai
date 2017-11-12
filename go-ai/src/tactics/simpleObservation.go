// 分析战场局势
package tactics

import (
	f "framework"
)

type Observation struct {
    // TotalSteps   int
    // CurSteps     int
    // MySocre      int
    // EmyScore     int
    // PreState     GameState
    CurState     *f.GameState
    HasFlag      bool
    Flag         Flag
	Kps          []f.Position  // 位置 提前分析地图
	Fcnt, Scnt, Kcnt int     // 各类角色分配数量
	// Forests      []Forest    //
    // StartSteps, EndSteps   int
}

type Flag struct {
    Pos            f.Position
    Next           int
    Occupied       bool
    MyTank         f.Tank    // 我方是否有坦克占领
    EmyTank        f.Tank    // 敌方是否有坦克占领
}

func NewObservation(state *f.GameState) (obs *Observation) {
    obs = &Observation{ CurState: state}

	// 观察苟点
	obs.observeKps(state)


	// 观察战旗
	obs.observeFlag(state)


	// 分配角色，有旗就分配一个旗手
    if obs.HasFlag {
        obs.Fcnt = 1
    } else {
        obs.Fcnt = 0
    }
    obs.Scnt = int(len(state.MyTank) / 2)
    obs.Kcnt = len(state.MyTank) - obs.Scnt - obs.Fcnt
	return obs
}

// 必须每回合都调用，因为记录 steps
func (o *Observation) makeObservation(state *f.GameState) {
	o.observeFlag(state)
}

func (o *Observation) observeFlag(state *f.GameState) {
	if state.FlagPos == (f.Position{}) {
		o.HasFlag = false
		return
	}

	o.HasFlag = true
	o.Flag = Flag { Pos: state.FlagPos, Next: state.FlagWait, Occupied: false }
	for _, tank := range state.MyTank {
		if tank.Pos.X == o.Flag.Pos.X && tank.Pos.Y == o.Flag.Pos.Y {
			o.Flag.Occupied = true
			o.Flag.MyTank   = tank
		}
	}
	for _, tank := range state.EnemyTank {
		if tank.Pos.X == o.Flag.Pos.X && tank.Pos.Y == o.Flag.Pos.Y {
			o.Flag.Occupied = true
			o.Flag.EmyTank  = tank
		}
	}
}
// 观察苟点
func (o *Observation) observeKps(state *f.GameState) {
	o.Kps = []f.Position{}
	o.Kps = append(o.Kps, f.Position { X:8,   Y:6})
	o.Kps = append(o.Kps, f.Position { X:10,  Y:12})
}

// func (o *Observation) observeForest(state *f.GameState) {
// }

// func (o *Observation) caculateScore(state *f.GameState) {
// }
