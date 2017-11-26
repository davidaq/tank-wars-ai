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
	noStop := make(map[Position]bool)
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
				noPass[pos] = true
				pos = nPos
			}
		}
		noStop[pos] = true
	}
	fmt.Println("shoot self", noPass, noStop)
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
			if noPass[pos] || noStop[pos] {
				movements[tank.Id] = ActionStay
				continue tankloop
			}
			for i := 0; i < state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
				if state.Terain.Get(pos.X, pos.Y) == 1 {
					continue tankloop
				}
				if noPass[pos] {
					movements[tank.Id] = ActionStay
					continue tankloop
				}
			}
			if noStop[pos] {
				movements[tank.Id] = ActionStay
				continue tankloop
			}
			for i := 0; i < state.Params.BulletSpeed * 2; i++ {
				pos = pos.step(dir)
				if state.Terain.Get(pos.X, pos.Y) == 1 {
					continue tankloop
				}
				if noPass[pos] || noStop[pos] {
					movements[tank.Id] = ActionStay
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
	noPass := make(map[Position]bool)
	noStop := make(map[Position]bool)
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
				noPass[pos] = true
				pos = nPos
			}
		}
		noStop[pos] = true
	}
	for _, eTank := range state.EnemyTank {
		for _, dir := range directions {
			preferVertical := true
			if dir == DirectionDown || dir == DirectionUp {
				preferVertical = false
			}
			pos := eTank.Pos.NoDirection()
			for i := 0; i < 1 + state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
				if state.Terain.Get(pos.X, pos.Y) == 1 {
					continue
				}
			}
			for i := 0; i < state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
				if state.Terain.Get(pos.X, pos.Y) == 1 {
					continue
				}
				dangerous[pos] = preferVertical
			}
		}
	}
	for _, tank := range state.MyTank {
		if preferVertical, danger := dangerous[tank.Pos.NoDirection()]; danger {
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
				if noPass[nPos] {
					blocked = true
				}
				if !blocked {
					preferDirection[dir] = true
				}
			}
			if len(preferDirection) > 0 {
				oldMove := movements[tank.Id]
				fixMove(state, radar, movements, tank, preferDirection)
				fmt.Println("Fix Danger Zone", tank.Id, oldMove, movements[tank.Id])
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
	
	if !dirtIsRight {
		movements[tank.Id] = nextDirt - DirectionUp + ActionTurnUp
	} else if movement >= ActionTurnUp && movement <= ActionTurnRight {
		nextActionDirt := movement - ActionTurnUp + DirectionUp
		moveIsRight := preferDirection[nextActionDirt]
		if !moveIsRight {
			movements[tank.Id] = ActionMove
		}
	}
}