package main

import (
	"strings"
	"testing"
	"time"
)

var all2GUI []string

func testTell(text ...string) {
	theCmd := ""
	for ix, txt := range text {
		_ = ix
		theCmd += txt
	}
	all2GUI = append(all2GUI, theCmd)
}

func Test_Uci(t *testing.T) {
	tell = testTell
	input := make(chan string)
	go uci(input) // if not 'go' we be blocked here

	tests := []struct {
		name   string
		cmd    string
		wanted []string
	}{
		{"uci", "uci", []string{"id name GoBit", "id author Carokanns", "option name Hash type spin default", "option name Threads type spin default", "uciok"}},
		{"isready", "isready", []string{"readyok"}},
		//{"set Hash", "setoption name Hash value 256", []string{"info string setoption not implemented"}},
		{"skit", "skit", []string{"info string unknown cmd skit"}},
		{"pos skit", "position skit", []string{"info string Error\"skit\" must be \"fen\" or \"startpos\""}},
		{"position no cmd", "position", []string{"info string Error[] wrong length=1"}},
		{"pos incorrect move 1", "position startpos moves e2j4", []string{"info string e2j4 in the position has an incorrect to square"}},
		{"pos incorrect move 2", "position startpos moves e3e4", []string{"info string e3e4 in the position command. fr_sq is an empty square"}},
		//	{"pos incorrect move 3", "position startpos moves e2e4 e7e5 e4e5", []string{"info string e4e5 in moves within the postion commad is not a corect move"}},
		{"ponderhit", "ponderhit", []string{"info string ponderhit not implemented"}},
		{"debug", "debug on", []string{"info string debug not implemented"}},
		//	{"go movetime", "go movetime 1000", []string{"info string engine got go! X"}},
		{"go movestogo", "go movestogo 20", []string{"info string go movestogo not implemented"}},
		{"go wtime", "go wtime 10000", []string{"info string go wtime not implemented"}},
		{"go btime", "go btime 11000", []string{"info string go btime not implemented"}},
		{"go winc", "go winc 500", []string{"info string go winc not implemented"}},
		{"go binc", "go binc 500", []string{"info string go binc not implemented"}},
		//	{"go depth", "go depth 7", []string{"info string go depth not implemented"}},
		{"go nodes", "go nodes 11000", []string{"info string go nodes not implemented"}},
		{"go mate", "go mate 11000", []string{"info string go mate not implemented"}},
		{"go ponder", "go ponder", []string{"info string go ponder not implemented"}},
		//	{"go infinte", "go infinite", []string{"info string go infinite not implemented"}},
		{"wrong cmd", "skit", []string{"info string unknown cmd"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			all2GUI = []string{}
			input <- tt.cmd
			time.Sleep(10 * time.Millisecond)
			for ix, want := range tt.wanted {
				if len(all2GUI) <= ix {
					t.Errorf("%v: we want %#v in ix=%v but got nothing", tt.name, want, ix)
					continue
				}
				if len(want) > len(all2GUI[ix]) {
					t.Errorf("%v: we want %#v (in index %v) but we got %#v", tt.name, want, ix, all2GUI[ix])
					continue
				}
				if all2GUI[ix][:len(want)] != want {
					t.Errorf("%v: Error. Should be %#v but we got %#v", tt.name, want, all2GUI[ix])
				}
			}

		})
	}
}

func Test_handleStop(t *testing.T) {
	tests := []struct {
		name     string
		saveBm   string
		infinite bool
		want     bool
	}{
		{"stop1", "bestmove a1h8", true, false},
		{"stop2", "", true, false},
		{"stop3", "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveBm = tt.saveBm
			limits.setInfinite(tt.infinite)
			handleStop()
			if limits.infinite != tt.want {
				t.Errorf("%v: limits.infinite should be %v but is %v", tt.name, tt.infinite, limits.infinite)
			}
			if limits.stop != true {
				t.Errorf("%v: limits.stop should be %v but is %v", tt.name, true, limits.stop)
			}
			if saveBm != "" {
				t.Errorf("%v: saveBm should be %v but is %v", tt.name, "", saveBm)
			}
		})
	}
}

func Test_handlePosition(t *testing.T) {
	type arg struct{ sq, pc int }
	tests := []struct {
		name  string
		cmd   string
		args  []arg
		castl castlings
	}{
		{"fen", "position fen rnbqkb1r/ppp1pp1p/5np1/3p4/3P1B2/2N1P3/PPP2PPP/R2QKBNR b KQq - 1 4", []arg{{A1, wR}}, castlings(shortW | longW | longB)},
		{"startpos", "position startpos", []arg{{A1, wR}, {A8, bR}, {E5, empty}}, castlings(shortW | longW | shortB | longB)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlePosition(tt.cmd)
			for _, arg := range tt.args {
				if board.sq[arg.sq] != arg.pc {
					t.Errorf("%v: sq=%v should be %v but is %v", tt.name, arg.sq, arg.pc, board.sq[arg.sq])
				}
			}
			if board.castlings != tt.castl {
				t.Errorf("%v: castlings should be %v but is %v", tt.name, tt.castl, board.castlings)
			}
		})
	}
}

func Test_handleNewgame(t *testing.T) {
	t.Run("ucinewgame", func(t *testing.T) {
		handleNewgame()
		if board.stm != WHITE {
			t.Errorf("%v: stm should be %v but we got %v", "ucinewgame", WHITE, board.stm)
		}
		if board.sq[A1] != wR {
			t.Errorf("%v: sq=%v should be %v but we got %v", "ucinewgame", A1, wR, board.sq[A1])
		}
		if board.sq[E3] != empty {
			t.Errorf("%v: sq=%v should be %v but we got %v", "ucinewgame", E3, empty, board.sq[E3])
		}
		if board.King[WHITE] != E1 {
			t.Errorf("%v: wKing sq should be %v but we got %v", "ucinewgame", E1, board.King[WHITE])
		}
		if board.ep != 0 {
			t.Errorf("%v: ep sq should be %v but we got %v", "ucinewgame", 0, board.ep)
		}
		if board.rule50 != 0 {
			t.Errorf("%v: 50 move rule should be %v but we got %v", "ucinewgame", 0, board.rule50)
		}
	})
}

func Test_handleSetOption(t *testing.T) {
	tests := []struct {
		name string
		cmd string
		wantEntries int64
	}{
	// Setoption Hash is tested in Test_transx_new(...) in position_test.go
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := strings.Split(tt.cmd," ")
			handleSetOption(words)
		})
	}
}
