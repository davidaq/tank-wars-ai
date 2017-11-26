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
                role.Target     = c.nfEntrance(role.Tank.Pos, forest)
                i++
            }
        }
    }
}

// 草丛最近入口
func (c *Catty) nfEntrance(pos f.Position, forest f.Forest) f.Position {
    var target f.Position
    dist := -1
    for p, _ := range forest.ForestMap {
        if nd := pos.SDist(p); dist < 0 || nd < dist {
            dist   = nd
            target = p
        }
    }
    return target
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
                    role.Target     = c.nfEntrance(role.Tank.Pos, forest)
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
    fmt.Printf("forestinfo: %+v\n", forestinfo)
    return forestinfo
}

func (c *Catty) Init(state *f.GameState) {
    // 初始化角色
    c.obs = NewObservation(state)
	for _, tank := range state.MyTank {
		c.Roles[tank.Id] = &CattyRole { obs: c.obs, Tank: tank}
	}

    c.forestmap = forestGrouping(len(c.obs.MyTank), c.obs.State.Terain, c.obs.mapanalysis)
    if len(c.forestmap) > 0 {
        c.dispatch()
    }
}

func (c *Catty) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    fmt.Printf("forestmap: %+v\n", c.forestmap)

    for tankid := range objective {
        delete(objective, tankid)
    }

    c.obs.makeObservation(state, radar, objective)

    c.updateRole()

    if len(c.forestmap) > 0{
        c.redispatch()
    }

    // 无草
    if len(c.forestmap) == 0 {
        for _, role := range c.Roles {
            if role.obs.Flag.Exist && role.obs.Flag.Next <= 5 {
                role.occupyFlag()
                continue
            }
            role.hunt()
            role.act()
            fmt.Printf("catty role: %+v\n", role)
        }

    // 有草
    } else {
        for _, role := range c.Roles {
            if !role.gotoforest {
                role.hunt()
                role.act()
            } else {
                role.patrol()
            }
            fmt.Printf("catty role: %+v\n", role)
        }
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
