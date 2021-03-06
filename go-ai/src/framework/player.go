package framework

import (
	"fmt"
    "time"
)

type Dodge struct {
	round int
	pos Position
}

type Player struct {
	inited bool
	tactics Tactics
	objectives map[string]Objective
	dodge map[string]Dodge
	radar *Radar
	traveller *Traveller
	differ *Diff
	round int
	firstFlagGenerated bool
	initTank int
	nextFlag int
	rotated bool
	rotatedTerain *Terain
	prevMove map[string]int
}

func NewPlayer(tactics Tactics) *Player {
	inst := &Player {
		tactics: tactics,
		objectives: make(map[string]Objective),
		dodge: make(map[string]Dodge),
		inited: false,
		radar: nil,
		traveller: nil,
		differ: NewDiff(),
		round: 0,
		firstFlagGenerated: false,
		initTank: 0,
		nextFlag: 0,
		prevMove: make(map[string]int),
	}
	return inst
}

func (self *Player) Play(state *GameState) map[string]int {
	start := time.Now()
	if self.initTank == 0 {
		self.initTank = len(state.MyTank)
		if state.MyTank[0].Pos.X > state.Terain.Width / 2 {
			self.rotated = true
		}
	}
	// 自动翻转
	if self.rotated {
		state.Terain = self.rotateTerain(state.Terain, state)
		state.MyTank = self.rotateTank(state.MyTank, state)
		state.EnemyTank = self.rotateTank(state.EnemyTank, state)
		state.MyBullet = self.rotateBullet(state.MyBullet, state)
		state.EnemyBullet = self.rotateBullet(state.EnemyBullet, state)
	}
	fmt.Println("State", state.MyTank, state.EnemyTank)
	// 预测旗子等待时间
	if state.FlagWait > 0 {
		state.FlagWait = 999999
		if self.round <= state.Params.MaxRound / 2 {
			if len(state.MyTank) == self.initTank {
				state.FlagWait = state.Params.MaxRound / 2 - self.round
			}
		} else if self.firstFlagGenerated {
			if self.nextFlag == 0 {
				base := state.Params.MaxRound / 2 / self.initTank + 1
				self.nextFlag = ((self.round - state.Params.MaxRound / 2) / base + 1) * base + state.Params.MaxRound / 2
			}
			state.FlagWait = self.nextFlag - self.round
		}
		if state.FlagWait <= 0 {
			state.FlagWait = 1
		}
	} else {
		self.firstFlagGenerated = true
		self.nextFlag = 0
		state.FlagWait = 0
	}
	if !self.inited {
		self.inited = true
		self.tactics.Init(state)
		self.radar = NewRadar()
		self.traveller = NewTraveller()
	}
	diff := self.differ.Compare(state, self.prevMove, self.traveller.CollidedTankInForest(state))
	radarResult := self.radar.Scan(state, diff)
	radarResult.ForestThreat = diff.ForestThreat
	self.tactics.Plan(state, radarResult, self.objectives)
	self.forestColideShoot(state, radarResult, self.objectives)

	movement := make(map[string]int)
	travel := make(map[string]*Position)
	var noForward []string
	for _, tank  := range state.MyTank {
		objective := self.objectives[tank.Id]
		if objective.Action == ActionTravel {
			travel[tank.Id] = &objective.Target
		} else if objective.Action == ActionTravelWithDodge {
			dodge, ok := radarResult.DodgeBullet[tank.Id]
			if ok && dodge.Threat > 0 {
				// if ododge, ok := self.dodge[tank.Id]; ok && self.round - ododge.round < 2 && ododge.pos.SDist(tank.Pos) > 0 {
				// 	travel[tank.Id] = &ododge.pos
				// } else 
				if dodge.SafePos.SDist(tank.Pos) == 0 {
					movement[tank.Id] = ActionTravelWithDodge
					noForward = append(noForward, tank.Id)
					travel[tank.Id] = &objective.Target
				} else {
					movement[tank.Id] = ActionTravelWithDodge
					travel[tank.Id] = &dodge.SafePos
					self.dodge[tank.Id] = Dodge {
						round: self.round,
						pos: dodge.SafePos,
					}
				}
			} else {
				travel[tank.Id] = &objective.Target
			}
		} else {
			movement[tank.Id] = objective.Action
			fireDirection := DirectionNone
			switch objective.Action {
				case ActionFireUp: fallthrough
				case ActionFireDown: fallthrough
				case ActionFireLeft: fallthrough
				case ActionFireRight:
					fireDirection = objective.Action - ActionFireUp + DirectionUp
			}
			if fireDirection != DirectionNone {
				pos := tank.Pos
				for i, c := 0, state.Params.BulletSpeed * 2 + 2; i < c; i++ {
					pos = pos.step(pos.Direction)
					radarResult.FullMapThreat[pos.NoDirection()] = 5
				}
			}
		}
	}
	self.traveller.Search(travel, state, radarResult.FullMapThreat, movement)
	for _, tankId := range noForward {
		action, _ := movement[tankId]
		if action == ActionMove {
			movement[tankId] = ActionStay
		}
	}
	for _, tank  := range state.MyTank {
		action, _ := movement[tank.Id]
		dir := 0
		isTurn := false
		switch action {
		case ActionLeft:
			isTurn = true
			dir = 1
		case ActionRight:
			isTurn = true
			dir = 3
		case ActionBack:
			isTurn = true
			dir = 2
		}
		if isTurn && dir > 0 {
			movement[tank.Id] = (tank.Pos.Direction + dir - DirectionUp + 4) % 4 + ActionTurnUp
		}
	}
	BadCase(state, radarResult, movement)
	self.prevMove = make(map[string]int)
	for k, v := range movement {
		self.prevMove[k] = v
	}
	if self.rotated {
		for tankId, action := range movement {
			switch action {
			case ActionTurnUp: fallthrough
			case ActionTurnLeft: fallthrough
			case ActionTurnDown: fallthrough
			case ActionTurnRight:
				movement[tankId] = (action - ActionTurnUp + 2) % 4 + ActionTurnUp
			case ActionFireUp: fallthrough
			case ActionFireLeft: fallthrough
			case ActionFireDown: fallthrough
			case ActionFireRight:
				movement[tankId] = (action - ActionFireUp + 2) % 4 + ActionFireUp
			}
		}
	}
	self.round++
	elapsed := time.Since(start)
	fmt.Println("Play function took ", elapsed)
	return movement
}

func (self *Player) forestColideShoot(state *GameState, radar *RadarResult, objective map[string]Objective) {
	type FireForest struct {
		tankId string
		action int
	}
	fireForest := make(map[Position]FireForest)
	for _, tank := range state.MyTank {
		radarFire := radar.Fire[tank.Id]
		if tank.Bullet == "" {
			if radarFire.Left != nil && radarFire.Left.Sin < 0.6 {
				fireForest[Position { X: tank.Pos.X - 1, Y: tank.Pos.Y }] = FireForest { tank.Id, ActionFireLeft }
			}
			if radarFire.Right != nil && radarFire.Right.Sin < 0.6 {
				fireForest[Position { X: tank.Pos.X + 1, Y: tank.Pos.Y }] = FireForest { tank.Id, ActionFireRight }
			}
			if radarFire.Up != nil && radarFire.Up.Sin < 0.6 {
				fireForest[Position { X: tank.Pos.X, Y: tank.Pos.Y - 1 }] = FireForest { tank.Id, ActionFireUp }
			}
			if radarFire.Down != nil && radarFire.Down.Sin < 0.6 {
				fireForest[Position { X: tank.Pos.X, Y: tank.Pos.Y + 1}] = FireForest { tank.Id, ActionFireDown }
			}
		}
	}
	for position, posibility := range radar.ForestThreat {
		if posibility > 0.9 {
			if fire, ok := fireForest[Position { X: position.X, Y: position.Y }]; ok {
				objective[fire.tankId] = Objective {
					Action: fire.action,
				}
			}
		}
	}
}

func (self *Player) End(state *GameState) {
	self.tactics.End(state)
}

func (self *Player) rotateTerain (terain *Terain, state *GameState) *Terain {
	// terain is asymetric, no need to rotate
	return terain
	// if self.rotatedTerain == nil {
	// 	self.rotatedTerain = &Terain {
	// 		Width: terain.Width,
	// 		Height: terain.Height,
	// 		Data: make([][]int, terain.Height),
	// 	}
	// 	for y, line := range terain.Data {
	// 		xline := make([]int, terain.Width)
	// 		for x, val := range line {
	// 			xline[terain.Width - x - 1] = val
	// 		}
	// 		self.rotatedTerain[terain.Height - y - 1] = xline
	// 	}
	// }
	// return self.rotatedTerain
}

func (self *Player) rotateTank (tank []Tank, state *GameState) []Tank {
	ret := make([]Tank, len(tank))
	for i, ot := range tank {
		nt := ot
		self.rotatePos(&nt.Pos, state)
		ret[i] = nt
	}
	return ret
}

func (self *Player) rotateBullet (bullet []Bullet, state *GameState) []Bullet {
	ret := make([]Bullet, len(bullet))
	for i, ot := range bullet {
		nt := ot
		self.rotatePos(&nt.Pos, state)
		ret[i] = nt
	}
	return ret
}

func (self *Player) rotatePos (pos *Position, state *GameState) {
	pos.X = state.Terain.Width - pos.X - 1
	pos.Y = state.Terain.Height - pos.Y - 1
	pos.Direction = (pos.Direction - DirectionUp + 2) % 4 + DirectionUp
}
