package tactics

import (
	f "framework"
	"math"
	"math/rand"
	"fmt"
)


func forestPartol(pos f.Position, terain *f.Terain, tankSpeed int) f.Position {
	if n := nextPos(pos.Direction, pos, terain, tankSpeed); n != pos && (terain.Get(n.X, n.Y) == 2 && rand.Int() % 8 != 0) {
		fmt.Println("Patrol", n)
		return n
	} else {
		var canDirection []f.Position
		for _, dir := range []int { f.DirectionUp, f.DirectionLeft, f.DirectionDown, f.DirectionRight } {
			if n := nextPos(dir, pos, terain, tankSpeed); n != pos && terain.Get(n.X, n.Y) == 2 {
				canDirection = append(canDirection, n)
			}
		}
		fmt.Println("Patrol", canDirection)
		if len(canDirection) > 0 {
			return canDirection[rand.Int() % len(canDirection)]
		}
	}
	return pos
}

func nextPos (direction int, pos f.Position, terain *f.Terain, tankSpeed int) f.Position {
	nextPos := pos
	for i := 0; i < tankSpeed; i++ {
		tPos := nextPos.Step(direction)
		if terain.Get(tPos.X, tPos.Y) == 1 {
			break
		}
		nextPos = tPos
	}
	return nextPos
}

func forestPartolA(pos f.Position, terain *f.Terain, tankSpeed int) f.Position {
	inputPos := pos
	forestExist, posRes := judgeDirectionGrove(inputPos, terain, tankSpeed)

	if forestExist > 0 {
		return posRes
	} else {
		inputPos.Direction = inputPos.Direction%4+1
		forestExist1, posRes1 := judgeDirectionGrove(inputPos, terain, tankSpeed)
		if forestExist1 > 0 {
			return posRes1
		} else {
			inputPos.Direction = inputPos.Direction%4+3
			forestExist2, posRes2 := judgeDirectionGrove(inputPos, terain, tankSpeed)
			if forestExist2 > 0 {
				return posRes2
			} else {
				inputPos.Direction = inputPos.Direction%4+2
				forestExist3, posRes3 := judgeDirectionGrove(inputPos, terain, tankSpeed)
				if forestExist3 > 0 {
					return posRes3
				} else {
					return pos
				}
			}
		}
	}
}

func judgeDirectionGrove (pos f.Position, terain *f.Terain, tankSpeed int) (int, f.Position) {
	posRes := f.Position {
		X: pos.X,
		Y: pos.Y,
		Direction: pos.Direction,
	}
	forestExist := 0
	step := 0
	switch pos.Direction {
	case f.DirectionUp:
		if pos.Y - tankSpeed < 0 {
			step = pos.Y
		} else {
			step = tankSpeed
		}
		for i:= 1; i<=step; i++ {
			if terain.Data[pos.Y-i][pos.X] == 0 {
				forestExist = 0
				break
			}
			if terain.Data[pos.Y-i][pos.X] == 1 {
				break
			}
			if terain.Data[pos.Y-i][pos.X] == 2 {
				forestExist = i
			}
		}
		if forestExist > 0 {
			posRes.Y = posRes.Y - forestExist
		}
	case f.DirectionDown:
		if pos.Y + tankSpeed > terain.Height-1 {
			step = terain.Height-pos.Y-1
		} else {
			step = tankSpeed
		}
		for i:= 1; i<=step; i++ {
			if terain.Data[pos.Y+i][pos.X] == 0 {
				forestExist = 0
				break
			}
			if terain.Data[pos.Y+i][pos.X] == 1 {
				break
			}
			if terain.Data[pos.Y+i][pos.X] == 2 {
				forestExist = i
			}
		}
		if forestExist > 0 {
			posRes.Y = posRes.Y + forestExist
		}
	case f.DirectionLeft:
		if pos.X - tankSpeed < 0 {
			step = pos.X
		} else {
			step = tankSpeed
		}
		for i:= 1; i<=step; i++ {
			if terain.Data[pos.Y][pos.X-i] == 0 {
				forestExist = 0
				break
			}
			if terain.Data[pos.Y][pos.X-i] == 1 {
				break
			}
			if terain.Data[pos.Y][pos.X-i] == 2 {
				forestExist = i
			}
		}
		if forestExist > 0 {
			posRes.X = posRes.X - forestExist
		}
	case f.DirectionRight:
		if pos.X + tankSpeed > terain.Width-1 {
			step = terain.Width-pos.X-1
		} else {
			step = tankSpeed
		}
		for i:= 1; i<=step; i++ {
			if terain.Data[pos.Y][pos.X+i] == 0 {
				forestExist = 0
				break
			}
			if terain.Data[pos.Y][pos.X+i] == 1 {
				break
			}
			if terain.Data[pos.Y][pos.X+i] == 2 {
				forestExist = i
			}
		}
		if forestExist > 0 {
			posRes.X = posRes.X + forestExist
		}
	}
	return forestExist, posRes
}

func caculateEnemyCost(bullet f.Bullet, terain *f.Terain, bulletSpeed int) float64 {
	count := 0
	status := true
	pos := bullet.Pos
	switch pos.Direction {
	case f.DirectionDown:
		for status {
			count++
			if terain.Data[pos.Y+count][pos.X] == 1 {
				status = false
			}
		}
	case f.DirectionUp:
		for status {
			count++
			if terain.Data[pos.Y-count][pos.X] == 1 {
				status = false
			}
		}
	case f.DirectionLeft:
		for status {
			count++
			if terain.Data[pos.Y][pos.X-count] == 1 {
				status = false
			}
		}
	case f.DirectionRight:
		for status {
			count++
			if terain.Data[pos.Y][pos.X+count] == 1 {
				status = false
			}
		}
	}
	return math.Ceil(float64(count/bulletSpeed))
}

// func forestGrouping (tankNum int, terain f.Terain, mapAnalysis f.MapAnalysis) (int, f.Forest) {
func forestGrouping (tankNum int, terain *f.Terain, mapAnalysis *f.MapAnalysis) map[int]int {
	// o := mapAnalysis.Ocnt
	// f := mapAnalysis.Fcnt
	// w := mapAnalysis.Wcnt
	forests := mapAnalysis.Forests
	width := terain.Width
	height := terain.Height
	// tData := terain.Data
	mapArea := width*height
	var large f.Forest
	res := make(map[int]int)

	if len(forests) == 0 {
		return res
	} else {
		large = forests[0]
		// 寻找最大草丛
		for _, forest := range forests {
			if forest.Area >large.Area {
				large = forest
			} else if forest.Area == large.Area {
				if forest.Center.X + forest.Center.Y < large.Center.X + large.Center.Y {
					large = forest
				}
			}
		}

		// if float64(float64(large.Area)/float64(mapArea)) > 0.3 {
		// 	return tankNum, large
		// } else if float64(float64(large.Area)/float64(mapArea)) > 0.15 {
		// 	return tankNum/2, large
		// }	else {
		// 	return 1, large
		// }

		if float64(float64(large.Area)/float64(mapArea)) > 0.15 {
			res[large.Id] = tankNum
		} else if float64(float64(large.Area)/float64(mapArea)) > 0.6 {
			res[large.Id] = tankNum/2
		}	else {
			res[large.Id] = 1
		}

		// 旗子非最大草丛
		for _, forest := range forests {
			if forest.Center.X == width/2 && forest.Center.Y == height/2 {
				if res[forest.Id] == 0 {
					res[forest.Id] = 1
				}
			}
		}

		return res
	}
}
