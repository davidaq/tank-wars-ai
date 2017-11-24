package tactics

import (
	f "framework"
	"math"
	// "fmt"
)

type CattyRole struct {
    obs       *Observation
    Tank      f.Tank
    Target    CattyTarget
    Dodge     f.RadarDodge     // 躲避建议
    Fire      f.RadarFireAll   // 开火建议
}

type CattyTarget struct {
	Tank    f.Tank        // 目标坦克
    Pos     f.Position    // 目标地点
}

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
	dist := -1
	var tpos f.Position
	var ttank f.Tank
	for pos, tankid := range r.obs.ShotPos {
		nd := r.Tank.Pos.SDist(pos)
		tank := r.obs.EmyTank[tankid]
		if dist < 0 {
			dist  = nd
			ttank = tank
			tpos  = pos
		} else if nd > r.obs.State.Params.BulletSpeed * 2 || (r.obs.pathReachable(pos, r.nextPos(pos))) {  // && r.canHunt(tank)) {
			if nd < dist {
				dist  = nd
				ttank = tank
				tpos  = pos
			}
		}
	}
	r.Target.Tank = ttank
	r.Target.Pos  = tpos
	delete(r.obs.ShotPos, tpos)
}

func (r *CattyRole) checkDone() bool {
    return r.Dodge.Threat == -1 || r.fireAction() != -1 || (r.Tank.Pos.X == r.Target.Pos.X && r.Tank.Pos.Y == r.Target.Pos.Y)
}

func (r *CattyRole) move() {
	if r.Tank.Bullet == "" {
		if r.Dodge.Threat == 1 {
			r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionTravelWithDodge, Target: r.Target.Pos }
		} else {
			r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionTravel, Target: r.Target.Pos }
		}
	} else {
		r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionTravel, Target: r.Dodge.SafePos }
	}
}

// 行动
func (r *CattyRole) act() {
	if r.Dodge.Threat == -1 {
		r.obs.Objs[r.Tank.Id] = f.Objective { Action: r.fireBeforeDying() }   // 光辉弹
	} else if r.fireAction() != -1 {
    	r.obs.Objs[r.Tank.Id] = f.Objective { Action: r.fireAction() }
	} else if r.Tank.Pos.X == r.Target.Pos.X && r.Tank.Pos.Y == r.Target.Pos.Y {
		r.hunt()
		r.move()
	}
}

// 光辉弹
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
        if rf == nil { continue }
        if mrf == nil || mrf.Faith - mrf.Sin < rf.Faith - rf.Sin {
            mrf = rf
        }
    }
	if mrf == nil || mrf.Faith < 0.4 || mrf.Sin >= 0.5 {
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

// 不要追击有援军的敌人
func (r *CattyRole) canHunt(tank f.Tank) bool {
	x := tank.Pos.X
	y := tank.Pos.Y
	for _, emytank := range r.obs.EmyTank {
		if int(math.Abs(float64(emytank.Pos.X - x))) < r.obs.State.Params.BulletSpeed && int(math.Abs(float64(emytank.Pos.Y - y))) < r.obs.State.Params.BulletSpeed {
			return false
		}
	}
	return true
}
