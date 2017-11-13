package main;

import (
	"lib/thrift"
	"thrift-player"
)

type PlayerServer struct {
}

func NewPlayerServer() *PlayerServer {
	return &PlayerServer {
	}
}

func (self *PlayerServer) UploadMap(gamemap [][]int32) (err error) {
	return nil
}

func (self *PlayerServer) UploadParamters(arguments *player.Args_) (err error) {
	return nil
}

func (self *PlayerServer) AssignTanks(tanks []int32) (err error) {
	return nil
}

func (self *PlayerServer) LatestState(state *GameState) (err error) {
	return nil
}

func (self *PlayerServer) GetNewOrders() (r []*Order, err error) {
	return nil, nil
}

func main() {
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory  := thrift.NewTBinaryProtocolFactory(false, false)
	
	serverTransport, err := thrift.NewTServerSocket("0.0.0.0:8787")
	if err != nil {
		panic(err)
	}
	processor := player.NewPlayerServerProcessor(NewPlayerServer())
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Thrift player server start")
	server.Serve()
}
