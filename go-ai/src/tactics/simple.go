/*
Simple
调度中心, 负责观察战场、确定模式、分配角色

SimpleSniper
占据苟点/草丛的坦克
    - 占据苟点/草丛
    - 除非受到威胁，否则不移动
    - 受到威胁移动后，按情况尝试重新占领草丛
    - 攻击范围内出现敌人，计算成功率的时机并开火

SimpleKiller
自由行动的坦克杀手
    - 负责追击敌方旗点附近坦克
    - 负责掩护或保护我方坦克

SimpleFlagman
负责扛旗的坦克
    - 从坦克中选出血量/距离最合适的坦克
    - 活动范围是旗点附近

SimpleObservation
战场局势观察分析

SimplePolicy
战术集合：占领、守卫、火力试探、追击、齐射等

模式说明
1) flyingstart:  开局，狙击手占领草丛，旗手赶往旗点附近，killer 定位到最近的地方坦克
2) flagfirst:    如果有旗or即将刷新，旗手赶往旗点，killer 在旗点附近巡游
3) tankfirst:    如果没有旗，旗手在旗点附近巡游，killer 追杀附近地方坦克
4) rabiddog:     即将结束且大分差落后时，开启疯狗模式
*/

package tactics

import (
	f "framework"
	"fmt"
)

type Simple struct {
    obs        *Observation         // 对局势的观察分析
    // mode       string               // 模式: flyingstart, flagfirst, tankfirst, rabiddog
    policy     *SimplePolicy        // 战术动作
    flagmen    *SimpleFlagMan
    snipers    *SimpleSniper
    killers    *SimpleKiller
}

func NewSimple() *Simple{
    return &Simple {}
}

func (s *Simple) NewFlagMan(tanks []f.Tank) {
    s.flagmen = &SimpleFlagMan { superior: s, policy: s.policy, tanks: tanks, obs: s.obs}
    s.flagmen.Init()
}

func (s *Simple) NewSniper(tanks []f.Tank) {
    s.snipers = &SimpleSniper { superior: s, policy: s.policy, tanks: tanks, obs: s.obs }
    s.snipers.Init()
}

func (s *Simple) NewKiller(tanks []f.Tank) {
    s.killers = &SimpleKiller { superior: s, policy: s.policy, tanks: tanks, obs: s.obs }
    s.killers.Init()
}

// 分配角色
func (s *Simple) initRole(state *f.GameState, fcnt int, scnt int, kcnt int) {
    s.NewFlagMan(state.MyTank[0:fcnt])          // 应选血厚、离Flag最近的坦克
    s.NewSniper(state.MyTank[fcnt:fcnt + scnt]) // 应选血薄、离草丛近的坦克
    s.NewKiller(state.MyTank[fcnt + scnt:])      // 剩余都作为 killer
}

func (s *Simple) Init(state *f.GameState) {
    s.obs    = NewObservation(state)
    s.policy = NewSimplePolicy()

    // s.mode   = "flyingstart"

    // 初始化角色
    s.initRole(state, s.obs.Fcnt, s.obs.Scnt, s.obs.Kcnt)
}

// 制定整体计划
func (s *Simple) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    // 分析局势，输出 mode
    s.makeObservation(state)

    // 设定模式
    // s.setMode(state)

	// 分析雷达建议
	s.checkRadar(radar, objective)

    // 子单位执行计划
    s.flagmen.Plan(state, objective)
    s.snipers.Plan(state, objective)
    s.killers.Plan(state, objective)

	fmt.Println("%+v", objective)
}

func (self *Simple) End(state *f.GameState) {
}

// 局势分析
func (s *Simple) makeObservation(state *f.GameState) {
    s.obs.makeObservation(state)
}

// 设定模式
// func (s *Simple) setMode() {
//     if s.obs.CurSteps < s.obs.StartSteps {
//         s.mode = "flyingstart"
//
//     // 局尾大比分落后，开启疯狗模式
//     } else if (s.obs.TotalSteps - s.obs.CurSteps < s.obs.EndSteps && s.obs.emyScore - s.obs.myScore > 0) {
//         s.mode = "rabiddog"
//
//     } else {
//         // 旗点已刷新 or 即将刷新
//         if s.obs.Flag.Exist || s.obs.Flag.Next <= s.MinStepsToFlag {
//             s.mode = "flagfirst"
//         } else {
//             s.mode = "tankfirst"
//         }
//     }
// }

// 检查雷达输出结果，决定躲避 or 开火
func (s *Simple) checkRadar(radar *f.RadarResult, objs map[string]f.Objective) {
	var mrf *f.RadarFire
	var rfs []*f.RadarFire
	for _, tank := range s.obs.CurState.MyTank {
		if radar.Dodge[tank.Id].Threat >= 0.7 {
			objs[tank.Id] = f.Objective{ Action: f.ActionTravel, Target: radar.Dodge[tank.Id].SafePos }
		} else {
			mrf = nil
			rfs = []*f.RadarFire{ radar.Fire[tank.Id].Up, radar.Fire[tank.Id].Down, radar.Fire[tank.Id].Left, radar.Fire[tank.Id].Right }
			for _, rf := range rfs {
				if mrf == nil || mrf.Faith - mrf.Sin < rf.Faith - rf.Sin {
					mrf = rf
				}
			}
			if mrf == nil {
				continue
			}
			if mrf.Faith > 0.5 && mrf.Sin < 0.3 {
				objs[tank.Id] = f.Objective{ Action: mrf.Action }
			}
		}
	}
}
