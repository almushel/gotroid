package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type EntityType int32

const (
	ENTITY_TYPE_UNDEFINED EntityType = iota
	ENTITY_TYPE_PLAYER
	ENTITY_TYPE_WALL
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

func NewEntity(newType EntityType) Entity {
	result := Entity{
		Type: newType,
	}

	switch newType {

	case ENTITY_TYPE_PLAYER:
		result.Collider = rl.Rectangle{
			X: 0, Y: 0,
			Width: 25, Height: 50,
		}
		result.PhysicsEnabled = true
		result.CollisionEnabled = true
		result.Color = rl.Green

	case ENTITY_TYPE_WALL:
		result.Collider = rl.Rectangle{
			X: 0, Y: 0,
			Width: 50, Height: 50,
		}
		result.CollisionEnabled = true
		result.Color = rl.Yellow

	default:
		result.Type = ENTITY_TYPE_UNDEFINED
		return result
	}

	return result
}

const playerAcceleration float32 = 1000
const jumpAcceleration float32 = 700
const airResistance float32 = 1
const friction float32 = 3
const gravity float32 = 2.5

func UpdateEntities(entities []Entity, dt float32) {
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

						if collision.Width < collision.Height {
							overlap.X -= collisionDirection.X
							newPosition.X -= collision.Width * collisionDirection.X
						} else {
							if collisionDirection.Y > 0 {
								entity.Grounded = true
							}
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
	}
}

func DrawEntities(entities []Entity) {
	for _, entity := range entities {
		rl.DrawRectanglePro(
			TranslateRect(entity.Collider, entity.Position),
			rl.Vector2{},
			entity.Rotation, entity.Color,
		)
	}
}
