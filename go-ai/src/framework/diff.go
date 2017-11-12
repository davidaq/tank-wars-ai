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
		if status {
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

func searchForest(preState *GameState, state *GameState, ret *DiffResult, watchList *ObservationList) {
	bulletSpeed := state.Params.BulletSpeed
	tankSpeed := state.Params.TankSpeed
	terain := state.Terain.Data

	tempList := ObservationList {
		tank: make(map[string]Tank),
		bullet: make(map[string]Bullet),
	}

	// 检查观察列表中子弹的状态
	if watchList.bullet != nil {
		for k,v := range(watchList.bullet) {
			live := false
			var curBullet Bullet 
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
				forestRange := caculateForestRange(terain, v.Pos)
				// 子弹消失
				if forestRange == 0 {
					// 草丛边缘，即将离开草丛
					delete(watchList.bullet, k)
				} else if forestRange < bulletSpeed {
					// 有干扰项障碍物/坦克
					maybeHit := false
					var chance int

					switch v.Pos.Direction {
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
							chance = bulletSpeed+1
						} else {
							chance = bulletSpeed
						}
						for i:=1;i<=forestRange;i++ {
							tempPos := Position {
								X: v.Pos.X,
								Y: v.Pos.Y-i,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = float64(1./chance)
						}
					case DirectionLeft:
						for i:=forestRange+1;i<=bulletSpeed;i++ {
							if terain[v.Pos.Y][v.Pos.X-i] == 1 {
								maybeHit = true
								break
							}
						}
						if maybeHit {
							chance = bulletSpeed+1
						} else {
							chance = bulletSpeed
						}
						for i:=1;i<=forestRange;i++ {
							tempPos := Position {
								X: v.Pos.X-i,
								Y: v.Pos.Y,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = float64(1./chance)
						}
					case DirectionDown:
						for i:=forestRange+1;i<=bulletSpeed;i++ {
							if terain[v.Pos.Y+i][v.Pos.X] == 1 {
								maybeHit = true
								break
							}
						}
						if maybeHit {
							chance = bulletSpeed+1
						} else {
							chance = bulletSpeed
						}
						for i:=1;i<=forestRange;i++ {
							tempPos := Position {
								X: v.Pos.X,
								Y: v.Pos.Y+i,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = float64(1./chance)
						}
					case DirectionRight:
						for i:=forestRange+1;i<=bulletSpeed;i++ {
							if terain[v.Pos.Y][v.Pos.X+i] == 1 {
								maybeHit = true
								break
							}
						}
						if maybeHit {
							chance = bulletSpeed+1
						} else {
							chance = bulletSpeed
						}
						for i:=1;i<=forestRange;i++ {
							tempPos := Position {
								X: v.Pos.X+i,
								Y: v.Pos.Y,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = float64(1./chance)
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
							ret.ForestThreat[tempPos] = float64(1./bulletSpeed)
						}
					case DirectionLeft:
						for i:=1;i<=bulletSpeed;i++ {
							tempPos := Position {
								X: v.Pos.X-i,
								Y: v.Pos.Y,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = float64(1./bulletSpeed)
						}
					case DirectionDown:
						for i:=1;i<=bulletSpeed;i++ {
							tempPos := Position {
								X: v.Pos.X,
								Y: v.Pos.Y+i,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = float64(1./bulletSpeed)
						}
					case DirectionRight:
						for i:=1;i<=bulletSpeed;i++ {
							tempPos := Position {
								X: v.Pos.X+i,
								Y: v.Pos.Y,
								Direction: v.Pos.Direction,
							}
							ret.ForestThreat[tempPos] = float64(1./bulletSpeed)
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
		step := 0
		
		switch curItemPos.Direction {
		case DirectionUp:
			if curItemPos.Y - bulletSpeed < 0 {
				step = 0
			} else {
				step = bulletSpeed
			}
			for j:= 1; j<=step; j++ {
				if terain[curItemPos.Y-j][curItemPos.X] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionLeft:
			if curItemPos.X - bulletSpeed < 0 {
				step = 0
			} else {
				step = bulletSpeed
			}
			for j:= 1; j<=step; j++ {
				if terain[curItemPos.Y][curItemPos.X-j] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionDown:
			if curItemPos.Y + bulletSpeed > state.Terain.Height {
				step = state.Terain.Height - curItemPos.Y
			} else {
				step = bulletSpeed
			}
			for j:= 1; j<=step; j++ {
				if terain[curItemPos.Y+j-1][curItemPos.X] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		case DirectionRight:
			if curItemPos.X + bulletSpeed > state.Terain.Width {
				step = state.Terain.Width - curItemPos.X
			} else {
				step = bulletSpeed
			}
			for j:= 1; j<=step; j++ {
				if terain[curItemPos.Y][curItemPos.X+j-1] == 2 {
					tempList.bullet[state.MyBullet[i].Id] = state.MyBullet[i]
					break
				}
			}
		}
	}

	// 检测地方可能进入草丛的坦克，进入监视名单
	for i:= 0; i<len(state.EnemyTank); i++ {
		curItemPos := state.EnemyTank[i].Pos
		step := 0

		switch curItemPos.Direction {
		case DirectionUp:
			if curItemPos.Y - tankSpeed < 0 {
				step = 0
			} else {
				step = tankSpeed
			}
			for j:= 1; j<=step; j++ {
				if terain[curItemPos.Y-j][curItemPos.X] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionLeft:
			if curItemPos.X - tankSpeed < 0 {
				step = 0
			} else {
				step = tankSpeed
			}
			for j:= 1; j<=step; j++ {
				if terain[curItemPos.Y][curItemPos.X-j] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionDown:
			if curItemPos.Y + tankSpeed > state.Terain.Height {
				step = state.Terain.Height - curItemPos.Y
			} else {
				step = tankSpeed
			}
			for j:= 1; j<=step; j++ {
				if terain[curItemPos.Y+j-1][curItemPos.X] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		case DirectionRight:
			if curItemPos.X + tankSpeed > state.Terain.Width {
				step = state.Terain.Width - curItemPos.X
			} else {
				step = tankSpeed
			}
			for j:= 1; j<=step; j++ {
				if terain[curItemPos.Y][curItemPos.X+j-1] == 2 {
					tempList.tank[state.EnemyTank[i].Id] = state.EnemyTank[i]
					break
				}
			}
		}
	}

	// return ret
}

func (self *Diff) Compare(newState *GameState) DiffResult {
	ret := DiffResult {
		ForestThreat: make(map[Position]float64),
	}
	if self.watchList == nil {
		self.watchList = &ObservationList {
			tank: make(map[string]Tank),
			bullet: make(map[string]Bullet),
		}
	}
	if self.prevState != nil {
		// TODO
		searchForest(self.prevState, newState, &ret, self.watchList)
	}
	self.prevState = newState
	fmt.Println(ret)
	return ret
}
