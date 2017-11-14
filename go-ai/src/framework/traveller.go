// 行走寻路行动子系统
package framework

import (
	"lib/go-astar"
	"sync"
)

type PathCache struct {
	path []Position
	target Position
	expect *Position
}

type Traveller struct {
	astar astar.AStar
	cache map[string]*PathCache
	colide map[string]int
}

func NewTraveller() *Traveller {
	inst := &Traveller {
		astar: nil,
		cache: make(map[string]*PathCache),
		colide: make(map[string]int),
	}
	return inst
}

func (self *Traveller) Search(travel map[string]*Position, state *GameState, movements map[string]int) {
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
	waits := 0
	waitchan := make(chan bool)
	var lock sync.Mutex
	occupy := make(map[Position]bool)
	a := self.astar.Clone()
	lw := state.Terain.Width
	for _, tank := range state.MyTank {
		w := 8 + lw * self.colide[tank.Id]
		a.FillTile(astar.Point{ Col: tank.Pos.X, Row: tank.Pos.Y }, w)
	}
	for _, tank := range state.EnemyTank {
		a.FillTile(astar.Point{ Col: tank.Pos.X, Row: tank.Pos.Y }, lw)
	}
	lock.Lock()
	for _, tank := range state.MyTank {
		if target, exists := travel[tank.Id]; exists {
			waits += 1
			id := tank.Id
			from := tank.Pos
			to := *target
			go (func () {
				nextPoint := from
				lock.Lock()
				cache, hasCache := self.cache[id]
				lock.Unlock()
				if from.X != to.X || from.Y != to.Y {
					if !hasCache || cache.target.X != to.X || cache.target.Y != to.Y {
						cache = &PathCache {
							path: self.path(a, from, to, state.Params.TankSpeed, state.Terain),
							target: to,
						}
						hasCache = true
						lock.Lock()
						self.cache[id] = cache
						lock.Unlock()
					}
					if cache.expect != nil {
						colide := self.colide[tank.Id]
						if cache.expect.Y != from.Y || cache.expect.X != from.X {
							self.colide[tank.Id] = colide + 10
							if abs(from.X - to.X) + abs(from.Y - to.Y) > state.Params.TankSpeed {
								cache.path = nil
							} else {
								cache.path = []Position { *cache.expect }
							}
						} else if colide > 0 {
							self.colide[tank.Id] = colide - 1
						}
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
						cache.path = self.path(a, from, to, state.Params.TankSpeed, state.Terain)
					}
					if len(cache.path) == 0 {
						nextPoint = to
					} else {
						nextPoint = cache.path[0]
					}
				}
				action := toAction(from, nextPoint)
				if hasCache {
					cache.expect = nil
				}
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
					if _, exists = occupy[p]; exists {
						action = ActionStay
						p = Position { Y: from.Y, X: from.X }
					} else if hasCache {
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
		} else {
			occupy[Position{ Y: tank.Pos.Y, X: tank.Pos.X }] = true
		}
	}
	lock.Unlock()
	for i := 0; i < waits; i++ {
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
