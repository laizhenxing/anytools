package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	peakFalloff         = 8.0
	defaultWindowWidth  = 800
	defaultWindowHeight = 450
	spectrumSize        = 80
)

var (
	fftOutputLock sync.RWMutex
	fftOutput     []complex128
	isDropped     = false
	done          = make(chan bool)
	isPlayer      = false
	playingText string
	f *os.File
	streamer beep.StreamSeekCloser 
)

func main() {
	freqSpectrum := make([]float64, spectrumSize)
	var (
		windowWidth  int32 = defaultWindowWidth
		windowHeight int32 = defaultWindowHeight
	)
	// 窗口调整大小
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(windowWidth, windowHeight, "xmp3")
	defer rl.CloseWindow()
	defer closeHandler()

	rl.SetTargetFPS(60)

	var droppedFiles []string

	for !rl.WindowShouldClose() {
		windowWidth = int32(rl.GetScreenWidth())
		windowHeight = int32(rl.GetScreenHeight())
		columnWidth := int32(windowWidth / spectrumSize)
		if rl.IsFileDropped() {
			droppedFiles = rl.LoadDroppedFiles()
			rl.UnloadDroppedFiles()
			err := handleFileDrop(droppedFiles[0])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fs, _ := f.Stat()
			playingText = "Now playing " + fs.Name() 
			fmt.Println(droppedFiles[0])
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		if !isDropped {
			drawDropzone(windowWidth, windowHeight)
		} else {
			select {
			case <-done:
				closeHandler()
				isDropped = false
			default:
			}

			if isPlayer {
				fftOutputLock.RLock()
				updateSpectrumValues(float64(windowHeight), freqSpectrum)
				fftOutputLock.RUnlock()
				for i, s := range freqSpectrum {
					rl.DrawRectangleGradientV(int32(i)*columnWidth, windowHeight-int32(s), columnWidth, int32(s), rl.Orange, rl.Green)
					rl.DrawRectangleLines(int32(i)*columnWidth, windowHeight-int32(s), columnWidth, int32(s), rl.Black)
				}
				rl.DrawText(playingText, 40, 40, 20, rl.White)
			}
		}

		rl.EndDrawing()
	}
}

func handleFileDrop(path string) (err error) {
	if !strings.HasSuffix(path, ".mp3") {
		return errors.New("is not mp3 file.")
	}
	f, err = os.Open(path)
	if err != nil {
		return err
	}

	var format beep.Format
	streamer, format, err = mp3.Decode(f)
	if err != nil {
		return err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(Streamer{Streamer: streamer}, beep.Callback(func() {
		done <- true
	})))
	isDropped = true
	return nil
}

func drawDropzone(windowWidth, windowHeight int32) {
	var fontSize float32 = 16.0
	font := rl.GetFontDefault()
	message := "Drop mp3 file to this window"
	textPos := rl.Vector2{
		X: float32(windowWidth)/2.0 - rl.MeasureTextEx(font, message, fontSize, 2).X/2.0,
		Y: float32(windowHeight)/2.0 - fontSize/2.0,
	}
	rl.DrawTextEx(font, message, textPos, fontSize, 2, rl.Yellow)
	rl.DrawRectangleLines(20, 20, windowWidth-40, windowHeight-40, rl.LightGray)
}

func updateSpectrumValues(maxValue float64, freqSpectrum []float64) {
	for i := 0; i < spectrumSize; i++ {
		fr := real(fftOutput[i])
		fi := imag(fftOutput[i])
		magnitude := math.Sqrt(fr*fr + fi*fi)
		val := math.Min(maxValue, math.Abs(magnitude))
		if freqSpectrum[i] > val {
			freqSpectrum[i] = math.Max(freqSpectrum[i]-peakFalloff, 0.0)
		} else {
			freqSpectrum[i] = (val + freqSpectrum[i]) / 2.0
		}
	}

}

func closeHandler() {
	if f != nil {
		f.Close()
	}
}