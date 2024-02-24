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
	Type                             EntityType
	Position, Velocity               rl.Vector2
	Rotation                         float32
	Collider                         rl.Rectangle
	Color                            rl.Color
	PhysicsEnabled, CollisionEnabled bool
	Grounded                         bool
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

func SubtractVector2(v1, v2 rl.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
	}
}

func ScaleVector2(v rl.Vector2, scale float32) rl.Vector2 {
	return rl.Vector2{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

func MagnitudeVector2(v rl.Vector2) float32 {
	var x float64 = float64(v.X)
	var y float64 = float64(v.Y)

	return float32(math.Hypot(x, y))
}

func NormalizeVector2(v rl.Vector2) rl.Vector2 {
	if v.X == 0 && v.Y == 0 {
		return v
	}
	magnitude := MagnitudeVector2(v)

	return rl.Vector2{
		X: v.X / (magnitude),
		Y: v.Y / (magnitude),
	}
}

func AbsVector2(v rl.Vector2) rl.Vector2 {
	if v.X < 0 {
		v.X *= -1
	}
	if v.Y < 0 {
		v.Y *= -1
	}

	return v
}

func main() {
	rl.InitWindow(1280, 720, "Gotroid")

	var entities []Entity
	entities = append(entities, Entity{
		Type:     ENTITY_TYPE_PLAYER,
		Position: rl.Vector2{X: 640, Y: 360},
		Collider: rl.Rectangle{
			X: 0, Y: 0,
			Width: 25, Height: 50,
		},
		PhysicsEnabled: true, CollisionEnabled: true,
		Color: rl.Green,
	})

	tilePosition := rl.Vector2{X: 640 - 200, Y: 360 + 50}
	for i := 0; i < 8; i++ {
		entities = append(entities, Entity{
			Position: tilePosition,
			Collider: rl.Rectangle{
				X: 0, Y: 0,
				Width: 50, Height: 50,
			},
			PhysicsEnabled: false, CollisionEnabled: true,
			Color: rl.Yellow,
		})

		entities = append(entities, Entity{
			Position: tilePosition,
			Collider: rl.Rectangle{
				X: 0, Y: -150,
				Width: 50, Height: 50,
			},
			PhysicsEnabled: false, CollisionEnabled: true,
			Color: rl.Yellow,
		})

		if i == 0 || i == 7 {
			for e := 0; e < 2; e++ {
				entities = append(entities, Entity{
					Position: tilePosition,
					Collider: rl.Rectangle{
						X: 0, Y: -50 - 50*float32(e),
						Width: 50, Height: 50,
					},
					PhysicsEnabled: false, CollisionEnabled: true,
					Color: rl.Yellow,
				})
			}
		}
		tilePosition.X += 50
	}

	// playerIndex := len(entities)-1

	const playerAcceleration float32 = 1000
	const jumpAcceleration float32 = 700
	const airResistance float32 = 1
	const friction float32 = 3
	const gravity float32 = 2.5

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)

		var dPad rl.Vector2

		for i := range entities {
			entity := &entities[i]
			if entity.Type == ENTITY_TYPE_PLAYER {
				accel := playerAcceleration * dt
				if !entity.Grounded {
					accel /= 2
				}

				movement := NormalizeVector2(dPad)
				if rl.IsKeyDown(rl.KeyA) {
					movement.X -= accel
				}
				if rl.IsKeyDown(rl.KeyD) {
					movement.X += accel
				}

				if rl.IsKeyPressed(rl.KeyW) && entity.Grounded {
					movement.Y -= jumpAcceleration
				}

				entity.Velocity = AddVector2(
					entity.Velocity,
					movement,
				)
			}

			if entity.PhysicsEnabled {
				entity.Velocity = AddVector2(entity.Velocity, rl.Vector2{X: 0, Y: gravity})
				entity.Grounded = false

				velocityRemaining := ScaleVector2(entity.Velocity, dt)
				newPosition := AddVector2(entity.Position, velocityRemaining)

				var overlap rl.Vector2
				if entity.CollisionEnabled {
					for e := 0; e < len(entities); e++ {
						if e == i {
							continue
						}

						if rl.CheckCollisionRecs(
							TranslateRect(entity.Collider, newPosition),
							TranslateRect(entities[e].Collider, entities[e].Position),
						) {
							r1 := TranslateRect(entity.Collider, newPosition)
							r2 := TranslateRect(entities[e].Collider, entities[e].Position)
							collision := rl.GetCollisionRec(r1, r2)

							c1 := rl.Vector2{X: r1.X + r1.Width/2, Y: r1.Y + r1.Height/2}
							c2 := rl.Vector2{X: r2.X + r2.Width/2, Y: r2.Y + r2.Height/2}
							delta := SubtractVector2(c2, c1)

							collisionDirection := NormalizeVector2(delta)
							if collisionDirection.Y > 0 {
								entity.Grounded = true
							}

							if collision.Width < collision.Height {
								overlap.X -= collisionDirection.X
								newPosition.X -= collision.Width * collisionDirection.X
							} else {
								overlap.Y -= collisionDirection.Y
								newPosition.Y -= collision.Height * collisionDirection.Y
							}
						}
					}
				}

				entity.Position = newPosition
				if overlap.X != 0 || overlap.Y != 0 {
					v := AbsVector2(entity.Velocity)
					entity.Velocity = AddVector2(entity.Velocity, rl.Vector2{X: overlap.X * v.X / 2, Y: overlap.Y * v.Y / 2})
					entity.Velocity = ScaleVector2(entity.Velocity, 1.0-(friction*dt))
				} else {
					entity.Velocity = ScaleVector2(entity.Velocity, 1.0-(airResistance*dt))
				}
			}

			rl.DrawRectanglePro(
				TranslateRect(entity.Collider, entity.Position),
				rl.Vector2{},
				entity.Rotation, entity.Color,
			)
		}

		rl.EndDrawing()
	}

}
