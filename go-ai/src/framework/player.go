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
			group := reactors {
				dodger: NewDodger(),
				attacker: NewAttacker(),
				traveller: NewTraveller(),
			}
			group.dodger.Init(state, tank.Id)
			group.attacker.Init(state, tank.Id)
			group.traveller.Init(state, tank.Id)
			self.reactors[tank.Id] = group
		}
	}
	self.tactics.Plan(state, &self.objectives)
	
	movement := make(map[string]int)
	for _, tank  := range state.MyTank {
		objective := self.objectives[tank.Id]
		reactors := self.reactors[tank.Id]
		suggestion := Suggestion {
			dodge: reactors.dodger.Suggest(state, &objective),
			attack: reactors.attacker.Suggest(state, &objective),
			travel: reactors.traveller.Suggest(state, &objective),
		}
		movement[tank.Id] = self.tactics.Decide(tank.Id, suggestion)
	}
	return movement
}

func (self *Player) Reset() {
}
