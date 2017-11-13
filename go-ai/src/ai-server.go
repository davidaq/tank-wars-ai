package main;

import (
	f "framework"
	t "tactics"
	"os"
	"fmt"
	"strconv"
	"lib/thrift"
	"thrift-player"
)

type PlayerServer struct {
	player *f.Player
	latestState *f.GameState
	params *f.Params
	terain *f.Terain
	flagX, flagY, flagWait int
	myTank map[int32]bool
	bulletTrack map[string]f.Position
	play chan []*player.Order
}

func NewPlayerServer() *PlayerServer {
	return &PlayerServer {}
}

func (self *PlayerServer) reset() {
	if self.player != nil {
		self.latestState.Ended = true
		go self.player.End(self.latestState)
		self.player = nil
		self.latestState = nil
	}
}

func (self *PlayerServer) UploadMap(gamemap [][]int32) (err error) {
	self.terain = &f.Terain {
		Width: len(gamemap[0]),
		Height: len(gamemap),
		Data: make([][]int, len(gamemap)),
	}
	self.flagX = 0
	self.flagY = 0
	self.flagWait = 1
	for y, lineIn := range gamemap {
		line := make([]int, len(lineIn))
		for x, val := range lineIn {
			if val == 3 {
				line[x] = 0
				self.flagX = x
				self.flagY = y
				self.flagWait = 0
			} else {
				line[x] = int(val)
			}
		}
		self.terain.Data[y] = line
	}
	return nil
}

func (self *PlayerServer) UploadParamters(args *player.Args_) (err error) {
	self.reset()
	self.params = &f.Params {
		TankSpeed: int(args.TankSpeed),
		BulletSpeed: int(args.ShellSpeed),
		TankScore: int(args.TankScore),
		FlagScore: int(args.FlagScore),
		MaxRound: int(args.MaxRound),
		Timeout: int(args.RoundTimeoutInMs),
	}
	return nil
}

func (self *PlayerServer) AssignTanks(tanks []int32) (err error) {
	self.reset()
	self.myTank = make(map[int32]bool)
	for _, id := range tanks {
		self.myTank[id] = true
	}
	return nil
}

func (self *PlayerServer) LatestState(raw *player.GameState) (err error) {
	if self.player == nil {
		tactics := t.StartTactics(os.Getenv("TACTICS"))
		self.player = f.NewPlayer(tactics)
		self.params.FlagTime = 50
		self.params.FlagX = self.terain.Width / 2
		self.params.FlagY = self.terain.Height / 2
	}
	state := &f.GameState {
		Params: self.params,
		Terain: self.terain,
		FlagWait: self.flagWait,
		MyFlag: int(raw.YourFlags),
		EnemyFlag: int(raw.EnemyFlags),
	}
	shotTank := make(map[int32]bool)
	for _, bulletIn := range raw.Shells {
		shotTank[bulletIn.ID] = true
		id := strconv.Itoa(int(bulletIn.ID))
		bullet := f.Bullet {
			Id: "B" + id,
			From: "T" + id,
			Pos: f.Position {
				X: int(bulletIn.Pos.X),
				Y: int(bulletIn.Pos.Y),
				Direction: (func () int {
					switch bulletIn.Dir {
					case player.Direction_UP: return f.DirectionUp
					case player.Direction_DOWN: return f.DirectionDown
					case player.Direction_LEFT: return f.DirectionLeft
					case player.Direction_RIGHT: return f.DirectionRight
					default: return f.DirectionNone
					}
				})(),
			},
		}
		var bulletSet *[]f.Bullet
		if self.myTank[bulletIn.ID] {
			bulletSet = &state.MyBullet
		} else {
			bulletSet = &state.EnemyBullet
		}
		*bulletSet = append(*bulletSet, bullet)
	}
	for _, tankIn := range raw.Tanks {
		tank := f.Tank {
			Id: "T" + strconv.Itoa(int(tankIn.ID)),
			Hp: int(tankIn.Hp),
			Pos: f.Position {
				X: int(tankIn.Pos.X),
				Y: int(tankIn.Pos.Y),
				Direction: (func () int {
					switch tankIn.Dir {
					case player.Direction_UP: return f.DirectionUp
					case player.Direction_DOWN: return f.DirectionDown
					case player.Direction_LEFT: return f.DirectionLeft
					case player.Direction_RIGHT: return f.DirectionRight
					default: return f.DirectionNone
					}
				})(),
			},
			Bullet: "",
		}
		if shotTank[tankIn.ID] {
			tank.Bullet = "B" + tank.Id
		}
		var tankSet *[]f.Tank
		if self.myTank[tankIn.ID] {
			tankSet = &state.MyTank
		} else {
			tankSet = &state.EnemyTank
		}
		*tankSet = append(*tankSet, tank)
	}
	self.latestState = state
	if state.Ended {
		go self.player.End(state)
	} else {
		go (func () {
			var orders []*player.Order
			for tankId, action := range self.player.Play(state, true) {
				numId, _ := strconv.Atoi(tankId[1:])
				order := &player.Order {
					TankId: int32(numId),
				}
				switch (action) {
				case f.ActionMove:
					order.Order = "move"
				case f.ActionTurnUp:
					order.Order = "turnTo"
					order.Dir = player.Direction_UP
				case f.ActionTurnLeft:
					order.Order = "fire"
					order.Dir = player.Direction_LEFT
				case f.ActionTurnDown:
					order.Order = "fire"
					order.Dir = player.Direction_DOWN
				case f.ActionTurnRight:
					order.Order = "fire"
					order.Dir = player.Direction_RIGHT
				case f.ActionFireUp:
					order.Order = "fire"
					order.Dir = player.Direction_UP
				case f.ActionFireLeft:
					order.Order = "fire"
					order.Dir = player.Direction_LEFT
				case f.ActionFireDown:
					order.Order = "fire"
					order.Dir = player.Direction_DOWN
				case f.ActionFireRight:
					order.Order = "fire"
					order.Dir = player.Direction_RIGHT
				default:
					order = nil
				}
				if order != nil {
					orders = append(orders, order)
				}
			}
			self.play <- orders
		})()
	}
	return nil
}

func (self *PlayerServer) GetNewOrders() (r []*player.Order, err error) {
	orders := <- self.play
	return orders, nil
}

func main() {
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory  := thrift.NewTBinaryProtocolFactory(false, false)
	port := os.Getenv("PORT")
	
	serverTransport, err := thrift.NewTServerSocket("0.0.0.0:" + port)
	if err != nil {
		panic(err)
	}
	processor := player.NewPlayerServerProcessor(NewPlayerServer())
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Thrift player server starting on port", port)
	server.Serve()
}
