package tactics

import (
	f "framework"
	"os"
	"fmt"
	"net/http"
	"encoding/json"
	"time"
)

type coplan struct {
	obj f.Objective
	dodge, attack, travel float64
	force int
}

type Proxy struct {
	server *http.Server
	setupCb, stateCb func(func (w http.ResponseWriter))
	coplan map[string]coplan
	coplanp map[string]coplan
	first bool
}

func NewProxy() *Proxy {
	inst := &Proxy {
		setupCb: nil,
		stateCb: nil,
		coplan: nil,
		coplanp: nil,
		first: true,
	}
	go inst.startServer();
	return inst
}

func (self *Proxy) startServer () {
	port := os.Getenv("PROXY_PORT")
	fmt.Println("Starting on", port)
	http.HandleFunc("/", func (w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			self.setupCb = func (fn func (w http.ResponseWriter)) {
				fn(w)
				self.setupCb = nil
			}
			for self.setupCb != nil {
				time.Sleep(100 * time.Millisecond)
			}
		} else if (req.Method == http.MethodPost) {
			defer req.Body.Close()
			decoder := json.NewDecoder(req.Body)
			var dat map[string]interface{}
			decoder.Decode(&dat)
			plan := make(map[string]coplan)
			for key, ival := range dat {
				if val, ok := ival.(map[string]interface{}); ok {
					x, _ := val["x"].(float64)
					y, _ := val["y"].(float64)
					dir, _ := val["direction"].(string)
					action, _ := val["action"].(string)
					plan[key] = coplan {
						obj: f.Objective {
							Target: f.Position {
								X: int(x),
								Y: int(y),
								Direction: f.DirectionFromStr(dir),
							},
							Action: f.ActionFromStr(action),
						},
					}
				} else if val, ok := ival.(string); ok {
					plan[key] = coplan {
						obj: f.Objective {
							Action: f.ActionFromStr(val),
						},
					}
				}
			}
			self.coplan = plan
			self.stateCb = func (fn func (w http.ResponseWriter)) {
				fn(w)
				self.stateCb = nil
			}
			for self.stateCb != nil {
				time.Sleep(100 * time.Millisecond)
			}
		}
	})
	http.ListenAndServe(":" + port, nil)
}

func (self *Proxy) Init(state *f.GameState) {
	for self.setupCb == nil {
		time.Sleep(100 * time.Millisecond)
	}
	self.setupCb(func (w http.ResponseWriter) {
		w.Write(state.Raw)
	})
	for self.setupCb != nil {
		time.Sleep(100 * time.Millisecond)
	}
}

func (self *Proxy) Plan(state *f.GameState, radar *f.RadarResult, objective map[string]f.Objective) {
	if self.first {
		self.first = false
	} else {
		for self.stateCb == nil {
			time.Sleep(100 * time.Millisecond)
		}
		self.stateCb(func (w http.ResponseWriter) {
			w.Write(state.Raw[:len(state.Raw) - 1])
			w.Write([]byte(",\"radar\":"))
			// TODO marshal radar
			radarjson, _ := json.Marshal(radar)
			w.Write(radarjson)
			w.Write([]byte("}"))
		})
		for self.stateCb != nil {
			time.Sleep(100 * time.Millisecond)
		}
	}
	for self.coplan == nil {
		time.Sleep(100 * time.Millisecond)
	}
	self.coplanp = self.coplan
	self.coplan = nil
	for k,v := range self.coplanp {
		objective[k] = v.obj
	}
}
