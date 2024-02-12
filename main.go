package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	rl.InitWindow(1280, 720, "Gotroid")

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Gotroid!", int32(rl.GetScreenWidth())/2, int32(rl.GetScreenHeight())/2, 32, rl.Black)

		rl.EndDrawing()
	}

}
