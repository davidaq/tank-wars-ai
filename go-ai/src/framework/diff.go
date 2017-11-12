// 检查当前state和之前state，判断草丛威胁
package framework

import (
	"fmt"
)

type Diff struct {
	prevState *GameState
	watchList *ObservationList
}

type ObservationList struct {
	tank map[string]Tank
	bullet map[string]Bullet
}

func NewDiff() *Diff {
	return &Diff {
		prevState: nil,
		watchList: nil,
	}
}

func caculateForestRange(terain [][]int, pos Position) int {
	status := true
	count := 0
	for {
		if !status {
			break
		} else {
			count++
			switch pos.Direction {
			case DirectionUp:
				if terain[pos.Y-count][pos.X] == 2 {
					count++
				} else {
					status = false
				}
			case DirectionLeft:
				if terain[pos.Y][pos.X-count] == 2 {
					count++
				} else {
					status = false
				}
			case DirectionDown:
				if terain[pos.Y+count][pos.X] == 2 {
					count++
				} else {
					status = false
				}
			case DirectionRight:
				if terain[pos.Y][pos.X+count] == 2 {
					count++
				} else {
					status = false
				}
			}
		}
	}
	return count
}

func searchForest(preState GameState, state GameState, ret DiffResult, watchList ObservationList) DiffResult {
	bulletSpeed := state.params.bulletSpeed
	tankSpeed := state.params.tankSpeed
	terain := state.Terain.Data
	events := state.Events

	tempList := ObservationList {
		tank: make(map[string]Position),
		bullet: make(map[string]Position),
	}

	// 检查观察列表中子弹的状态
	if watchList.bullet != nil {
		for k,v := range(watchList.bullet) {
			live := false
			for i:=0;i<len(state.MyBullet);i++ {
				curItemPos := state.MyBullet[i]
				if curItemPos.Id == k {
					live = true
				}
			}
			if live && terain[v.Pos.Y][v.Pos.X] != 2 {
				// 子弹存活，并飞出草丛
				delete(watchList.bullet, k)
			} else {
				forestRange := caculateForestRange(terain, v)
				maybeHit := false
	
				// 前方存在可能击中的坦克
				for i:=0;i<len(events);i++ {
					if events[i].typ == "me-hit-enemy" && events[i].From == v.From {
						maybeHit = true
						delete(watchList.bullet, k)
						break
					}
				}

				if !maybeHit {
					// 子弹消失
					if forestRange == 0 {
						// 草丛边缘，即将离开草丛
						delete(watchList.bullet, k)
					} else if forestRange < bulletSpeed {
						// 有干扰项障碍物/坦克
						switch v.Direction {
						case DirectionUp:
							// 判断是否撞在障碍物上
							for i:=forestRange+1;i<=bulletSpeed;i++ {
								// 前方存在可能击中的障碍物
								if terain[v.Pos.Y-i][v.Pos.X] == 1 {
									maybeHit = true
									break
								}
							}
							if maybeHit {
								chance := bulletSpeed+1
							} else {
								chance := bulletSpeed
							}
							for i:=1;i<=forestRange;i++ {
								tempPos := Position {
									X: v.Pos.X,
									Y: v.Pos.Y-i,
									Direction: v.Direction,
								}
								ret.ForestThreat[tempPos] = 1/chance
							}
						case DirectionLeft:
							for i:=forestRange+1;i<=bulletSpeed;i++ {
								if terain[v.Pos.Y][v.Pos.X-i] == 1 {
									maybeHit = true
									break
								}
							}
							if maybeHit {
								chance := bulletSpeed+1
							} else {
								chance := bulletSpeed
							}
							for i:=1;i<=forestRange;i++ {
								tempPos := Position {
									X: v.Pos.X-i,
									Y: v.Pos.Y,
									Direction: v.Direction,
								}
								ret.ForestThreat[tempPos] = 1/chance
							}
						case DirectionDown:
							for i:=forestRange+1;i<=bulletSpeed;i++ {
								if terain[v.Pos.Y+i][v.Pos.X] == 1 {
									maybeHit = true
									break
								}
							}
							if maybeHit {
								chance := bulletSpeed+1
							} else {
								chance := bulletSpeed
							}
							for i:=1;i<=forestRange;i++ {
								tempPos := Position {
									X: v.Pos.X,
									Y: v.Pos.Y+i,
									Direction: v.Direction,
								}
								ret.ForestThreat[tempPos] = 1/chance
							}
						case DirectionRight:
							for i:=forestRange+1;i<=bulletSpeed;i++ {
								if terain[v.Pos.Y][v.Pos.X+i] == 1 {
									maybeHit = true
									break
								}
							}
							if maybeHit {
								chance := bulletSpeed+1
							} else {
								chance := bulletSpeed
							}
							for i:=1;i<=forestRange;i++ {
								tempPos := Position {
									X: v.Pos.X+i,
									Y: v.Pos.Y,
									Direction: v.Direction,
								}
								ret.ForestThreat[tempPos] = 1/chance
							}
						}
					} else {
						// 纯草丛击中
						switch v.Pos.Direction {
						case DirectionUp:
							for i:=1;i<=bulletSpeed;i++ {
								tempPos := Position {
									X: v.Pos.X,
									Y: v.Pos.Y-i,
									Direction: v.Pos.Direction,
								}
								ret.ForestThreat[tempPos] = 1/bulletSpeed
							}
						case DirectionLeft:
							for i:=1;i<=bulletSpeed;i++ {
								tempPos := Position {
									X: v.Pos.X-i,
									Y: v.Pos.Y,
									Direction: v.Pos.Direction,
								}
								ret.ForestThreat[tempPos] = 1/bulletSpeed
							}
						case DirectionDown:
							for i:=1;i<=bulletSpeed;i++ {
								tempPos := Position {
									X: v.Pos.X,
									Y: v.Pos.Y+i,
									Direction: v.Pos.Direction,
								}
								ret.ForestThreat[tempPos] = 1/bulletSpeed
							}
						case DirectionRight:
							for i:=1;i<=bulletSpeed;i++ {
								tempPos := Position {
									X: v.Pos.X+i,
									Y: v.Pos.Y,
									Direction: v.Pos.Direction,
								}
								ret.ForestThreat[tempPos] = 1/bulletSpeed
							}
						}
					}
				}

				
			}
		}
	}

	// 检查观察列表中坦克的状态
	if watchList.tank != nil {
		for k,v := range(watchList.tank) {
			disappear := false
			increaseHP := 0

			// 判断坦克是否消失
			for i:=0;i<len(state.EnemyTank);i++ {
				if k==state.EnemyTank[i].Id {
					delete(watchList.tank, k)
					break
				}
			}

			if disappear {
				// 判断坦克是否被击杀
				for i:=0;i<len(events);i++ {
					if events[i].typ == "me-hit-enemy" || events[i].typ == "enemy-hit-enemy" {
						if v == events[i].target {
							increaseHP++
						}
					}
				}
				if v.Hp <= increaseHP {
					delete(watchList.tank, k)
				} else {
					// 草丛无障碍物
					switch v.Pos.Direction {
					case DirectionUp:
						tempPos := Position {
							X: v.Pos.X,
							Y: v.Pos.Y-tankSpeed,
							Direction: v.Pos.Direction,
						}
						ret.ForestThreat[tempPos] = 1
					case DirectionLeft:
						tempPos := Position {
							X: v.Pos.X-tankSpeed,
							Y: v.Pos.Y,
							Direction: v.Pos.Direction,
						}
						ret.ForestThreat[tempPos] = 1
					case DirectionDown:
						tempPos := Position {
							X: v.Pos.X,
							Y: v.Pos.Y+tankSpeed,
							Direction: v.Pos.Direction,
						}
						ret.ForestThreat[tempPos] = 1
					case DirectionRight:
						tempPos := Position {
							X: v.Pos.X+tankSpeed,
							Y: v.Pos.Y,
							Direction: v.Pos.Direction,
						}
						ret.ForestThreat[tempPos] = 1
					}
				}
			}
		}
	}

	// 检测我方即将进入草丛的子弹，列入监视名单
	for i:= 0; i<len(state.MyBullet); i++ {
		curItemPos := state.MyBullet[i].Pos

		switch curItemPos.Direction {
		case DirectionUp:
			for j:= 1; j<=tankSpeed; j++ {
				if state.Terain.Data[curItemPos.Y-j][curtItemPos.X] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionLeft:
			for j:= 1; j<=tankSpeed; j++ {
				if state.Terain.Data[curItemPos.Y][curtItemPos.X-j] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionDown:
			for j:= 1; j<=tankSpeed; j++ {
				if state.Terain.Data[curItemPos.Y+j][curtItemPos.X] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionRight:
			for j:= 1; j<=tankSpeed; j++ {
				if state.Terain.Data[curItemPos.Y][curtItemPos.X+j] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		}
	}

	// 检测地方可能进入草丛的坦克，进入监视名单
	for i:=0; i<len(state.EnemyTank); i++ {
		curItemPos := state.EnemyTank[i].Pos
		
		switch curItemPos.Direction {
		case DirectionUp:
			for j:= 1; j<=bulletSpeed; j++ {
				if state.Terain.Data[curItemPos.Y-j][curtItemPos.X] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionLeft:
			for j:= 1; j<=bulletSpeed; j++ {
				if state.Terain.Data[curItemPos.Y][curtItemPos.X-j] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionDown:
			for j:= 1; j<=bulletSpeed; j++ {
				if state.Terain.Data[curItemPos.Y+j][curtItemPos.X] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionRight:
			for j:= 1; j<=bulletSpeed; j++ {
				if state.Terain.Data[curItemPos.Y][curtItemPos.X+j] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		}
	}

	return ret
}

func searchForest(preState GameState, state GameState, ret DiffResult, watchList ObservationList) DiffResult {
	bulletSpeed := state.params.bulletSpeed
	tankSpeed := state.params.tankSpeed
	terain := state.Terain.Data

	tempList := ObservationList {
		tank: make(map[string]Position),
		bullet: make(map[string]Position),
	}

	// 检查观察列表中子弹的状态
	if watchList.bullet != nil {
		for k,v := range(watchList.bullet) {
			live := false
			curBullet := nil
			for i:=0;i<len(state.MyBullet);i++ {
				if state.MyBullet[i].Id == k {
					live = true
					curBullet = state.MyBullet[i]
				}
			}
			if live {
				if terain[v.Pos.Y][v.Pos.X] != 2 {
					// 子弹存活，并飞出草丛
					delete(watchList.bullet,k)
				} else {
					// 子弹存活，仍在草丛
					watchList.bullet[k] = curBullet
				}
			} else {
				forestRange := caculateForestRange(terain, v)
				
				// 子弹消失
				if forestRange == 0 {
					// 草丛边缘，即将离开草丛
					delete(watchList.bullet, k)
				} else if forestRange < bulletSpeed {
					// 有干扰项障碍物/坦克
					switch v.Direction {
					case DirectionUp:
						// 判断是否撞在障碍物上
						for i:=forestRange+1;i<=bulletSpeed;i++ {
							// 前方存在可能击中的障碍物
							if terain[v.Pos.Y-i][v.Pos.X] == 1 {
								maybeHit = true
								break
							}
						}
						if maybeHit {
							chance := bulletSpeed+1
						} else {
							chance := bulletSpeed
						}
						for i:=1;i<=forestRange;i++ {
							tempPos := Position {
								X: v.Pos.X,
								Y: v.Pos.Y-i,
								Direction: v.Direction,
							}
							ret.ForestThreat[tempPos] = 1/chance
						}
					case DirectionLeft:
						for i:=forestRange+1;i<=bulletSpeed;i++ {
							if terain[v.Pos.Y][v.Pos.X-i] == 1 {
								maybeHit = true
								break
							}
						}
						if maybeHit {
							chance := bulletSpeed+1
						} else {
							chance := bulletSpeed
						}
						for i:=1;i<=forestRange;i++ {
							tempPos := Position {
								X: v.Pos.X-i,
								Y: v.Pos.Y,
								Direction: v.Direction,
							}
							ret.ForestThreat[tempPos] = 1/chance
						}
					case DirectionDown:
						for i:=forestRange+1;i<=bulletSpeed;i++ {
							if terain[v.Pos.Y+i][v.Pos.X] == 1 {
								maybeHit = true
								break
							}
						}
						if maybeHit {
							chance := bulletSpeed+1
						} else {
							chance := bulletSpeed
						}
						for i:=1;i<=forestRange;i++ {
							tempPos := Position {
								X: v.Pos.X,
								Y: v.Pos.Y+i,
								Direction: v.Direction,
							}
							ret.ForestThreat[tempPos] = 1/chance
						}
					case DirectionRight:
						for i:=forestRange+1;i<=bulletSpeed;i++ {
							if terain[v.Pos.Y][v.Pos.X+i] == 1 {
								maybeHit = true
								break
							}
						}
						if maybeHit {
							chance := bulletSpeed+1
						} else {
							chance := bulletSpeed
						}
						for i:=1;i<=forestRange;i++ {
							tempPos := Position {
								X: v.Pos.X+i,
								Y: v.Pos.Y,
								Direction: v.Direction,
							}
							ret.ForestThreat[tempPos] = 1/chance
						}
					}
					delete(watchList.bullet, k)
				} else {
					// 纯草丛击中
					switch v.Pos.Direction {
					case DirectionUp:
						for i:=1;i<=bulletSpeed;i++ {
							tempPos := Position {
								X: v.Pos.X,
								Y: v.Pos.Y-i,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = 1/bulletSpeed
						}
					case DirectionLeft:
						for i:=1;i<=bulletSpeed;i++ {
							tempPos := Position {
								X: v.Pos.X-i,
								Y: v.Pos.Y,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = 1/bulletSpeed
						}
					case DirectionDown:
						for i:=1;i<=bulletSpeed;i++ {
							tempPos := Position {
								X: v.Pos.X,
								Y: v.Pos.Y+i,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = 1/bulletSpeed
						}
					case DirectionRight:
						for i:=1;i<=bulletSpeed;i++ {
							tempPos := Position {
								X: v.Pos.X+i,
								Y: v.Pos.Y,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = 1/bulletSpeed
						}
					}
					delete(watchList.bullet, k)
				}
			}
		}
	}

	// 检查观察列表中坦克的状态
	if watchList.tank != nil {
		for k,v := range(watchList.tank) {

			// 判断坦克是否消失
			for i:=0;i<len(state.EnemyTank);i++ {
				if k==state.EnemyTank[i].Id {
					delete(watchList.tank, k)
					break
				}
			}

			switch v.Pos.Direction {
			case DirectionUp:
				tempPos := Position {
					X: v.Pos.X,
					Y: v.Pos.Y-tankSpeed,
					Direction: v.Pos.Direction,
				}
				ret.ForestThreat[tempPos] = 1
			case DirectionLeft:
				tempPos := Position {
					X: v.Pos.X-tankSpeed,
					Y: v.Pos.Y,
					Direction: v.Pos.Direction,
				}
				ret.ForestThreat[tempPos] = 1
			case DirectionDown:
				tempPos := Position {
					X: v.Pos.X,
					Y: v.Pos.Y+tankSpeed,
					Direction: v.Pos.Direction,
				}
				ret.ForestThreat[tempPos] = 1
			case DirectionRight:
				tempPos := Position {
					X: v.Pos.X+tankSpeed,
					Y: v.Pos.Y,
					Direction: v.Pos.Direction,
				}
				ret.ForestThreat[tempPos] = 1
			}
			delete(watchList.tank, k)
		}
	}

	// 检测我方即将进入草丛的子弹，列入监视名单
	for i:=0; i<len(state.MyBullet); i++ {
		curItemPos := state.MyBullet[i].Pos
		
		switch curItemPos.Direction {
		case DirectionUp:
			for j:= 1; j<=bulletSpeed; j++ {
				if state.Terain.Data[curItemPos.Y-j][curtItemPos.X] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionLeft:
			for j:= 1; j<=bulletSpeed; j++ {
				if state.Terain.Data[curItemPos.Y][curtItemPos.X-j] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionDown:
			for j:= 1; j<=bulletSpeed; j++ {
				if state.Terain.Data[curItemPos.Y+j][curtItemPos.X] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionRight:
			for j:= 1; j<=bulletSpeed; j++ {
				if state.Terain.Data[curItemPos.Y][curtItemPos.X+j] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		}
	}

	// 检测地方可能进入草丛的坦克，进入监视名单
	for i:= 0; i<len(state.EnemyTank); i++ {
		curItemPos := state.EnemyTank[i].Pos

		switch curItemPos.Direction {
		case DirectionUp:
			for j:= 1; j<=tankSpeed; j++ {
				if state.Terain.Data[curItemPos.Y-j][curtItemPos.X] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionLeft:
			for j:= 1; j<=tankSpeed; j++ {
				if state.Terain.Data[curItemPos.Y][curtItemPos.X-j] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionDown:
			for j:= 1; j<=tankSpeed; j++ {
				if state.Terain.Data[curItemPos.Y+j][curtItemPos.X] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionRight:
			for j:= 1; j<=tankSpeed; j++ {
				if state.Terain.Data[curItemPos.Y][curtItemPos.X+j] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		}
	}
}

func (self *Diff) compare(newState *GameState) DiffResult {
	ret := DiffResult {
		ForestThreat: make(map[Position]float64),
	}
	res := nil
	if self.prevState != nil {
		// TODO
		res = searchForest(self.prevState, newState, ret, self.watchList)
		// ret.ForestThreat[Position { X: 0, Y: 0 }] = 1.
	}
	self.prevState = newState
	fmt.println(res)
	return res
}
