// 行走寻路行动子系统
package framework

import (
	"lib/go-astar"
	"sync"
	"fmt"
)

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

type PathCache struct {
	path *astar.PathPoint
	target Position
	expect *astar.Point
}

type Traveller struct {
	astar astar.AStar
	cache map[string]*PathCache
}

func NewTraveller() *Traveller {
	inst := &Traveller {
		astar: nil,
		cache: make(map[string]*PathCache),
	}
	return inst
}

func (self *Traveller) Search(travel map[string]*Position, state *GameState, movements map[string]int) {
	if self.astar == nil {
		self.astar = astar.NewAStar(state.Terain.Height, state.Terain.Width)
		for y := 0; y < state.Terain.Height; y++ {
			for x := 0; x < state.Terain.Width; x++ {
				if state.Terain.Get(x, y) != 0 {
					self.astar.FillTile(astar.Point{ Col: x, Row :y }, -1)
				}
			}
		}
	}
	waits := 0
	waitchan := make(chan bool)
	var lock sync.Mutex
	occupy := make(map[astar.Point]bool)
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
							path: self.path(self.astar.Clone(), from, to, state.Params.TankSpeed, &state.Terain),
							target: to,
						}
						hasCache = true
						lock.Lock()
						self.cache[id] = cache
						lock.Unlock()
					}
					if cache.expect != nil && (cache.expect.Row != from.Y || cache.expect.Col != from.X) {
						cache.path = nil
					}
					for cache.path != nil {
						if cache.path.Col == from.X && cache.path.Row == from.Y {
							cache.path = cache.path.Parent
						} else {
							break
						}
					}
					if cache.path == nil {
						fmt.Println("Recalc")
						cache.path = self.path(self.astar.Clone(), from, to, state.Params.TankSpeed, &state.Terain)
					}
					nextPoint.X = cache.path.Col
					nextPoint.Y = cache.path.Row
				}
				if id == state.MyTank[0].Id {
					fmt.Println(from, nextPoint)
				}
				action := toAction(from, nextPoint)
				if hasCache {
					cache.expect = nil
				}
				lock.Lock()
				if action == ActionMove {
					p := astar.Point{ Row: nextPoint.Y, Col: nextPoint.X }
					if _, exists = occupy[p]; exists {
						action = ActionStay
						p = astar.Point{ Row: from.Y, Col: from.X }
					} else if hasCache {
						cache.expect = &p
					}
					occupy[p] = true
				} else {
					p := astar.Point{ Row: from.Y, Col: from.X }
					occupy[p] = true
				}
				movements[id] = action
				lock.Unlock()
				waitchan <- true
			})()
		} else {
			occupy[astar.Point{ Row: tank.Pos.Y, Col: tank.Pos.X }] = true
		}
	}
	lock.Unlock()
	for i := 0; i < waits; i++ {
		_ = <- waitchan
	}
}

func (self *Traveller) path(a astar.AStar, source Position, target Position, movelen int, terain *Terain) *astar.PathPoint {
	p2p := astar.NewPointToPoint()
	// cols := terain.Width
	// rows := terain.Height
	// switch target.Direction {
	// case DirectionUp: 
	// 	tmpI := target.Y + 1
	// 	for ; tmpI < target.Y + 5; tmpI++ {
	// 		if tmpI == cols {
	// 			break;
	// 		}
	// 		if terain.Get(target.X, tmpI) == 1 {
	// 			break;
	// 		}
	// 	}
	// 	if source.X == target.X && source.Y > target.Y {
	// 		break;
	// 	}
	// 	target.Y = tmpI - 1
	// 	break;
	// case DirectionLeft:
	// 	tmpI := target.X + 1
	// 	for ;tmpI < target.X + 5; tmpI++ {
	// 		if tmpI == rows {
	// 			break;
	// 		}
	// 		if terain.Get(tmpI, target.Y) == 1 {
	// 			break;
	// 		}
	// 	}
	// 	if source.Y == target.Y && source.X > target.X {
	// 		break;
	// 	}
	// 	target.X = tmpI - 1
	// 	break;
	// case DirectionDown:
	// 	tmpI := target.Y - 1
	// 	for ;tmpI > target.Y - 5; tmpI-- {
	// 		if tmpI == -1 {
	// 			break;
	// 		}
	// 		if terain.Get(target.X, tmpI) == 1 {
	// 			break;
	// 		}
	// 	}
	// 	if source.X == target.X && source.Y < target.Y {
	// 		break;
	// 	}
	// 	target.Y = tmpI + 1
	// 	break;
	// case DirectionRight:
	// 	tmpI := target.X - 1
	// 	for ;tmpI > target.X - 5; tmpI-- {
	// 		if tmpI == -1 {
	// 			break;
	// 		}
	// 		if terain.Get(tmpI, target.Y) == 1 {
	// 			break;
	// 		}
	// 	}
	// 	if source.Y == target.Y && source.X < target.X {
	// 		break;
	// 	}
	// 	target.X = tmpI + 1
	// 	break;
	// }

	sourcePoint := []astar.Point{ astar.Point{ Row: source.Y, Col: source.X } }
	targetPoint := []astar.Point{ astar.Point{ Row: target.Y, Col: target.X } }

	return a.FindPath(p2p, sourcePoint, targetPoint, movelen)
}
