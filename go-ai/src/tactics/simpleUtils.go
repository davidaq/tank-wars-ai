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

// 寻找血最厚的坦克
func TankyByHp(tanks map[string]f.Tank) f.Tank {
	var tanky f.Tank
	for _, tank := range tanks {
        if tanky == (f.Tank{}) || tanky.Hp < tank.Hp {
            tanky = tank
        }
	}
	return tanky
}

// 按距离match一组坦克与坐标(坐标数不小于坦克数)
func MatchPosTank(ps []f.Position, tanks []f.Tank) map[string]f.Position {
	tankpos := make(map[string]f.Position)
	if len(ps) == 0 || len(tanks) == 0 || len(ps) < len(tanks) {
		return nil
	}
	for _, tank := range tanks {
		ps = SortByPos(tank.Pos, ps)
		tankpos[tank.Id] = ps[0]
	}
	return tankpos
}

// 地点是否可达（是否超出地图范围、是否墙壁）
func IsReachable(p f.Position, terain f.Terain) bool {
    // 超出地图范围
    if p.X < 0 || p.X >= terain.Width || p.Y < 0 || p.Y >= terain.Height {
        return false
    }
    // 是否墙壁
    if terain.Data[p.Y][p.X] == 1 {
        return false
    }
    // TODO 暂时不考虑是否有坦克
    return true
}

// 子弹是否可达（路径上是否有墙壁）
func IsBulletReachable(startpos f.Position, endpos f.Position, terain f.Terain) bool {
    var min, max int
    if startpos.X == endpos.X {
        min = int(math.Min(float64(startpos.Y), float64(endpos.Y)))
        max = int(math.Max(float64(startpos.Y), float64(endpos.Y)))
        for i := min+1; i < max; i++ {
            if terain.Data[i][startpos.X] == 1 {
                return false
            }
        }
    } else if startpos.Y == endpos.Y {
        min = int(math.Min(float64(startpos.X), float64(endpos.X)))
        max = int(math.Max(float64(startpos.X), float64(endpos.X)))
        for i := min+1; i < max; i++ {
            if terain.Data[startpos.Y][i] == 1 {
                return false
            }
        }
    } else {
        return false
    }
    return true
}

// 寻找合适开火地点
func FindShootPos(emypos f.Position, terain f.Terain, bspeed int) (result []f.Position) {
    // 周围合适的攻击地点
    positions := []f.Position {
        f.Position{X: emypos.X - (bspeed + 1), Y:emypos.Y},
        f.Position{X: emypos.X + (bspeed + 1), Y:emypos.Y},
        f.Position{X: emypos.X, Y:emypos.Y - (bspeed + 1)},
        f.Position{X: emypos.X, Y:emypos.Y + (bspeed + 1)},
    }
    for _, position := range positions {
        // 地点是否可达
        if IsReachable(position, terain) {
            // 射程内是否有障碍物
            if IsBulletReachable(position, emypos,  terain) {
                result = append(result, position)
            }
        }
    }
    return result
}
