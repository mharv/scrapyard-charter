package basics

import "math"

//Floats
type Vector2f struct {
	X float64
	Y float64
}

type Vector3f struct {
	X float64
	Y float64
	Z float64
}

type FloatRect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

//Ints
type Vector2i struct {
	X int
	Y int
}

type Vector3i struct {
	X int
	Y int
	Z int
}

type IntRect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func FloatLerp(a, b, t float64) float64 {
	return (1-t)*a + t*b
}

func Vec2FLerp(v1, v2 Vector2f, t float64) Vector2f {
	return Vector2f{X: FloatLerp(v1.X, v2.X, t), Y: FloatLerp(v1.Y, v2.Y, t)}
}

func FloatDistance(v1, v2 Vector2f) float64 {
	return math.Sqrt(((v2.X - v1.X) * (v2.X - v1.X)) + ((v2.Y - v1.Y) * (v2.Y - v1.Y)))
}

func FloatNormalise(v Vector2f) Vector2f {
	return Vector2f{}
}

func FloatClamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func FloatRotAroundPoint(point, center Vector2f, angle float64) Vector2f {
	vec1 := Vector2f{X: point.X - center.X, Y: point.Y - center.Y}
	vec2 := Vector2f{
		X: vec1.X*math.Cos(angle) - vec1.Y*math.Sin(angle),
		Y: vec1.X*math.Sin(angle) + vec1.Y*math.Cos(angle),
	}

	// translate point back:
	value := Vector2f{X: vec2.X + center.X, Y: vec2.Y + center.Y}
	return value
}
