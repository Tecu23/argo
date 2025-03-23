# UCI Protocol Implementation

## Overview

The Universal Chess Interface (UCI) is a standard protocol for communication
between chess engines and user interfaces. ArGO implements UCI to enable
compatibility with popular chess GUIs such as Arena, Fritz, and Cutechess.
This document details ArGO's UCI implementation.

## Protocol Architecture

The UCI implementation in ArGO is encapsulated in the `uci` package, with the
core functionality in the `Protocol` struct:

```go
type Protocol struct {
    name         string              // Engine name
    author       string              // Engine author
    version      string              // Engine version
    options      []Option            // Configurable engine options
    engine       Engine              // The underlying chess engine instance
    boards       []board.Board       // Stack of board positions
    thinking     bool                // Indicates if engine is currently searching
    engineOutput chan SearchInfo     // Channel for search updates
    cancel       context.CancelFunc  // For canceling the search
}
```

## UCI Command Loop

The protocol's main loop processes commands from standard input and handles search
output from the engine:

```go
func (uci *Protocol) Run(logger *log.Logger) {
    commands := make(chan string)

    // Goroutine for reading commands from stdin
    go func() {
        defer close(commands)
        readCommands(commands)
    }()

    var searchResult SearchInfo
    for {
        select {
        // Handle search results from engine
        case si, ok := <-uci.engineOutput:
            if ok {
                // Print intermediate search info
                fmt.Println(searchInfoToUci(si))
                searchResult = si
            } else {
                // Engine finished searching
                if len(searchResult.MainLine) != 0 {
                    fmt.Printf("bestmove %v\n", searchResult.MainLine[0])
                }
                // Reset state
                uci.thinking = false
                uci.cancel = nil
                uci.engineOutput = nil
                searchResult = SearchInfo{}
            }

        // Handle incoming commands
        case commandLine, ok := <-commands:
            if !ok {
                return
            }
            err := uci.handle(commandLine)
            if err != nil {
                logger.Println(err)
            }
        }
    }
}
```

## Command Handling

The `handle` function parses and dispatches UCI commands:

```go
func (uci *Protocol) handle(commandLine string) error {
    fields := strings.Fields(commandLine)
    if len(fields) == 0 {
        return nil // Empty line
    }

    commandName := fields[0]
    fields = fields[1:]

    // If engine is thinking, only allow 'stop'
    if uci.thinking {
        if commandName == "stop" {
            uci.cancel()
            return nil
        }
        return errors.New("search still running")
    }

    // Map command to handler
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
```

## Core UCI Commands

### 1. UCI Command

The `uci` command identifies the engine and lists available options:

```go
func (uci *Protocol) uciCommand(_ []string) error {
    fmt.Printf("id name %s %s\n", uci.name, uci.version)
    fmt.Printf("id author %s\n", uci.author)

    // Print all available options
    for _, option := range uci.options {
        fmt.Println(option.UciString())
    }

    fmt.Println("uciok")
    return nil
}
```

### 2. Set Option Command

The `setoption` command allows configuring engine parameters:

```go
func (uci *Protocol) setOptionCommand(fields []string) error {
    if len(fields) < 4 {
        return errors.New("invalid setoption arguments")
    }

    name, value := fields[1], fields[3]

    // Find and set the matching option
    for _, option := range uci.options {
        if strings.EqualFold(option.UciName(), name) {
            return option.Set(value)
        }
    }

    return errors.New("unhandled option")
}
```

### 3. IsReady Command

The `isready` command checks if the engine is ready to receive commands:

```go
func (uci *Protocol) isReadyCommand(_ []string) error {
    uci.engine.Prepare()
    fmt.Println("readyok")
    return nil
}
```

### 4. Position Command

The `position` command sets up a chess position:

```go
func (uci *Protocol) positionCommand(fields []string) error {
    args := fields
    token := args[0]
    var fen string
    movesIndex := findIndexString(args, "moves")

    // Handle "startpos" or "fen" positions
    if token == "startpos" {
        fen = StartPosition
    } else if token == "fen" {
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

    // Apply moves if specified
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
```

### 5. Go Command

The `go` command starts the search:

```go
func (uci *Protocol) goCommand(fields []string) error {
    limits := parseLimits(fields)
    ctx, cancel := context.WithCancel(context.TODO())
    uci.cancel = cancel
    uci.thinking = true
    uci.engineOutput = make(chan SearchInfo, 3)

    // Run search async
    go func() {
        searchResult := uci.engine.Search(ctx, SearchParams{
            Boards: uci.boards,
            Limits: limits,
            Progress: func(si SearchInfo) {
                // Send intermediate info
                select {
                case uci.engineOutput <- si:
                default:
                }
            },
        })

        // Send final result
        uci.engineOutput <- searchResult
        close(uci.engineOutput)
    }()

    return nil
}
```

### 6. UCI New Game Command

The `ucinewgame` command resets the engine for a new game:

```go
func (uci *Protocol) uciNewGameCommand(_ []string) error {
    uci.engine.Clear()
    return nil
}
```

## Search Limits Parsing

The engine supports various search constraints specified in the `go` command:

```go
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
```

## Search Info Formatting

The `searchInfoToUci` function formats search results in UCI protocol format:

```go
func searchInfoToUci(si SearchInfo) string {
    sb := &strings.Builder{}

    // Basic info: depth and score
    fmt.Fprintf(sb, "info depth %v", si.Depth)
    if si.Score.Mate != 0 {
        fmt.Fprintf(sb, " score mate %v", si.Score.Mate)
    } else {
        fmt.Fprintf(sb, " score cp %v", si.Score.Centipawns)
    }

    // Performance metrics
    timeMs := si.Time.Milliseconds()
    nps := si.Nodes * 1000 / (timeMs + 1)
    fmt.Fprintf(sb, " nodes %v time %v nps %v", si.Nodes, timeMs, nps)

    // Principal variation
    if len(si.MainLine) != 0 {
        fmt.Fprintf(sb, " pv")
        for _, move := range si.MainLine {
            sb.WriteString(" ")
            sb.WriteString(move.String())
        }
    }

    return sb.String()
}
```

## Engine Options

ArGO implements a flexible option system that can handle different option types:

```go
type Option interface {
    UciName() string
    UciString() string
    Set(s string) error
}

// Boolean option implementation
type BoolOption struct {
    Name  string
    Value *bool
}

func (opt *BoolOption) UciName() string {
    return opt.Name
}

func (opt *BoolOption) UciString() string {
    return fmt.Sprintf("option name %v type %v default %v",
        opt.Name, "check", *opt.Value)
}

func (opt *BoolOption) Set(s string) error {
    v, err := strconv.ParseBool(s)
    if err != nil {
        return err
    }
    *opt.Value = v
    return nil
}
```

Additional option types can be implemented (int, string, button, combo) following
the same interface.

## UCI Engine Interface

ArGO defines an `Engine` interface that must be implemented by the actual chess engine:

```go
type Engine interface {
    Prepare()
    Clear()
    Search(ctx context.Context, searchParams SearchParams) SearchInfo
}
```

This abstraction allows for potentially swapping different engine implementations
while maintaining the same UCI interface.

## Context-Based Cancellation

ArGO uses Go's context package for search cancellation:

```go
ctx, cancel := context.WithCancel(context.TODO())
uci.cancel = cancel

// Later, to stop the search:
uci.cancel()
```

This provides a clean mechanism for terminating searches when requested by the user.

## Timing and Search Parameters

The `SearchParams` struct encapsulates all information needed for a search:

```go
type SearchParams struct {
    Boards   []board.Board            // Board positions (current + history)
    Limits   LimitsType               // Time and depth constraints
    Progress func(si SearchInfo)      // Callback for search progress
}
```

And `LimitsType` defines the search constraints:

```go
type LimitsType struct {
    Ponder         bool    // Search in ponder mode
    Infinite       bool    // Search until stopped
    WhiteTime      int     // White's remaining time in milliseconds
    BlackTime      int     // Black's remaining time in milliseconds
    WhiteIncrement int     // White's increment per move in milliseconds
    BlackIncrement int     // Black's increment per move in milliseconds
    MoveTime       int     // Fixed time per move in milliseconds
    MovesToGo      int     // Moves until next time control
    Depth          int     // Maximum search depth
    Nodes          int     // Maximum nodes to search
    Mate           int     // Search for mate in x moves
}
```

## Pondering Support

Pondering allows the engine to think during the opponent's time. While
ArGO includes a `ponderhit` command handler, the implementation is currently minimal:

```go
func (uci *Protocol) ponderhitCommand(_ []string) error {
    return errors.New("not implemented")
}
```

Full pondering support would require additional logic to:

1. Continue searching after a "ponderhit" command
2. Adjust time management based on pondering results
3. Handle ponder hits for predicted vs. actual opponent moves

## Typical UCI Session

A typical UCI session follows this flow:

1. **Initialization**:

   ```txt
   uci
   id name ArGO 1.0
   id author Tecu23
   option name Hash type spin default 32 min 1 max 1024
   option name Threads type spin default 1 min 1 max 128
   uciok
   isready
   readyok
   ```

2. **Game Setup**:

   ```txt
   ucinewgame
   position startpos
   ```

3. **Search and Move**:

   ```txt
   go wtime 300000 btime 300000 winc 2000 binc 2000
   info depth 1 score cp 50 nodes 30 time 5 nps 6000 pv e2e4
   info depth 2 score cp 32 nodes 115 time 10 nps 11500 pv e2e4 e7e5
   ...
   bestmove e2e4
   ```

4. **Position Update and Next Search**:

   ```txt
   position startpos moves e2e4 e7e5
   go wtime 299000 btime 299000 winc 2000 binc 2000
   ...
   bestmove g1f3
   ```

## Limitations and Future Improvements

1. **Full Pondering Support**: Implementing complete pondering capabilities
2. **Multi-PV Support**: Reporting multiple principal variations
3. **Improved Options**: Adding more configurable engine parameters
4. **Analysis Features**: Providing detailed position analysis beyond
   just the best move
5. **Custom Commands**: Supporting engine-specific commands for debugging
   or advanced features
