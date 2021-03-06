package framework

import (
	"fmt"
)

// 写死防止干蠢事

func BadCase(state *GameState, radar *RadarResult, movements map[string]int) {
	badCaseControlEnemy(state, radar, movements)
	badCaseDangerZone(state, radar, movements)
	badCaseShootSelf(state, radar, movements)
}

func badCaseControlEnemy(state *GameState, radar *RadarResult, movements map[string]int) {
	safe := make(map[string]bool)
	for _, tank  := range state.MyTank {
		safe[tank.Id] = true
	}
	for key, _ := range movements {
		if !safe[key] {
			delete(movements, key)
		}
	}
}

func badCaseShootSelf(state *GameState, radar *RadarResult, movements map[string]int) {
	noPass := make(map[Position]bool)
	ePos := make(map[Position]bool)
	for _, tank  := range state.MyTank {
		pos := tank.Pos.NoDirection()
		noPass[pos] = true
		action := movements[tank.Id]
		if action == ActionMove {
			for i := 0; i < state.Params.TankSpeed; i++ {
				nPos := pos.step(tank.Pos.Direction)
				if state.Terain.Get(nPos.X, nPos.Y) == 1 {
					break
				}
				noPass[nPos] = true
				pos = nPos
			}
		}
	}
	fmt.Println("shoot self", noPass)
	for _, etank  := range state.EnemyTank {
		ePos[etank.Pos] = true
	}
	tankloop: for _, tank  := range state.MyTank {
		action := movements[tank.Id]
		if action >= ActionFireUp && action <= ActionFireRight {
			dir := action - ActionFireUp + DirectionUp
			pos := tank.Pos.NoDirection().step(dir)
			if state.Terain.Get(pos.X, pos.Y) == 1 {
				continue tankloop
			}
			if ePos[pos] {
				continue tankloop
			}
			for i := 0; i <= state.Params.BulletSpeed * 2; i++ {
				if noPass[pos] {
					movements[tank.Id] = ActionStay
					fmt.Println("Self shoot prevented", tank.Id)
					continue tankloop
				}
				pos = pos.step(dir)
				if state.Terain.Get(pos.X, pos.Y) == 1 {
					continue tankloop
				}
			}
		}
	}
}

func badCaseDangerZone(state *GameState, radar *RadarResult, movements map[string]int) {
	dangerous := make(map[Position]bool)
	directions := []int { DirectionUp, DirectionLeft, DirectionDown, DirectionRight }
	vDirections := []int { DirectionUp, DirectionDown }
	hDirections := []int { DirectionLeft, DirectionRight }
	for _, eTank := range state.EnemyTank {
		dirloop: for _, dir := range directions {
			preferVertical := true
			if dir == DirectionDown || dir == DirectionUp {
				preferVertical = false
			}
			pos := eTank.Pos.NoDirection()
			for i := 0; i < 1 + state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
				if state.Terain.Get(pos.X, pos.Y) == 1 {
					continue dirloop
				}
			}
			for i := 0; i < state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
				if state.Terain.Get(pos.X, pos.Y) == 1 {
					continue dirloop
				}
				dangerous[pos] = preferVertical
			}
		}
	}
	for _, tank := range state.MyTank {
		noPass := make(map[Position]string)
		for _, oTank  := range state.MyTank {
			if oTank.Id != tank.Id {
				pos := oTank.Pos.NoDirection()
				noPass[pos] = oTank.Id
				action := movements[oTank.Id]
				if action == ActionMove {
					for i := 0; i < state.Params.TankSpeed; i++ {
						nPos := pos.step(oTank.Pos.Direction)
						if state.Terain.Get(nPos.X, nPos.Y) == 1 {
							break
						}
						noPass[nPos] = oTank.Id
						pos = nPos
					}
				}
			}
		}
		getPreferB := func (preferVertical bool, ignoreAlly bool) map[int]bool {
			preferDirection := make(map[int]bool)
			dirs := hDirections
			if preferVertical {
				dirs = vDirections
			}
			for _, dir := range dirs {
				nPos := tank.Pos.step(dir).NoDirection()
				blocked := false
				if state.Terain.Get(nPos.X, nPos.Y) == 1 {
					blocked = true
				}
				if ignoreAlly {
					if id, ok := noPass[nPos]; ok && id != tank.Id {
						blocked = true
					}
				}
				if !blocked {
					preferDirection[dir] = true
				}
			}
			return preferDirection
		}
		getPrefer := func (preferVertical bool) map[int]bool {
			ret := getPreferB(preferVertical, false)
			if len(ret) == 0 {
				ret = getPreferB(preferVertical, true)
			}
			return ret
		}
		if preferVertical, danger := dangerous[tank.Pos.NoDirection()]; danger {  // 位于危险区域
			preferDirection := getPrefer(preferVertical)
			if len(preferDirection) > 0 {
				oldMove := movements[tank.Id]
				fixMove(state, radar, movements, tank, preferDirection)
				fmt.Println("Fix Danger Zone", tank.Id, oldMove, movements[tank.Id])
			}
		} else if movements[tank.Id] == ActionMove {
			pos := tank.Pos.NoDirection()
			for i := 0; i < state.Params.TankSpeed; i++ {
				nPos := pos.step(tank.Pos.Direction)
				if state.Terain.Get(nPos.X, nPos.Y) == 1 {
					break
				}
				pos = nPos
			}
			if preferVertical, danger := dangerous[tank.Pos.NoDirection()]; danger { // 下一步走进危险区
				preferDirection := getPrefer(preferVertical)
				if !preferDirection[tank.Pos.Direction] {
					movements[tank.Id] = ActionStay
				}
			}
		}
	}
}

// 矫正在危险位置的坦克行为
func fixMove (state *GameState, radar *RadarResult, movements map[string]int, tank Tank, preferDirection map[int]bool) {

	dirtIsRight := preferDirection[tank.Pos.Direction]

	nextDirt := 9999
	for dirt, isRight := range preferDirection {
		if !isRight {
			continue
		}
		nextDirt = dirt
		break
	}
	if len(preferDirection) > 1 {
		dirt := nextDirt
		revDirt := (dirt - DirectionUp + 2) % 4 + DirectionUp
		curDirtThreat, revDirtThreat := 0., 0.
		pos := tank.Pos.NoDirection()
		for i := 0; i < state.Params.TankSpeed; i++ {
			pos = pos.step(dirt)
			if state.Terain.Get(pos.X, pos.Y) == 1 {
				break
			}
			curDirtThreat += radar.FullMapThreat[pos]
		}
		pos = pos.step(dirt)
		if state.Terain.Get(pos.X, pos.Y) == 1 {
			curDirtThreat += 0.1
		}
		pos = tank.Pos.NoDirection()
		for i := 0; i < state.Params.TankSpeed; i++ {
			pos = pos.step(revDirt)
			if state.Terain.Get(pos.X, pos.Y) == 1 {
				break
			}
			revDirtThreat += radar.FullMapThreat[pos]
		}
		pos = pos.step(dirt)
		if state.Terain.Get(pos.X, pos.Y) == 1 {
			revDirtThreat += 0.1
		}

		if curDirtThreat > revDirtThreat {
			nextDirt = revDirt
		} else {
			nextDirt = dirt
		}
	}

	movement := movements[tank.Id]
	fmt.Println("Danger zone threat", radar.FullMapThreat[tank.Pos.NoDirection()], preferDirection, tank.Pos.Direction, tank.Id)
	if !dirtIsRight {
		movements[tank.Id] = nextDirt - DirectionUp + ActionTurnUp
	} else if radar.FullMapThreat[tank.Pos.NoDirection()] > 0.2 {
		movements[tank.Id] = ActionMove
	} else if movement >= ActionTurnUp && movement <= ActionTurnRight {
		nextActionDirt := movement - ActionTurnUp + DirectionUp
		moveIsRight := preferDirection[nextActionDirt]
		if !moveIsRight {
			movements[tank.Id] = ActionMove
		}
	}
}