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

	// 观察战旗
	obs.observeFlag()

	return obs
}

func (o *Observation) makeObservation(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    o.Steps  += 1
    o.State  = state
    o.Radar  = radar
    o.Objs   = objective

    // 观察战旗
    o.observeFlag()
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

func (o *Observation) observeFlag() {
    if o.Steps < o.TotalSteps / 2 || o.State.FlagWait == 999999 {
        o.Flag = Flag { Pos: f.Position{ X: o.State.Params.FlagX, Y:o.State.Params.FlagY }, Exist: false, Next: 0 }
    } else {
        o.Flag = Flag{ Pos: f.Position{ X: o.State.Params.FlagX, Y:o.State.Params.FlagY }, Exist: true, Next: o.State.FlagWait }
    }
}
