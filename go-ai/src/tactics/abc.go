package tactics

import (
	f "framework"
)

func StartTactics (name string) f.Tactics {
	switch (name) {
	case "random":
		return NewRandom()
	case "proxy":
		return NewProxy()
	case "nearst":
		return NewNearest()
	case "simple":
		return NewSimple()
	default:
		return NewRandom()
	}
}
