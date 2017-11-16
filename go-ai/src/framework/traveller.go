// 行走寻路行动子系统
package framework

import (
	"lib/go-astar"
	"sync"
	"math/rand"
)

type PathCache struct {
	path []Position
	target Position
	expect *Position
}

type Traveller struct {
	astar astar.AStar
	cache map[string]*PathCache
	collide map[string]int
}

func NewTraveller() *Traveller {
	inst := &Traveller {
		astar: nil,
		cache: make(map[string]*PathCache),
		collide: make(map[string]int),
	}
	return inst
}

func (self *Traveller) CollidedTankInForest(state *GameState) []Position {
	var candidate []Position
	myTankPos := make(map[Position]bool)
	for _, tank := range state.MyTank {
		p := Position {
			X: tank.Pos.X,
			Y: tank.Pos.Y,
		}
		myTankPos[p] = true
		if cache, ok := self.cache[tank.Id]; ok {
			from := &tank.Pos
			if cache.expect != nil && (cache.expect.Y != from.Y || cache.expect.X != from.X) {
				pos := Position {
					X: from.X + sign(cache.expect.X - from.X),
					Y: from.Y + sign(cache.expect.Y - from.Y),
				}
				if state.Terain.Get(pos.X, pos.Y) == 2 {
					candidate = append(candidate, pos)
				}
			}
		}
	}
	var ret []Position
	for _, pos := range candidate {
		if !myTankPos[pos] {
			ret = append(ret, pos)
		}
	}
	return ret
}

func (self *Traveller) Search(travel map[string]*Position, state *GameState, movements map[string]int) {
	maxPathCalc := 9
	if self.astar == nil {
		self.astar = astar.NewAStar(state.Terain.Height, state.Terain.Width)
		for y := 0; y < state.Terain.Height; y++ {
			for x := 0; x < state.Terain.Width; x++ {
				switch state.Terain.Get(x, y) {
				case 1:
					self.astar.FillTile(astar.Point{ Col: x, Row: y }, -1)
				case 2:
					self.astar.FillTile(astar.Point{ Col: x, Row: y }, 1)
				}
			}
		}
	}
	waitchan := make(chan bool)
	var lock sync.Mutex
	occupy := make(map[Position]bool)
	a := self.astar.Clone()
	lw := state.Terain.Width
	for _, tank := range state.MyTank {
		w := 8 + lw * self.collide[tank.Id]
		a.FillTile(astar.Point{ Col: tank.Pos.X, Row: tank.Pos.Y }, w)
	}
	for _, tank := range state.EnemyTank {
		a.FillTile(astar.Point{ Col: tank.Pos.X, Row: tank.Pos.Y }, lw)
	}
	lock.Lock()
	var myTanks []*Tank
	for _, tank := range state.MyTank {
		if _, exists := travel[tank.Id]; exists {
			t := tank
			myTanks = append(myTanks, &t)
		} else {
			occupy[Position{ Y: tank.Pos.Y, X: tank.Pos.X }] = true
		}
	}
	if len(myTanks) > maxPathCalc {
		full := myTanks
		myTanks = make([]*Tank, len(full))
		for i, j := range rand.Perm(len(myTanks)) {
			myTanks[i] = full[j]
		}
	}
	for _, tank := range myTanks {
		id := tank.Id
		from := tank.Pos
		to := *travel[tank.Id]
		go (func () {
			nextPoint := to
			lock.Lock()
			cache, hasCache := self.cache[id]
			if !hasCache {
				cache = &PathCache {}
				self.cache[id] = cache
			}
			lock.Unlock()
			if from.X != to.X || from.Y != to.Y {
				if cache.target.X != to.X || cache.target.Y != to.Y {
					cache.path = nil
				}
				cache.target = to
				if cache.expect != nil {
					lock.Lock()
					collide := self.collide[tank.Id]
					if cache.expect.Y != from.Y || cache.expect.X != from.X {
						self.collide[tank.Id] = collide + 10
						if abs(from.X - to.X) + abs(from.Y - to.Y) > state.Params.TankSpeed {
							cache.path = nil
						} else {
							cache.path = []Position { *cache.expect }
						}
					} else if collide > 0 {
						self.collide[tank.Id] = collide - 1
					}
					lock.Unlock()
				}
				for len(cache.path) > 0 {
					p := cache.path[0]
					if p.X == from.X && p.Y == from.Y {
						cache.path = cache.path[1:]
					} else {
						break
					}
				}
				if len(cache.path) == 0 {
					lock.Lock()
					allowCalc := false
					if maxPathCalc > 0 {
						maxPathCalc--
						allowCalc = true
					}
					lock.Unlock()
					if allowCalc {
						cache.path = self.path(a, from, to, state.Params.TankSpeed, state.Terain)
						for len(cache.path) > 0 {
							p := cache.path[0]
							if p.X == from.X && p.Y == from.Y {
								cache.path = cache.path[1:]
							} else {
								break
							}
						}
					}
				}
				if len(cache.path) == 0 {
					nextPoint = to
				} else {
					nextPoint = cache.path[0]
				}
			}
			action := toAction(from, nextPoint)
			cache.expect = nil
			lock.Lock()
			if action == ActionMove {
				p := Position { Y: from.Y, X: from.X }
				if nextPoint.Y > from.Y {
					p.Y++
				} else if nextPoint.Y < from.Y {
					p.Y--
				} else if nextPoint.X > from.X {
					p.X++
				} else if nextPoint.X < from.X {
					p.X--
				}
				if _, exists := occupy[p]; exists {
					action = ActionStay
					p = Position { Y: from.Y, X: from.X }
				} else {
					cache.expect = &nextPoint
				}
				occupy[p] = true
			} else {
				p := Position { Y: from.Y, X: from.X }
				occupy[p] = true
			}
			movements[id] = action
			lock.Unlock()
			waitchan <- true
		})()
	}
	lock.Unlock()
	for _, _ = range myTanks {
		_ = <- waitchan
	}
}

func (self *Traveller) path(a astar.AStar, source Position, target Position, movelen int, terain *Terain) []Position {
	p2p := astar.NewPointToPoint()

	sourcePoint := []astar.Point{ astar.Point{ Row: source.Y, Col: source.X } }
	targetPoint := []astar.Point{ astar.Point{ Row: target.Y, Col: target.X } }

	p := a.FindPath(p2p, targetPoint, sourcePoint, movelen, source.Direction)
	
	var ret []Position
	for p != nil {
		ret = append(ret, Position {
			X: p.Col,
			Y: p.Row,
		})
		p = p.Parent
	}
	c := len(ret)
	for i, n := 0, c / 2; i < n; i++ {
		j := c - i - 1
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret
}

func abs (val int) int {
	if val < 0 {
		return -val
	} else {
		return val
	}
}

func sign (val int) int {
	if val > 0 {
		return 1
	} else if val < 0 {
		return -1
	} else {
		return 0
	}
}

func toAction (source Position, target Position) int {
	targetDirection := DirectionNone
	if source.X < target.X {
		targetDirection = DirectionRight
	} else if source.X > target.X {
		targetDirection = DirectionLeft
	} else if source.Y < target.Y {
		targetDirection = DirectionDown
	} else if source.Y > target.Y {
		targetDirection = DirectionUp
	} else {
		targetDirection = target.Direction
		if targetDirection == DirectionNone || source.Direction == target.Direction {
			return ActionStay	
		}
	}
	if targetDirection == source.Direction {
		return ActionMove
	}
	switch ((targetDirection - 1) - (source.Direction - 1) + 4) % 4 {
	case 1:
		return ActionLeft
	case 2:
		return ActionBack
	default:
		return ActionRight
	}
}
