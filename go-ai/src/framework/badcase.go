package framework

// 写死防止干蠢事

func BadCase(state *GameState, radar *RadarResult, movements map[string]int) {
	badCaseDangerZone(state, radar, movements)
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
			fixMove(state, radar, movements, tank, preferDirection)
		}
	}
}

// 矫正在危险位置的坦克行为
func fixMove (state *GameState, radar *RadarResult, movements map[string]int, tank Tank, preferDirection map[int]bool) {
}
