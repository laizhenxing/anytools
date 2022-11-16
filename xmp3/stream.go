package main

import (
	"github.com/faiface/beep"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
)

type Streamer struct {
	Streamer beep.Streamer
}

func (s Streamer) Stream(samples [][2]float64) (n int, ok bool) {
	var ware = make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		ware[i] = samples[i][0] * 35
	}
	window.Apply(ware, window.Bartlett)
	fftOutputLock.Lock()
	fftOutput = fft.FFTReal(ware)
	fftOutputLock.Unlock()
	if !isPlayer {
		isPlayer = true
	}
	return s.Streamer.Stream(samples)
}

func (s Streamer) Err() error {
	return s.Streamer.Err()
}
