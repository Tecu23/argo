// Package nnue keeps the NNUE (Efficiently Updated Neural Network) responsible for
// evaluation the current position
package nnue

import (
	"encoding/binary"
	"fmt"
	"os"
)

// LoadWeights loads network weights from a binary file.
// The file is expected to contain input weights, input bias, hidden weights, and hidden bias in sequence.
func LoadWeights(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	// Read input weights
	for i := 0; i < InputSize; i++ {
		for j := 0; j < HiddenSize; j++ {

			var val int16
			if err := binary.Read(file, binary.LittleEndian, &val); err != nil {
				return fmt.Errorf("error reading input weights: %v", err)
			}
			InputWeights[i][j] = val
		}
	}

	// Read input bias
	for i := 0; i < HiddenSize; i++ {
		var val int16
		if err := binary.Read(file, binary.LittleEndian, &val); err != nil {
			return fmt.Errorf("error reading input weights: %v", err)
		}
		InputBias[i] = val
	}

	// Read hidden weights
	for i := 0; i < OutputSize; i++ {
		for j := 0; j < HiddenDSize; j++ {
			var val int16
			if err := binary.Read(file, binary.LittleEndian, &val); err != nil {
				return fmt.Errorf("error reading input weights: %v", err)
			}
			HiddenWeights[i][j] = val
		}
	}

	// Read hidden bias
	for i := 0; i < OutputSize; i++ {
		var val int32
		if err := binary.Read(file, binary.LittleEndian, &val); err != nil {
			return fmt.Errorf("error reading input weights: %v", err)
		}
		HiddenBias[i] = val
	}

	return nil
}
