// // 分析战场局势
package tactics
//
// import (
// 	f "framework"
//     "fmt"
// )
//
// type Observation struct {
//     Side         string
//     CurSteps     int
//     CurState     *f.GameState
//     PreState     *f.GameState
//     MyTank       map[string]f.Tank // 我方存活坦克
//     EmyTank      map[string]f.Tank // 敌方存活坦克
//     Radar        *f.RadarResult
//     HasFlag      bool
//     Flag         Flag
// 	Kps          []f.Position  // snipers 的苟点（提前分析地图）
//     FlagKps      []f.Position  // flagmen 的苟点（提前分析地图）
//     ShootPos     []f.Position  // killers 的射击目标点
// 	Fcnt, Scnt, Kcnt int       // 各类角色分配数量
// }
//
// type Flag struct {
//     Pos            f.Position
//     Exist          bool
//     // Next           int    // 距离刷新的回合数
//     Occupied       bool
//     MyTank         f.Tank    // 我方是否有坦克占领
//     EmyTank        f.Tank    // 敌方是否有坦克占领
// }
//
// func NewObservation(state *f.GameState) (obs *Observation) {
//     obs = &Observation{ CurState: state, CurSteps: 0}
//
//     // 观察坦克
//     obs.observeTank(state)
//
// 	// 观察苟点
// 	obs.observeKps(state)
//
// 	// 观察战旗
// 	obs.observeFlag(state)
//
//     // 分配角色
//     obs.assignRole(state)
//
// 	return obs
// }
//
// // 必须每回合都调用，因为记录 steps
// func (o *Observation) makeObservation(state *f.GameState, radar *f.RadarResult) {
//     // step
//     o.CurSteps += 1
//
//     o.CurState  = state
//     o.Radar     = radar
//
//     // 更新坦克状态
//     o.observeTank(state)
//
//     // 更新旗点状态
// 	o.observeFlag(state)
//
//     fmt.Printf("CurSteps: %+v\n", o.CurSteps)
// }
//
// func (o *Observation) observeTank(state *f.GameState) {
//     o.MyTank = make(map[string]f.Tank)
//     for _, tank := range state.MyTank {
//         o.MyTank[tank.Id] = tank
//     }
//     o.EmyTank = make(map[string]f.Tank)
//     for _, tank := range state.EnemyTank {
//         o.EmyTank[tank.Id] = tank
//     }
// }
//
// func (o *Observation) observeFlag(state *f.GameState) {
//     // TODO 这个数值靠谱么
//     if state.FlagWait > 50 {
//         o.HasFlag = false
//         o.Flag = Flag { Pos: f.Position{ X: state.Params.FlagX, Y:state.Params.FlagY }, Exist: false, Occupied: false }
//         return
//     }
// 	o.HasFlag = true
// 	o.Flag = Flag { Pos: f.Position{ X: state.Params.FlagX, Y:state.Params.FlagY }, Exist: state.FlagWait == 0, Occupied: false }
// 	for _, tank := range state.MyTank {
// 		if tank.Pos.X == o.Flag.Pos.X && tank.Pos.Y == o.Flag.Pos.Y {
// 			o.Flag.Occupied = true
// 			o.Flag.MyTank   = tank
// 		}
// 	}
// 	for _, tank := range state.EnemyTank {
// 		if tank.Pos.X == o.Flag.Pos.X && tank.Pos.Y == o.Flag.Pos.Y {
// 			o.Flag.Occupied = true
// 			o.Flag.EmyTank  = tank
// 		}
// 	}
// }
//
// // 观察苟点
// func (o *Observation) observeKps(state *f.GameState) {
//     o.FlagKps  = []f.Position{}
//     o.Kps      = []f.Position{}
//     // 判断红蓝方
//     if state.MyTank[0].Pos.X < state.Terain.Width / 2 {
//         o.Side    = "blue"
//         o.FlagKps = append(o.Kps, f.Position { X:9, Y:12})
//         o.Kps     = append(o.Kps, f.Position { X:10, Y:12})
//     } else {
//         o.Side    = "red"
//         o.FlagKps = append(o.Kps, f.Position { X:9, Y:6})
//         o.Kps     = append(o.Kps, f.Position { X:8, Y:6})
//     }
//     o.ShootPos = []f.Position{
//         f.Position { X: 5, Y: 8},
//         f.Position { X: 5, Y: 9},
//         f.Position { X: 13, Y: 10},
//         f.Position { X: 13, Y: 9},
//     }
// }
//
// // 分配角色
// func (o *Observation) assignRole(state *f.GameState) {
//     if o.HasFlag {
//         o.Fcnt = 1
//         o.Scnt = len(o.Kps)
//     } else {
//         o.Fcnt = 0
//         o.Scnt = 0
//     }
//     o.Kcnt = len(state.MyTank) - o.Scnt - o.Fcnt
// }
