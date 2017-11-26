
package tactics

import (
	f "framework"
    "fmt"
)

type Simple struct {
    state         *f.GameState
    radar         *f.RadarResult
    mapanalysis   *f.MapAnalysis
    tanks         map[string]f.Tank
    flaggroup     *FlagGroup
    battlegroup   *BattleGroup
    objectives    map[string]f.Objective
}

func NewSimple() *Simple{
    return &Simple {
        mapanalysis:  &f.MapAnalysis{},
        tanks:        make(map[string]f.Tank),
    }
}

type SimpleRole struct {
    Tank         f.Tank
    Target       f.Position
    Dodge        f.RadarDodge     // 躲避建议
    Fire         f.RadarFireAll   // 开火建议
    ExtDangerSrc []f.ExtDangerSrc // 躲不掉和火线上的威胁源
}

func (s *Simple) NewFlagGroup() {
    s.flaggroup = &FlagGroup { parent: s }
}

func (s *Simple) NewBattleGroup() {
    s.battlegroup = &BattleGroup { parent: s }
}

func (s *Simple) initgroup(state *f.GameState, fcnt int) {
    if fcnt > 0 {
        s.flaggroup.Init(state.MyTank[0:fcnt])                  // 距离flag最近的坦克
    }
    if len(state.MyTank) > fcnt {
        s.battlegroup.Init(state.MyTank[fcnt:])     // 剩余都作为 killer
    }
}


func (s *Simple) Init(state *f.GameState) {
    s.NewFlagGroup()
    s.NewBattleGroup()

    // 启用地图分析
    s.mapanalysis.Analysis(state)

    // 决定分组数量
    fcnt, forest := s.analysis()
    fmt.Printf("forest:%+v\n", forest)

    // 执行分组
    s.initgroup(state, fcnt)

}

// 制定整体计划
func (s *Simple) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    s.state = state
    s.radar = radar

    // 清空上一步的 objective
    for tankid := range objective {
        delete(objective, tankid)
    }

    // 更新信息
    s.update()

    // 重分配角色
    s.reassign()

    // 子单位执行计划
    s.flaggroup.Plan()
    s.battlegroup.Plan()

	fmt.Printf("objective: %+v\n", objective)
}

func (self *Simple) End(state *f.GameState) {}


func (s *Simple) analysis() (fcnt int, forest f.Forest){
    fcnt   = 1
    forest = f.Forest{}
    return
}

// 更新角色信息
func (s *Simple) update() {
    for _, tank := range s.state.MyTank {
        s.tanks[tank.Id] = tank
    }
    for id, role := range s.flaggroup.roles {
        if s.tanks[id] != (f.Tank{}) {
            role.Tank  = s.tanks[id]
            role.Dodge = s.radar.DodgeBullet[id]
            role.Fire  = s.radar.Fire[id]
            role.ExtDangerSrc = s.radar.ExtDangerSrc[id]
        } else {
            delete(s.flaggroup.roles, id)
        }
    }
    for id, role := range s.battlegroup.roles {
        if s.tanks[id] != (f.Tank{}) {
            role.Tank  = s.tanks[id]
            role.Dodge = s.radar.DodgeBullet[id]
            role.Fire  = s.radar.Fire[id]
            role.ExtDangerSrc = s.radar.ExtDangerSrc[id]
        } else {
            delete(s.battlegroup.roles, id)
        }
    }
}

// 重分组
func (s *Simple) reassign() {

}
