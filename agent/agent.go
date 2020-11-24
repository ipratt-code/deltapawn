package agent

import (
	"fmt"
	"github.com/notnil/chess"
)

//DoTheChess is a test function
func DoTheChess() {
	game := chess.NewGame()
	moves := game.ValidMoves()
	game.Move(moves[0])
	fmt.Println(game.Position().Board().Draw()) // b1a3
}
