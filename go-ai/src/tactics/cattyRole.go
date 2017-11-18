package tactics

import (
	f "framework"
	// "fmt"
)

type CattyRole struct {
    obs       *Observation
    Tank      f.Tank
    Objective f.Objective
}

// 更新坦克信息
func (r *CattyRole) SetTank(tank f.Tank)  {
    r.Tank = tank
}

// 设置行动目标
func (r *CattyRole) SetTarget(pos f.Position)  {
    r.Objective = f.Objective { Action: pos.Direction, Target: pos }
}

// 检查目标是否完成
func (r *CattyRole) checkArrive() bool {
    if r.Tank.Pos.X ==  r.Objective.Target.X && r.Tank.Pos.Y ==  r.Objective.Target.Y {
        return true
    } else {
        return false
    }
}

// 移动
func (r *CattyRole) Move() {
    r.obs.Objs[r.Tank.Id] = f.Objective { Action: f.ActionTravel, Target: r.Objective.Target }
}

// 行动
func (r *CattyRole) Act() {
    radarFire := r.obs.Radar.Fire[r.Tank.Id]
    rfs       := []*f.RadarFire{ radarFire.Up, radarFire.Down, radarFire.Left, radarFire.Right }
    for _, rf := range rfs {
        if rf.Action == r.Objective.Action {
            if rf.Sin < 0.3 {
                r.obs.Objs[r.Tank.Id] = f.Objective { Action: r.Objective.Action }   // 开火
            }
        }
    }
}
