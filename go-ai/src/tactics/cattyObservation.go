// 观察者
package tactics

import (
	f "framework"
    // "fmt"
)

type Observation struct {
    TotalSteps   int
    Steps        int
    PreState     *f.GameState
    State        *f.GameState
    Radar        *f.RadarResult
    Objs         map[string]f.Objective
    MyTank       map[string]f.Tank // 我方存活坦克
    EmyTank      map[string]f.Tank // 敌方存活坦克
    Flag         Flag
	SNook, FNook   []f.Position  // sniper, flagman的苟点（提前分析地图）
	Fcnt, Scnt   int       // 各类角色分配数量
}

type Flag struct {
    Pos            f.Position
    Exist          bool      // 是否有旗定时刷新
    Next           int       // 距离刷新的回合数
}

func NewObservation(state *f.GameState) (obs *Observation) {
    obs = &Observation{ TotalSteps: state.Params.MaxRound, Steps: 0, State: state}

    // 观察坦克
    obs.observeTank()

	// 观察苟点
	obs.observePos()

	// 观察战旗
	obs.observeFlag()

    // 分配角色
    // obs.assignRole()

	return obs
}

func (o *Observation) makeObservation(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    o.Steps  += 1
    o.State  = state
    o.Radar  = radar
    o.Objs   = objective

    // 观察苟点
    o.observePos()

    // 观察战旗
    o.observeFlag()

    // 分配角色
    // o.assignRole()
}

func (o *Observation) observeTank() {
    o.MyTank = make(map[string]f.Tank)
    for _, tank := range o.State.MyTank {
        o.MyTank[tank.Id] = tank
    }
    o.EmyTank = make(map[string]f.Tank)
    for _, tank := range o.State.EnemyTank {
        o.EmyTank[tank.Id] = tank
    }
}

func (o *Observation) observePos() {
    o.SNook = []f.Position {
        f.Position { X: 5, Y:6, Direction: f.ActionFireLeft },
        f.Position { X: 5, Y:7, Direction: f.ActionFireLeft },
        f.Position { X: 3, Y:4, Direction: f.ActionFireUp },
        f.Position { X: 2, Y:1, Direction: f.ActionFireUp },
    }
    o.FNook = []f.Position {
        f.Position { X: 9, Y:9, Direction: f.DirectionDown },
    }
}

func (o *Observation) observeFlag() {
    if o.Steps < o.TotalSteps / 2 || o.State.FlagWait == 999999 {
        o.Flag = Flag { Pos: f.Position{ X: o.State.Params.FlagX, Y:o.State.Params.FlagY }, Exist: false, Next: 0}
    } else {
        o.Flag = Flag { Pos: f.Position{ X: o.State.Params.FlagX, Y:o.State.Params.FlagY }, Exist: true, Next: o.State.FlagWait}
    }
}

// 分配角色
// func (o *Observation) assignRole() {
//     if o.Flag.Exist {
//         o.Fcnt = 1
//     } else {
//         o.Fcnt = 0
//     }
//     o.Scnt = len(o.State.MyTank) - o.Fcnt
// }
