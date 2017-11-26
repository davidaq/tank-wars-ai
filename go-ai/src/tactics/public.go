package tactics

import (
	f "framework"
	"math"
	// "fmt"
)


func forestPartol(pos f.Position, terain f.Terain, tankSpeed int) f.Position {
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

func judgeDirectionGrove (pos f.Position, terain f.Terain, tankSpeed int) (int, f.Position) {
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