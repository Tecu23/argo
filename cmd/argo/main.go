// Package main is the entry point of the program
package main

import (
	"flag"
	"log"
	"os"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/engine"
	"github.com/Tecu23/argov2/pkg/nnue"
	"github.com/Tecu23/argov2/pkg/uci"
	"github.com/Tecu23/argov2/pkg/util"
)

const (
	name   = "ArGO"
	author = "Tecu23"
)

var version = "?"

var debug bool

func main() {
	flag.BoolVar(&debug, "debug", false, "specifies if engine ran on debug mode")
	flag.Parse()
	initHelpers()

	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	options := engine.NewOptions()
	engine := engine.NewEngine(options)

	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	logger.Fatalf("cpu profiling file could not be created: %v", err)
	// }
	//
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	logger.Fatalf("could not start profiling: %v", err)
	// }
	//
	// defer pprof.StopCPUProfile()
	//
	// if debug {
	// 	ev := nnue.NewEvaluator()
	// 	b, _ := board.ParseFEN(constants.StartPosition)
	// 	ev.Reset(&b)
	//
	// 	score := ev.Evaluate(&b)
	//
	// 	fmt.Printf("Initial score for initial position: %d\n", score)
	//
	// 	mvs := b.GenerateMoves()
	// 	maxScore := -1_000_000
	// 	bestMove := move.NoMove
	// 	for _, mv := range mvs {
	// 		cpy := b.CopyBoard()
	// 		if !cpy.MakeMove(mv, board.AllMoves) {
	// 			continue
	// 		}
	// 		ev.ProcessMove(&cpy, mv)
	//
	// 		s1 := ev.Evaluate(&cpy)
	// 		if s1 > maxScore {
	// 			maxScore = s1
	// 			bestMove = mv
	// 		}
	//
	// 		mvs2 := cpy.GenerateMoves()
	// 		minScore := 1_000_000
	// 		bm2 := move.NoMove
	// 		for _, mv2 := range mvs2 {
	// 			cpy2 := cpy.CopyBoard()
	// 			if !cpy2.MakeMove(mv2, board.AllMoves) {
	// 				continue
	// 			}
	// 			ev.ProcessMove(&cpy2, mv2)
	//
	// 			s2 := ev.Evaluate(&cpy2)
	// 			if s2 < minScore {
	// 				minScore = s2
	// 				bm2 = mv2
	// 			}
	//
	// 			ev.PopAccumulation()
	// 		}
	// 		fmt.Printf(
	// 			"Min Score at depth 2 for initial position, after move %s: %d, %s\n",
	// 			mv,
	// 			minScore,
	// 			bm2,
	// 		)
	// 		ev.PopAccumulation()
	// 	}
	// 	fmt.Printf("Max Score at depth 1 for initial position: %d, %s\n", maxScore, bestMove)
	//
	// 	return
	// }
	//
	protocol := uci.New(name, author, version, engine, []uci.Option{})
	protocol.Run(logger)
}

func initHelpers() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(constants.Bishop)
	attacks.InitSliderPiecesAttacks(constants.Rook)

	util.InitFen2Sq()

	hash.Init()

	if err := nnue.InitializeNNUE(); err != nil {
		log.Fatalf("Error initializing NNUE: %v", err)
	}
}
