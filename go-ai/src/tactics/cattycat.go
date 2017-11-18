package tactics

import (
	f "framework"
	"fmt"
)

type Catty struct {
    obs      *Observation
    Roles    map[string]*CattyRole
}

func NewCatty() *Catty{
    return &Catty { Roles: make(map[string]*CattyRole) }
}

func (c *Catty) Init(state *f.GameState) {
    c.obs     = NewObservation(state)

    var target f.Position
    if len(state.MyTank) > 0 {
        for _, tank := range state.MyTank {
            c.Roles[tank.Id] = &CattyRole { obs: c.obs}
            c.Roles[tank.Id].SetTank(tank)
            if tank.Pos.X == 1 && tank.Pos.Y == 2 {
                target = f.Position { X: 2, Y:1, Direction: f.ActionFireRight }
            } else if tank.Pos.X == 1 && tank.Pos.Y == 1 {
                target = f.Position { X: 3, Y:4, Direction: f.ActionFireRight }
            } else if tank.Pos.X == 2 && tank.Pos.Y == 1 {
                target = f.Position { X: 5, Y:6, Direction: f.ActionFireLeft }
            } else {
                target = f.Position { X: 5, Y:7, Direction: f.ActionFireLeft }
            }
            c.Roles[tank.Id].SetTarget(target)
        }
    }
}

// 执行计划
func (c *Catty) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    // 清空上一步的 objective
    for tankid := range objective {
        delete(objective, tankid)
    }

    // 观察局势
    c.obs.makeObservation(state, radar, objective)

    // 更新Tank信息
    c.updateRole()

    // 检查雷达
    c.checkRadar()

    for _, role := range c.freeRole() {
        if role.checkArrive() {
            role.Act()
        } else {
            role.Move()
        }
    }

    fmt.Printf("catty objective: %+v\n", c.obs.Objs)
}

func (c *Catty) updateRole() {
    for _, tank := range c.obs.State.MyTank {
        if c.Roles[tank.Id] != nil {
            c.Roles[tank.Id].SetTank(tank)
        }
    }
}

func (c *Catty) freeRole() (freeRole []*CattyRole) {
    for id, role := range c.Roles {
        if c.obs.Objs[id] == (f.Objective{}) {
            freeRole = append(freeRole, role)
        }
    }
    return freeRole
}

func (c *Catty) checkRadar() {

}

func (c *Catty) End(state *f.GameState) {

}
