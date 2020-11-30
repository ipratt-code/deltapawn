package main

func main() {
	tell("info string initializing deltapawn...")

	uci(input())

    tell("info string shutting down deltapawn...")
}

func init() {
    initFenSq2Int()
	initMagic()
    initAtksKings()
    initAtksKnights()
}