// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package uci

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// Engine is the interface that any chess engine implementation must follow.
// It provides methods to prepare the engine, clear internal state, and run a search.
type Engine interface {
	Prepare()
	Clear()
	Search(ctx context.Context, searchParams SearchParams) SearchInfo
}

// Protocol represents the UCI (Universal Chess Interface) protocol implementation.
// It manages communication between the UI (like a chess GUI) and the Engine.
type Protocol struct {
	name         string
	author       string
	version      string
	options      []Option           // a list of engine options that can be set via the UCI commands
	engine       Engine             // The underlying chess engine instance
	boards       []board.Board      // The stack of boards representing the current game state
	thinking     bool               // Indicates if the engine is currently searching
	engineOutput chan SearchInfo    // Channel used to receive SearchInfo updates from the engine
	cancel       context.CancelFunc // Used to cancel ongoing searches
}

// New creates a new Protocol instance with given engine name, author, version, and options.
// It also initializes the board to the standard chess starting position.
func New(name, author, version string, engine Engine, options []Option) *Protocol {
	initBoard, err := board.ParseFEN(StartPosition)
	if err != nil {
		panic(err)
	}
	return &Protocol{
		name:    name,
		author:  author,
		version: version,
		engine:  engine,
		options: options,
		boards:  []board.Board{initBoard},
	}
}

// Run starts the main UCI loop, listening for incoming commands from stdin
// and handling them. It also listens for the engine's search output
func (uci *Protocol) Run(logger *log.Logger) {
	commands := make(chan string)

	// This goroutine coninuously reads lines from stdin and sends them to the commands channel
	go func() {
		defer close(commands)
		readCommands(commands)
	}()

	var searchResult SearchInfo
	for {
		select {
		// If the engine sends a SearchInfo on engineOutput:
		case si, ok := <-uci.engineOutput:
			if ok {
				// Print the intermediate search info in UCI format
				fmt.Println(searchInfoToUci(si))
				searchResult = si
			} else {
				// Engine finished searching (channel closed)
				if len(searchResult.MainLine) != 0 {
					// Print the best move found
					fmt.Printf("bestmove %v\n", searchResult.MainLine[0])
				}
				// Reset state
				uci.thinking = false
				uci.cancel = nil
				uci.engineOutput = nil
				searchResult = SearchInfo{}
			}
		// When a new command arrives from stdin:
		case commandLine, ok := <-commands:
			if !ok {
				// uci quit
				return
			}
			// Handle the incomming command line
			err := uci.handle(commandLine)
			if err != nil {
				logger.Println(err)
			}

		}
	}
}

// readCommands reads input lines from stdin and sends them to the provided channel.
// It stops reading if the "quit" command is encountered
func readCommands(commands chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		commandLine := scanner.Text()
		if commandLine == "quit" {
			// Stop reading commands on "quit"
			return
		}

		if commandLine != "" {
			commands <- commandLine
		}
	}
}

// handle takes a single UCI command line, parses it, and executes the corresponding handler.
func (uci *Protocol) handle(commandLine string) error {
	fields := strings.Fields(commandLine)
	if len(fields) == 0 {
		return nil // Empty line, do nothing
	}

	commandName := fields[0]
	fields = fields[1:]

	// If engine is currently searching (thinking), only certain commands like "stop" are allowed
	if uci.thinking {
		if commandName == "stop" {
			uci.cancel()
			// Stop the ongoing search
			return nil
		}
		return errors.New("search still run")
	}

	/// Map commandName to the appropriate handler function
	var h func(fields []string) error

	switch commandName {
	case "uci":
		h = uci.uciCommand
	case "setoption":
		h = uci.setOptionCommand
	case "isready":
		h = uci.isReadyCommand
	case "position":
		h = uci.positionCommand
	case "go":
		h = uci.goCommand
	case "ucinewgame":
		h = uci.uciNewGameCommand
	case "ponderhit":
		h = uci.ponderhitCommand
	}

	if h == nil {
		return errors.New("command not found")
	}

	return h(fields)
}

// uciCommand handles the "uci" command, which requests engine identification and available options.
func (uci *Protocol) uciCommand(_ []string) error {
	fmt.Printf("id name %s %s\n", uci.name, uci.version)
	fmt.Printf("id author %s\n", uci.author)
	// Print all available options in UCI format
	for _, option := range uci.options {
		fmt.Println(option.UciString())
	}
	fmt.Println("uciok")
	return nil
}

// setOptionCommand handle the "setoption", allowing the GUI to change engine output
func (uci *Protocol) setOptionCommand(fields []string) error {
	if len(fields) < 4 {
		// Expected output: setoption name <name> value <value>
		return errors.New("invalid setoption arguments")
	}

	name, value := fields[1], fields[3]
	// Try to find and set the matching option
	for _, option := range uci.options {
		if strings.EqualFold(option.UciName(), name) {
			return option.Set(value)
		}
	}

	return errors.New("unhandled option")
}

// isReadyCommand handles the "isready" command.
// The engine should do any necessary initialization and then print "readyok".
func (uci *Protocol) isReadyCommand(_ []string) error {
	uci.engine.Prepare()
	fmt.Println("readyok")
	return nil
}

// positionCommand sets up a position in the engine. It can either be "startpos" or a custom "fen"
// followed by "moves" for a sequence of moves. After processing, the engine's internal board state is updated
func (uci *Protocol) positionCommand(fields []string) error {
	args := fields
	token := args[0]
	var fen string
	movesIndex := findIndexString(args, "moves")

	// Handle "startpos" or "fen" positions
	if token == "startpos" {
		fen = StartPosition
	} else if token == "fen" {
		// If "fen" is specified, parse everything until "moves" as the FEN string
		if movesIndex == -1 {
			fen = strings.Join(args[1:], " ")
		} else {
			fen = strings.Join(args[1:movesIndex], " ")
		}
	} else {
		return errors.New("unknown position command")
	}

	// Parse the given FEN into a Board structure
	b, err := board.ParseFEN(fen)
	if err != nil {
		return err
	}

	boards := []board.Board{b}
	// If there are moves following "moves", apply them sequentially to reach the final position
	if movesIndex >= 0 && movesIndex+1 < len(args) {
		for _, smove := range args[movesIndex+1:] {
			newBoard, ok := boards[len(boards)-1].ParseMove(smove)
			if !ok {
				return errors.New("parse move failed")
			}

			boards = append(boards, newBoard)
		}
	}
	uci.boards = boards
	return nil
}

// goCommand starts the search with the given time/move constraints (limits).
// It creates a cancellable context and runs the search in a separate goroutine.
// Intermediate and final results are sent to engineOutput channel.
func (uci *Protocol) goCommand(fields []string) error {
	limits := parseLimits(fields)
	ctx, cancel := context.WithCancel(context.TODO())
	uci.cancel = cancel
	uci.thinking = true
	uci.engineOutput = make(chan SearchInfo, 3)

	// Run the search async
	go func() {
		searchResult := uci.engine.Search(ctx, SearchParams{
			Boards: uci.boards,
			Limits: limits,
			Progress: func(si SearchInfo) {
				// Send intermediate search info, but don't block if channel is full
				select {
				case uci.engineOutput <- si:
				default:
				}
			},
		})
		// After the search completes, send the final result and close the channel
		uci.engineOutput <- searchResult
		close(uci.engineOutput)
	}()
	return nil
}

// uciNewGameCommand signals that a new game is starting, so the engine should reset it internal state.
func (uci *Protocol) uciNewGameCommand(_ []string) error {
	uci.engine.Clear()
	return nil
}

// ponderhitCommand is a UCI Command that indicates the opponent has made a move and
// the engine should start pondering and start searching. Not yet implemented
func (uci *Protocol) ponderhitCommand(_ []string) error {
	return errors.New("not implemented")
}

// searchInfoToUci converts a SearchInfo structure into a string that follows UCI's "info" line format.
// It includes depth, score, nodes, time, nps, and the principal variation (pv)
func searchInfoToUci(si SearchInfo) string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "info depth %v", si.Depth)
	if si.Score.Mate != 0 {
		fmt.Fprintf(sb, " score mate %v", si.Score.Mate)
	} else {
		fmt.Fprintf(sb, " score cp %v", si.Score.Centipawns)
	}

	timeMs := si.Time.Milliseconds()
	nps := si.Nodes * 1000 / (timeMs + 1)
	fmt.Fprintf(sb, " nodes %v time %v nps %v", si.Nodes, timeMs, nps)
	if len(si.MainLine) != 0 {
		fmt.Fprintf(sb, " pv")
		for _, move := range si.MainLine {
			sb.WriteString(" ")
			sb.WriteString(move.String())
		}
	}

	return sb.String()
}

// parseLimits parses the arguments from "go" command to extract time controls, depth, nodes, etc.,
// and returns them in a LimitsType struct.
func parseLimits(args []string) (result LimitsType) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "ponder":
			result.Ponder = true
		case "wtime":
			result.WhiteTime, _ = strconv.Atoi(args[i+1])
		case "btime":
			result.BlackTime, _ = strconv.Atoi(args[i+1])
		case "winc":
			result.WhiteIncrement, _ = strconv.Atoi(args[i+1])
		case "binc":
			result.BlackIncrement, _ = strconv.Atoi(args[i+1])
		case "movestogo":
			result.MovesToGo, _ = strconv.Atoi(args[i+1])
		case "depth":
			result.Depth, _ = strconv.Atoi(args[i+1])
		case "nodes":
			result.Nodes, _ = strconv.Atoi(args[i+1])
		case "mate":
			result.Mate, _ = strconv.Atoi(args[i+1])
		case "movetime":
			result.MoveTime, _ = strconv.Atoi(args[i+1])
		case "infinite":
			result.Infinite = true
		}
	}
	return
}

// findIndexString searches for a specific string in a slice and returns its index.
// If not found, returns -1
func findIndexString(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
