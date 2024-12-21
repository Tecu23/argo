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

func TestPieceValueMg(t *testing.T) {
	tests := []struct {
		name            string
		fen             string
		pieceValueBlack int // pieceValue evaluation in centipawns
		pieceValueWhite int // pieceValue evaluation in centipawns
		psqtWhite       int
		psqtBlack       int
		imbalanceWhite  int
		imbalanceBlack  int
		imbalanceTotal  int
		desc            string
	}{
		// Material Balance Tests
		{
			name:            "Even Material",
			fen:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       12,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "Starting position - equal material",
		},
		{
			name:            "Pawn Up",
			fen:             "rnbqkbnr/ppp1pppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9170,
			pieceValueWhite: 9294,
			psqtWhite:       12,
			psqtBlack:       -7,
			imbalanceWhite:  12794,
			imbalanceBlack:  11844,
			imbalanceTotal:  59,
			desc:            "White up a pawn",
		},

		// Pawn Structure Tests
		{
			name:            "Doubled Pawns",
			fen:             "rnbqkbnr/pppppppp/8/8/8/P7/1PPPPPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       0,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White has doubled a-pawns",
		},
		{
			name:            "Tripled Pawns",
			fen:             "rnbqkbnr/pppppppp/8/8/P7/P7/1PPPPPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9418,
			psqtWhite:       -4,
			psqtBlack:       12,
			imbalanceWhite:  14688,
			imbalanceBlack:  13662,
			imbalanceTotal:  64,
			desc:            "White has tripled a-pawns",
		},
		{
			name:            "Isolated Pawn",
			fen:             "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       13,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White has isolated d-pawn",
		},
		{
			name:            "Multiple Isolated Pawns",
			fen:             "rnbqkbnr/pppppppp/8/8/3P1P2/8/PPP2PPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       14,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White has two isolated pawns",
		},
		{
			name:            "Backward Pawn",
			fen:             "rnbqkbnr/pppppppp/8/2P5/3P4/8/PP2PPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       -10,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White has backward c-pawn",
		},
		{
			name:            "Passed Pawn",
			fen:             "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9170,
			pieceValueWhite: 9418,
			psqtWhite:       13,
			psqtBlack:       -7,
			imbalanceWhite:  14254,
			imbalanceBlack:  12278,
			imbalanceTotal:  123,
			desc:            "White has passed d-pawn",
		},
		{
			name:            "Protected Passed Pawn",
			fen:             "rnbqkbnr/ppp1pppp/8/2PP4/8/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9170,
			pieceValueWhite: 9418,
			psqtWhite:       -19,
			psqtBlack:       -7,
			imbalanceWhite:  14254,
			imbalanceBlack:  12278,
			imbalanceTotal:  123,
			desc:            "White has protected passed pawn",
		},
		{
			name:            "Connected Passed Pawns",
			fen:             "rnbqkbnr/ppp1pppp/8/2PP4/8/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9170,
			pieceValueWhite: 9418,
			psqtWhite:       -19,
			psqtBlack:       -7,
			imbalanceWhite:  14254,
			imbalanceBlack:  12278,
			imbalanceTotal:  123,
			desc:            "White has connected passed pawns",
		},

		// Piece Mobility Tests
		{
			name:            "Trapped Knight",
			fen:             "rnbqkbnr/pppppppp/8/8/8/7N/PPPPPPPP/RNBQKB1R w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       43,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White knight trapped in corner",
		},
		{
			name:            "Trapped Bishop",
			fen:             "rnbqkbnr/ppppp1pp/5p2/8/8/7B/PPPPPPPP/RNBQK1NR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       13,
			psqtBlack:       15,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White bishop trapped by pawn",
		},
		{
			name:            "Central Knight",
			fen:             "rnbqkbnr/pppppppp/8/8/4N3/8/PPPPPPPP/RNBQKB1R w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       153,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White knight on strong central square",
		},
		{
			name:            "Open File Rook",
			fen:             "rnbqkbnr/pppp1ppp/8/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9170,
			pieceValueWhite: 9170,
			psqtWhite:       -4,
			psqtBlack:       -4,
			imbalanceWhite:  11410,
			imbalanceBlack:  11410,
			imbalanceTotal:  0,
			desc:            "White rook on open e-file",
		},

		// King Safety Tests
		{
			name:            "Exposed King",
			fen:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 7688,
			psqtWhite:       112,
			psqtBlack:       12,
			imbalanceWhite:  8657,
			imbalanceBlack:  12436,
			imbalanceTotal:  -326,
			desc:            "White king exposed with missing kingside pawns",
		},
		{
			name:            "Castled King",
			fen:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1RK1 w kq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 7688,
			psqtWhite:       258,
			psqtBlack:       12,
			imbalanceWhite:  8657,
			imbalanceBlack:  12436,
			imbalanceTotal:  -326,
			desc:            "White king safely castled",
		},
		{
			name:            "Fianchettoed Bishop",
			fen:             "rnbqkbnr/pppppppp/8/8/8/6P1/PPPPPP1P/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       10,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White preparing kingside fianchetto",
		},

		// Space and Control Tests
		{
			name:            "Central Control",
			fen:             "rnbqkbnr/pppppppp/8/8/3PP3/8/PPP2PPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       37,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White controlling center with pawns",
		},
		{
			name:            "Advanced Center",
			fen:             "rnbqkbnr/pppppppp/8/3PP3/8/8/PPP2PPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       -11,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White with advanced center pawns",
		},
		{
			name:            "Space Advantage",
			fen:             "rnbqkbnr/pppppppp/8/8/2PPP3/2N5/PP3PPP/R1BQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       131,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White with significant space advantage",
		},

		// Piece Coordination Tests
		{
			name:            "Bishop Pair",
			fen:             "rnbqk1nr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 8469,
			pieceValueWhite: 9294,
			psqtWhite:       12,
			psqtBlack:       20,
			imbalanceWhite:  12526,
			imbalanceBlack:  11239,
			imbalanceTotal:  170,
			desc:            "White has bishop pair vs bishop+knight",
		},
		{
			name:            "Doubled Rooks",
			fen:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3RK2 w kq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 3544,
			psqtWhite:       307,
			psqtBlack:       12,
			imbalanceWhite:  2572,
			imbalanceBlack:  12257,
			imbalanceTotal:  -695,
			desc:            "White has doubled rooks on open file",
		},
		{
			name:            "Knight Outpost",
			fen:             "rnbqkbnr/ppp1pppp/8/3p4/4N3/8/PPPPPPPP/RNBQKB1R w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       153,
			psqtBlack:       13,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White knight on protected outpost",
		},

		// Development and Tempo Tests
		{
			name:            "Better Development",
			fen:             "r1bqkbnr/pppppppp/2n5/8/4P3/2N1B3/PPPP1PPP/R1BQK1NR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       159,
			psqtBlack:       110,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White better developed in opening",
		},
		{
			name:            "Undeveloped Position",
			fen:             "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1",
			pieceValueBlack: 9294,
			pieceValueWhite: 9294,
			psqtWhite:       36,
			psqtBlack:       12,
			imbalanceWhite:  13228,
			imbalanceBlack:  13228,
			imbalanceTotal:  0,
			desc:            "White with poor development",
		},

		// Complex Positions
		{
			name:            "Complex Center",
			fen:             "r1bqk2r/ppp2ppp/2n2n2/1B1pp3/4P3/2PP1N2/PP3PPP/RNBQK2R w KQkq - 0 1",
			pieceValueBlack: 8469,
			pieceValueWhite: 9294,
			psqtWhite:       168,
			psqtBlack:       241,
			imbalanceWhite:  12526,
			imbalanceBlack:  11239,
			imbalanceTotal:  170,
			desc:            "Complex central tension",
		},
		{
			name:            "Attacking Position",
			fen:             "rnbqk2r/ppp2ppp/5n2/3Pp3/1b2P3/2N2N2/PPP2PPP/R1BQKB1R w KQkq - 0 1",
			pieceValueBlack: 9170,
			pieceValueWhite: 9294,
			psqtWhite:       214,
			psqtBlack:       152,
			imbalanceWhite:  12794,
			imbalanceBlack:  11844,
			imbalanceTotal:  59,
			desc:            "White with attacking chances",
		},

		// Special Cases
		{
			name:            "Rook vs Three Minors",
			fen:             "8/8/8/8/8/2nbn3/8/4R3 w - - 0 1",
			pieceValueBlack: 2387,
			pieceValueWhite: 1276,
			psqtWhite:       -5,
			psqtBlack:       122,
			imbalanceWhite:  -184,
			imbalanceBlack:  -240,
			imbalanceTotal:  3,
			desc:            "Material imbalance special case",
		},
		{
			name:            "Queen vs Three Minors",
			fen:             "8/8/8/8/8/2nbn3/8/4Q3 w - - 0 1",
			pieceValueBlack: 2387,
			pieceValueWhite: 2538,
			psqtWhite:       4,
			psqtBlack:       122,
			imbalanceWhite:  47,
			imbalanceBlack:  -240,
			imbalanceTotal:  17,
			desc:            "Roughly equal material imbalance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.ParseFEN(tt.fen)
			got := PieceValueMg(&b)
			psqt := PsqtMg(&b)
			imb := Imbalance(&b)

			if got != tt.pieceValueWhite {
				t.Errorf(
					"\nPosition: %s\nDescription: %s\nGot eval: %d\nWant eval: %d\nDiff: %d\nFEN: %s",
					tt.name,
					tt.desc,
					got,
					tt.pieceValueWhite,
					got-tt.pieceValueWhite,
					tt.fen,
				)
			}

			if psqt != tt.psqtWhite {
				t.Errorf(
					"\nPosition: %s\nDescription: %s\nGot psqt: %d\nWant psqt: %d\nDiff: %d\nFEN: %s",
					tt.name,
					tt.desc,
					psqt,
					tt.psqtWhite,
					psqt-tt.psqtWhite,
					tt.fen,
				)
			}

			if imb != tt.imbalanceWhite {
				t.Errorf(
					"\nPosition: %s\nDescription: %s\nGot imbalance: %d\nWant imbalance: %d\nDiff: %d\nFEN: %s",
					tt.name,
					tt.desc,
					imb,
					tt.imbalanceWhite,
					imb-tt.imbalanceWhite,
					tt.fen,
				)
			}

			// Test evaluation symmetry
			mirroredBoard := b.Mirror()
			mirroredEval := PieceValueMg(mirroredBoard)
			psqt = PsqtMg(mirroredBoard)
			imb = Imbalance(mirroredBoard)

			if mirroredEval != tt.pieceValueBlack {
				t.Errorf(
					"\nPosition: %s\nDescription: %s\nGot eval: %d\nWant eval: %d\nDiff: %d\nFEN: %s",
					tt.name,
					tt.desc,
					mirroredEval,
					tt.pieceValueBlack,
					mirroredEval-tt.pieceValueBlack,
					tt.fen,
				)
			}

			if psqt != tt.psqtBlack {
				t.Errorf(
					"\nPosition: %s\nDescription: %s\nGot psqt: %d\nWant psqt: %d\nDiff: %d\nFEN: %s",
					tt.name,
					tt.desc,
					psqt,
					tt.psqtBlack,
					psqt-tt.psqtBlack,
					tt.fen,
				)
			}

			if imb != tt.imbalanceBlack {
				t.Errorf(
					"\nPosition: %s\nDescription: %s\nGot imbalance: %d\nWant imbalance: %d\nDiff: %d\nFEN: %s",
					tt.name,
					tt.desc,
					imb,
					tt.imbalanceBlack,
					imb-tt.imbalanceBlack,
					tt.fen,
				)
			}

			imbTotal := ImbalanceTotal(&b, mirroredBoard)

			if imbTotal != tt.imbalanceTotal {
				t.Errorf(
					"\nPosition: %s\nDescription: %s\nGot imbalance total: %d\nWant imbalance total: %d\nDiff: %d\nFEN: %s",
					tt.name,
					tt.desc,
					imbTotal,
					tt.imbalanceTotal,
					imbTotal-tt.imbalanceTotal,
					tt.fen,
				)
			}
		})
	}
}
