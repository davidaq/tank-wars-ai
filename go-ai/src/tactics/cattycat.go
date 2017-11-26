package tactics

import (
	f "framework"
	"fmt"
)

type Catty struct {
    obs           *Observation
    Roles         map[string]*CattyRole
    forestmap     map[int]int
}

func NewCatty() *Catty{
    return &Catty {
        Roles: make(map[string]*CattyRole),
    }
}

// 分析草丛人员分配
func (c *Catty) analysis() {

}

// 分配
func (c *Catty) dispatch() {
    // 分配一些去旗子，其它在草外
    for forestid, cnt := range c.forestmap {
        i := 0
        for _, role := range c.Roles {
            if i == cnt {
                break
            }
            if role.gotoforest == false {
                forest          := c.obs.Forests[forestid]
                role.gotoforest = true
                role.forest     = forest
                role.Target     = forest.Center
                i++
            }
        }
    }

}

// 有旗的草丛才补充人员
func (c *Catty) redispatch() {
    forestinfo := c.forestinfo()
    for forestid, cnt := range c.forestmap {
        forest := c.obs.Forests[forestid]
        if forest.HasFlag {
            curcnt := forestinfo[forestid]
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
}

func (c *Catty) forestinfo() map[int]int{
    forestinfo := make(map[int]int)
    for _, role := range c.Roles {
        if role.gotoforest {
            forestinfo[role.forest.Id] += 1
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

    c.forestmap = forestGrouping(len(c.obs.MyTank), c.obs.State.Terain, c.obs.mapanalysis)

    c.dispatch()
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
