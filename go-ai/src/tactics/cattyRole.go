package tactics

import (
	f "framework"
	"math"
	"fmt"
)

type CattyRole struct {
    obs       *Observation
    Tank      f.Tank
    Target    CattyTarget
    Dodge     f.RadarDodge     // 躲避建议
    Fire      f.RadarFireAll   // 开火建议
}

type CattyTarget struct {
    Pos     f.Position    // 目标地点
    Action  int           // 达成目标后执行动作
    Tank    f.Tank        // 目标坦克
    After   *CattyTarget  // 达成目标后的下一步动作
}

func (r *CattyRole) refreshTarget() {
    if r.Target.Tank != (f.Tank{}) {
		if r.Target.Tank.Bullet == "" {
			r.Target.Pos   = r.shotPos()
	        r.Target.After = &CattyTarget { Pos: r.safePos(), After: &CattyTarget{ Tank: r.Target.Tank } }
			fmt.Printf("=============== r.Target: %+v\n", r.Target)
		} else {
			r.Target.Pos   = r.Dodge.SafePos
		}
    }
}

func (r *CattyRole) occupyFlag() {
    travel := f.ActionTravel
    if r.Dodge.Threat > 0.9 {
        travel = f.ActionTravelWithDodge
    }
    r.obs.Objs[r.Tank.Id] = f.Objective {
        Action: travel,
        Target: r.obs.Flag.Pos,
    }
}

func (r *CattyRole) hunt() {
    var ttank f.Tank
	dist := -1
	for _, tank := range r.obs.EmyTank {
        if nd := r.Tank.Pos.SDist(tank.Pos); dist < 0 || nd < dist {
            dist  = nd
            ttank = tank
        }
	}
    r.Target.Tank = ttank
}

func (r *CattyRole) checkDone() bool {
	// fmt.Printf("r.Fire: %+v\n", r.Fire)
    return r.Dodge.Threat == -1 || r.fireAction() != -1 || (r.Tank.Pos.X == r.Target.Pos.X || r.Tank.Pos.X == r.Target.Pos.X)
}

func (r *CattyRole) move() {
    travel := f.ActionTravel
	// if r.Dodge.Threat >= 0.8 {
	// 	travel = f.ActionTravelWithDodge
	// }
	// if r.Target.Tank.Bullet != "" {
    //     travel = f.ActionTravelWithDodge
    // }
	r.obs.Objs[r.Tank.Id] = f.Objective { Action: travel, Target: r.Target.Pos }
}

// 行动
func (r *CattyRole) act() {
	if r.Dodge.Threat == -1 {
		r.obs.Objs[r.Tank.Id] = f.Objective { Action: r.fire() }
		if r.Target.After != nil { // && *r.Target.After != (CattyTarget{}) {
			r.Target = *r.Target.After
		}
	} else if r.fireAction() != -1 {
    	r.obs.Objs[r.Tank.Id] = f.Objective { Action: r.fireAction() }
		if r.Target.After != nil { // && *r.Target.After != (CattyTarget{}) {
			r.Target = *r.Target.After
		}
	} else if (r.Tank.Pos.X == r.Target.Pos.X || r.Tank.Pos.X == r.Target.Pos.X) {
		if r.Target.After != nil { // && *r.Target.After != (CattyTarget{}) {
			r.Target = *r.Target.After
			if r.Target.Pos == (f.Position{}) {
				r.refreshTarget()
			}
			r.move()
		}
	}
}

// 根据友伤判断能否开火
// func (r *CattyRole) canTakeAction() bool {
//     canTakeAction  := true
//     radarFire := r.obs.Radar.Fire[r.Tank.Id]
//     rfs       := []*f.RadarFire{ radarFire.Up, radarFire.Down, radarFire.Left, radarFire.Right }
//     for _, rf := range rfs {
//         if rf == nil {
//             continue
//         }
//         // 有自伤的情况下，取消行动
//         if rf.Action == r.Target.Action && rf.Sin >= 0.3 {
//             canTakeAction = false
//         }
//     }
//     return canTakeAction
// }

// 根据地形判断能否开火（子弹是否可达，将己方坦克位置视为障碍）
// func (r *CattyRole) canFireByTerain(startpos f.Position, endpos f.Position) bool {
// 	 terain := r.obs.Terain.Data
//      var min, max int
// 	 // XY路径上若有障碍物则不可达
//      if startpos.X == endpos.X {
//          min = int(math.Min(float64(startpos.Y), float64(endpos.Y)))
//          max = int(math.Max(float64(startpos.Y), float64(endpos.Y)))
//          for i := min+1; i < max; i++ {
//              if terain[i][startpos.X] == 1 || terain[i][startpos.X] == 4 {
//                  return false
//              }
//          }
//      } else if startpos.Y == endpos.Y {
//          min = int(math.Min(float64(startpos.X), float64(endpos.X)))
//          max = int(math.Max(float64(startpos.X), float64(endpos.X)))
//          for i := min+1; i < max; i++ {
//              if terain[startpos.Y][i] == 1 || terain[startpos.Y][i] == 4 {
//                  return false
//              }
//          }
//      } else {
// 		 return false
// 	 }
// 	 return true
// }

//
func (r *CattyRole) fire() int {
    var mrf *f.RadarFire
    for _, rf := range []*f.RadarFire{ r.Fire.Up, r.Fire.Down, r.Fire.Left, r.Fire.Right } {
        if rf == nil { continue }
        if mrf == nil || mrf.Faith - mrf.Sin < rf.Faith - rf.Sin {
            mrf = rf
        }
    }
	if mrf == nil || mrf.Faith == 0 || mrf.Sin >= 0.5 {
		return 0
	} else {
		return mrf.Action
	}
}

// 可以开火的方向
func (r *CattyRole) fireAction() int {
    var mrf *f.RadarFire
    for _, rf := range []*f.RadarFire{ r.Fire.Up, r.Fire.Down, r.Fire.Left, r.Fire.Right } {
        if rf == nil { continue }
		// fmt.Printf("rf: %+v\n", rf)
        if mrf == nil || mrf.Faith - mrf.Sin < rf.Faith - rf.Sin {
            mrf = rf
        }
    }
	if mrf == nil || mrf.Faith < 0.5 || mrf.Sin > 0.3 {
		return -1
	} else {
		return mrf.Action
	}
}

// 只考虑前进方向
func (r *CattyRole) safePos() f.Position {
	if r.Tank.Pos.Direction == f.DirectionUp {
		return f.Position { X: r.Tank.Pos.X, Y: r.Tank.Pos.Y + r.obs.State.Params.TankSpeed	 }
	} else if r.Tank.Pos.Direction == f.DirectionDown {
		return f.Position { X: r.Tank.Pos.X, Y: r.Tank.Pos.Y - r.obs.State.Params.TankSpeed}
	} else if r.Tank.Pos.Direction == f.DirectionRight {
		return f.Position { X: r.Tank.Pos.X - r.obs.State.Params.TankSpeed, Y: r.Tank.Pos.Y }
	} else {
		return f.Position { X: r.Tank.Pos.X + r.obs.State.Params.TankSpeed, Y: r.Tank.Pos.Y }
	}
}

func (r *CattyRole) shotPos() f.Position {
    var positions []f.Position
    for i := r.obs.State.Params.BulletSpeed; i > 0; i-- {
        if r.Target.Tank.Pos.Direction == f.DirectionUp || r.Target.Tank.Pos.Direction == f.DirectionDown {
            positions = append(positions, f.Position { X: r.Target.Tank.Pos.X, Y: r.Target.Tank.Pos.Y + r.obs.State.Params.BulletSpeed + i + 1})
            positions = append(positions, f.Position { X: r.Target.Tank.Pos.X, Y: r.Target.Tank.Pos.Y - r.obs.State.Params.BulletSpeed - i - 1})
        } else {
            positions = append(positions, f.Position { X: r.Target.Tank.Pos.X + r.obs.State.Params.BulletSpeed+ i + 1, Y: r.Target.Tank.Pos.Y})
            positions = append(positions, f.Position { X: r.Target.Tank.Pos.X - r.obs.State.Params.BulletSpeed - i - 1, Y: r.Target.Tank.Pos.Y})
        }
    }
	fmt.Println("r.Tank.id: ", r.Tank.Id)
	fmt.Printf("positions: %+v\n", positions)

	dist := -1
	var ret f.Position
	for _, pos := range positions {
		if r.reachable(pos) { //&& r.bulletReachable(r.Tank.Pos, pos) {
			if nd := r.Tank.Pos.SDist(pos); dist < 0 || nd < dist {
				dist = nd
				ret  = pos
			}
		}
	}
	fmt.Printf("ret: %+v\n", ret)
	return ret
}

// 地点是否可达（是否超出地图范围、是否墙壁）
func (r *CattyRole) reachable(p f.Position) bool {
    // 超出地图范围
    if p.X < 0 || p.X >= r.obs.Terain.Width || p.Y < 0 || p.Y >= r.obs.Terain.Height {
        return false
    }
    // 是否墙壁
    if r.obs.Terain.Data[p.Y][p.X] == 1 {
        return false
    }
    return true
}

// 子弹是否可达（
func (r *CattyRole) bulletReachable(startpos f.Position, endpos f.Position) bool {
    var min, max int
    if startpos.X == endpos.X {
        min = int(math.Min(float64(startpos.Y), float64(endpos.Y)))
        max = int(math.Max(float64(startpos.Y), float64(endpos.Y)))
        for i := min+1; i < max; i++ {
            if r.obs.Terain.Data[i][startpos.X] == 1 {
                return false
            }
        }
    } else if startpos.Y == endpos.Y {
        min = int(math.Min(float64(startpos.X), float64(endpos.X)))
        max = int(math.Max(float64(startpos.X), float64(endpos.X)))
        for i := min+1; i < max; i++ {
            if r.obs.Terain.Data[startpos.Y][i] == 1 {
                return false
            }
        }
    } else {
        return false
    }
    return true
}
