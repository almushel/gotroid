package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
