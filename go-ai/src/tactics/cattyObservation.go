// 观察者
package tactics

import (
	f "framework"
	"math"
    // "fmt"
)

type Observation struct {
    TotalSteps   int
    Steps        int
    PreState     *f.GameState
    State        *f.GameState
    Radar        *f.RadarResult
    Objs         map[string]f.Objective
    MyTank       map[string]f.Tank // 我方存活坦克
    EmyTank      map[string]f.Tank // 敌方存活坦克
    Flag         Flag
	Terain       *f.Terain
	ShotPos      map[f.Position]string
	Q            int   // 安全系数 Q * bulletspeed
}

type Flag struct {
    Pos            f.Position
    Exist          bool      // 是否有旗定时刷新
    Next           int       // 距离刷新的回合数
}

func NewObservation(state *f.GameState) (obs *Observation) {
    obs = &Observation{ TotalSteps: state.Params.MaxRound, Steps: 0, State: state, Terain: state.Terain}

	obs.Q = 1

    // 观察坦克
    obs.observeTank()

	// 观察战旗
	obs.observeFlag()

	return obs
}

func (o *Observation) makeObservation(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
    o.Steps  += 1
    o.State  = state
    o.Radar  = radar
    o.Objs   = objective
	o.Terain = state.Terain

	// 观察坦克
    o.observeTank()

    // 观察战旗
    o.observeFlag()

	// 观察地图
	o.observeTerain()

	if len(o.MyTank) < len(o.EmyTank) {
		o.Q = 2
	} else {
		o.Q = 1
	}

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
	o.ShotPos = make(map[f.Position]string)
    var pos f.Position
	for _, tank := range o.EmyTank {
		for i := o.State.Params.BulletSpeed; i > 0; i-- {
	        if tank.Pos.Direction == f.DirectionUp || tank.Pos.Direction == f.DirectionDown {
				pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y + o.State.Params.BulletSpeed * (o.Q+1) + i + 2}
				if o.reachable(pos) {
					o.ShotPos[pos] = tank.Id
				}
				pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y - o.State.Params.BulletSpeed * o.Q - i - 2}
				if o.reachable(pos) {
					o.ShotPos[pos] = tank.Id
				}
	        } else {
				pos = f.Position { X: tank.Pos.X + o.State.Params.BulletSpeed * (o.Q+1) + i + 2, Y: tank.Pos.Y}
				if o.reachable(pos) {
					o.ShotPos[pos] = tank.Id
				}
				pos = f.Position { X: tank.Pos.X - o.State.Params.BulletSpeed * o.Q - i - 2, Y: tank.Pos.Y}
				if o.reachable(pos) {
					o.ShotPos[pos] = tank.Id
				}
	        }
	    }
	}
}

// 换个思路定位追击点
// func (o *Observation) observeShotPos() {
// 	o.ShotPos = make(map[f.Position]string)
//     var pos f.Position
// 	for _, tank := range o.EmyTank {
//         if tank.Pos.Direction == f.DirectionUp {
// 			pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y + o.State.Params.TankSpeed }
// 			if o.reachable(pos) {
// 				o.ShotPos[pos] = tank.Id
// 			}
// 		} else if tank.Pos.Direction == f.DirectionDown {
// 			pos = f.Position { X: tank.Pos.X, Y: tank.Pos.Y - o.State.Params.TankSpeed }
// 			if o.reachable(pos) {
// 				o.ShotPos[pos] = tank.Id
// 			}
//         } else if tank.Pos.Direction == f.DirectionRight {
// 			pos = f.Position { X: tank.Pos.X + o.State.Params.TankSpeed, Y: tank.Pos.Y}
// 			if o.reachable(pos) {
// 				o.ShotPos[pos] = tank.Id
// 			}
// 		} else if tank.Pos.Direction == f.DirectionLeft {
// 			pos = f.Position { X: tank.Pos.X - o.State.Params.TankSpeed, Y: tank.Pos.Y}
// 			if o.reachable(pos) {
// 				o.ShotPos[pos] = tank.Id
// 			}
//         }
// 	}
// }

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
