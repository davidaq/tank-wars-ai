package framework

type Player struct {
	inited bool
	tactics Tactics
	objectives map[string]Objective
	radar *Radar
	traveller *Traveller
}

func NewPlayer(tactics Tactics) *Player {
	inst := &Player {
		tactics: tactics,
		objectives: make(map[string]Objective),
		inited: false,
		radar: nil,
		traveller: nil,
	}
	return inst
}

func (self *Player) Play(state *GameState) map[string]int {
	if !self.inited {
		self.inited = true
		self.tactics.Init(state)
		self.radar = NewRadar()
		self.traveller = NewTraveller()
	}
	radarResult := self.radar.Scan(state)
	self.tactics.Plan(state, radarResult, self.objectives)

	movement := make(map[string]int)
	travellingTank := make(map[string]*Position)
	var noForward []string
	for _, tank  := range state.MyTank {
		objective := self.objectives[tank.Id]
		if objective.Action == ActionTravel {
			travellingTank[tank.Id] = &objective.Target
		} else if objective.Action == ActionTravelWithDodge {
			dodge, ok := radarResult.Dodge[tank.Id]
			if ok {
				if dodge.SafePos.X == tank.Pos.X && dodge.SafePos.Y == tank.Pos.Y {
					noForward = append(noForward, tank.Id)
					travellingTank[tank.Id] = &objective.Target
				} else {
					travellingTank[tank.Id] = &dodge.SafePos
				}
			} else {
				travellingTank[tank.Id] = &objective.Target
			}
		} else {
			movement[tank.Id] = objective.Action
		}
	}
	self.traveller.Search(travellingTank, state, movement)
	for _, tankId := range noForward {
		action, _ := movement[tankId]
		if action == ActionMove {
			movement[tankId] = ActionStay
		}
	}
	return movement
}

func (self *Player) End(state *GameState) {
	self.tactics.End(state)
}
