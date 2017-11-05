package framework

type reactors struct {
	dodger, attacker, traveller reactor
}

type Player struct {
	inited bool
	tactics Tactics
	objectives map[string]Objective
	reactors map[string]reactors
}

func NewPlayer(tactics Tactics) *Player {
	inst := &Player {
		tactics: tactics,
		objectives: make(map[string]Objective),
		inited: false,
		reactors: make(map[string]reactors),
	}
	return inst
}

func (self *Player) Play(state *GameState) map[string]int {
	if !self.inited {
		self.inited = true
		self.tactics.Init(state)
		for _, tank  := range state.MyTank {
			self.reactors[tank.Id] = reactors {
				dodger: NewDodger(),
				attacker: NewAttacker(),
				traveller: NewTraveller(),
			}
		}
	}
	self.tactics.Plan(state, &self.objectives)
	
	movement := make(map[string]int)
	for _, tank  := range state.MyTank {
		objective := self.objectives[tank.Id]
		reactors := self.reactors[tank.Id]
		suggestion := Suggestion {
			Dodge: reactors.dodger.Suggest(&tank, state, &objective),
			Attack: reactors.attacker.Suggest(&tank, state, &objective),
			Travel: reactors.traveller.Suggest(&tank, state, &objective),
		}
		movement[tank.Id] = self.tactics.Decide(&tank, state, suggestion)
	}
	return movement
}

func (self *Player) Reset() {
}
