// 雷达系统，用于侦测威胁以及可以开火的目标
package framework;

type Radar struct {
}

func NewRadar() *Radar {
	inst := &Radar {
	}
	return inst
}

// 检查自己当前的位置以及面朝方向前进的位置周围是否安全
// 如果当前位置受到威胁，给出所有可行躲避步骤
// 如果面朝方向受到威胁，给出标记即可
func (self *Radar) ScanThreat(tank *Tank, state *GameState) {
}

func (self *Radar) ScanAttack(tank *Tank, state *GameState) {
}
