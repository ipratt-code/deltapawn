
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)
// setting variables to be named something shorter
var tell = mainTell
var trim = strings.TrimSpace
var low = strings.ToLower

var saveBm = ""

// handles all the uci protocol cases
func uci(inp chan string) {
	tell("info string Hello from uci")
	toEng, frEng := engine()
	bInfinite := false
	var cmd string
	var bm string
	quit := false
	for !quit {
		select {
		case cmd = <-inp:
			tell("info string uci got ", cmd)
		case bm = <-frEng:
			handleBm(bm, bInfinite)
			continue
		}
		words := strings.Split(cmd, " ")
		words[0] = trim(low(words[0]))

		switch words[0] {
		case "uci":
			handleUci()
		case "setoption":
			handleSetOption(words)
		case "isready":
			handleIsReady()
		case "ucinewgame":
			handleNewgame()
		case "position":
			handlePosition(words)
		case "debug":
			handleDebug(words)
		case "register":
			handleRegister(words)
		case "go":
			handleGo(words)
		case "ponderhit":
			handlePonderhit()

		case "stop":
			handleStop(toEng, &bInfinite)
		case "quit", "q":
			handleQuit(toEng)
			quit = true
			continue

		default:
        // defaults to telling us that this command is unknown
			fmt.Println("unknown cmd", cmd)
			tell("info string unknown cmd")
		}
	}
	
	tell("info string leaving uci()")
}

func handleUci() {
    // tells the client that this is deltapawn
	tell("id name deltapawn")
	tell("id author Ian Pratt")

    // the options we can set
	tell("option name Hash type spin default 128 min 16 max 1024")
	tell("option name Threads type spin default 1 min 1 max 16")
	
    // tells the client we are alright
    tell("uciok")
}

func handleIsReady() {

    // tell the client we are ready to go!
	tell("readyok")
}

// this handles the case when we may have to stop the engine
func handleStop(toEng chan string, bInfinite *bool) {
	if *bInfinite {
		if saveBm != "" {
			tell(saveBm)
			saveBm = ""
		}

		toEng <- "stop"
		*bInfinite = false
	}
}

// handleQuit not really necessary
func handleQuit(toEng chan string) {
	toEng <- "stop"
}

// handles the bestmove output from the engine
func handleBm(bm string, bInfinite bool) {
	if bInfinite {
		saveBm = bm
		return
	}
	tell(bm)
}

// not implemented

func handleSetOption(words []string) {
	// setoption name Hash value 256
	fmt.Println("handleSetOption starting", words)
	tell("info string setoption not implemented")
}
func handleNewgame() {
	fmt.Println("handleNewgame starting")
	tell("info string ucinewgame not implemented")
}

func handlePosition(words []string) {
	// position [fen <fenstring> | startpos ]  moves <move1> .... <movei>
	
	if len(words) > 1 {
		words[1] = trim(low(words[1]))
		switch words[1] {
		case "startpos":
			tell("info string position startpos not implemented")
		case "fen":
			tell("info string position fen not implemented")
		default:
			tell("info string position ", words[1], " not implemented")
		}
	} else {
		tell("info string position not implemented")
	}
}
func handleGo(words []string) {
	// go  searchmoves <move1-moveii>/ponder/wtime <ms>/ btime <ms>/winc <ms>/binc <ms>/movestogo <x>/depth <x>/nodes <x>/movetime <ms>/mate <x>/infinite
	fmt.Println("handleGo starting")
	if len(words) > 1 {
		words[1] = trim(low(words[1]))
		switch words[1] {
		case "searchmoves":
			tell("info string go searchmoves not implemented")
		case "ponder":
			tell("info string go ponder not implemented")
		case "wtime":
			tell("info string go wtime not implemented")
		case "btime":
			tell("info string go btime not implemented")
		case "winc":
			tell("info string go winc not implemented")
		case "binc":
			tell("info string go binc not implemented")
		case "movestogo":
			tell("info string go movestogo not implemented")
		case "depth":
			tell("info string go depth not implemented")
		case "nodes":
			tell("info string go nodes not implemented")
		case "movetime":
			tell("info string go movetime not implemented")
		case "mate":
			tell("info string go mate not implemented")
		case "infinite":
			tell("info string go infinite not implemented")
		default:
			tell("info string go ", words[1], " not implemented")
		}
	} else {
		tell("info string go not implemented")
	}
}

func handlePonderhit() {
	fmt.Println("handlePonderhit starting")
	tell("info string ponderhit not implemented")
}

func handleDebug(words []string) {
	// debug [ on | off ]
	fmt.Println("handleDebug starting")
	tell("info string debug not implemented")
}
func handleRegister(words []string) {
	// register later/name <x>/code <y>
	fmt.Println("handleRegister starting")
	tell("info string register not implemented")
}

func mainTell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println(toGUI)
}

func input() chan string {
	line := make(chan string)
	var reader *bufio.Reader
	reader = bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				line <- text
			}
		}
	}()
	return line
}