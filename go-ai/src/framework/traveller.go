// 行走寻路行动子系统
package framework

import (
	"lib/go-astar";
)

func path(env [][]int, source Pos, target Pos) int {
	rows := len(env)
	cols := len(env[0])

	a := astar.NewAStar(rows, cols)
	p2p := astar.NewPointToPoint()

	
	a.FillTile(astar.Point{ Row: source.y, Col: source.x + 1}, 1);
	a.FillTile(astar.Point{ Row: source.y, Col: source.x - 1}, 1);
	a.FillTile(astar.Point{ Row: source.y + 1, Col: source.x}, 1);
	a.FillTile(astar.Point{ Row: source.y - 1, Col: source.x}, 1);

	switch source.direction {
	case DirectionLeft: 
		a.FillTile(astar.Point{ Row: source.y, Col: source.x - 1}, 0);
		break;
	case DirectionRight: 
		a.FillTile(astar.Point{ Row: source.y, Col: source.x + 1}, 0);
		break;
	case DirectionUp: 
		a.FillTile(astar.Point{ Row: source.y - 1, Col: source.x}, 0);
		break;
	case DirectionDown: 
		a.FillTile(astar.Point{ Row: source.y + 1, Col: source.x}, 0);
		break;
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if env[i][j]!=0 {
				a.FillTile(astar.Point{ Row: i, Col: j }, -1) 
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
			if env[tmpI][target.x] == 1 {
				break;
			}
		}
		if source.x == target.x && source.y > target.y {
			break;
		}
		target.y = tmpI - 1
		break;
	case 2:
		tmpI := target.x + 1
		for ;tmpI < target.x + 5; tmpI++ {
			if tmpI == rows {
				break;
			}
			if env[target.y][tmpI] == 1 {
				break;
			}
		}
		if source.y == target.y && source.x > target.x {
			break;
		}
		target.x = tmpI - 1
		break;
	case 3:
		tmpI := target.y - 1
		for ;tmpI > target.y - 5; tmpI-- {
			if tmpI == -1 {
				break;
			}
			if env[tmpI][target.x] == 1 {
				break;
			}
		}
		if source.x == target.x && source.y < target.y {
			break;
		}
		target.y = tmpI + 1
		break;
	case 4:
		tmpI := target.x - 1
		for ;tmpI > target.x - 5; tmpI-- {
			if tmpI == -1 {
				break;
			}
			if env[target.y][tmpI] == 1 {
				break;
			}
		}
		if source.y == target.y && source.x < target.x {
			break;
		}
		target.x = tmpI + 1
		break;
	}

	sourcePoint := []astar.Point{astar.Point{Row: source.y, Col: source.x}}
	targetPoint := []astar.Point{astar.Point{Row: target.y, Col: target.x}}

	pathoutput := a.FindPath(p2p, sourcePoint, targetPoint)

	nextPoint := Pos{
		x: source.x,
		y: source.y,
	}

	count := 0
	for pathoutput != nil {
			count++
			if count == 2 {
				nextPoint.x = pathoutput.Col
				nextPoint.y = pathoutput.Row
				break
			}
			pathoutput = pathoutput.Parent
	}

	return transDirection(source, nextPoint)
}

func transDirection (source Pos, target Pos) int {
	targetDirection := DirectionNone
	if source.x < target.x {
		targetDirection = DirectionRight
	} else if source.x > target.x {
		targetDirection = DirectionLeft
	} else if source.y < target.y {
		targetDirection = DirectionDown
	} else if source.y > target.y {
		targetDirection = DirectionUp
	} else {
		targetDirection = target.direction
		if targetDirection == DirectionNone || source.direction == target.direction {
			return ActionStay	
		} 
	}
	if targetDirection == source.direction {
		return ActionMove
	}
	if ((targetDirection - 1) - (source.direction - 1) + 4) % 4 == 1 {
		return ActionLeft
	} else if ((targetDirection - 1) - (source.direction - 1) + 4) % 4 == 2 {
		return ActionBack
	} else {
		return ActionRight
	}
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

// 每一回合第一个坦克寻路前调用
func (self *Traveller) BeforeSearch(state* GameState) {
}

// 对于每一个坦克依次调用
func (self *Traveller) Search(tank *Tank, state *GameState, target *Position) int {
	source := Pos {
		x: tank.Pos.X,
		y: tank.Pos.Y,
		direction: tank.Pos.Direction,
	}
	ntarget := Pos {
		x: target.X,
		y: target.Y,
		direction: target.Direction,
	}
	return path(state.Terain.Data, source, ntarget)
}
