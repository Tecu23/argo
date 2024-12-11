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
	"time"

	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/move"
)

type UciScore struct {
	Centipawns int
	Mate       int
}

type SearchInfo struct {
	Score    UciScore
	Depth    int
	Nodes    int64
	Time     time.Duration
	MainLine []move.Move
}

type LimitsType struct {
	Ponder         bool
	Infinite       bool
	WhiteTime      int
	BlackTime      int
	WhiteIncrement int
	BlackIncrement int
	MoveTime       int
	MovesToGo      int
	Depth          int
	Nodes          int
	Mate           int
}

type SearchParams struct {
	Boards   []board.Board
	Limits   LimitsType
	Progress func(si SearchInfo)
}

type Engine interface {
	Prepare()
	Clear()
	Search(ctx context.Context, searchParams SearchParams) SearchInfo
}

type Protocol struct {
	name         string
	author       string
	version      string
	options      []Option
	engine       Engine
	boards       []board.Board
	thinking     bool
	engineOutput chan SearchInfo
	cancel       context.CancelFunc
}

func New(name, author, version string, engine Engine, options []Option) *Protocol {
	initBoard, err := board.ParseFEN(constants.StartPosition)
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

func (uci *Protocol) Run(logger *log.Logger) {
	commands := make(chan string)

	go func() {
		defer close(commands)
		readCommands(commands)
	}()

	var searchResult SearchInfo
	for {
		select {
		case si, ok := <-uci.engineOutput:
			if ok {
				fmt.Println(searchInfoToUci(si))
				searchResult = si
			} else {
				if len(searchResult.MainLine) != 0 {
					fmt.Printf("bestmove %v\n", searchResult.MainLine[0])
				}
				uci.thinking = false
				uci.cancel = nil
				uci.engineOutput = nil
				searchResult = SearchInfo{}
			}
		case commandLine, ok := <-commands:
			if !ok {
				// uci quit
				return
			}
			err := uci.handle(commandLine)
			if err != nil {
				logger.Println(err)
			}

		}
	}
}

func readCommands(commands chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		commandLine := scanner.Text()
		if commandLine == "quit" {
			return
		}

		if commandLine != "" {
			commands <- commandLine
		}
	}
}

func (uci *Protocol) handle(commandLine string) error {
	fields := strings.Fields(commandLine)
	if len(fields) == 0 {
		return nil
	}

	commandName := fields[0]
	fields = fields[1:]

	if uci.thinking {
		if commandName == "stop" {
			uci.cancel()
			return nil
		}
		return errors.New("search still run")
	}

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

func (uci *Protocol) uciCommand(fields []string) error {
	fmt.Printf("id name %s %s\n", uci.name, uci.version)
	fmt.Printf("id author %s\n", uci.author)
	for _, option := range uci.options {
		fmt.Println(option.UciString())
	}
	fmt.Println("uciok")
	return nil
}

func (uci *Protocol) setOptionCommand(fields []string) error {
	if len(fields) < 4 {
		return errors.New("invalid setoption arguments")
	}

	name, value := fields[1], fields[3]
	for _, option := range uci.options {
		if strings.EqualFold(option.UciName(), name) {
			return option.Set(value)
		}
	}

	return errors.New("unhandled option")
}

func (uci *Protocol) isReadyCommand(fields []string) error {
	uci.engine.Prepare()
	fmt.Println("readyok")
	return nil
}

func (uci *Protocol) positionCommand(fields []string) error {
	args := fields
	token := args[0]
	var fen string
	movesIndex := findIndexString(args, "moves")
	if token == "startpos" {
		fen = constants.StartPosition
	} else if token == "fen" {
		if movesIndex == -1 {
			fen = strings.Join(args[1:], " ")
		} else {
			fen = strings.Join(args[1:movesIndex], " ")
		}
	} else {
		return errors.New("unknown position command")
	}

	b, err := board.ParseFEN(fen)
	if err != nil {
		return err
	}

	boards := []board.Board{b}
	if movesIndex >= 0 && movesIndex+1 < len(args) {
		for _, smove := range args[movesIndex+1:] {
			newBoard, ok := boards[len(boards)-1].MakeMoveLAN(smove)
			if !ok {
				return errors.New("parse move failed")
			}
			boards = append(boards, newBoard)
		}
	}
	uci.boards = boards
	return nil
}

func (uci *Protocol) goCommand(fields []string) error {
	limits := parseLimits(fields)
	ctx, cancel := context.WithCancel(context.TODO())
	uci.cancel = cancel
	uci.thinking = true
	uci.engineOutput = make(chan SearchInfo, 3)
	go func() {
		searchResult := uci.engine.Search(ctx, SearchParams{
			Boards: uci.boards,
			Limits: limits,
			Progress: func(si SearchInfo) {
				select {
				case uci.engineOutput <- si:
				default:
				}
			},
		})
		uci.engineOutput <- searchResult
		close(uci.engineOutput)
	}()
	return nil
}

func (uci *Protocol) uciNewGameCommand(fields []string) error {
	uci.engine.Clear()
	return nil
}

func (uci *Protocol) ponderhitCommand(fields []string) error {
	return errors.New("not implemented")
}

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

func findIndexString(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
