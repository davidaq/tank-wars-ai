package tactics

import (
	f "framework"
	// "math"
	"fmt"
)

type CattyRole struct {
    obs       *Observation
    Tank      f.Tank
    Target    f.Position
    Dodge     f.RadarDodge     // 躲避建议
    Fire      f.RadarFireAll   // 开火建议
}

// type CattyTarget struct {
//     Pos     f.Position    // 目标地点
// }

func (r *CattyRole) occupyFlag() {
    travel := f.ActionTravel
    if r.Dodge.Threat == 1 {
        travel = f.ActionTravelWithDodge
    }
    r.obs.Objs[r.Tank.Id] = f.Objective {
        Action: travel,
        Target: r.obs.Flag.Pos,
    }
}

// 距离最近的可攻击位置
// 如果敌方坦克密度较大，放弃那个位置
func (r *CattyRole) hunt() {
    fmt.Println("--------in hunt--------")
    fmt.Printf("r.obs.ShotPos: %+v\n", r.obs.ShotPos)

    // 如果没有绝杀点
    if len(r.obs.ShotPos) == 0 {
        ttank := r.neareastEmy()
        // 距离最近的坦克很远
        if nd := r.Tank.Pos.SDist(ttank.Pos); nd > r.obs.State.Params.BulletSpeed * 2 {
            r.Target = ttank.Pos
            return
        }
    }

    // 如果有绝杀点
	dist := -1
	var tpos f.Position
	for pos, _ := range r.obs.ShotPos {
		nd   := r.Tank.Pos.SDist(pos)
        // 很接近目标，且逃跑路线不顺畅，不攻击
        if nd < r.obs.State.Params.BulletSpeed * 2 && !r.obs.pathReachable(pos, r.nextPos(pos)) {
            continue
        }
        if dist < 0 || nd < dist{
            dist  = nd
            tpos  = pos
        }
	}

    // tpos 可能为空
    if tpos != (f.Position{}) {
        r.Target = tpos
        delete(r.obs.ShotPos, tpos)

    } else {
        r.Target = r.Tank.Pos
    }
}

func (r *CattyRole) neareastEmy() f.Tank {
    dist := -1
    var ttank f.Tank
    for _, tank := range r.obs.EmyTank {
        if nd:= r.Tank.Pos.SDist(tank.Pos); dist < 0 || nd < dist {
            dist  = nd
            ttank = tank
        }
    }
    return ttank
}

func (r *CattyRole) move() {
    r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionTravelWithDodge, Target: r.Target }
}

// 行动
func (r *CattyRole) act() {
    // 必死
	if r.Dodge.Threat == -1 {
		r.obs.Objs[r.Tank.Id] = f.Objective { Action: r.fireBeforeDying() }

    // 可开火
	} else if r.fireAction() != -1 && r.Dodge.Threat < 1 {
		r.obs.Objs[r.Tank.Id] = f.Objective { Action: r.fireAction() }

    // 可朝旗开火
    } else if r.canFireToFlag() {
        r.fireFlag()

    // 其余情况寻路
	} else {
		r.move()
	}
}

// 光荣弹
func (r *CattyRole) fireBeforeDying() int {
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

// 选择最适合的开火方向，没有可开火方向返回- 1
func (r *CattyRole) fireAction() int {
    var mrf *f.RadarFire
    for _, rf := range []*f.RadarFire{ r.Fire.Up, r.Fire.Down, r.Fire.Left, r.Fire.Right } {
        if rf == nil || !rf.IsStraight { continue }
        if mrf == nil || mrf.Faith - mrf.Sin < rf.Faith - rf.Sin {
            mrf = rf
        }
    }
	if mrf == nil || mrf.Faith < 0.5 || mrf.Sin >= 0.5 {
		return -1
	} else {
		return mrf.Action
	}
}

// 下一步位置
func (r *CattyRole) nextPos(pos f.Position) f.Position {
	if r.Tank.Pos.Direction == f.DirectionUp {
		return f.Position { X: pos.X, Y: pos.Y + r.obs.State.Params.TankSpeed	 }
	} else if pos.Direction == f.DirectionDown {
		return f.Position { X: pos.X, Y: pos.Y - r.obs.State.Params.TankSpeed}
	} else if pos.Direction == f.DirectionRight {
		return f.Position { X: pos.X - r.obs.State.Params.TankSpeed, Y: pos.Y }
	} else {
		return f.Position { X: pos.X + r.obs.State.Params.TankSpeed, Y: pos.Y }
	}
}

// 是否可以朝旗子开火
func (r *CattyRole) canFireToFlag() bool {
    if r.Tank.Pos.X == r.obs.Flag.Pos.X || r.Tank.Pos.Y == r.obs.Flag.Pos.Y {
        // 自己在旗子中不开火
        if (r.Tank.Pos.X == r.obs.Flag.Pos.X && r.Tank.Pos.Y == r.obs.Flag.Pos.Y){
            return false
        }
        // 可以向旗子开火
        if r.Dodge.Threat == 0 && r.obs.pathReachable(r.Tank.Pos, r.obs.Flag.Pos) {
            // 判断友伤
            var rf *f.RadarFire
            if r.Tank.Pos.X == r.obs.Flag.Pos.X {
                if r.Tank.Pos.Y > r.obs.Flag.Pos.Y {
                     rf = r.Fire.Up
                } else {
                    rf = r.Fire.Down
                }
            } else {
                if r.Tank.Pos.X > r.obs.Flag.Pos.X {
                    rf = r.Fire.Left
                } else {
                    rf = r.Fire.Right
                }
            }
            if rf == nil || rf.Sin <= 0.0 {
                return true
            }
        }
    }
    return false
}

// 向旗子开火
func (r *CattyRole) fireFlag() {
    if r.Tank.Pos.X == r.obs.Flag.Pos.X {
        if r.Tank.Pos.Y > r.obs.Flag.Pos.Y {
            r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionFireUp }
        } else {
            r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionFireDown }
        }
    } else {
        if r.Tank.Pos.X > r.obs.Flag.Pos.X {
            r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionFireLeft }
        } else {
            r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionFireRight }
        }
    }
}
