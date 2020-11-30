package main

import (
	"math"
)

type searchLimits struct {
	depth    int
	nodes    uint64
	moveTime int // in milliseconds
	infinite  bool
}

var limits searchLimits

func (s *searchLimits) init() {
	s.depth = 9999
	s.nodes = math.MaxUint64
	s.moveTime = 9999999999
    s.infinite = false
}

func (s *searchLimits) setDepth(depth int) {
	s.depth = depth
}

func (s *searchLimits) setMoveTime(time int) {
	s.moveTime = time
}
func engine() (frEng, toEng chan string) {
	tell("info string initializing deltapawn engine...")
	frEng = make(chan string)
	toEng = make(chan string)
	go func() {
		for cmd := range toEng {
			switch cmd {
			case "stop":
			case "quit":
            case "go": tell("info string deltapawn is thinking!...")
			}
		}
	}()
	return
}
