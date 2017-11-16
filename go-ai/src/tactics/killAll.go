package tactics

import (
	f "framework"
)

type KillAll struct {
}

//func NewKillAll() *KillALL {
//	return &KillAll {}
//}

func (self *KillAll) Init(state *f.GameState) {
}

func (self *KillAll) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
}

func (self *KillAll) End(state *f.GameState) {
}