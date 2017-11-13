// 分析战场局势
package tactics

import (
	f "framework"
    "fmt"
)

type Observation struct {
    // TotalSteps   int
    CurSteps     int
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
    Exist          bool
    // Next           int
    Occupied       bool
    MyTank         f.Tank    // 我方是否有坦克占领
    EmyTank        f.Tank    // 敌方是否有坦克占领
}

func NewObservation(state *f.GameState) (obs *Observation) {
    obs = &Observation{ CurState: state, CurSteps: 0}

	// 观察苟点
	obs.observeKps(state)

	// 观察战旗
	obs.observeFlag(state)

    // 分配角色
    obs.assignRole(state)

	return obs
}

// 必须每回合都调用，因为记录 steps
func (o *Observation) makeObservation(state *f.GameState) {
    o.CurSteps += 1
	o.observeFlag(state)
    fmt.Printf("CurSteps: %+v\n", o.CurSteps)
}

func (o *Observation) observeFlag(state *f.GameState) {
    // TODO 判断条件暂时不明确，暂时当做始终有旗
	// if state.FlagPos == (f.Position{}) {
	// 	o.HasFlag = false
	// 	return
	// }
    // 始终有旗
	o.HasFlag = true
	o.Flag = Flag { Pos: f.Position{ X: state.Params.FlagX, Y:state.Params.FlagY }, Exist: state.FlagWait == 0, Occupied: false }
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
    fmt.Printf("obs.Flag: %+v\n", o.Flag)
}

// 分配角色
func (o *Observation) assignRole(state *f.GameState) {
    // 分配角色，有旗就分配一个旗手
    if o.HasFlag {
        o.Fcnt = 2
    } else {
        o.Fcnt = 0
    }
    o.Scnt = 0
    // o.Scnt = int(len(state.MyTank) / 2)
    o.Kcnt = len(state.MyTank) - o.Scnt - o.Fcnt
}

// 观察苟点
func (o *Observation) observeKps(state *f.GameState) {
	o.Kps = []f.Position{}
	o.Kps = append(o.Kps, f.Position { X:8,   Y:6})
	// o.Kps = append(o.Kps, f.Position { X:10,  Y:12})
    o.Kps = append(o.Kps, f.Position { X:15,  Y:17})
}

// func (o *Observation) observeForest(state *f.GameState) {
// }

// func (o *Observation) caculateScore(state *f.GameState) {
// }
