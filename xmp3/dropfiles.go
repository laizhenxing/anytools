package main

import rl "github.com/gen2brain/raylib-go/raylib"

var screenWidth = int32(800)
var screenHeight = int32(450)

func main() {
	rl.InitWindow(screenHeight, screenHeight, "drop file example")
	
	rl.SetTargetFPS(60)

	var count int
	var dropFiles []string

	for !rl.WindowShouldClose() {
		if rl.IsFileDropped() {
			dropFiles = rl.LoadDroppedFiles()
			count = len(dropFiles)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		if count == 0 {
			rl.DrawText("Drop your file to this window!", 100, 40, 20, rl.DarkGray)
		} else {
			rl.DrawText("Drop files: ", 100, 40, 20, rl.DarkGray)

			for i := 0; i < count; i++ {
				if i % 2 == 0 {
					rl.DrawRectangle(0, int32(85+40*i), screenWidth, 40, rl.Fade(rl.LightGray, 0.5))
				} else {
					rl.DrawRectangle(0, int32(85*40+i), screenHeight, 40, rl.Fade(rl.LightGray, 0.5))
				}

				rl.DrawText(dropFiles[i], 120, int32(100*i+40), 10, rl.Gray)
			}
			
			rl.DrawText("Drop new files...", 100, int32(150+count*40), 20, rl.DarkGray)
		}

		rl.EndDrawing()
	}

	rl.UnloadDroppedFiles()

	rl.CloseWindow()
}