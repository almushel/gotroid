package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1280, 720, "Gotroid")

	var entities []Entity
	player := NewEntity(ENTITY_TYPE_PLAYER)
	player.Position = rl.Vector2{X: 640, Y: 360}

	entities = append(entities, player)

	tilePosition := rl.Vector2{X: 640 - 200, Y: 360 + 50}
	for i := 0; i < 8; i++ {
		tile := NewEntity(ENTITY_TYPE_WALL)
		tile.Position = tilePosition
		entities = append(entities, tile)

		tile.Position.Y -= 150
		entities = append(entities, tile)

		if i == 0 || i == 7 {
			for e := 0; e < 2; e++ {
				tile.Position.Y = tilePosition.Y - 50 - 50*float32(e)
				entities = append(entities, tile)
			}
		}
		tilePosition.X += 50
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		UpdateEntities(entities, dt)

		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)
		DrawEntities(entities)

		rl.EndDrawing()
	}

}
