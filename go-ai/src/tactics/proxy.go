package tactics

import (
	f "framework"
	"os"
	"fmt"
	"net/http"
	"encoding/json"
	"time"
)

type Proxy struct {
	server *http.Server
	stateChan chan []byte
	planChan chan map[string]f.Objective
}

func NewProxy() *Proxy {
	inst := &Proxy {
		stateChan: make(chan []byte),
		planChan: make(chan map[string]f.Objective),
	}
	go inst.startServer();
	return inst
}

func (self *Proxy) startServer () {
	port := os.Getenv("PROXY_PORT")
	fmt.Println("Starting on", port)
	http.HandleFunc("/", func (w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			w.Write(<- self.stateChan)
		} else if (req.Method == http.MethodPost) {
			defer req.Body.Close()
			decoder := json.NewDecoder(req.Body)
			var dat map[string]interface{}
			decoder.Decode(&dat)
			plan := make(map[string]f.Objective)
			for key, ival := range dat {
				if val, ok := ival.(map[string]interface{}); ok {
					x, _ := val["x"].(float64)
					y, _ := val["y"].(float64)
					dir, _ := val["direction"].(string)
					action, _ := val["action"].(string)
					plan[key] = f.Objective {
						Target: f.Position {
							X: int(x),
							Y: int(y),
							Direction: f.DirectionFromStr(dir),
						},
						Action: f.ActionFromStr(action),
					}
				} else if val, ok := ival.(string); ok {
					plan[key] = f.Objective {
						Action: f.ActionFromStr(val),
					}
				} else if val, ok := ival.(float64); ok {
					plan[key] = f.Objective {
						Action: int(val),
					}
				}
			}
			self.planChan <- plan
			w.Write(<- self.stateChan)
		}
	})
	http.ListenAndServe(":" + port, nil)
}

func (self *Proxy) Init(state *f.GameState) {
}

func (self *Proxy) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	var stateByte []byte
	stateByte = append(state.Raw[:len(state.Raw) - 1], []byte(",\"radar\":")...)
	radarjson, _ := json.Marshal(radar)
	stateByte = append(stateByte, radarjson...)
	stateByte = append(stateByte, byte('}'))
	self.stateChan <- stateByte
	recv := <- self.planChan
	for k, v := range recv {
		objective[k] = v
	}
}

func (self *Proxy) End(state *f.GameState) {
	self.stateChan <- state.Raw
	time.Sleep(300 * time.Millisecond)
}
