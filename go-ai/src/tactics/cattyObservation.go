// 观察者
package tactics

import (
	f "framework"
	"math"
    "sort"
    "fmt"
)

type Observation struct {
    TotalSteps   int
    Steps        int
    PreState     *f.GameState
    State        *f.GameState
    Radar        *f.RadarResult
    mapanalysis  *f.MapAnalysis
    Objs         map[string]f.Objective
    MyTank       map[string]f.Tank // 我方存活坦克
    EmyTank      map[string]f.Tank // 敌方存活坦克
    Flag         Flag
	Terain       *f.Terain
	// ShotPos      []f.Position
    ShotPos      map[f.Position]string
    Forests      map[int]f.Forest
    TankCnt      int   // 初始坦克数量
}

type Flag struct {
    Pos            f.Position
    Exist          bool      // 是否有旗定时刷新
    Next           int       // 距离刷新的回合数
}

type Pair struct {
    Key    f.Position
    Value  int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

func NewObservation(state *f.GameState) (obs *Observation) {
    obs = &Observation{
        TotalSteps: state.Params.MaxRound,
        Steps: 0,
        State: state,
        Terain: state.Terain,
        mapanalysis: &f.MapAnalysis{},
    }

    obs.TankCnt = len(state.MyTank)

    // 地图分析
    obs.mapanalysis.Analysis(state)

    obs.Forests = make(map[int]f.Forest)
    for _, forest := range obs.mapanalysis.Forests {
        obs.Forests[forest.Id] = forest
    }

    // 观察坦克
    obs.observeTank()

	// 观察战旗
	obs.observeFlag()

	return obs
}

func (o *Observation) makeObservation(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    o.Steps   += 1
    o.State   = state
    o.Radar   = radar
    o.Objs    = objective
	o.Terain  = state.Terain

	// 观察坦克
    o.observeTank()

    // 观察战旗
    o.observeFlag()

	// 观察地图
	o.observeTerain()

	// 观察适合射击的地点
	o.observeShotPos()
}

func (o *Observation) observeTank() {
    o.MyTank = make(map[string]f.Tank)
    for _, tank := range o.State.MyTank {
        o.MyTank[tank.Id] = tank
    }
    o.EmyTank = make(map[string]f.Tank)
    for _, tank := range o.State.EnemyTank {
        o.EmyTank[tank.Id] = tank
    }
}

func (o *Observation) observeFlag() {
    if o.Steps < o.TotalSteps / 2 || o.State.FlagWait == 999999 {
        o.Flag = Flag { Pos: f.Position{ X: o.State.Params.FlagX, Y:o.State.Params.FlagY }, Exist: false, Next: 0 }
    } else {
        o.Flag = Flag{ Pos: f.Position{ X: o.State.Params.FlagX, Y:o.State.Params.FlagY }, Exist: true, Next: o.State.FlagWait }
    }
}

// 将坦克和子弹渲染到地图上
func (o *Observation) observeTerain() {
	for _, tank := range o.MyTank {
		o.Terain.Data[tank.Pos.Y][tank.Pos.X] = 4 // 我方坦克
	}
	for _, tank := range o.EmyTank {
		o.Terain.Data[tank.Pos.Y][tank.Pos.X] = 5 // 敌方坦克
	}
	for _, bullet := range o.State.MyBullet {
		o.Terain.Data[bullet.Pos.Y][bullet.Pos.X] = 6 // 子弹
	}
	for _, bullet := range o.State.EnemyBullet {
		o.Terain.Data[bullet.Pos.Y][bullet.Pos.X] = 6 // 子弹
	}
}

func (o *Observation) observeShotPos() {
    // 清空上一步结果
    // o.ShotPos = []f.Position{}
    o.ShotPos = make(map[f.Position]string)
    shotPos := o.findShotPos()
    if len(shotPos) > 0 {
        // 如果有绝杀点，筛选掉其中不能去的点
        for pos, tankid := range shotPos {
            if o.huntable(pos, tankid) {
                o.ShotPos[pos] = tankid
            }
        }
        // 按距离我方中心点排序
        // avgPos := o.avgpos()
        // fmt.Printf("avgPos: %+v\n", avgPos)

        // 按离中心点距离，给 positions 排序
        // o.ShotPos = o.sortByPos(avgPos, o.ShotPos)

        // 多选一些点，避免集中
        // if len(o.ShotPos) - len(o.MyTank) >= len(o.MyTank) {
        //     o.ShotPos = o.ShotPos[0:len(o.MyTank)+len(o.MyTank)]
        // }
    }
    fmt.Printf("observeShotPos: %+v\n", o.ShotPos)
}

// 寻找攻击地点【前后选点】
// func (o *Observation) findShotPos() map[f.Position]string {
// 	shotPos := make(map[f.Position]string)
//     var pos f.Position
// 	for _, tank := range o.EmyTank {
//         if tank.Pos.Direction == f.DirectionUp {
//             pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y + 2*o.State.Params.BulletSpeed + 1 + o.State.Params.TankSpeed }
//             if o.reachable(pos) && (pos.Y - tank.Pos.Y <= 2 * o.State.Params.BulletSpeed) {
// 				shotPos[pos] = tank.Id
// 			}
// 			pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y - 2*o.State.Params.BulletSpeed - 1}
// 			if o.reachable(pos) {
// 				shotPos[pos] = tank.Id
// 			}
//
//         } else if tank.Pos.Direction == f.DirectionDown {
// 			pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y + 2*o.State.Params.BulletSpeed + 1}
// 			if o.reachable(pos) {
// 				shotPos[pos] = tank.Id
// 			}
// 			pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y - 2*o.State.Params.BulletSpeed - 1 - o.State.Params.TankSpeed }
// 			if o.reachable(pos) && tank.Pos.Y - pos.Y <= 2 * o.State.Params.BulletSpeed {
// 				shotPos[pos] = tank.Id
// 			}
//
//         } else if tank.Pos.Direction == f.DirectionLeft {
// 			pos = f.Position { X: tank.Pos.X + 2*o.State.Params.BulletSpeed + 1, Y: tank.Pos.Y}
// 			if o.reachable(pos) {
// 				shotPos[pos] = tank.Id
// 			}
// 			pos = f.Position { X: tank.Pos.X - 2*o.State.Params.BulletSpeed - 1 - o.State.Params.TankSpeed, Y: tank.Pos.Y}
// 			if o.reachable(pos) && tank.Pos.X - pos.X <= 2 * o.State.Params.BulletSpeed{
// 				shotPos[pos] = tank.Id
// 			}
//         } else {
//             pos = f.Position { X: tank.Pos.X + 2*o.State.Params.BulletSpeed + 1 + o.State.Params.TankSpeed, Y: tank.Pos.Y}
// 			if o.reachable(pos) && pos.X - tank.Pos.X <= 2 * o.State.Params.BulletSpeed {
// 				shotPos[pos] = tank.Id
// 			}
// 			pos = f.Position { X: tank.Pos.X - 2*o.State.Params.BulletSpeed - 1, Y: tank.Pos.Y}
// 			if o.reachable(pos) {
// 				shotPos[pos] = tank.Id
// 			}
//         }
// 	}
//     return shotPos
// }


// 我方坦克中心(含旗点)
func (o *Observation) avgpos() f.Position {
    sumX  := 0
	sumY  := 0
	count := 0
    for _, tank := range o.MyTank {
        sumX += tank.Pos.X
        sumY += tank.Pos.Y
        count++
    }
    sumX  += o.Flag.Pos.X
    sumY  += o.Flag.Pos.Y
    count += 1

    // 中心点
    return f.Position { X: sumX / count, Y: sumY / count }
}

// 按距离给一组点排序
func (o *Observation) sortByPos(pos f.Position, ps []f.Position) (positions []f.Position) {
    pl := make(PairList, len(ps))
    i  := 0
    for _, p := range ps {
        dist := pos.SDist(p)
        pl[i] = Pair{ p, dist }
        i++
    }
    sort.Sort(pl)
    for _, p := range pl {
        positions = append(positions, p.Key)
    }
    return positions
}


// 寻找攻击地点【十字围杀】
func (o *Observation) findShotPos() map[f.Position]string{
    shotPos := make(map[f.Position]string)
    var pos f.Position
	for _, tank := range o.EmyTank {
        if tank.Pos.Direction == f.DirectionUp {
			pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y - 2*o.State.Params.BulletSpeed - 2}
			if o.reachable(pos) {
				shotPos[pos] = tank.Id
			}
            pos = f.Position { X: tank.Pos.X + 2*o.State.Params.BulletSpeed + 2, Y: tank.Pos.Y}
			if o.reachable(pos) {
				shotPos[pos] = tank.Id
			}
            pos = f.Position { X: tank.Pos.X - 2*o.State.Params.BulletSpeed - 2, Y: tank.Pos.Y}
            if o.reachable(pos) {
                shotPos[pos] = tank.Id
            }

        } else if tank.Pos.Direction == f.DirectionDown {
			pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y + 2*o.State.Params.BulletSpeed + 2}
			if o.reachable(pos) {
				shotPos[pos] = tank.Id
			}
            pos = f.Position { X: tank.Pos.X + 2*o.State.Params.BulletSpeed + 2, Y: tank.Pos.Y}
            if o.reachable(pos) {
                shotPos[pos] = tank.Id
            }
            pos = f.Position { X: tank.Pos.X - 2*o.State.Params.BulletSpeed - 2, Y: tank.Pos.Y}
            if o.reachable(pos) {
                shotPos[pos] = tank.Id
            }

        } else if tank.Pos.Direction == f.DirectionLeft {
            pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y + 2*o.State.Params.BulletSpeed + 2}
            if o.reachable(pos) {
                shotPos[pos] = tank.Id
            }
            pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y - 2*o.State.Params.BulletSpeed - 2}
            if o.reachable(pos) {
                shotPos[pos] = tank.Id
            }
			pos = f.Position { X: tank.Pos.X + 2*o.State.Params.BulletSpeed + 2, Y: tank.Pos.Y}
			if o.reachable(pos) {
				shotPos[pos] = tank.Id
			}
        } else {
            pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y + 2*o.State.Params.BulletSpeed + 2}
            if o.reachable(pos) {
                shotPos[pos] = tank.Id
            }
            pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y - 2*o.State.Params.BulletSpeed - 2}
            if o.reachable(pos) {
                shotPos[pos] = tank.Id
            }
            // pos = f.Position { X: tank.Pos.X + 2*o.State.Params.BulletSpeed + 2 + o.State.Params.TankSpeed, Y: tank.Pos.Y}
			// if o.reachable(pos) && pos.X - tank.Pos.X <= 2 * o.State.Params.BulletSpeed {
			// 	shotPos[pos] = tank.Id
			// }
			pos = f.Position { X: tank.Pos.X - 2*o.State.Params.BulletSpeed - 2, Y: tank.Pos.Y}
			if o.reachable(pos) {
				shotPos[pos] = tank.Id
			}
        }
	}
    return shotPos
}

// 追击点能不能去
func (o *Observation) huntable(pos f.Position, tankid string) bool{
    tank := o.EmyTank[tankid]
    return o.pathReachable(pos, tank.Pos)

    // // 追击点两侧的逃生点
    // tank := o.EmyTank[tankid]
    // positions := make([]f.Position, 2)
    // if tank.Pos.Direction == f.DirectionUp || tank.Pos.Direction == f.DirectionDown {
    //     positions = []f.Position {
    //         f.Position { X: pos.X - o.State.Params.TankSpeed, Y: pos.Y },
    //         f.Position { X: pos.X + o.State.Params.TankSpeed, Y: pos.Y },
    //     }
    // } else {
    //     positions = []f.Position {
    //         f.Position { X: pos.X, Y: pos.Y - o.State.Params.TankSpeed},
    //         f.Position { X: pos.X, Y: pos.Y + o.State.Params.TankSpeed},
    //     }
    // }
    // // 若无可达逃生点，那么不能去
    // if !o.pathReachable(pos, positions[0]) && !o.pathReachable(pos, positions[1]){
    //     return false
    // } else {
    //     return true
    // }
}


// 地点是否可达（是否超出地图范围、是否墙壁）
func (o *Observation) reachable(pos f.Position) bool {
    // 超出地图范围
    if pos.X < 0 || pos.X >= o.Terain.Width || pos.Y < 0 || pos.Y >= o.Terain.Height {
        return false
    }
    // 是否墙壁
    if o.Terain.Data[pos.Y][pos.X] == 1 {
        return false
    }
    return true
}

// 路径是否可达
func (o *Observation) pathReachable(startpos f.Position, endpos f.Position) bool {
	terain := o.Terain.Data
	// 两点是否可达
	if !o.reachable(startpos) || !o.reachable(endpos) {
		return false
	}
	// 路径是否可达（路径无墙壁、无坦克、无子弹）
    var min, max int
    if startpos.X == endpos.X {
        min = int(math.Min(float64(startpos.Y), float64(endpos.Y)))
        max = int(math.Max(float64(startpos.Y), float64(endpos.Y)))
        for i := min+1; i < max; i++ {
            if terain[i][startpos.X] != 0 && terain[i][startpos.X] != 2 && terain[i][startpos.X] != 3 {
                return false
            }
        }
    } else if startpos.Y == endpos.Y {
        min = int(math.Min(float64(startpos.X), float64(endpos.X)))
        max = int(math.Max(float64(startpos.X), float64(endpos.X)))
        for i := min+1; i < max; i++ {
            if terain[startpos.Y][i] != 0 && terain[startpos.Y][i] != 2 && terain[startpos.Y][i] != 3 {
                return false
            }
        }
    } else {
        return false
    }
    return true
}
