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
	case "nearest":
		return NewNearest()
	case "simple":
		return NewSimple()
	default:
		return NewRandom()
	}
}
