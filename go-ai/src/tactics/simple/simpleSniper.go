package tactics
//
// import (
// 	f "framework"
//     // "fmt"
// )
//
// type SimpleSniper struct {
//     superior *Simple
//     Tanks    map[string]f.Tank
//     policy   *SimplePolicy
// 	obs      *Observation
// }
//
// // 初始化
// func (s *SimpleSniper) Init(tanks []f.Tank) {
//     s.Tanks = make(map[string]f.Tank)
//     if len(tanks) == 0 { return }
//     for _, tank := range tanks {
//         s.Tanks[tank.Id] = tank
//     }
// }
//
// func (s *SimpleSniper) Plan(state *f.GameState, objs map[string]f.Objective){
//     // 本次可支配的坦克
//     var tanks []f.Tank
//     for id, tank := range s.Tanks {
//         if objs[id] == (f.Objective{}) {
//             tanks = append(tanks, tank)
//         }
//     }
// 	if len(tanks) > 0 {
// 		s.HideAndFire(tanks, objs)
// 	}
// }
//
// func (s *SimpleSniper) HideAndFire(tanks []f.Tank, objs map[string]f.Objective) {
//     // 分派狙击手到苟点
//     ftanks := make(map[string]f.Tank)
//     if len(s.obs.Kps) < len(tanks) {
//         ftanks = s.policy.Dispatch(tanks[0:len(s.obs.Kps)],  s.obs.Kps, objs)
//         for _, tank := range tanks[len(s.obs.Kps):]{
//             ftanks[tank.Id] = tank
//         }
//     } else {
//         ftanks = s.policy.Dispatch(tanks, s.obs.Kps, objs)
//     }
//     // 已抵达目的地坦克，根据友伤选择是否开火
//     if len(ftanks) > 0 {
//         s.policy.FireToFlag(ftanks, s.obs, objs)
//     }
// }
