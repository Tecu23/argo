package evaluation

import (
	"testing"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(Bishop)
	attacks.InitSliderPiecesAttacks(Rook)

	util.InitFen2Sq()

	hash.Init()
}

// func TestPieceValueMg(t *testing.T) {
// 	tests := []struct {
// 		name            string
// 		fen             string
// 		pieceValueBlack int // pieceValue evaluation in centipawns
// 		pieceValueWhite int // pieceValue evaluation in centipawns
// 		psqtWhite       int
// 		psqtBlack       int
// 		imbalanceWhite  int
// 		imbalanceBlack  int
// 		imbalanceTotal  int
// 		desc            string
// 	}{
// 		// Material Balance Tests
// 		{
// 			name:            "Even Material",
// 			fen:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       12,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "Starting position - equal material",
// 		},
// 		{
// 			name:            "Pawn Up",
// 			fen:             "rnbqkbnr/ppp1pppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9170,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       12,
// 			psqtBlack:       -7,
// 			imbalanceWhite:  12794,
// 			imbalanceBlack:  11844,
// 			imbalanceTotal:  59,
// 			desc:            "White up a pawn",
// 		},
//
// 		// Pawn Structure Tests
// 		{
// 			name:            "Doubled Pawns",
// 			fen:             "rnbqkbnr/pppppppp/8/8/8/P7/1PPPPPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       0,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White has doubled a-pawns",
// 		},
// 		{
// 			name:            "Tripled Pawns",
// 			fen:             "rnbqkbnr/pppppppp/8/8/P7/P7/1PPPPPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9418,
// 			psqtWhite:       -4,
// 			psqtBlack:       12,
// 			imbalanceWhite:  14688,
// 			imbalanceBlack:  13662,
// 			imbalanceTotal:  64,
// 			desc:            "White has tripled a-pawns",
// 		},
// 		{
// 			name:            "Isolated Pawn",
// 			fen:             "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       13,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White has isolated d-pawn",
// 		},
// 		{
// 			name:            "Multiple Isolated Pawns",
// 			fen:             "rnbqkbnr/pppppppp/8/8/3P1P2/8/PPP2PPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       14,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White has two isolated pawns",
// 		},
// 		{
// 			name:            "Backward Pawn",
// 			fen:             "rnbqkbnr/pppppppp/8/2P5/3P4/8/PP2PPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       -10,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White has backward c-pawn",
// 		},
// 		{
// 			name:            "Passed Pawn",
// 			fen:             "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9170,
// 			pieceValueWhite: 9418,
// 			psqtWhite:       13,
// 			psqtBlack:       -7,
// 			imbalanceWhite:  14254,
// 			imbalanceBlack:  12278,
// 			imbalanceTotal:  123,
// 			desc:            "White has passed d-pawn",
// 		},
// 		{
// 			name:            "Protected Passed Pawn",
// 			fen:             "rnbqkbnr/ppp1pppp/8/2PP4/8/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9170,
// 			pieceValueWhite: 9418,
// 			psqtWhite:       -19,
// 			psqtBlack:       -7,
// 			imbalanceWhite:  14254,
// 			imbalanceBlack:  12278,
// 			imbalanceTotal:  123,
// 			desc:            "White has protected passed pawn",
// 		},
// 		{
// 			name:            "Connected Passed Pawns",
// 			fen:             "rnbqkbnr/ppp1pppp/8/2PP4/8/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9170,
// 			pieceValueWhite: 9418,
// 			psqtWhite:       -19,
// 			psqtBlack:       -7,
// 			imbalanceWhite:  14254,
// 			imbalanceBlack:  12278,
// 			imbalanceTotal:  123,
// 			desc:            "White has connected passed pawns",
// 		},
//
// 		// Piece Mobility Tests
// 		{
// 			name:            "Trapped Knight",
// 			fen:             "rnbqkbnr/pppppppp/8/8/8/7N/PPPPPPPP/RNBQKB1R w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       43,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White knight trapped in corner",
// 		},
// 		{
// 			name:            "Trapped Bishop",
// 			fen:             "rnbqkbnr/ppppp1pp/5p2/8/8/7B/PPPPPPPP/RNBQK1NR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       13,
// 			psqtBlack:       15,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White bishop trapped by pawn",
// 		},
// 		{
// 			name:            "Central Knight",
// 			fen:             "rnbqkbnr/pppppppp/8/8/4N3/8/PPPPPPPP/RNBQKB1R w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       153,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White knight on strong central square",
// 		},
// 		{
// 			name:            "Open File Rook",
// 			fen:             "rnbqkbnr/pppp1ppp/8/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9170,
// 			pieceValueWhite: 9170,
// 			psqtWhite:       -4,
// 			psqtBlack:       -4,
// 			imbalanceWhite:  11410,
// 			imbalanceBlack:  11410,
// 			imbalanceTotal:  0,
// 			desc:            "White rook on open e-file",
// 		},
//
// 		// King Safety Tests
// 		{
// 			name:            "Exposed King",
// 			fen:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 7688,
// 			psqtWhite:       112,
// 			psqtBlack:       12,
// 			imbalanceWhite:  8657,
// 			imbalanceBlack:  12436,
// 			imbalanceTotal:  -326,
// 			desc:            "White king exposed with missing kingside pawns",
// 		},
// 		{
// 			name:            "Castled King",
// 			fen:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1RK1 w kq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 7688,
// 			psqtWhite:       258,
// 			psqtBlack:       12,
// 			imbalanceWhite:  8657,
// 			imbalanceBlack:  12436,
// 			imbalanceTotal:  -326,
// 			desc:            "White king safely castled",
// 		},
// 		{
// 			name:            "Fianchettoed Bishop",
// 			fen:             "rnbqkbnr/pppppppp/8/8/8/6P1/PPPPPP1P/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       10,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White preparing kingside fianchetto",
// 		},
//
// 		// Space and Control Tests
// 		{
// 			name:            "Central Control",
// 			fen:             "rnbqkbnr/pppppppp/8/8/3PP3/8/PPP2PPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       37,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White controlling center with pawns",
// 		},
// 		{
// 			name:            "Advanced Center",
// 			fen:             "rnbqkbnr/pppppppp/8/3PP3/8/8/PPP2PPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       -11,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White with advanced center pawns",
// 		},
// 		{
// 			name:            "Space Advantage",
// 			fen:             "rnbqkbnr/pppppppp/8/8/2PPP3/2N5/PP3PPP/R1BQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       131,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White with significant space advantage",
// 		},
//
// 		// Piece Coordination Tests
// 		{
// 			name:            "Bishop Pair",
// 			fen:             "rnbqk1nr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 8469,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       12,
// 			psqtBlack:       20,
// 			imbalanceWhite:  12526,
// 			imbalanceBlack:  11239,
// 			imbalanceTotal:  170,
// 			desc:            "White has bishop pair vs bishop+knight",
// 		},
// 		{
// 			name:            "Doubled Rooks",
// 			fen:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3RK2 w kq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 3544,
// 			psqtWhite:       307,
// 			psqtBlack:       12,
// 			imbalanceWhite:  2572,
// 			imbalanceBlack:  12257,
// 			imbalanceTotal:  -695,
// 			desc:            "White has doubled rooks on open file",
// 		},
// 		{
// 			name:            "Knight Outpost",
// 			fen:             "rnbqkbnr/ppp1pppp/8/3p4/4N3/8/PPPPPPPP/RNBQKB1R w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       153,
// 			psqtBlack:       13,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White knight on protected outpost",
// 		},
//
// 		// Development and Tempo Tests
// 		{
// 			name:            "Better Development",
// 			fen:             "r1bqkbnr/pppppppp/2n5/8/4P3/2N1B3/PPPP1PPP/R1BQK1NR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       159,
// 			psqtBlack:       110,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White better developed in opening",
// 		},
// 		{
// 			name:            "Undeveloped Position",
// 			fen:             "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1",
// 			pieceValueBlack: 9294,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       36,
// 			psqtBlack:       12,
// 			imbalanceWhite:  13228,
// 			imbalanceBlack:  13228,
// 			imbalanceTotal:  0,
// 			desc:            "White with poor development",
// 		},
//
// 		// Complex Positions
// 		{
// 			name:            "Complex Center",
// 			fen:             "r1bqk2r/ppp2ppp/2n2n2/1B1pp3/4P3/2PP1N2/PP3PPP/RNBQK2R w KQkq - 0 1",
// 			pieceValueBlack: 8469,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       168,
// 			psqtBlack:       241,
// 			imbalanceWhite:  12526,
// 			imbalanceBlack:  11239,
// 			imbalanceTotal:  170,
// 			desc:            "Complex central tension",
// 		},
// 		{
// 			name:            "Attacking Position",
// 			fen:             "rnbqk2r/ppp2ppp/5n2/3Pp3/1b2P3/2N2N2/PPP2PPP/R1BQKB1R w KQkq - 0 1",
// 			pieceValueBlack: 9170,
// 			pieceValueWhite: 9294,
// 			psqtWhite:       214,
// 			psqtBlack:       152,
// 			imbalanceWhite:  12794,
// 			imbalanceBlack:  11844,
// 			imbalanceTotal:  59,
// 			desc:            "White with attacking chances",
// 		},
//
// 		// Special Cases
// 		{
// 			name:            "Rook vs Three Minors",
// 			fen:             "8/8/8/8/8/2nbn3/8/4R3 w - - 0 1",
// 			pieceValueBlack: 2387,
// 			pieceValueWhite: 1276,
// 			psqtWhite:       -5,
// 			psqtBlack:       122,
// 			imbalanceWhite:  -184,
// 			imbalanceBlack:  -240,
// 			imbalanceTotal:  3,
// 			desc:            "Material imbalance special case",
// 		},
// 		{
// 			name:            "Queen vs Three Minors",
// 			fen:             "8/8/8/8/8/2nbn3/8/4Q3 w - - 0 1",
// 			pieceValueBlack: 2387,
// 			pieceValueWhite: 2538,
// 			psqtWhite:       4,
// 			psqtBlack:       122,
// 			imbalanceWhite:  47,
// 			imbalanceBlack:  -240,
// 			imbalanceTotal:  17,
// 			desc:            "Roughly equal material imbalance",
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			got := PieceValueMg(&b)
// 			psqt := PsqtMg(&b)
// 			imb := Imbalance(&b)
//
// 			if got != tt.pieceValueWhite {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot eval: %d\nWant eval: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					got,
// 					tt.pieceValueWhite,
// 					got-tt.pieceValueWhite,
// 					tt.fen,
// 				)
// 			}
//
// 			if psqt != tt.psqtWhite {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot psqt: %d\nWant psqt: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					psqt,
// 					tt.psqtWhite,
// 					psqt-tt.psqtWhite,
// 					tt.fen,
// 				)
// 			}
//
// 			if imb != tt.imbalanceWhite {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot imbalance: %d\nWant imbalance: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					imb,
// 					tt.imbalanceWhite,
// 					imb-tt.imbalanceWhite,
// 					tt.fen,
// 				)
// 			}
//
// 			// Test evaluation symmetry
// 			mirroredBoard := b.Mirror()
// 			mirroredEval := PieceValueMg(mirroredBoard)
// 			psqt = PsqtMg(mirroredBoard)
// 			imb = Imbalance(mirroredBoard)
//
// 			if mirroredEval != tt.pieceValueBlack {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot eval: %d\nWant eval: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					mirroredEval,
// 					tt.pieceValueBlack,
// 					mirroredEval-tt.pieceValueBlack,
// 					tt.fen,
// 				)
// 			}
//
// 			if psqt != tt.psqtBlack {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot psqt: %d\nWant psqt: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					psqt,
// 					tt.psqtBlack,
// 					psqt-tt.psqtBlack,
// 					tt.fen,
// 				)
// 			}
//
// 			if imb != tt.imbalanceBlack {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot imbalance: %d\nWant imbalance: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					imb,
// 					tt.imbalanceBlack,
// 					imb-tt.imbalanceBlack,
// 					tt.fen,
// 				)
// 			}
//
// 			imbTotal := ImbalanceTotal(&b, mirroredBoard)
//
// 			if imbTotal != tt.imbalanceTotal {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot imbalance total: %d\nWant imbalance total: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					imbTotal,
// 					tt.imbalanceTotal,
// 					imbTotal-tt.imbalanceTotal,
// 					tt.fen,
// 				)
// 			}
// 		})
// 	}
// }
//
// func TestPawnMethods(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		fen          string
// 		desc         string
// 		whitePawnsMg int
// 		blackPawnsMg int
// 	}{
// 		{
// 			name:         "Double Pawns",
// 			fen:          "8/1p6/p3pp1p/8/8/1PP2P1P/1P3P1P/8 w - -",
// 			whitePawnsMg: 2,
// 			blackPawnsMg: 72,
// 			desc:         "Different Types of doubled pawns",
// 		},
// 		{
// 			name:         "Crippled Majority",
// 			fen:          "8/1pp2ppp/1p6/8/8/5P2/PPP2PPP/8 w - -",
// 			whitePawnsMg: 120,
// 			blackPawnsMg: 99,
// 			desc:         "Pawn structures with a crippled majority by doubled pawns",
// 		},
// 		{
// 			name:         "Backward Pawn 1",
// 			fen:          "8/8/1p6/8/p1P1p3/3pP3/1P1P4/8 w - -",
// 			whitePawnsMg: 11,
// 			blackPawnsMg: 37,
// 			desc:         "Backward pawn",
// 		},
// 		{
// 			name:         "Backward Pawn 2",
// 			fen:          "8/1pp5/8/2P5/8/8/8/1R6 w - -",
// 			whitePawnsMg: -5,
// 			blackPawnsMg: 35,
// 			desc:         "Stops with negative SEE",
// 		},
// 		// Weak Pawns
// 		{
// 			name:         "Weak Pawns",
// 			fen:          "1k6/2p5/1p6/1P6/p1P5/P7/3K4/8 w - -",
// 			whitePawnsMg: 21,
// 			blackPawnsMg: 9,
// 			desc:         "Overly advanced pawns",
// 		},
//
// 		{
// 			name:         "Hanging Pawns",
// 			fen:          "8/pp3ppp/4p3/8/2PP4/8/P4PPP/8 w - -",
// 			whitePawnsMg: 109,
// 			blackPawnsMg: 114,
// 			desc:         "Hanging Pawns Formation",
// 		},
//
// 		// Chain Detection
// 		{
// 			name:         "Pawn Chains",
// 			fen:          "6k1/5p2/4p3/3pP3/2pP4/2P5/8/6K1 w - -",
// 			whitePawnsMg: 63,
// 			blackPawnsMg: 79,
// 			desc:         "Pawn Chains Formation",
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			pawnEval := PawnsMg(&b)
//
// 			if pawnEval != tt.whitePawnsMg {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot eval: %d\nWant eval: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					pawnEval,
// 					tt.whitePawnsMg,
// 					pawnEval-tt.whitePawnsMg,
// 					tt.fen,
// 				)
// 			}
//
// 			// Test evaluation symmetry
// 			mirroredBoard := b.Mirror()
// 			mirroredBoard.PrintBoard()
// 			pawnEval = PawnsMg(mirroredBoard)
//
// 			if pawnEval != tt.blackPawnsMg {
// 				t.Errorf(
// 					"\nPosition: %s\nDescription: %s\nGot eval: %d\nWant eval: %d\nDiff: %d\nFEN: %s",
// 					tt.name,
// 					tt.desc,
// 					pawnEval,
// 					tt.blackPawnsMg,
// 					pawnEval-tt.blackPawnsMg,
// 					tt.fen,
// 				)
// 			}
// 		})
// 	}
// }

// func TestMinorBehindPawn(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int
// 		expected int
// 	}{
// 		{
// 			name:     "White Minor Behind White Pawn",
// 			fen:      "8/8/8/8/4P3/4N3/8/8 w - - 0 1",
// 			square:   E3,
// 			expected: 1,
// 		},
// 		{
// 			name:     "White Minor Behind Black Pawn",
// 			fen:      "8/8/8/4p3/4N3/8/8/8 w - - 0 1",
// 			square:   E4,
// 			expected: 1,
// 		},
// 		{
// 			name:     "No Pawn in Front",
// 			fen:      "8/8/8/8/8/4N3/8/8 w - - 0 1",
// 			square:   E3,
// 			expected: 0,
// 		},
// 		{
// 			name:     "First Rank",
// 			fen:      "8/8/8/8/8/8/4P3/4N3 w - - 0 1",
// 			square:   E1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Eighth Rank",
// 			fen:      "4N3/4P3/8/8/8/8/8/8 w - - 0 1",
// 			square:   E8,
// 			expected: 0,
// 		},
// 		{
// 			name:     "File A Edge",
// 			fen:      "8/8/8/8/P7/N7/8/8 w - - 0 1",
// 			square:   A3,
// 			expected: 1,
// 		},
// 		{
// 			name:     "File H Edge",
// 			fen:      "8/8/8/8/7P/7N/8/8 w - - 0 1",
// 			square:   H3,
// 			expected: 1,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := MinorBehindPawn(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("got %v, want %v", result, tt.expected)
// 			}
// 		})
// 	}
// }

func TestBishopPawns(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		square   int // bishop square
		expected int
	}{
		// Basic Tests
		{
			name:     "Empty Board with Bishop",
			fen:      "8/8/8/8/4B3/8/8/8 w - - 0 1",
			square:   E4,
			expected: 0,
		},
		{
			name:     "Bishop with Same-Color Pawns",
			fen:      "8/8/8/8/4B3/3P1P2/8/8 w - - 0 1",
			square:   E4,
			expected: 0,
		},
		{
			name:     "Bishop with Blocked Center Pawn",
			fen:      "8/8/8/3p4/4B3/3P4/8/8 w - - 0 1",
			square:   E4,
			expected: 0,
		},
		{
			name:     "Bishop Under Pawn Attack",
			fen:      "8/8/3p4/4B3/8/8/8/8 w - - 0 1",
			square:   E5,
			expected: 0,
		},

		// Edge Cases
		{
			name:     "Bishop on Edge with Edge Pawns",
			fen:      "8/8/8/8/B7/P7/8/8 w - - 0 1",
			square:   A4,
			expected: 0,
		},
		{
			name:     "Bishop with Multiple Blocked Center Pawns",
			fen:      "8/8/8/3ppp2/4B3/3PPP2/8/8 w - - 0 1",
			square:   E4,
			expected: 2,
		},

		// Complex Positions
		{
			name:     "Complex Position 1",
			fen:      "8/8/2pppp2/8/2PPPP2/4B3/8/8 w - - 0 1",
			square:   E3,
			expected: 2,
		},
		{
			name:     "Complex Position 2",
			fen:      "8/8/8/3p4/2P1P3/4B3/3P4/8 w - - 0 1",
			square:   E3,
			expected: 0,
		},
		//
		// Special Cases
		{
			name:     "No Bishop on Square",
			fen:      "8/8/8/8/3P4/8/8/8 w - - 0 1",
			square:   E4,
			expected: 0,
		},
		{
			name:     "Bishop with All Edge Pawns",
			fen:      "8/8/8/8/B7/P7/1P6/8 w - - 0 1",
			square:   A4,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.ParseFEN(tt.fen)
			result := BishopPawns(&b, tt.square)
			if result != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}
