// 行走寻路行动子系统
package framework

import "lib/go-astar"

func path(env [][]int, source Pos, target Pos, ret SuggestionItem) (SuggestionItem) {
	rows := len(env)
	cols := len(env[0])

	a := astar.NewAStar(rows, cols)
	p2p := astar.NewPointToPoint()

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if env[i][j]!=0 {
				a.FillTile(astar.Point{j, i}, -1) 
			}
		}
  	}

	switch target.direction {
	case 0: break;
	case 1: 
		tmpI := target.y + 1
		for ;tmpI < target.y + 5; tmpI++ {
			if tmpI == cols {
				break;
			}
			if env[target.x][tmpI] == 1 {
				break;
			}
		}
		target.y = tmpI - 1
		break;
	case 2:
		tmpI := target.x + 1
		for ;tmpI < target.x + 5; tmpI++ {
			if tmpI == rows {
				break;
			}
			if env[tmpI][target.y] == 1 {
				break;
			}
		}
		target.x = tmpI - 1
		break;
	case 3:
		tmpI := target.y - 1
		for ;tmpI > target.y - 5; tmpI-- {
			if tmpI == -1 {
				break;
			}
			if env[target.x][tmpI] == 1 {
				break;
			}
		}
		target.y = tmpI + 1
		break;
	case 4:
		tmpI := target.x - 1
		for ;tmpI > target.x - 5; tmpI-- {
			if tmpI == -1 {
				break;
			}
			if env[tmpI][target.y] == 1 {
				break;
			}
		}
		target.x = tmpI + 1
		break;
	}

	sourcePoint := []astar.Point{astar.Point{source.x,source.y}}
	targetPoint := []astar.Point{astar.Point{target.x,target.y}}

	pathoutput := a.FindPath(p2p, sourcePoint, targetPoint)

	firstPoint := Pos{
		x: source.x,
		y: source.y,
	}

	count := 0
	for pathoutput != nil {
			count++
			// fmt.Printf("At (%d, %d)\n", pathoutput.Col, pathoutput.Row)
			if count == 2 {
				firstPoint.x = pathoutput.Col
				firstPoint.y = pathoutput.Row
			}
			pathoutput = pathoutput.Parent
	}

	action := transDirection(source, firstPoint)

	ret.Action = action
	ret.Urgent = count
	return ret
}

func transDirection (source Pos, target Pos) int {
	res := 0 
	if source.x == target.x && source.y == target.y {
		res = 1
		return res
	}
	switch source.direction {
	case 1:
		if source.x < target.x {
			res = 4
		} else if source.x > target.x {
			res = 3
		} else if source.y < target.y {
			res = 4
		} else if source.y > target.y {
			res = 2
		}
	case 2:
		if source.x < target.x {
			res = 4
		} else if source.x > target.x {
			res = 2
		} else if source.y < target.y {
			res = 3
		} else if source.y > target.y {
			res = 4
		}
	case 3:
		if source.x < target.x {
			res = 3
		} else if source.x > target.x {
			res = 4
		} else if source.y < target.y {
			res = 2
		} else if source.y > target.y {
			res = 4
		}
	case 4:
		if source.x < target.x {
			res = 2
		} else if source.x > target.x {
			res = 4
		} else if source.y < target.y {
			res = 4
		} else if source.y > target.y {
			res = 3
		}
	}
	return res
}

type Pos struct {
	x int
	y int
	direction int
}

type Traveller struct {
}

func NewTraveller() *Traveller {
	inst := &Traveller {
	}
	return inst
}
func (self *Traveller) Suggest(tank *Tank, state *GameState, objective *Objective) SuggestionItem {
	ret := SuggestionItem {
		Action: ActionLeft,
		Urgent: 1,
	}

	source := Pos{
		x:tank.Pos.X,
		y:tank.Pos.Y,
		direction:tank.Pos.Direction,
	}
	target := Pos{
		x:objective.Target.X,
		y:objective.Target.Y,
		direction:objective.Target.Direction,
	}

	return path(state.Terain.Data, source, target, ret)
}
