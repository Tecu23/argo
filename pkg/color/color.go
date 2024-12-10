// Package color defines the representation of chess piece colors and
// provides utility functions to work with them. The main colors are
// WHITE and BLACK, with an optional BOTh state for special scenarios.
package color

// Color represents the color of a chess piece or side.
// It is typically either WHITE or BLACK, but BOTh is
// provided as a special case.
type Color int

const (
	// WHITE represents the White side in chess.
	WHITE = Color(0)
	// BLACK represents the Black side in chess.
	BLACK = Color(1)
	// BOTH can represent a scenario where both sides are considered,
	// for example in certain combined calculations or testing.
	BOTH = Color(2)
)

// Opp returns the opposite Color. If the Color is WHITE, it returns BLACK.
// If it is BLACK, it returns WHITE. If it is BOTh, it returns BLACK due to
// the XOR operation, but generally BOTh is not used in normal gameplay logic.
func (c Color) Opp() Color {
	if c == BOTH {
		return BOTH
	}
	return c ^ 0x1
}

// String returns a single-character string representing the Color.
// "W" for White, "B" for Black. If called on BOTh, it currently returns "B",
// as there is no special handling for BOTh in the string representation.
func (c Color) String() string {
	if c == WHITE {
		return "W"
	}
	return "B"
}
