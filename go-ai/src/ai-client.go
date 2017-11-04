package main;

import (
	f "framework"
	t "tactics"
	"net/http"
	"io/ioutil"
	//"fmt"
	"encoding/json"
	"strings"
)

func main() {
	host := "ml.niven.cn:8777"
	gameid := "ryImt-sAW"
	side := "red"

	player := f.NewPlayer(t.NewRandom())
	var state *f.GameState = setup(host, gameid, side)
	for !state.Ended {
		state = act(host, gameid, side, player.Play(state))
	}
	player.Reset()
}

func setup(host string, gameid string, side string) *f.GameState {
	url := "http://" + host + "/game/" + gameid + "/match/" + side
	for {
		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				state, err := f.ParseGameState(body)
				if (err == nil) {
					return state
					break
				}
			}
		}
		resp, err = http.Get("http://" + host + "/game/" + gameid + "/interrupt/")
		if err == nil {
			resp.Body.Close()
		}
	}
	return nil
}

func act(host string, gameid string, side string, move map[string]int) *f.GameState {
	send := make(map[string]string)
	for k,v := range move {
		switch v {
		case f.ActionMove:
			send[k] = "move"
			break
		case f.ActionLeft:
			send[k] = "left"
			break
		case f.ActionRight:
			send[k] = "right"
			break
		case f.ActionFire:
			send[k] = "fire"
			break
		}
	}
	url := "http://" + host + "/game/" + gameid + "/match/" + side
	encoded, _ := json.Marshal(send)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(encoded)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	state, err := f.ParseGameState(body)
	if err != nil {
		panic(err)
	}
	return state
}
