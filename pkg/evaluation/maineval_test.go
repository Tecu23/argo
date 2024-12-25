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
			result = MiddleGameEvaluation(&b, false)
			if result != tt.eval {
				t.Errorf("Middle Game Evaluation, %s: got %v, want %v", tt.name, result, tt.eval)
			}
		})
	}
}
