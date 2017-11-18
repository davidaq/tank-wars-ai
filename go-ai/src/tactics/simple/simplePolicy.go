// // 各类战术动作
package tactics
//
// import (
// 	f "framework"
//     // "fmt"
// )
//
// type SimplePolicy struct {
// }
//
// func NewSimplePolicy() *SimplePolicy {
//     return &SimplePolicy{ }
// }
//
// // 移动到某地点
// func (p *SimplePolicy) MoveTo(pos f.Position, tank f.Tank, objs map[string]f.Objective){
//     objs[tank.Id] = f.Objective{ Action: f.ActionTravel, Target: pos}
// }
//
// // 将一组坦克派到一组地点，并返回空闲的坦克
// func (p *SimplePolicy) Dispatch(tanks []f.Tank, pos []f.Position, objs map[string]f.Objective) (ftanks map[string]f.Tank){
//     ftanks = make(map[string]f.Tank)
//     match  := MatchPosTank(pos[0:len(tanks)], tanks)
// 	for _, tank := range tanks {
//         if tank.Pos.X == match[tank.Id].X && tank.Pos.Y == match[tank.Id].Y {
//             ftanks[tank.Id] = tank
//         } else {
//             objs[tank.Id] = f.Objective{ Action: f.ActionTravel, Target: match[tank.Id]}
//         }
// 	}
//     return ftanks
// }
//
// // FireToFlag(地图针对策略，有固定朝向的射击)
// func (p *SimplePolicy) FireToFlag(ftanks map[string]f.Tank, obs *Observation, objs map[string]f.Objective ) {
//     var radarFire f.RadarFireAll
//     for id, _ := range ftanks {
//         radarFire = obs.Radar.Fire[id]
//         // 无友伤则开火
//         if obs.Side == "red" {
//             if radarFire == (f.RadarFireAll{}) || radarFire.Down.Sin < 0.3 {
//                 objs[id] = f.Objective{ Action: f.ActionFireDown}
//             } else {
//                 objs[id] = f.Objective{ Action: f.ActionStay}
//             }
//         } else {
//             if radarFire == (f.RadarFireAll{}) || radarFire.Up.Sin < 0.3 {
//                 objs[id] = f.Objective{ Action: f.ActionFireUp}
//             } else {
//                 objs[id] = f.Objective{ Action: f.ActionStay}
//             }
//         }
//     }
// }
//
// // FireTowardsFlag(地图针对策略，需要计算朝向后射击)
// func (p *SimplePolicy) FireTowardsFlag(ftanks map[string]f.Tank, obs *Observation, objs map[string]f.Objective ) {
//     var radarFire f.RadarFireAll
//     for id, tank := range ftanks {
//         radarFire = obs.Radar.Fire[id]
//         if radarFire == (f.RadarFireAll{}) || radarFire.Up.Sin < 0.3 {
//             objs[id] = f.Objective{ Action: FireDirection(tank.Pos, obs.Flag.Pos) }
//         }
//     }
// }
//
// // FreeFire
// func (p *SimplePolicy) FreeFire(ftanks map[string]f.Tank, obs *Observation, objs map[string]f.Objective ) {
//     var radarFire f.RadarFireAll
//     var rfs       []*f.RadarFire
//     var mrf       *f.RadarFire
//     for id, _ := range ftanks {
//         radarFire = obs.Radar.Fire[id]
//         rfs       = []*f.RadarFire{ radarFire.Up, radarFire.Down, radarFire.Left, radarFire.Right }
//         for _, rf := range rfs {
//             if rf == nil { continue }
//             if mrf == nil || mrf.Faith - mrf.Sin < rf.Faith - rf.Sin {
//                 mrf = rf
//             }
//         }
//         if mrf == nil { continue; }
//         if mrf.Faith > 0.5 && mrf.Sin <= 0.3 {
//             objs[id] = f.Objective{ Action: mrf.Action }
//         }
//     }
// }
