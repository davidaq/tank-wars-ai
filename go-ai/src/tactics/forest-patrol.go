package tactics

import (
	f "framework"
	"math/rand"
)

type ForestPatrol struct {
	stateMachine map[string]*StateMachine
	role map[string]string
	round int
}

func NewForestPatrol() *ForestPatrol {
	inst := &ForestPatrol {
		stateMachine: make(map[string]*StateMachine),
		role: make(map[string]string),
		round: 0,
	}
	return inst
}

func (self *ForestPatrol) Init(state *f.GameState) {
	for _, tank := range state.MyTank {
		self.stateMachine[tank.Id] = NewStateMachine("init")
	}
}

type FireForest struct {
	tankId string
	action int
}

func (self *ForestPatrol) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	self.round++
	alive := make(map[string]bool)
	for _, tank := range state.MyTank {
		alive[tank.Id] = true
	}
	for role, id := range self.role {
		if !alive[id] {
			delete(self.role, role)
		}
	}
	fireForest := make(map[f.Position]FireForest)
	tankloop: for _, tank := range state.MyTank {
		fireForest[f.Position { X: tank.Pos.X - 1, Y: tank.Pos.Y }] = FireForest { tank.Id, f.ActionFireLeft }
		fireForest[f.Position { X: tank.Pos.X + 1, Y: tank.Pos.Y }] = FireForest { tank.Id, f.ActionFireRight }
		fireForest[f.Position { X: tank.Pos.X, Y: tank.Pos.Y - 1 }] = FireForest { tank.Id, f.ActionFireUp }
		fireForest[f.Position { X: tank.Pos.X, Y: tank.Pos.Y + 1}] = FireForest { tank.Id, f.ActionFireDown }
		fireRadar := radar.Fire[tank.Id]
		for _, fire := range []*f.RadarFire { fireRadar.Up, fireRadar.Down, fireRadar.Left, fireRadar.Right } {
			if fire != nil && fire.Sin < 0.1 && fire.Faith > 0.7 {
				objective[tank.Id] = f.Objective {
					Action: fire.Action,
				}
				continue tankloop
			}
		}
		if obj := self.stateMachine[tank.Id].Run(self, &tank, state, radar); obj != nil {
			switch obj.Action {
			case f.ActionFireUp:
				if fireRadar.Up.Sin > 0.1 || self.round < 35 {
					obj.Action = f.ActionStay
				}
			case f.ActionFireLeft:
				if fireRadar.Left.Sin > 0.1 || self.round < 35 {
					obj.Action = f.ActionStay
				}
			case f.ActionFireDown:
				if fireRadar.Down.Sin > 0.1 || self.round < 35 {
					obj.Action = f.ActionStay
				}
			case f.ActionFireRight:
				if fireRadar.Right.Sin > 0.1 || self.round < 35 {
					obj.Action = f.ActionStay
				}
			}
			objective[tank.Id] = *obj
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

func (self *ForestPatrol) End(state *f.GameState) {
}

type StateMachine struct {
	currentState string
	counter int
}

func NewStateMachine (initialState string) *StateMachine {
	return &StateMachine {
		currentState: initialState,
		counter: 0,
	}
}

func (self *StateMachine) Run (ctx *ForestPatrol, tank *f.Tank, state *f.GameState, radar *f.RadarResult) *f.Objective {
	switch self.currentState {
	case "init":
		target := f.Position { X: 3, Y: 8 }
		if target.SDist(tank.Pos) < 2 {
			if _, ok := ctx.role["center-shooter"]; !ok {
				self.currentState = "center-shooter-init"
				ctx.role["center-shooter"] = tank.Id
				return self.Run(ctx, tank, state, radar)
			}
			if _, ok := ctx.role["center-patrol"]; !ok {
				self.currentState = "center-patrol-init"
				ctx.role["center-patrol"] = tank.Id
				return self.Run(ctx, tank, state, radar)
			}
			if _, ok := ctx.role["side-patrol"]; !ok {
				self.currentState = "side-patrol-init"
				ctx.role["side-patrol"] = tank.Id
				return self.Run(ctx, tank, state, radar)
			}
			if _, ok := ctx.role["side-shooter"]; !ok {
				self.currentState = "side-shooter-init"
				ctx.role["side-shooter"] = tank.Id
				return self.Run(ctx, tank, state, radar)
			}
			if _, ok := ctx.role["double-shooter"]; !ok {
				self.currentState = "double-shooter-init"
				ctx.role["double"] = tank.Id
				return self.Run(ctx, tank, state, radar)
			}
			self.currentState = "patrol-1"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}
	case "center-shooter-init":
		target := f.Position { X: 9, Y: 7, Direction: f.DirectionDown }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			self.currentState = "center-shooter"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}
	case "center-shooter":
		if state.FlagWait < 4 {
			self.currentState = "flagger"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Action: f.ActionFireDown,
		}
	case "flagger":
		if state.FlagWait >= 4 {
			self.currentState = "center-shooter-init"
			return self.Run(ctx, tank, state, radar)
		}
		target := f.Position { X: 9, Y: 9 }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			return &f.Objective {
				Action: f.ActionFireUp + rand.Int() % 4,
			}
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}

	case "side-shooter-init":
		target := f.Position { X: 10, Y: 12 }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			self.counter = 0
			self.currentState = "side-shooter-A"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}
	case "side-shooter-A":
		if self.counter == 1 {
			self.currentState = "side-shooter-B"
			return self.Run(ctx, tank, state, radar)
		}
		self.counter++
		return &f.Objective {
			Action: f.ActionFireUp,
		}
	case "side-shooter-B":
		target := f.Position { X: 10, Y: 11 }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			self.counter = 0
			self.currentState = "side-shooter-init"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}


	case "center-patrol-init":
		target := f.Position { X: 8, Y: 7, Direction: f.DirectionUp }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			self.counter = 0
			self.currentState = "center-patrol-A"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}
	case "center-patrol-A":
		if self.counter == 1 {
			self.currentState = "center-patrol-B"
			return self.Run(ctx, tank, state, radar)
		}
		self.counter++
		return &f.Objective {
			Action: f.ActionFireDown,
		}
	case "center-patrol-B":
		target := f.Position { X: 9, Y: 6, Direction: f.DirectionLeft }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			self.counter = 0
			self.currentState = "center-patrol-C"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}
	case "center-patrol-C":
		if self.counter == 2 {
			self.currentState = "center-patrol-init"
			return self.Run(ctx, tank, state, radar)
		}
		self.counter++
		return &f.Objective {
			Action: f.ActionFireRight,
		}

	case "side-patrol-init":
		target := f.Position { X: 6, Y: 6 }
		if target.SDist(tank.Pos) == 0 {
			self.counter = 0
			self.currentState = "side-patrol-A"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}
	case "side-patrol-A":
		if self.counter == 1 {
			self.currentState = "side-patrol-B"
			return self.Run(ctx, tank, state, radar)
		}
		self.counter++
		return &f.Objective {
			Action: f.ActionFireDown,
		}
	case "side-patrol-B":
		target := f.Position { X: 5, Y: 6 }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			self.counter = 0
			self.currentState = "side-patrol-init"
			return self.Run(ctx, tank, state, radar)
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}

	case "double-shooter-init":
		target := f.Position { X: 7, Y: 9, Direction: f.DirectionDown }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			if state.FlagWait > 2 {
				self.currentState = "double-shooter-A"
				return &f.Objective {
					Action: f.ActionFireRight,
				}
			}
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}
	case "double-shooter-A":
		target := f.Position { X: 7, Y: 12, Direction: f.DirectionUp }
		if target.SDist(tank.Pos) < state.Params.TankSpeed {
			self.currentState = "double-shooter-init"
			return &f.Objective {
				Action: f.ActionFireLeft,
			}
		}
		return &f.Objective {
			Target: target,
			Action: f.ActionTravel,
		}
		
	}
	return nil
}
