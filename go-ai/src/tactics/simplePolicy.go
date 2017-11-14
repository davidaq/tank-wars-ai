// 各类战术动作
package tactics

import (
	f "framework"
    // "fmt"
)

type SimplePolicy struct {

}

func NewSimplePolicy() *SimplePolicy {
    return &SimplePolicy{ }
}

// 占领
func (p *SimplePolicy) Occupy(pos f.Position, tank f.Tank, objs map[string]f.Objective){
    objs[tank.Id] = f.Objective{ Action: f.ActionTravel, Target: pos}
}

// 巡查某个位置附近的坦克
func (p *SimplePolicy) Patrol(pos f.Position, tanks []f.Tank, emytanks []f.Tank,  objs map[string]f.Objective) (ftanks map[string]f.Tank){
	emypos := make([]f.Position, len(emytanks))
	for i, emytank := range emytanks {
		emypos[i] = emytank.Pos
	}
	arrPos := SortByPos(pos, emypos)  // 按远近排序
    // 敌方坦克数量小于我方坦克
    if (len(arrPos) < len(tanks)) {
        return p.Dispatch(tanks[0:len(arrPos)], arrPos, objs)
    } else {
        return p.Dispatch(tanks, arrPos[0:len(tanks)], objs)
    }
}

// 将一组坦克派到一组地点，并返回空闲的坦克
func (p *SimplePolicy) Dispatch(tanks []f.Tank, pos []f.Position, objs map[string]f.Objective) (ftanks map[string]f.Tank){
    ftanks = make(map[string]f.Tank)
    match  := MatchPosTank(pos[0:len(tanks)], tanks)
	for _, tank := range tanks {
        if tank.Pos.X == match[tank.Id].X && tank.Pos.Y == match[tank.Id].Y {
            ftanks[tank.Id] = tank
        } else {
            objs[tank.Id] = f.Objective{ Action: f.ActionTravel, Target: match[tank.Id]}
        }
	}
    return ftanks
}
