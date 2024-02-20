package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EntityType int32

const (
	ENTITY_TYPE_UNDEFINED EntityType = iota
	ENTITY_TYPE_PLAYER
)

type Entity struct {
	Type               EntityType
	Position, Velocity rl.Vector2
	Rotation           float32
}

func TranslateRect(r rl.Rectangle, offset rl.Vector2) rl.Rectangle {
	return rl.Rectangle{
		X:      r.X + offset.X,
		Y:      r.Y + offset.Y,
		Width:  r.Width,
		Height: r.Height,
	}
}

func AddVector2(v1, v2 rl.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

func ScaleVector2(v rl.Vector2, scale float32) rl.Vector2 {
	return rl.Vector2{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

func NormalizeVector2(v rl.Vector2) rl.Vector2 {
	if v.X == 0 && v.Y == 0 {
		return v
	}
	magnitude := math.Sqrt(float64(v.X*v.X) + float64(v.Y*v.Y))

	return rl.Vector2{
		X: v.X / float32(magnitude),
		Y: v.Y / float32(magnitude),
	}
}

func main() {
	rl.InitWindow(1280, 720, "Gotroid")

	var entities []Entity
	entities = append(entities, Entity{
		Type:     ENTITY_TYPE_PLAYER,
		Position: rl.Vector2{X: 640, Y: 360},
	})
	// playerIndex := len(entities)-1

	entityRect := rl.Rectangle{
		X: 0, Y: 0,
		Width: 25, Height: 50,
	}
	var playerAcceleration float32 = 4
	var friction float32 = 3

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)

		var dPad rl.Vector2

		if rl.IsKeyDown(rl.KeyA) {
			dPad.X--
		}
		if rl.IsKeyDown(rl.KeyD) {
			dPad.X++
		}
		if rl.IsKeyDown(rl.KeyW) {
			dPad.Y--
		}
		if rl.IsKeyDown(rl.KeyS) {
			dPad.Y++
		}

		for i := range entities {
			entity := &entities[i]
			if entity.Type == ENTITY_TYPE_PLAYER {
				movement := NormalizeVector2(dPad)
				entity.Velocity = AddVector2(
					entity.Velocity,
					ScaleVector2(movement, playerAcceleration),
				)
			}

			entity.Position = AddVector2(
				entity.Position,
				ScaleVector2(entity.Velocity, dt),
			)

			entity.Velocity = ScaleVector2(entity.Velocity, 1.0-(friction*dt))

			rl.DrawRectanglePro(
				TranslateRect(entityRect, entity.Position),
				rl.Vector2{X: -entityRect.Width / 2, Y: -entityRect.Height / 2},
				entity.Rotation, rl.Green,
			)
		}
		rl.EndDrawing()
	}

}
