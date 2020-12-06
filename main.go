package main

import (
	"strings"
)

func main() {
	tell("info string Starting deltapawn...")

	uci(input())

	tell("info string quits deltapawn...")
}

func init() {
	initFen2Sq()
	initMagic()
	initKeys()
	initAtksKings()
	initAtksKnights()
	initCastlings()
	pcSqInit()
	board.newGame()
	handleSetOption(strings.Split("setoption name hash value 32", " "))

}
