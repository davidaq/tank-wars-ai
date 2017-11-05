package framework

import (
	"encoding/json"
)

type GameState struct {
	Raw []byte
	Ended bool
	Events []Event
	Terain Terain
	MyTank, EnemyTank []Tank
	MyBullet, EnemyBullet []Bullet
}

type Terain struct {
	Width int
	Height int
	Data [][]int
}

func (self Terain) Get(x int, y int) int {
	if (x < 0 || x >= self.Width || y < 0 || y >= self.Height) {
		return 1
	}
	return self.Data[y][x]
}

type Tank struct {
	Id string
	Pos Position
}

type Bullet struct {
	Id string
	From string
	Pos Position
}

type Event struct {
	Typ string
	Target string
	From string
}

// 位置（可携带方向）
type Position struct {
	X, Y, Direction int
}

var directionMapToInt map[string]int
func DirectionFromStr (str string) int {
	if directionMapToInt == nil {
		directionMapToInt := make(map[string]int)
		directionMapToInt["none"] = DirectionNone
		directionMapToInt["up"] = DirectionUp
		directionMapToInt["left"] = DirectionLeft
		directionMapToInt["down"] = DirectionDown
		directionMapToInt["right"] = DirectionRight
	}
	return directionMapToInt[str]
}

func ParseGameState (bytes []byte) (*GameState, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(bytes, &dat); err != nil {
		return nil, err
	}
	ret := &GameState {
		Raw: bytes,
		Ended: dat["ended"].(bool),
		Events: nil,
		MyTank: nil,
		EnemyTank: nil,
		MyBullet: nil,
		EnemyBullet: nil,
		Terain: Terain {
			Width: 0,
			Height: 0,
			Data: nil,
		},
	}
	for _, ievent := range dat["events"].([]interface{}) {
		event := ievent.(map[string]interface{})
		from, _ := event["from"].(string)
		ret.Events = append(ret.Events, Event {
			Typ: event["type"].(string),
			Target: event["target"].(string),
			From: from,
		})
	}
	parseTank(dat["myTank"].([]interface{}), &ret.MyTank)
	parseTank(dat["enemyTank"].([]interface{}), &ret.EnemyTank)
	parseBullet(dat["myBullet"].([]interface{}), &ret.MyBullet)
	parseBullet(dat["enemyBullet"].([]interface{}), &ret.EnemyBullet)

	for _, iline := range dat["terain"].([]interface{}) {
		line := iline.([]interface{})
		ret.Terain.Width = len(line)
		oline := make([]int, ret.Terain.Width)
		for _, v := range line {
			oline = append(oline, int(v.(float64)))
		}
		ret.Terain.Data = append(ret.Terain.Data, oline)
	}
	ret.Terain.Height = len(ret.Terain.Data)
	return ret, nil
}

func parseTank(dat []interface{}, tanks *[]Tank) {
	for _, itank := range dat {
		tank := itank.(map[string]interface{})
		*tanks = append(*tanks, Tank {
			Id: tank["id"].(string),
			Pos: Position {
				X: int(tank["x"].(float64)),
				Y: int(tank["y"].(float64)),
				Direction: DirectionFromStr(tank["direction"].(string)),
			},
		})
	}
}

func parseBullet(dat []interface{}, bullets *[]Bullet) {
	for _, ibullet := range dat {
		bullet := ibullet.(map[string]interface{})
		*bullets = append(*bullets, Bullet {
			Id: bullet["id"].(string),
			From: bullet["from"].(string),
			Pos: Position {
				X: int(bullet["x"].(float64)),
				Y: int(bullet["y"].(float64)),
				Direction: DirectionFromStr(bullet["direction"].(string)),
			},
		})
	}
}
