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
	case "less":
		return NewLess()
	// case "simple":
	// 	return NewSimple()
	case "killall":
		return NewKillAll()
    case "cattycat":
        return NewCatty()
	case "forest-patrol":
		return NewForestPatrol()
	case "sweep":
		return NewSweep()
	case "waitsweep":
		return NewWaitSweep()
	default:
		return NewRandom()
	}
}
