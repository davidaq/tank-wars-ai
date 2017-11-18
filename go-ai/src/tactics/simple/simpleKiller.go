package tactics
//
// import (
// 	f "framework"
//     // "fmt"
// )
//
// type SimpleKiller struct {
//     superior *Simple
//     Tanks    map[string]f.Tank
//     policy   *SimplePolicy
// 	obs      *Observation
// }
//
// // 初始化
// func (s *SimpleKiller) Init(tanks []f.Tank) {
//     s.Tanks = make(map[string]f.Tank)
//     if len(tanks) == 0 { return }
//     for _, tank := range tanks {
//         s.Tanks[tank.Id] = tank
//     }
// }
//
// // 制定执行计划
// func (s *SimpleKiller) Plan(state *f.GameState, objs map[string]f.Objective){
//     // 本次可支配的坦克
//     var tanks []f.Tank
//     for id, tank := range s.Tanks {
//         if objs[id] == (f.Objective{}) {
//             tanks = append(tanks, tank)
//         }
//     }
// 	if len(tanks) > 0 {
//         if len(s.obs.EmyTank) > 0 {
//             s.huntEmyTank(tanks, objs)
//         // 无可见敌方坦克
//         } else {
//             s.fireToForest(tanks, objs)
//         }
// 	}
// }
//
// // 追击敌方坦克
// func (s *SimpleKiller) huntEmyTank(tanks []f.Tank, objs map[string]f.Objective) {
//     // 适合射击的地点
//     var shootPos []f.Position
//     ftanks   := make(map[string]f.Tank)
//
//     // // 寻找适合追击的坦克
//     // targetTanks := FindTargetTanks(s.obs.EmyTank)
//
//     // 寻找适合射击的地点
//     for _, emytank := range s.obs.EmyTank {
//         shootPos = append(shootPos, FindShootPos(emytank.Pos, *s.obs.CurState.Terain, s.obs.CurState.Params.BulletSpeed)...)
//     }
//
//     // match 附近的坦克
//     if len(shootPos) < len(tanks) {
//         ftanks = s.policy.Dispatch(tanks[0:len(shootPos)], shootPos, objs)
//         for _, tank := range tanks[len(shootPos):]{
//             ftanks[tank.Id] = tank
//         }
//     } else {
//         ftanks = s.policy.Dispatch(tanks, shootPos, objs)
//     }
//
//     if len(ftanks) > 0 {
//         s.policy.FreeFire(ftanks, s.obs, objs)
//     }
// }
//
// // 向草丛内定点射击
// func (s *SimpleKiller) fireToForest(tanks []f.Tank, objs map[string]f.Objective) {
//     // 适合射击的地点
//     shootPos := s.obs.ShootPos
//     ftanks   := make(map[string]f.Tank)
//     // match 附近的坦克
//     if len(shootPos) < len(tanks) {
//         ftanks = s.policy.Dispatch(tanks[0:len(shootPos)], shootPos, objs)
//         for _, tank := range tanks[len(shootPos):]{
//             ftanks[tank.Id] = tank
//         }
//     } else {
//         ftanks = s.policy.Dispatch(tanks, shootPos, objs)
//     }
//     if len(ftanks) > 0 {
//         s.policy.FireTowardsFlag(ftanks, s.obs, objs)
//     }
// }
