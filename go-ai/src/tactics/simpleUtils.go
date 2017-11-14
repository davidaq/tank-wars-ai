package tactics

import (
	f "framework"
	"math"
)

// N个点坐标不能重复，
// x、y 也尽量互不相同（避免阻碍友方攻击）
// 在可能的情况下应离我方大本营更近
// 暂时写死
// func FindNearByPos(pos f.Position, n int, tspeed int) ([]f.Position){
// 	var arrPos []f.Position
// 	arrPos = append(arrPos, f.Position{ X:6, Y:6})
// 	arrPos = append(arrPos, f.Position{ X:6, Y:12})
// 	arrPos = append(arrPos, f.Position{ X:12,Y:12})
// 	arrPos = append(arrPos, f.Position{ X:12,Y:6})
//     arrPos = append(arrPos, f.Position{ X:5,Y:9})
// 	return arrPos[0:n]
// }

// 按离P点距离给一组位置排序
func SortByPos(p f.Position, ps []f.Position) (arrPos []f.Position) {
	var distance int
	mapPos := make(map[int][]f.Position)
	for _, pitem := range ps {
        // 有墙壁阻碍时，X或Y轴的距离是不准确的
		distance  = int(math.Min(math.Abs(float64(pitem.X-p.X)), math.Abs(float64(pitem.Y - p.Y))))
		mapPos[distance] = append(mapPos[distance], pitem)
	}
	for _, pos := range mapPos {
		arrPos = append(arrPos, pos...)
	}
	return arrPos
}

// 为一组坐标和一组坦克，按距离远近匹配
func MatchPosTank(ps []f.Position, tanks []f.Tank) map[string]f.Position {
	tankpos := make(map[string]f.Position)
	if len(ps) <= 0 || len(tanks) <= 0 || len(ps) != len(tanks) {
		return nil
	}
	for _, tank := range tanks {
		ps = SortByPos(tank.Pos, ps)
		tankpos[tank.Id] = ps[0]
	}
	return tankpos
}
