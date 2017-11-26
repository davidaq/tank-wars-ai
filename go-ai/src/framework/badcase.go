// 写死防止干蠢事

func BadCase(state *f.GameState, radar *f.RadarResult, movements map[string]int) {
	for _, eTank := range state.EnemyTank {

	}
	for _, tank := range state.MyTank {

	}
}

// 矫正在危险位置的坦克行为
func fixMove (state *f.GameState, radar *f.RadarResult, movements map[string]int, tank Tank, preferVertical bool) int {
}
