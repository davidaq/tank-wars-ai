package tactics

import (
	f "framework"
	"fmt"
)

type Brute struct {
	round int
	relays []f.Position
}

func NewBrute() *Brute {
	inst := &Brute {
		round: 0,
	}
	inst.relays = append(inst.relays, f.Position { X: 6, Y: 15 }, f.Position { X: 15, Y: 6 }, f.Position { X: 24, Y: 15 }, f.Position { X: 15, Y: 24 })
	return inst
}

func (self *Brute) Init(state *f.GameState) {
}

func (self *Brute) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	self.round++
	if state.FlagWait < 5 {
		self.PlanCatchFlag(state, radar, objective)
	} else {
		self.PlanFarShoot(state, radar, objective)
	}
	fireForest := make(map[f.Position]FireForest)
	for _, tank := range state.MyTank {
		if tank.Bullet == "" {
			fireRadar := radar.Fire[tank.Id]
			fireForest[f.Position { X: tank.Pos.X - 1, Y: tank.Pos.Y }] = FireForest { tank.Id, f.ActionFireLeft }
			fireForest[f.Position { X: tank.Pos.X + 1, Y: tank.Pos.Y }] = FireForest { tank.Id, f.ActionFireRight }
			fireForest[f.Position { X: tank.Pos.X, Y: tank.Pos.Y - 1 }] = FireForest { tank.Id, f.ActionFireUp }
			fireForest[f.Position { X: tank.Pos.X, Y: tank.Pos.Y + 1}] = FireForest { tank.Id, f.ActionFireDown }
			for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
				if fire != nil && fire.Sin < 0.5 && fire.Faith > 0.1 {
					objective[tank.Id] = f.Objective {
						Action: fire.Action,
					}
					continue
				}
			}
		}
	}
	for position, posibility := range radar.ForestThreat {
		if posibility > 0.9 {
			if fire, ok := fireForest[f.Position { X: position.X, Y: position.Y }]; ok {
				objective[fire.tankId] = f.Objective {
					Action: fire.action,
				}
			}
		}
	}
}

func around(pos f.Position, counter *int) f.Position {
	ret := pos
	for i, j := 0, 0; i < *counter; i++ {
		switch j {
		case 0:
			ret.X++
			j = 1
		case 1:
			ret.X--
			ret.Y++
			j = 2
		case 2:
			ret.X++
			j = 0
		}
	}
	(*counter)++
	fmt.Println(ret)
	return ret
}

func (self *Brute) nearestRelay(x, y int) f.Position {
	p := f.Position { X: x, Y : y }
	ret := p
	dist := -1
	for _, c := range self.relays {
		if nd := c.SDist(p); dist < 0 || nd < dist {
			dist = nd
			ret = c
		}
	}
	return ret
}

func (self *Brute) PlanFarShoot(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	aroundRelay := 0
	sumX := 0
	sumY := 0
	count := 0
	for _, tank := range state.MyTank {
		sumX += tank.Pos.X
		sumY += tank.Pos.Y
		count++
	}
	for _, tank := range state.EnemyTank {
		sumX += state.Terain.Width - tank.Pos.X - 1
		sumY += state.Terain.Height - tank.Pos.Y - 1
		count++
	}
	avgX := sumX / count
	avgY := sumY / count
	for _, tank := range state.MyTank {
		travel := f.ActionTravel
		if radar.Dodge[tank.Id].Threat > 0.8 || tank.Bullet != "" {
			travel = f.ActionTravelWithDodge
		}
		target := around(self.nearestRelay((tank.Pos.X + avgX) / 2, (tank.Pos.Y + avgY) / 2), &aroundRelay)
		objective[tank.Id] = f.Objective {
			Action: travel,
			Target: target,
		}
	}
}

func (self *Brute) PlanCatchFlag(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	for _, tank := range state.MyTank {
		travel := f.ActionTravel
		if radar.Dodge[tank.Id].Threat > 0.9 {
			travel = f.ActionTravelWithDodge
		}
		objective[tank.Id] = f.Objective {
			Action: travel,
			Target: f.Position { X: state.Params.FlagX, Y: state.Params.FlagY },
		}
	}
}

func (self *Brute) End(state *f.GameState) {
}
