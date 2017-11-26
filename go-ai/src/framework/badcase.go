// 写死防止干蠢事

func BadCase(state *f.GameState, radar *f.RadarResult, movements map[string]int) {
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
		if prefer, danger := dangerous[tank.Pos.NoDirection()]; danger {
			preferAction := make(map[int]bool)
			dirs := hDirections
			if preferVertical {

			}
		}
	}
}

// 矫正在危险位置的坦克行为
func fixMove (state *f.GameState, radar *f.RadarResult, movements map[string]int, tank Tank, preferDirection map[int]bool) {
}
