package tactics

import (
	f "framework"
	"math"
)

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

// 按离P点距离给一组坦克排序
func SortTankByPos(p f.Position, tanks []f.Tank) (res []f.Tank) {
	var distance int
	mapTank := make(map[int][]f.Tank)
	for _, tank := range tanks {
        // 有墙壁阻碍时，X或Y轴的距离是不准确的
		distance  = int(math.Min(math.Abs(float64(tank.Pos.X-p.X)), math.Abs(float64(tank.Pos.Y - p.Y))))
		mapTank[distance] = append(mapTank[distance], tank)
	}
	for _, tank := range mapTank {
		res = append(res, tank...)
	}
	return res
}

// 按离P点距离给一组坦克排序
func TankyByHp(tanks map[string]f.Tank) f.Tank {
	var tanky f.Tank
	for _, tank := range tanks {
        if tanky == (f.Tank{}) || tanky.Hp < tank.Hp {
            tanky = tank
        }
	}
	return tanky
}

// 为一组坐标和一组坦克，按距离远近匹配
// 坐标点数量必须不小于坦克数量
func MatchPosTank(ps []f.Position, tanks []f.Tank) map[string]f.Position {
	tankpos := make(map[string]f.Position)
	if len(ps) <= 0 || len(tanks) <= 0 || len(ps) < len(tanks) {
		return nil
	}
	for _, tank := range tanks {
		ps = SortByPos(tank.Pos, ps)
		tankpos[tank.Id] = ps[0]
	}
	return tankpos
}
