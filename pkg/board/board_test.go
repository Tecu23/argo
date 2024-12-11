// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	util.InitFen2Sq() // Make sure square mappings are initialized
}
