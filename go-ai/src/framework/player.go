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
	for _, tank := range state.MyTank {
		self.radar.Scan(&tank, state)
	}
	self.tactics.Plan(state, &self.objectives)
	
	movement := make(map[string]int)
	for _, tank  := range state.MyTank {
		objective := self.objectives[tank.Id]
		if objective.Action == ActionTravel {
			movement[tank.Id] = self.traveller.Search(&tank, state, &objective.Target)
		} else {
			movement[tank.Id] = objective.Action
		}
	}
	return movement
}

func (self *Player) Reset() {
}
