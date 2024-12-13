package evaluation

const (
	// Piece base values in centipawns
	PawnValue   = 100
	KnightValue = 320
	BishopValue = 330
	RookValue   = 500
	QueenValue  = 900
	KingValue   = 20_000

	// Special position scores
	DrawScore     = 0
	MateScore     = 1_000_000
	InfiniteScore = 9_999_999

	// Phase scores for game stage detection
	PawnPhase   = 0
	KnightPhase = 1
	BishopPhase = 1
	RookPhase   = 2
	QueenPhase  = 4

	// Evaluation weights
	DefaultMaterialWeight      = 1.0
	DefaultMobilityWeight      = 0.1
	DefaultPawnStructureWeight = 0.3
	DefaultPositionalWeight    = 0.5

	// Common penalties and bonuses
	BishopPairBonus     = 50
	RookOpenFileBonus   = 25
	RookSemiOpenBonus   = 10
	DoublePawnPenalty   = -10
	IsolatedPawnPenalty = -20
	PassedPawnBonus     = 20
)
