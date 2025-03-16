package evaluation

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/board"
)

// TODO: FIX THESE TESTS

func TestFinalMiddleGameEvaluation(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		eval int
	}{
		// Complex Position
		{
			name: "Complex Position 1",
			fen:  "r1bq1rk1/ppp2ppp/1bn1p3/3pP3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - d6 0 10",
			eval: 1117,
		},
		{
			name: "Complex Position 2",
			fen:  "2rr2k1/pp2qppp/2n1pn2/8/3P4/2P1PNP1/PQ2PPBP/1RR3K1 w - - 0 20",
			eval: 631,
		},
		{
			name: "Complex Position 3",
			fen:  "r2q1rk1/1b2bppp/p1n1pn2/1p2P3/3P4/2NB1N2/PP3PPP/R1BQ1RK1 w - - 0 10",
			eval: 347 + 31,
		},
		{
			name: "Complex Position 4",
			fen:  "4rrk1/1ppq2bp/p1np2p1/4n3/2B1P3/2N2N2/PPP2PPP/2QRR1K1 w - - 0 17",
			eval: 86,
		},
		{
			name: "Complex Position 5",
			fen:  "r1bqkb1r/1ppn1ppp/p1n1pn2/8/2BP4/2N2N2/PPP2PPP/R1BQK2R w KQkq - 0 7",
			eval: -931,
		},
		{
			name: "Complex Position 6",
			fen:  "3r1rk1/1bqn1ppp/p2p1n2/1p6/3NP3/1BN1B3/PP2QPPP/R4RK1 w - - 0 17",
			eval: 1165 - 31,
		},
		{
			name: "Complex Position 7",
			fen:  "r2q1r1k/1p1bbp1p/p1n3p1/2P1p3/P1BP4/2NB1N2/1P3PPP/R1BQR1K1 w - - 0 16",
			eval: 2304 + 31,
		},
		{
			name: "Complex Position 8",
			fen:  "1rb1r1k1/p1q2pbp/1pn2np1/3p4/3P4/1PN1P1P1/PBQN1PBP/R4RK1 w - - 0 14",
			eval: 191,
		},
		{
			name: "Complex Position 9",
			fen:  "r1bq1rk1/pp1nppbp/2n3p1/2p5/2B1P3/5N2/PP1N1PPP/R1BQ1RK1 w - - 0 9",
			eval: -183 + 31,
		},
		{
			name: "Complex Position 10",
			fen:  "2rq1rk1/p3bppp/1pn1pn2/3p4/3P4/2N1PN2/PPQ1BPPP/2RR1RK1 w - - 0 13",
			eval: 1354,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result int
			b, _ := board.ParseFEN(tt.fen)
			ev := NewEvaluator()
			result = ev.Evaluate(&b)
			if result != tt.eval {
				t.Errorf("Middle Game Evaluation, %s: got %v, want %v", tt.name, result, tt.eval)
			}
		})
	}
}

func TestEndgameEvaluation(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		eval int
	}{
		// Complex Position
		{
			name: "Complex Position 1",
			fen:  "8/8/8/8/8/8/8/K6k w - - 0 1",
			eval: -38,
		},
		{
			name: "Complex Position 2",
			fen:  "8/8/8/8/8/8/4P3/k6K w - - 0 1",
			eval: 271,
		},
		{
			name: "Complex Position 3",
			fen:  "8/8/8/8/8/3P4/4K3/k7 w - - 0 1",
			eval: 491,
		},
		{
			name: "Complex Position 4",
			fen:  "8/8/8/8/8/2P5/8/k2K4 w - - 0 1",
			eval: 315,
		},
		{
			name: "Complex Position 5",
			fen:  "8/8/8/8/8/3K2Q1/8/k7 w - - 0 1",
			eval: 3009,
		},
		{
			name: "Complex Position 6",
			fen:  "8/8/8/8/8/K7/2R5/k7 w - - 1 1",
			eval: 1659,
		},
		{
			name: "Complex Position 8",
			fen:  "8/8/8/8/8/8/6N1/k5K1 w - - 0 1",
			eval: 795,
		},
		{
			name: "Complex Position 10",
			fen:  "8/8/8/8/3K4/8/8/4k3 w - - 0 1",
			eval: 0,
		},
		{
			name: "Complex Position 11",
			fen:  "8/3k4/3P4/8/8/8/3K4/8 w - - 0 1",
			eval: 0,
		},
		{
			name: "Complex Position 12",
			fen:  "8/8/8/5p2/8/2K5/8/6k1 w - - 0 1",
			eval: -282,
		},
		{
			name: "Complex Position 13",
			fen:  "6k1/8/8/8/6P1/8/8/6K1 w - - 0 1",
			eval: 113,
		},
		{
			name: "Complex Position 14",
			fen:  "8/8/8/8/3K4/2b5/8/4k3 w - - 0 1",
			eval: -625,
		},
		{
			name: "Complex Position 15",
			fen:  "8/8/8/8/2B5/3K4/8/4k3 w - - 0 1",
			eval: 930,
		},
		{
			name: "Complex Position 16",
			fen:  "8/8/8/8/1B6/6k1/8/6K1 w - - 0 1",
			eval: 708,
		},
		{
			name: "Complex Position 17",
			fen:  "8/8/8/8/8/2N5/K7/k7 w - - 0 1",
			eval: 833,
		},
		{
			name: "Complex Position 18",
			fen:  "8/8/2n5/8/8/8/K7/k7 w - - 0 1",
			eval: -674,
		},
		{
			name: "Complex Position 19",
			fen:  "8/4k3/8/3n4/8/8/3K4/8 w - - 0 1",
			eval: -692,
		},
		{
			name: "Complex Position 20",
			fen:  "8/8/8/5p2/6P1/8/6K1/5k2 w - - 0 1",
			eval: 154,
		},
		{
			name: "Complex Position 21",
			fen:  "2r3k1/1pp4p/p3pn2/1P6/2Q2P2/P2R4/5P1P/3R2K1 w - - 0 1",
			eval: 3587 + 22,
		},
		{
			name: "Complex Position 22",
			fen:  "7r/5pp1/4p2p/2k1n3/8/2PBK1P1/1P3P1P/R7 b - - 0 1",
			eval: 494,
		},
		{
			name: "Complex Position 23",
			fen:  "r1b3k1/pp1n1ppp/2p5/1P2p3/3P1N2/4B3/PP3PPP/R3K2R w - - 0 1",
			eval: 1341,
		},
		{
			name: "Complex Position 24",
			fen:  "8/p1r2p1k/1p2p1p1/n1p5/3P4/P1P1P1P1/1P3P1K/1BRN4 w - - 0 1",
			eval: 1080,
		},
		{
			name: "Complex Position 25",
			fen:  "r3k2r/1pqn1ppp/p1pbp3/8/3P4/PNP1P3/1B3PPP/R1BQK2R w - - 0 1",
			eval: 830,
		},
		{
			name: "Complex Position 26",
			fen:  "4r1k1/1q3pp1/p2bpn1p/1p2N3/3Q4/1P3P1P/PB3PP1/2RR2K1 w - - 0 1",
			eval: 1687 + 22,
		},
		{
			name: "Complex Position 27",
			fen:  "6k1/1p3ppp/p2b4/4p3/P2pP3/2P3P1/1P1N1P1P/4BK2 w - - 0 1",
			eval: 883,
		},
		{
			name: "Complex Position 28",
			fen:  "4rrk1/1bq2ppp/p3pn2/1pp1P3/3P4/PNP3P1/1B2QP1P/1R1R1NK1 b - - 0 1",
			eval: 880,
		},
		{
			name: "Complex Position 29",
			fen:  "6k1/2q3pp/2P1pn2/r7/2P5/1P2PNP1/P3QPKP/R7 w - - 0 1",
			eval: 1577,
		},
		{
			name: "Complex Position 30",
			fen:  "8/pbpp2p1/1pn1pk1p/4n3/2PPP3/1PN1BPP1/P5KP/3R4 w - - 0 1",
			eval: 923,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result int
			b, _ := board.ParseFEN(tt.fen)
			ev := NewEvaluator()
			result = ev.Evaluate(&b)
			if result != tt.eval {
				t.Errorf("Endgame Evaluation, %s: got %v, want %v", tt.name, result, tt.eval)
			}
		})
	}
}
