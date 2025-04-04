// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

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

var (
	gitCommit = "?"
	buildDate = "?"
)

var debug bool

func main() {
	flag.BoolVar(&debug, "debug", false, "specifies if engine ran on debug mode")
	flag.Parse()
	initHelpers()

	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

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

	if err := nnue.InitializeNNUE(); err != nil {
		log.Fatalf("Error initializing NNUE: %v", err)
	}
}
