package tactics

import (
	f "framework"
	"fmt"
)

type Catty struct {
    obs           *Observation
    Roles         map[string]*CattyRole
    
}

func NewCatty() *Catty{
    return &Catty {
        // mapanalysis:  &f.MapAnalysis{},
        Roles: make(map[string]*CattyRole),
    }
}

// 分析草丛人员分配
func (c *Catty) analysis() map[f.Forest]int {
    return make(map[f.Forest]int)
}

// 分配
func (c *Catty) dispatch() {
    forestinfo := c.forestinfo()

    // 分配一些去旗子，其它在草外
    for forest, cnt := range c.analysis() {
        curcnt := forestinfo[forest]
        for _, role := range c.Roles {
            if curcnt == cnt {
                break
            }
            if role.gotoforest == false {
                role.gotoforest = true
                role.forest     = forest
                role.Target     = forest.Center
                curcnt++
            }
        }
    }
}

func (c *Catty) forestinfo() map[f.Forest]int{
    forestinfo := make(map[f.Forest]int)
    for _, role := range c.Roles {
        if role.gotoforest {
            forestinfo[role.forest] += 1
        }
    }
    return forestinfo
}

func (c *Catty) Init(state *f.GameState) {

    // 初始化角色
    c.obs = NewObservation(state)
	for _, tank := range state.MyTank {
		c.Roles[tank.Id] = &CattyRole { obs: c.obs, Tank: tank}
	}

    c.dispatch()

    // 分配一个去旗子
    // for _, role := range c.Roles {
    //     role.gotoforest = true
    //     role.Target     = c.obs.Flag.Pos
    //     break
    // }
}

func (c *Catty) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    for tankid := range objective {
        delete(objective, tankid)
    }

    c.obs.makeObservation(state, radar, objective)

    c.updateRole()

    c.redispatch()

    for _, role := range c.Roles {
        if c.obs.Flag.Exist && c.obs.Flag.Next <= 5 {
            role.occupyFlag()
            continue
        }
        if !role.gotoforest {
            role.hunt()
            role.act()
        } else {
            role.patrol()
        }
        fmt.Printf("catty role target: %+v\n", role.Target)
    }
    fmt.Printf("catty objective: %+v\n", c.obs.Objs)
}

func (c *Catty) updateRole() {
    foreatcnt := 0
	for id, role := range c.Roles {
		if c.obs.MyTank[id] != (f.Tank{}) {
            if role.gotoforest {
                foreatcnt++
            }
			role.Tank         = c.obs.MyTank[id]
			role.Dodge        = c.obs.Radar.DodgeBullet[id]
            role.ExtDangerSrc = c.obs.Radar.ExtDangerSrc[id]
            role.Fire         = c.obs.Radar.Fire[id]
		} else {
			delete(c.Roles, id)
		}
	}
    if foreatcnt <= 0 {
        for _, role := range c.Roles {
            role.gotoforest = true
            role.Target     = c.obs.Flag.Pos
            break
        }
    }
}

func (c *Catty) End(state *f.GameState) { }
