// 雷达系统，用于侦测威胁以及可以开火的目标
package framework;

type Radar struct {
}

func NewRadar() *Radar {
	inst := &Radar {
	}
	return inst
}

func (self *Radar) Scan(tank *Tank, state *GameState) {
}
