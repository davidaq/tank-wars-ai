package framework

import (
	"fmt"
    "time"
)

type Player struct {
	inited bool
	tactics Tactics
	objectives map[string]Objective
	radar *Radar
	traveller *Traveller
	differ *Diff
}

func NewPlayer(tactics Tactics) *Player {
	inst := &Player {
		tactics: tactics,
		objectives: make(map[string]Objective),
		inited: false,
		radar: nil,
		traveller: nil,
		differ: NewDiff(),
	}
	return inst
}

func (self *Player) Play(state *GameState, absTurn bool) map[string]int {
	start := time.Now()

	if !self.inited {
		self.inited = true
		self.tactics.Init(state)
		self.radar = NewRadar()
		self.traveller = NewTraveller()
	}
	diff := self.differ.Compare(state, self.traveller.CollidedTankInForest(state))
	radarResult := self.radar.Scan(state, diff)
	radarResult.ForestThreat = diff.ForestThreat
	self.tactics.Plan(state, radarResult, self.objectives)

	movement := make(map[string]int)
	travel := make(map[string]*Position)
	var noForward []string
	for _, tank  := range state.MyTank {
		objective := self.objectives[tank.Id]
		if objective.Action == ActionTravel {
			travel[tank.Id] = &objective.Target
		} else if objective.Action == ActionTravelWithDodge {
			dodge, ok := radarResult.Dodge[tank.Id]
			if ok && dodge.Threat > 0.001 {
				if dodge.SafePos.X == tank.Pos.X && dodge.SafePos.Y == tank.Pos.Y {
					noForward = append(noForward, tank.Id)
					travel[tank.Id] = &objective.Target
				} else {
					travel[tank.Id] = &dodge.SafePos
				}
			} else {
				travel[tank.Id] = &objective.Target
			}
		} else {
			movement[tank.Id] = objective.Action
		}
	}

	self.traveller.Search(travel, state, movement)
	for _, tankId := range noForward {
		action, _ := movement[tankId]
		if action == ActionMove {
			movement[tankId] = ActionStay
		}
	}
	if absTurn {
		for _, tank  := range state.MyTank {
			action, _ := movement[tank.Id]
			dir := 0
			isTurn := true
			switch action {
			case ActionLeft:
				dir = 1
			case ActionRight:
				dir = 3
			case ActionBack:
				dir = 2
			default:
				isTurn = false
			}
			if isTurn {
				movement[tank.Id] = (tank.Pos.Direction + dir - DirectionUp + 4) % 4 + ActionTurnUp
			}
		}
	}
	elapsed := time.Since(start)
	fmt.Println("Play function took ", elapsed)
	return movement
}

func (self *Player) End(state *GameState) {
	self.tactics.End(state)
}
