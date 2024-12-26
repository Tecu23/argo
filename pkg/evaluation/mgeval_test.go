package evaluation

import (
	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
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
//
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
//
// func TestBishopPawns(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // bishop square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			square:   D3,
// 			expected: 8,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			square:   G2,
// 			expected: 6,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			square:   C1,
// 			expected: 15,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   C4,
// 			expected: 16,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   C4,
// 			expected: 9,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			square:   E3,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			square:   C4,
// 			expected: 6,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   B2,
// 			expected: 10,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   C4,
// 			expected: 6,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			square:   E2,
// 			expected: 6,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := BishopPawns(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestBishopXrayPawns(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // bishop square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			square:   D3,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			square:   G2,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			square:   D3,
// 			expected: 3,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   C4,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   C4,
// 			expected: 3,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			square:   E3,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			square:   D3,
// 			expected: 3,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   G2,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   C4,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			square:   E2,
// 			expected: 0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := BishopXrayPawns(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestRookOnQueenFile(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // rook square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			square:   A1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			square:   B1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			square:   F1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   D1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   H1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			square:   A1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			square:   E1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   A1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   F1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			square:   D1,
// 			expected: 1,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := RookOnQueenFile(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestRookOnKingRing(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // rook square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			square:   F1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			square:   B1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			square:   F1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   D1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   H1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			square:   A1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			square:   E1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   F1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   F1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			square:   D1,
// 			expected: 0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := RookOnKingRing(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestBishopOnKingRing(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // bishop square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			square:   D3,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			square:   G2,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p1P4/8/2N2N2/PB3PPP/R1BQ1RK1 w - - 0 10",
// 			square:   B2,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   C4,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   C4,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/4P3/5N2/PB1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   C1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			square:   D3,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   G2,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   C1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			square:   E2,
// 			expected: 0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := BishopOnKingRing(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestRookOnFile(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // rook square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			square:   F1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			square:   B1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			square:   F1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   D1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   H1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			square:   A1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			square:   E1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   F1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   F1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			square:   C1,
// 			expected: 2,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := RookOnFile(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestTrappedRook(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // rook square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			square:   F1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			square:   B1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPPBQPPP/RK5R w q - 2 8",
// 			square:   A1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   D1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   H1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N5/PPPBQPPP/RRK5 b q - 4 8",
// 			square:   A1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			square:   E1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPPBQPPP/R4K1R w q - 2 8",
// 			square:   H1,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   F1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N5/PPPBQPPP/RRK5 b q - 4 8",
// 			square:   B1,
// 			expected: 1,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := TrappedRook(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestWeakQueen(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // bishop square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp1bppp/2n1p3/3pP3/Q2P4/2NB1N2/PP3PPP/R1B2RK1 w - - 1 10",
// 			square:   A4,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/2Q5/3P4/2P1PNP1/P3PPBP/1RR3K1 w - - 1 20",
// 			square:   C5,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P3Q/2NB1N2/PP3PPP/R1B2RK1 w - - 0 10",
// 			square:   H4,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   C1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   D1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			square:   E2,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1pQ2/P1BP4/2NB1N2/1P3PPP/R1B1R1K1 w - - 1 16",
// 			square:   F5,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   C2,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/Q1B1P3/5N2/PP1N1PPP/R1B2RK1 w - - 1 9",
// 			square:   A4,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/p3bppp/1pn1pn2/3p4/3P3Q/2N1PN2/PP2BPPP/2RR1RK1 w - - 1 13",
// 			square:   H4,
// 			expected: 1,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := WeakQueen(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestQueenInfiltration(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // bishop square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/pppQbppp/2n1p3/3pP3/3P4/2NB1N2/PP3PPP/R1B2RK1 w - - 2 10",
// 			square:   D7,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pQ2qppp/2n1pn2/8/3P4/2P1PNP1/P3PPBP/1RR3K1 w - - 2 20",
// 			square:   B7,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P3Q/2NB1N2/PP3PPP/R1B2RK1 w - - 0 10",
// 			square:   H4,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   C1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "1Qbqkb1r/rppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1B1K2R w KQk - 2 8",
// 			square:   B8,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			square:   E2,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1pQ2/P1BP4/2NB1N2/1P3PPP/R1B1R1K1 w - - 1 16",
// 			square:   F5,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   C2,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1b2rk1/3nppbp/pQn1q1p1/1pp5/2B1P3/5N2/PP1N1PPP/R1B2RK1 w - - 2 9",
// 			square:   B6,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/pQ2bppp/1pn1pn2/3p4/3P4/2N1PN2/PP2BPPP/2RR1RK1 w - - 1 13",
// 			square:   B7,
// 			expected: 1,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := QueenInfiltration(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestLongDiagonalBishop(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int // bishop square
// 		expected int
// 	}{
// 		// Complex Position
// 		{
// 			name:     "Complex Position 1",
// 			fen:      "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			square:   D3,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 2",
// 			fen:      "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			square:   G2,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 3",
// 			fen:      "r2q1rk1/1b2bppp/p1n1pn2/1p1P4/8/2N2N2/PB3PPP/R1BQ1RK1 w - - 0 10",
// 			square:   B2,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 4",
// 			fen:      "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			square:   C4,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 5",
// 			fen:      "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			square:   C4,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 6",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/4P3/5N2/PB1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   B2,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 7",
// 			fen:      "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			square:   D3,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 8",
// 			fen:      "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			square:   G2,
// 			expected: 1,
// 		},
// 		{
// 			name:     "Complex Position 9",
// 			fen:      "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			square:   C1,
// 			expected: 0,
// 		},
// 		{
// 			name:     "Complex Position 10",
// 			fen:      "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			square:   E2,
// 			expected: 0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			result := LongDiagonalBishop(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
//
// func TestPiecesEvaluation(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fen       string
// 		whiteEval int
// 		blackEval int
// 	}{
// 		// Complex Position
// 		{
// 			name:      "Complex Position 1",
// 			fen:       "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			whiteEval: -144,
// 			blackEval: -147,
// 		},
// 		{
// 			name:      "Complex Position 2",
// 			fen:       "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			whiteEval: 44,
// 			blackEval: -10,
// 		},
// 		{
// 			name:      "Complex Position 3",
// 			fen:       "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			whiteEval: -131,
// 			blackEval: -45,
// 		},
// 		{
// 			name:      "Complex Position 4",
// 			fen:       "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			whiteEval: -99,
// 			blackEval: -15,
// 		},
// 		{
// 			name:      "Complex Position 5",
// 			fen:       "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			whiteEval: -174,
// 			blackEval: -140,
// 		},
// 		{
// 			name:      "Complex Position 6",
// 			fen:       "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			whiteEval: -68,
// 			blackEval: -23,
// 		},
// 		{
// 			name:      "Complex Position 7",
// 			fen:       "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			whiteEval: -160,
// 			blackEval: -74,
// 		},
// 		{
// 			name:      "Complex Position 8",
// 			fen:       "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			whiteEval: -41,
// 			blackEval: -54,
// 		},
// 		{
// 			name:      "Complex Position 9",
// 			fen:       "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			whiteEval: -70,
// 			blackEval: -11,
// 		},
// 		{
// 			name:      "Complex Position 10",
// 			fen:       "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			whiteEval: 16,
// 			blackEval: 1,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var result int
// 			b, _ := board.ParseFEN(tt.fen)
// 			result = PiecesMg(&b)
// 			if result != tt.whiteEval {
// 				t.Errorf("White %s: got %v, want %v", tt.name, result, tt.whiteEval)
// 			}
//
// 			mirror := b.Mirror()
// 			result = PiecesMg(mirror)
// 			if result != tt.blackEval {
// 				t.Errorf("Black %s: got %v, want %v", tt.name, result, tt.blackEval)
// 			}
// 		})
// 	}
// }
//
// func TestMiddlegameMobility(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fen       string
// 		whiteEval int
// 		blackEval int
// 	}{
// 		// Complex Position
// 		{
// 			name:      "Complex Position 1",
// 			fen:       "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			whiteEval: 167,
// 			blackEval: 48,
// 		},
// 		{
// 			name:      "Complex Position 2",
// 			fen:       "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			whiteEval: 186,
// 			blackEval: 127,
// 		},
// 		{
// 			name:      "Complex Position 3",
// 			fen:       "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			whiteEval: 159,
// 			blackEval: 165,
// 		},
// 		{
// 			name:      "Complex Position 4",
// 			fen:       "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			whiteEval: 189,
// 			blackEval: 220,
// 		},
// 		{
// 			name:      "Complex Position 5",
// 			fen:       "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			whiteEval: 128,
// 			blackEval: 78,
// 		},
// 		{
// 			name:      "Complex Position 6",
// 			fen:       "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			whiteEval: 191,
// 			blackEval: 170,
// 		},
// 		{
// 			name:      "Complex Position 7",
// 			fen:       "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			whiteEval: 265,
// 			blackEval: 201,
// 		},
// 		{
// 			name:      "Complex Position 8",
// 			fen:       "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			whiteEval: 164,
// 			blackEval: 190,
// 		},
// 		{
// 			name:      "Complex Position 9",
// 			fen:       "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			whiteEval: 135,
// 			blackEval: 108,
// 		},
// 		{
// 			name:      "Complex Position 10",
// 			fen:       "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			whiteEval: 209,
// 			blackEval: 114,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var result int
// 			b, _ := board.ParseFEN(tt.fen)
// 			result = MobilityMg(&b)
// 			if result != tt.whiteEval {
// 				t.Errorf("White %s: got %v, want %v", tt.name, result, tt.whiteEval)
// 			}
//
// 			mirror := b.Mirror()
// 			result = MobilityMg(mirror)
// 			if result != tt.blackEval {
// 				t.Errorf("Black %s: got %v, want %v", tt.name, result, tt.blackEval)
// 			}
// 		})
// 	}
// }
//
// func TestMiddlegameSpace(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fen       string
// 		whiteEval int
// 		blackEval int
// 	}{
// 		// Complex Position
// 		{
// 			name:      "Complex Position 1",
// 			fen:       "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			whiteEval: 192,
// 			blackEval: 96,
// 		},
// 		{
// 			name:      "Complex Position 2",
// 			fen:       "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			whiteEval: 75,
// 			blackEval: 45,
// 		},
// 		{
// 			name:      "Complex Position 3",
// 			fen:       "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			whiteEval: 159,
// 			blackEval: 73,
// 		},
// 		{
// 			name:      "Complex Position 4",
// 			fen:       "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			whiteEval: 83,
// 			blackEval: 56,
// 		},
// 		{
// 			name:      "Complex Position 5",
// 			fen:       "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			whiteEval: 99,
// 			blackEval: 84,
// 		},
// 		{
// 			name:      "Complex Position 6",
// 			fen:       "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			whiteEval: 83,
// 			blackEval: 56,
// 		},
// 		{
// 			name:      "Complex Position 7",
// 			fen:       "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			whiteEval: 147,
// 			blackEval: 56,
// 		},
// 		{
// 			name:      "Complex Position 8",
// 			fen:       "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			whiteEval: 122,
// 			blackEval: 105,
// 		},
// 		{
// 			name:      "Complex Position 9",
// 			fen:       "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			whiteEval: 83,
// 			blackEval: 81,
// 		},
// 		{
// 			name:      "Complex Position 10",
// 			fen:       "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			whiteEval: 122,
// 			blackEval: 105,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var result int
// 			b, _ := board.ParseFEN(tt.fen)
// 			result = Space(&b)
// 			if result != tt.whiteEval {
// 				t.Errorf("White %s: got %v, want %v", tt.name, result, tt.whiteEval)
// 			}
//
// 			mirror := b.Mirror()
// 			result = Space(mirror)
// 			if result != tt.blackEval {
// 				t.Errorf("Black %s: got %v, want %v", tt.name, result, tt.blackEval)
// 			}
// 		})
// 	}
// }
//
// func TestMiddleGamePassedPawnEvaluation(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fen       string
// 		whiteEval int
// 		blackEval int
// 	}{
// 		// Complex Position
// 		{
// 			name:      "Complex Position 1",
// 			fen:       "r1bq1rk1/5ppp/1bn1p3/1PppP3/p2P4/2NB1N2/P4PPP/R1BQ1RK1 w - - 0 11",
// 			whiteEval: 51,
// 			blackEval: 3,
// 		},
// 		{
// 			name:      "Complex Position 2",
// 			fen:       "2rr2k1/1p2qppp/2n2n2/4P3/p2P4/5NP1/1Q3PBP/1RR3K1 w - - 0 20",
// 			whiteEval: -18,
// 			blackEval: 96,
// 		},
// 		{
// 			name:      "Complex Position 3",
// 			fen:       "r2q1rk1/1b2bppp/p1n2n2/P3P3/3P4/1pPB1N2/4NPPP/R1BQ1RK1 b - - 1 10",
// 			whiteEval: -13,
// 			blackEval: 157,
// 		},
// 		{
// 			name:      "Complex Position 4",
// 			fen:       "4rrk1/1ppq2b1/p1np2p1/4nP1p/2B5/2N2N2/PPP5/2QRR1K1 w - - 0 17",
// 			whiteEval: 40,
// 			blackEval: 21,
// 		},
// 		{
// 			name:      "Complex Position 5",
// 			fen:       "r1bqkb1r/1p1n1ppp/p1n2n2/2pP4/2B1pP2/2N2N2/PPP3PP/R1BQK2R w KQkq - 0 7",
// 			whiteEval: 64,
// 			blackEval: 29,
// 		},
// 		{
// 			name:      "Complex Position 6",
// 			fen:       "3r1rk1/1bqn1ppp/1P1p1n2/Pp6/p2NP3/1BN1B3/4QPPP/R4RK1 w - - 0 17",
// 			whiteEval: 254,
// 			blackEval: 84,
// 		},
// 		{
// 			name:      "Complex Position 7",
// 			fen:       "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			whiteEval: 10,
// 			blackEval: 0,
// 		},
// 		{
// 			name:      "Complex Position 8",
// 			fen:       "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			whiteEval: 0,
// 			blackEval: 0,
// 		},
// 		{
// 			name:      "Complex Position 9",
// 			fen:       "r1bq1rk1/3nppbp/2n3p1/1Pp5/2B1P3/5N2/P2N1PPP/R1BQ1RK1 w - - 0 9",
// 			whiteEval: 61,
// 			blackEval: -7,
// 		},
// 		{
// 			name:      "Complex Position 10",
// 			fen:       "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			whiteEval: 0,
// 			blackEval: 0,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var result int
// 			b, _ := board.ParseFEN(tt.fen)
// 			result = PassedMg(&b)
// 			if result != tt.whiteEval {
// 				t.Errorf("White %s: got %v, want %v", tt.name, result, tt.whiteEval)
// 			}
//
// 			mirror := b.Mirror()
// 			result = PassedMg(mirror)
// 			if result != tt.blackEval {
// 				t.Errorf("Black %s: got %v, want %v", tt.name, result, tt.blackEval)
// 			}
// 		})
// 	}
// }
//
// func TestMiddlegameThreats(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fen       string
// 		whiteEval int
// 		blackEval int
// 	}{
// 		// Complex Position
// 		{
// 			name:      "Complex Position 1",
// 			fen:       "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			whiteEval: 100,
// 			blackEval: 68,
// 		},
// 		{
// 			name:      "Complex Position 2",
// 			fen:       "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			whiteEval: 119,
// 			blackEval: 7,
// 		},
// 		{
// 			name:      "Complex Position 3",
// 			fen:       "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			whiteEval: 201,
// 			blackEval: 102,
// 		},
// 		{
// 			name:      "Complex Position 4",
// 			fen:       "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			whiteEval: 92,
// 			blackEval: 307,
// 		},
// 		{
// 			name:      "Complex Position 5",
// 			fen:       "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			whiteEval: 28,
// 			blackEval: 83,
// 		},
// 		{
// 			name:      "Complex Position 6",
// 			fen:       "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			whiteEval: 44,
// 			blackEval: 102,
// 		},
// 		{
// 			name:      "Complex Position 7",
// 			fen:       "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			whiteEval: 117,
// 			blackEval: 61,
// 		},
// 		{
// 			name:      "Complex Position 8",
// 			fen:       "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			whiteEval: 79,
// 			blackEval: 104,
// 		},
// 		{
// 			name:      "Complex Position 9",
// 			fen:       "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			whiteEval: 14,
// 			blackEval: 47,
// 		},
// 		{
// 			name:      "Complex Position 10",
// 			fen:       "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			whiteEval: 21,
// 			blackEval: 58,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var result int
// 			b, _ := board.ParseFEN(tt.fen)
// 			result = ThreatsMg(&b)
// 			if result != tt.whiteEval {
// 				t.Errorf("White %s: got %v, want %v", tt.name, result, tt.whiteEval)
// 			}
//
// 			mirror := b.Mirror()
// 			result = ThreatsMg(mirror)
// 			if result != tt.blackEval {
// 				t.Errorf("Black %s: got %v, want %v", tt.name, result, tt.blackEval)
// 			}
// 		})
// 	}
// }
//
// func TestMiddleGameKingEvaluation(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fen       string
// 		whiteEval int
// 		blackEval int
// 	}{
// 		// Complex Position
// 		{
// 			name:      "Complex Position 1",
// 			fen:       "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
// 			whiteEval: 72,
// 			blackEval: -126,
// 		},
// 		{
// 			name:      "Complex Position 2",
// 			fen:       "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
// 			whiteEval: -140,
// 			blackEval: -100,
// 		},
// 		{
// 			name:      "Complex Position 3",
// 			fen:       "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
// 			whiteEval: -36,
// 			blackEval: -126,
// 		},
// 		{
// 			name:      "Complex Position 4",
// 			fen:       "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
// 			whiteEval: 156,
// 			blackEval: -70,
// 		},
// 		{
// 			name:      "Complex Position 5",
// 			fen:       "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
// 			whiteEval: -86,
// 			blackEval: -94,
// 		},
// 		{
// 			name:      "Complex Position 6",
// 			fen:       "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
// 			whiteEval: -50,
// 			blackEval: -118,
// 		},
// 		{
// 			name:      "Complex Position 7",
// 			fen:       "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
// 			whiteEval: -45,
// 			blackEval: -124,
// 		},
// 		{
// 			name:      "Complex Position 8",
// 			fen:       "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
// 			whiteEval: -106,
// 			blackEval: -58,
// 		},
// 		{
// 			name:      "Complex Position 9",
// 			fen:       "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
// 			whiteEval: -76,
// 			blackEval: -164,
// 		},
// 		{
// 			name:      "Complex Position 10",
// 			fen:       "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
// 			whiteEval: -86,
// 			blackEval: -118,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var result int
// 			b, _ := board.ParseFEN(tt.fen)
// 			result = KingMg(&b)
// 			if result != tt.whiteEval {
// 				t.Errorf("White %s: got %v, want %v", tt.name, result, tt.whiteEval)
// 			}
//
// 			mirror := b.Mirror()
// 			result = KingMg(mirror)
// 			if result != tt.blackEval {
// 				t.Errorf("Black %s: got %v, want %v", tt.name, result, tt.blackEval)
// 			}
// 		})
// 	}
// }
