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
	for _, tank := range state.MyTank {
		c.Roles[tank.Id] = &CattyRole { obs: c.obs, Tank: tank}
	}
}

func (c *Catty) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    for tankid := range objective {
        delete(objective, tankid)
    }

    c.obs.makeObservation(state, radar, objective)

    c.updateRole()

    for _, role := range c.Roles {
        if c.obs.Flag.Exist && c.obs.Flag.Next <= 5 {
            role.occupyFlag()
            continue
        }

        role.hunt()
        role.act()

        fmt.Printf("catty role target: %+v\n", role.Target)
    }
    fmt.Printf("catty objective: %+v\n", c.obs.Objs)
}

func (c *Catty) updateRole() {
	for id, role := range c.Roles {
		if c.obs.MyTank[id] != (f.Tank{}) {
			role.Tank         = c.obs.MyTank[id]
			role.Dodge        = c.obs.Radar.DodgeBullet[id]
            role.ExtDangerSrc = c.obs.Radar.ExtDangerSrc[id]
            role.Fire         = c.obs.Radar.Fire[id]
		} else {
			delete(c.Roles, id)
		}
	}
}

func (c *Catty) End(state *f.GameState) { }
