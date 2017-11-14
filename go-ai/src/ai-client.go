package main;

import (
	f "framework"
	t "tactics"
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"strings"
	"time"
)

func main() {
	host := os.Getenv("HOST")
	gameid := os.Getenv("GAME")
	side := os.Getenv("SIDE")
	tactics := t.StartTactics(os.Getenv("TACTICS"))
	// tactics := t.NewRandom()
	player := f.NewPlayer(tactics)
	var state *f.GameState = setup(host, gameid, side)
	i := 0
	for !state.Ended {
		state = act(host, gameid, side, player.Play(state, false))
		i++
		fmt.Print(i, "\tmy:", len(state.MyTank), "\tenemy:", len(state.EnemyTank), "\t\t\n")
	}
	player.End(state)
	fmt.Println("\nEnd")
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
	responded := false
	defer (func() { responded = true })()
	go (func () {
		for i := 0; i < 20 && !responded; i++ {
			time.Sleep(50 * time.Millisecond)
		}
		if !responded {
			fmt.Println("Game server not responding")
		}
	})()
	send := make(map[string]string)
	for k,v := range move {
		send[k] = f.ActionToStr(v)
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
