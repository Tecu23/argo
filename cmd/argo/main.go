// Package main is the entry point of the program
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/board"
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

	file, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}

	if err := pprof.StartCPUProfile(file); err != nil {
		log.Fatal(err)
	}

	defer pprof.StopCPUProfile()

	if debug {
		b, _ := board.ParseFEN(
			"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		)

		evaluator := nnue.NewEvaluator()
		evaluator.Reset(&b)

		score := evaluator.Evaluate(&b)
		fmt.Println("Initial Score", score)

		moves := b.GenerateMoves()
		for _, move := range moves {
			copyB := b.CopyBoard()
			if !copyB.MakeMove(move, board.AllMoves) {
				continue
			}
			evaluator.ProcessMove(&copyB, move)

			// Now let's evaluate this position
			score := evaluator.Evaluate(&copyB)
			fmt.Printf("Score After move: %s, %d\n", move.String(), score)

			evaluator.PopAccumulation()
		}

		return
	}

	options := engine.NewOptions()
	engine := engine.NewEngine(options)

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

	err := nnue.LoadWeights("./default.net")
	if err != nil {
		fmt.Printf("Error loading weights: %v\n", err)
		return
	}
}
