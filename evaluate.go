package main

import (
	//"fmt"
	"math/bits"
)

const (
	maxEval  = +10000
	minEval  = -maxEval
	mateEval = maxEval + 1
	noScore  = minEval - 1
)

// files for the chess board mainly for deducing pawn structure
const (
	aFile   = bitBoard(0x0101010101010101)
	bFile   = bitBoard(0x0202020202020202)
	cFile   = bitBoard(0x0404040404040404)
	dFile   = bitBoard(0x0808080808080808)
	eFile   = bitBoard(0x1010101010101010)
	fFile   = bitBoard(0x2020202020202020)
	gFile   = bitBoard(0x4040404040404040)
	hFile   = bitBoard(0x8080808080808080)
	fileNum = 8
)

var pieceVal = [16]int{100, -100, 325, -325, 350, -350, 500, -500, 950, -950, 10000, -10000, 0, 0, 0, 0}

var knightFile = [8]int{-4, -3, -2, +2, +2, 0, -2, -4}
var knightRank = [8]int{-15, 0, +5, +6, +7, +8, +2, -4}
var centerFile = [8]int{-8, -1, 0, +1, +1, 0, -1, -3}
var kingFile = [8]int{+1, +2, 0, -2, -2, 0, +2, +1}
var kingRank = [8]int{+1, 0, -2, -4, -6, -8, -10, -12}
var pawnRank = [8]int{0, 0, 0, 0, +2, +6, +25, 0}
var pawnFile = [8]int{0, 0, +1, +10, +10, +8, +10, +8}

const longDiag = 10

// Piece Square Table
var pSqTab [12][64]int

// count how many ones there are in the bitboard

//TODO: eval hash
//TODO: pawn hash
//TODO: pawn structures. isolated, backward, duo, passed (guarded and not), double and more...
//first punish many pawns on one file
func pawnStructEval(b *boardStruct) int {
	eval := 0
	wpBB := b.pieceBB[Pawn] & b.wbBB[WHITE]
	bpBB := b.pieceBB[Pawn] & b.wbBB[BLACK]

	eval += (bits.OnesCount64(uint64(wpBB&aFile)) - 1) * (pieceVal[0] / 2)
	eval += (bits.OnesCount64(uint64(wpBB&bFile)) - 1) * (pieceVal[0] / 2)
	eval += (bits.OnesCount64(uint64(wpBB&cFile)) - 1) * (pieceVal[0] / 2)
	eval += (bits.OnesCount64(uint64(wpBB&dFile)) - 1) * (pieceVal[0] / 2)
	eval += (bits.OnesCount64(uint64(wpBB&eFile)) - 1) * (pieceVal[0] / 2)
	eval += (bits.OnesCount64(uint64(wpBB&fFile)) - 1) * (pieceVal[0] / 2)
	eval += (bits.OnesCount64(uint64(wpBB&gFile)) - 1) * (pieceVal[0] / 2)
	eval += (bits.OnesCount64(uint64(wpBB&hFile)) - 1) * (pieceVal[0] / 2)

	eval += (bits.OnesCount64(uint64(bpBB&aFile)) - 1) * (pieceVal[1] / 2)
	eval += (bits.OnesCount64(uint64(bpBB&bFile)) - 1) * (pieceVal[1] / 2)
	eval += (bits.OnesCount64(uint64(bpBB&cFile)) - 1) * (pieceVal[1] / 2)
	eval += (bits.OnesCount64(uint64(bpBB&dFile)) - 1) * (pieceVal[1] / 2)
	eval += (bits.OnesCount64(uint64(bpBB&eFile)) - 1) * (pieceVal[1] / 2)
	eval += (bits.OnesCount64(uint64(bpBB&fFile)) - 1) * (pieceVal[1] / 2)
	eval += (bits.OnesCount64(uint64(bpBB&gFile)) - 1) * (pieceVal[1] / 2)
	eval += (bits.OnesCount64(uint64(bpBB&hFile)) - 1) * (pieceVal[1] / 2)

	return eval
}

//TODO: bishop pair
//TODO: King safety. pawn shelter, guarding pieces
//TODO: King attack. Attacking area surrounding the enemy king, closeness to the enemy king
//TODO: space, center control, knight outposts, connected rooks, 7th row and more
//TODO: combine middle game and end game values

// evaluate returns score from white pov
func evaluate(b *boardStruct) int {
	ev := 0
	for sq := A1; sq <= H8; sq++ {
		pc := b.sq[sq]
		if pc == empty {
			continue
		}
		ev += pieceVal[pc]
		ev += pcSqScore(pc, sq)
	}
	ev += pawnStructEval(b)
	return ev
}

// Score returns the piece square table value for a given piece on a given square. Stage = MG/EG
func pcSqScore(pc, sq int) int {
	return pSqTab[pc][sq]
}

// PstInit intits the pieces-square-tables when the program starts
func pcSqInit() {
	for pc := 0; pc < 12; pc++ {
		for sq := 0; sq < 64; sq++ {
			pSqTab[pc][sq] = 0
		}
	}

	for sq := 0; sq < 64; sq++ {

		fl := sq % 8
		rk := sq / 8

		pSqTab[wP][sq] = pawnFile[fl] + pawnRank[rk]

		pSqTab[wN][sq] = knightFile[fl] + knightRank[rk]
		pSqTab[wB][sq] = centerFile[fl] + centerFile[rk]*2

		pSqTab[wR][sq] = centerFile[fl] * 5

		pSqTab[wQ][sq] = centerFile[fl] + centerFile[rk]

		pSqTab[wK][sq] = (kingFile[fl] + kingRank[rk]) * 8
	}

	// bonus for e4 d5 and c4
	pSqTab[wP][E2], pSqTab[wP][D2], pSqTab[wP][E3], pSqTab[wP][D3], pSqTab[wP][E4], pSqTab[wP][D4], pSqTab[wP][C4] = 0, 0, 6, 6, 24, 20, 12

	// long diagonal
	for sq := A1; sq <= H8; sq += NE {
		pSqTab[wB][sq] += longDiag - 2
	}
	for sq := H1; sq <= A8; sq += NW {
		pSqTab[wB][sq] += longDiag
	}

	// for Black
	for pt := Pawn; pt <= King; pt++ {

		wPiece := pt2pc(pt, WHITE)
		bPiece := pt2pc(pt, BLACK)

		for bSq := 0; bSq < 64; bSq++ {
			wSq := oppRank(bSq)
			pSqTab[bPiece][bSq] = -pSqTab[wPiece][wSq]
		}
	}
}

// mirror the rank_sq
func oppRank(sq int) int {
	fl := sq % 8
	rk := sq / 8
	rk = 7 - rk
	return rk*8 + fl
}
