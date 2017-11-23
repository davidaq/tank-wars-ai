package tactics

import (
	f "framework"
	"fmt"
	"math/rand"
)

type Brute struct {
	round int
	relays []f.Position
}

func NewBrute() *Brute {
	inst := &Brute {
		round: 0,
	}
	return inst
}

func (self *Brute) Init(state *f.GameState) {
	fx, fy := state.Params.FlagX, state.Params.FlagY
	up, left, down, right := false, false, false, false
	upv, leftv, downv, rightv := 0, 0, 0, 0
	for i, d := state.Params.BulletSpeed * 2 + 1, state.Params.BulletSpeed * 4; i <= d; i++ {
		if !up {
			if state.Terain.Get(fx, fy - i) == 0 {
				upv = i
			} else {
				up = true
			}
		}
		if !down {
			if state.Terain.Get(fx, fy + i) == 0 {
				downv = i
			} else {
				down = true
			}
		}
		if !left {
			if state.Terain.Get(fx - i, fy) == 0 {
				leftv = i
			} else {
				left = true
			}
		}
		if !right {
			if state.Terain.Get(fx + i, fy) == 0 {
				leftv = i
			} else {
				left = true
			}
		}
	}
	upv -= 1
	leftv -= 1
	downv -= 1
	rightv -= 1
	if up {
		self.relays = append(self.relays, f.Position { X: fx, Y: fy - upv })
	}
	if left {
		self.relays = append(self.relays, f.Position { X: fx - leftv, Y: fy })
	}
	if down {
		self.relays = append(self.relays, f.Position { X: fx, Y: fy + downv })
	}
	if right {
		self.relays = append(self.relays, f.Position { X: fx + rightv, Y: fy })
	}
}

func (self *Brute) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	self.round++
	if state.FlagWait < 5 {
		self.PlanCatchFlag(state, radar, objective)
	} else if !self.PlanKill(state, radar, objective) {
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
			faith := 0.
			var pfire *f.RadarFire
			if radar.DodgeBullet[tank.Id].Threat > 0.7 {
				faith = 0.7
			}
			for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
				fmt.Println("RECV FIRE", tank.Id, fire)
				if fire.Sin < 1 && fire.Faith > faith {
					if (fire.Sin < 0.5 && (fire.Faith > 0.3 || fire.Cost < 5)) || (fire.Sin > 0.49 && fire.Faith > 0.8) {
						pfire = fire
						faith = fire.Faith
					}
				}
			}
			if pfire != nil {
				objective[tank.Id] = f.Objective {
					Action: pfire.Action,
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

func (self *Brute) PlanKill(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) bool {
	if self.round < state.Params.MaxRound / 2 || len(state.MyTank) > len(state.EnemyTank) {
		return false
	}
	sumX := 0
	sumY := 0
	count := 0
	for _, tank := range state.MyTank {
		sumX += tank.Pos.X
		sumY += tank.Pos.Y
		count++
	}
	if count < 2 {
		return false
	}
	mAvg := f.Position { X: sumX / count, Y: sumY / count }

	sumX = 0
	sumY = 0
	count = 0
	for _, tank := range state.EnemyTank {
		sumX += tank.Pos.X
		sumY += tank.Pos.Y
		count++
	}
	if count == 0 {
		return false
	}
	eAvg := f.Position { X: sumX / count, Y: sumY / count }

	least := 0
	var ttank *f.Tank
	for _, tank := range state.EnemyTank {
		if away := tank.Pos.SDist(eAvg); away > state.Params.TankSpeed * 3 {
			d := tank.Pos.SDist(mAvg) - away
			if ttank == nil || d < least {
				least = d
				ttank = &tank
			}
		}
	}
	if ttank == nil {
		return false
	}
	for _, tank := range state.MyTank {
		travel := f.ActionTravel
		if radar.Dodge[tank.Id].Threat > 0.8 || tank.Bullet != "" {
			travel = f.ActionTravelWithDodge
		}
		dx := ttank.Pos.X - tank.Pos.X
		dy := ttank.Pos.Y - tank.Pos.Y
		target := ttank.Pos
		if dx > dy {
			target.X = tank.Pos.X
		} else {
			target.Y = tank.Pos.Y
		}
		objective[tank.Id] = f.Objective {
			Action: travel,
			Target: target,
		}
	}
	fmt.Println("pursue")
	return true
}

func around(state *f.GameState, pos f.Position, counter *int) f.Position {
	ret := pos
	switch *counter {
	case 1:
		ret.X -= state.Params.TankSpeed
		ret.Y -= state.Params.TankSpeed
	case 2:
		ret.X -= state.Params.TankSpeed
		ret.Y += state.Params.TankSpeed
	case 3:
		ret.X += state.Params.TankSpeed
		ret.Y -= state.Params.TankSpeed
	case 4:
		ret.X += state.Params.TankSpeed
		ret.Y += state.Params.TankSpeed
	}
	(*counter)++
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
		} else if radar.Dodge[tank.Id].Threat < 0.1 && tank.Bullet == "" {
			fireRadar := radar.Fire[tank.Id]
			var tryFire *f.RadarFire
			if tank.Pos.X == state.Params.FlagX {
				if tank.Pos.Y > state.Params.FlagY {
					tryFire = fireRadar.Up
				} else {
					tryFire = fireRadar.Down
				}
			} else if tank.Pos.Y == state.Params.FlagY {
				if tank.Pos.X > state.Params.FlagX {
					tryFire = fireRadar.Left
				} else {
					tryFire = fireRadar.Right
				}
			}
			if tryFire != nil && tryFire.Sin < 0.1 {
				fmt.Println("far shoot")
				objective[tank.Id] = f.Objective {
					Action: tryFire.Action,
				}
				continue
			}
		}
		target := around(state, self.nearestRelay((tank.Pos.X * 3 + avgX * 2) / 5, (tank.Pos.Y * 3 + avgY * 2) / 5), &aroundRelay)

		objective[tank.Id] = f.Objective {
			Action: travel,
			Target: target,
		}
	}
}

func (self *Brute) PlanCatchFlag(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	for _, tank := range state.MyTank {
		travel := f.ActionTravel

		if rand.Int() % 3 == 0 {
			if radar.Dodge[tank.Id].Threat > 0.9 {
				travel = f.ActionTravelWithDodge
			}
		} else {
			if radar.DodgeBullet[tank.Id].Threat > 0.7 && radar.DodgeBullet[tank.Id].Threat <= 1 {
				travel = f.ActionTravelWithDodge
			}
		}
		objective[tank.Id] = f.Objective {
			Action: travel,
			Target: f.Position { X: state.Params.FlagX, Y: state.Params.FlagY },
		}
	}
}

func (self *Brute) End(state *f.GameState) {
}
