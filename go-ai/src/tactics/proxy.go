package tactics

// import (
// 	f "framework"
// 	"os"
// 	"fmt"
// 	"net/http"
// 	"encoding/json"
// 	"time"
// 	"sort"
// )

// type coplan struct {
// 	obj f.Objective
// 	dodge, attack, travel float64
// 	force int
// }

// type Proxy struct {
// 	server *http.Server
// 	setupCb, stateCb func(func () *f.GameState)
// 	coplan map[string]coplan
// 	coplanp map[string]coplan
// 	first bool
// }

// func NewProxy() *Proxy {
// 	inst := &Proxy {
// 		setupCb: nil,
// 		stateCb: nil,
// 		coplan: nil,
// 		coplanp: nil,
// 		first: true,
// 	}
// 	go inst.startServer();
// 	return inst
// }

// func (self *Proxy) startServer () {
// 	port := os.Getenv("PROXY_PORT")
// 	fmt.Println("Starting on", port)
// 	http.HandleFunc("/", func (w http.ResponseWriter, req *http.Request) {
// 		if req.Method == http.MethodGet {
// 			self.setupCb = func (fn func () *f.GameState) {
// 				state := fn()
// 				w.Write(state.Raw)
// 				self.setupCb = nil
// 			}
// 			for self.setupCb != nil {
// 				time.Sleep(100 * time.Millisecond)
// 			}
// 		} else if (req.Method == http.MethodPost) {
// 			defer req.Body.Close()
// 			decoder := json.NewDecoder(req.Body)
// 			var dat map[string]interface{}
// 			decoder.Decode(&dat)
// 			plan := make(map[string]coplan)
// 			for key, ival := range dat {
// 				val := ival.(map[string]interface{})
// 				x, _ := val["x"].(float64)
// 				y, _ := val["y"].(float64)
// 				dir, _ := val["direction"].(string)
// 				dodge, _ := val["dodge"].(float64)
// 				attack, _ := val["attack"].(float64)
// 				travel, _ := val["travel"].(float64)
// 				force, _ := val["force"].(string)
// 				plan[key] = coplan {
// 					obj: f.Objective {
// 						Target: f.Position {
// 							X: int(x),
// 							Y: int(y),
// 							Direction: f.DirectionFromStr(dir),
// 						},
// 					},
// 					dodge: dodge,
// 					attack: attack,
// 					travel: travel,
// 					force: f.ActionFromStr(force),
// 				}
// 			}
// 			self.coplan = plan
// 			self.stateCb = func (fn func () *f.GameState) {
// 				state := fn()
// 				w.Write(state.Raw)
// 				self.stateCb = nil
// 			}
// 			for self.stateCb != nil {
// 				time.Sleep(100 * time.Millisecond)
// 			}
// 		}
// 	})
// 	http.ListenAndServe(":" + port, nil)
// }

// func (self *Proxy) Init(state *f.GameState) {
// 	for self.setupCb == nil {
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	self.setupCb(func () *f.GameState {
// 		return state
// 	})
// 	for self.setupCb != nil {
// 		time.Sleep(100 * time.Millisecond)
// 	}
// }

// func (self *Proxy) Plan(state *f.GameState, objective *map[string]f.Objective) {
// 	if self.first {
// 		self.first = false
// 	} else {
// 		for self.stateCb == nil {
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 		self.stateCb(func () *f.GameState {
// 			return state
// 		})
// 		for self.stateCb != nil {
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}
// 	for self.coplan == nil {
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	self.coplanp = self.coplan
// 	self.coplan = nil
// 	for k,v := range self.coplanp {
// 		(*objective)[k] = v.obj
// 	}
// }

// type citem struct {
// 	urg float64
// 	action int
// }
// type citems []citem
// func (a citems) Len() int           { return len(a) }
// func (a citems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a citems) Less(i, j int) bool { return a[i].urg < a[j].urg }

// func (self *Proxy) Decide(tank *f.Tank, state *f.GameState, suggestion f.Suggestion) int {
// 	cp := self.coplanp[tank.Id]
// 	if cp.force != f.ActionNone {
// 		return cp.force
// 	}
// 	var candidate citems
// 	coef := []float64 { cp.dodge, cp.attack, cp.travel }
// 	for i, v := range []f.SuggestionItem { suggestion.Dodge, suggestion.Attack, suggestion.Travel } {
// 		if v.Action != f.ActionNone {
// 			candidate = append(candidate, citem {
// 				urg: coef[i] * float64(v.Urgent),
// 				action: v.Action,
// 			})
// 		}
// 	}
// 	if candidate == nil {
// 		return 0
// 	}
// 	sort.Sort(candidate)
// 	return candidate[0].action;
// }
