// Copyright (C) 2025 Tecu23
// Port of Koivisto evaluation, licensed under GNU GPL v3

//go:build embed
// +build embed

// Package nnue keeps the NNUE (Efficiently Updated Neural Network) responsible for
// evaluation the current position
package nnue

import (
	"embed"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

// Global cache for weight loading
var (
	initOnce          sync.Once
	initializationErr error
)

//go:embed default.net
var embeddedWeights embed.FS

// InitializeNNUE sets up the neural network weights using the embedded file
// This function is safe to call multiple times, weights will only be loaded once
func InitializeNNUE() error {
	initOnce.Do(func() {
		initializationErr = loadDefaultWeights()
	})
	return initializationErr
}

// loadDefaultWeights loads weights from the embedded file when using build tag 'embed'
func loadDefaultWeights() error {
	const filename = "default.net"

	// Open the embedded file
	weightFile, err := embeddedWeights.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open embedded weights file: %v", err)
	}
	defer weightFile.Close()

	// Load Weights from the file
	if err := loadWeights(weightFile); err != nil {
		return fmt.Errorf("error loading embedded weights: %v", err)
	}

	fmt.Println("Loaded embedded NNUE weights")
	return nil
}

// LoadWeights loads network weights from a binary file.
// The file is expected to contain input weights, input bias, hidden weights, and hidden bias in sequence.
func loadWeights(r io.Reader) error {
	// Read input weights
	for i := 0; i < InputSize; i++ {
		for j := 0; j < HiddenSize; j++ {

			var val int16
			if err := binary.Read(r, binary.LittleEndian, &val); err != nil {
				return fmt.Errorf("error reading input weights: %v", err)
			}
			InputWeights[i][j] = val
		}
	}

	// Read input bias
	for i := 0; i < HiddenSize; i++ {
		var val int16
		if err := binary.Read(r, binary.LittleEndian, &val); err != nil {
			return fmt.Errorf("error reading input weights: %v", err)
		}
		InputBias[i] = val
	}

	// Read hidden weights
	for i := 0; i < OutputSize; i++ {
		for j := 0; j < HiddenDSize; j++ {
			var val int16
			if err := binary.Read(r, binary.LittleEndian, &val); err != nil {
				return fmt.Errorf("error reading input weights: %v", err)
			}
			HiddenWeights[i][j] = val
		}
	}

	// Read hidden bias
	for i := 0; i < OutputSize; i++ {
		var val int32
		if err := binary.Read(r, binary.LittleEndian, &val); err != nil {
			return fmt.Errorf("error reading input weights: %v", err)
		}
		HiddenBias[i] = val
	}

	return nil
}
