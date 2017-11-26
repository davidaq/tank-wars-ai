package framework

// 写死防止干蠢事

func BadCase(state *GameState, radar *RadarResult, movements map[string]int) {
	badCaseDangerZone(state, radar, movements)
	badCaseShootSelf(state, radar, movements)
}

func badCaseShootSelf(state *GameState, radar *RadarResult, movements map[string]int) {
	noPass := make(map[Position]bool)
	noStop := make(map[Position]bool)
	ePos := make(map[Position]bool)
	for _, tank  := range state.MyTank {
		noPass[tank.Pos.NoDirection()] = true
		action := movements[tank.Id]
		if action == ActionMove {
			pos := tank.Pos.NoDirection()
			for i := 0; i < state.Params.TankSpeed; i++ {
				nPos := pos.step(tank.Pos.Direction)
				if state.Terain.Get(nPos.X, nPos.Y) == 1 {
					break
				}
				pos = nPos
			}
			noStop[pos] = true
		}
	}
	for _, etank  := range state.EnemyTank {
		ePos[etank.Pos] = true
	}
	tankloop: for _, tank  := range state.MyTank {
		action := movements[tank.Id]
		if action >= ActionFireUp && action <= ActionFireRight {
			dir := action - ActionFireUp + DirectionUp
			pos := tank.Pos.NoDirection().step(dir)
			if ePos[pos] {
				continue tankloop
			}
			if noPass[pos] || noStop[pos] {
				movements[tank.Id] = ActionStay
				continue tankloop
			}
			for i := 0; i < state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
				if noPass[pos] {
					movements[tank.Id] = ActionStay
					continue tankloop
				}
			}
			if noStop[pos] {
				movements[tank.Id] = ActionStay
				continue tankloop
			}
			for i := 0; i < state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
				if noPass[pos] {
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
	for _, eTank := range state.EnemyTank {
		for _, dir := range directions {
			preferVertical := true
			if dir == DirectionDown || dir == DirectionUp {
				preferVertical = false
			}
			pos := eTank.Pos.NoDirection()
			for i := 0; i < 2 + state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
			}
			for i := 0; i < state.Params.BulletSpeed; i++ {
				pos = pos.step(dir)
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
				nPos := tank.Pos.step(dir)
				if state.Terain.Get(nPos.X, nPos.Y) != 1 {
					preferDirection[dir] = true
				}
			}
			if len(preferDirection) > 0 {
				fixMove(state, radar, movements, tank, preferDirection)
			}
		}
	}
}

// 矫正在危险位置的坦克行为
func fixMove (state *GameState, radar *RadarResult, movements map[string]int, tank Tank, preferDirection map[int]bool) {
}
