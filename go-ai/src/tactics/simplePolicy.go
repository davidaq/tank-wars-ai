// 各类战术动作
package tactics

import (
	f "framework"
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

// 拱卫
// func (p *SimplePolicy) Defend(pos f.Position, tanks []f.Tank, tspeed int, objs map[string]f.Objective){
//     positions := FindNearByPos(pos, len(tanks), tspeed)
// 	p.Dispatch(tanks, positions, objs)
// }

// 开火后占领（靠策略）
// func (p *SimplePolicy) FireAndOccupy(pos Position, tanks []f.Tank, objs *map[string]f.Objective){
// }

// 开火后躲避（靠雷达）
// func (p *SimplePolicy) FireAndDodge(pos Position, tanks []f.Tank, objs *map[string]f.Objective){
// }

// 齐射（还没想好）
// func (p *SimplePolicy) Volley(pos Position, tanks []f.Tank, objs *map[string]f.Objective){
// }

// 巡查某个位置附近的坦克
func (p *SimplePolicy) Patrol(pos f.Position, tanks []f.Tank, emytanks []f.Tank,  objs map[string]f.Objective) (ftanks map[string]f.Tank){
	emypos := make([]f.Position, len(emytanks))
	for i, emytank := range emytanks {
		emypos[i] = emytank.Pos
	}
	// 将敌方坦克按照距某点位置排序
	arrPos := SortByPos(pos, emypos)
    // 敌方坦克数量小于我方坦克
    if (len(arrPos) < len(tanks)) {
        // TODO 会有坦克选择stay
        return p.Dispatch(tanks[0:len(arrPos)], arrPos, objs)
    } else {
        return p.Dispatch(tanks, arrPos[0:len(tanks)], objs)
    }
}

// 指定追击某坦克
// func (p *SimplePolicy) Hunt(tank f.Tank, pos f.Position, objs map[string]f.Objective){
// 	objs[tank.Id] = f.Objective{ Action: f.ActionTravel, Target: Pos}
// }

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
