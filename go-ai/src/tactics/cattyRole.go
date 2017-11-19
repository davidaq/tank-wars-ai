package tactics

import (
	f "framework"
	// "fmt"
)

type CattyRole struct {
    obs       *Observation
    Tank      f.Tank
    Target    CattyTarget
}

type CattyTarget struct {
    Pos     f.Position
    Action  int
    After   *CattyTarget
}

// 更新坦克信息
func (r *CattyRole) SetTank(tank f.Tank)  {
    r.Tank = tank
}

// 设置行动目标
func (r *CattyRole) SetTarget(target CattyTarget)  {
    r.Target = target
}

// 检查目标是否完成
func (r *CattyRole) checkArrive() bool {
    if r.Tank.Pos.X ==  r.Target.Pos.X && r.Tank.Pos.Y ==  r.Target.Pos.Y {
        return true
    } else {
        return false
    }
}

// 移动
func (r *CattyRole) Move() {
    r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionTravel, Target: r.Target.Pos }
}

// 行动
func (r *CattyRole) Act() {
    canfire   := true
    radarFire := r.obs.Radar.Fire[r.Tank.Id]
    rfs       := []*f.RadarFire{ radarFire.Up, radarFire.Down, radarFire.Left, radarFire.Right }
    for _, rf := range rfs {
        if rf == nil {
            continue
        }
        if rf.Action == r.Target.Action && rf.Sin >= 0.3 {
            canfire = false
        }
    }
    if canfire {
        r.obs.Objs[r.Tank.Id] = f.Objective { Action: r.Target.Action }   // 开火
    }
    // 设置下一步目标
    if r.Target.After != nil {
        r.SetTarget(*r.Target.After)
    }
}
