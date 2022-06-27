package basics

import "math"

//Floats
type Vector2f struct {
	X float64
	Y float64
}

type Vector2[T interface{}] struct {
	X T
	Y T
}

type Rect[T interface{}] struct {
	X      T
	Y      T
	Width  T
	Height T
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

// UI elements

type FloatRectUI struct {
	Name   string
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func (f *FloatRectUI) IsClicked(positionOfClick Vector2f) bool {
	if positionOfClick.X > f.X &&
		positionOfClick.X < (f.X+f.Width) &&
		positionOfClick.Y > f.Y &&
		positionOfClick.Y < (f.Y+f.Height) {
		return true
	} else {
		return false
	}
}

func (f *FloatRectUI) IsHoveredOver(positionOfMouse Vector2f) bool {
	if positionOfMouse.X > f.X &&
		positionOfMouse.X < (f.X+f.Width) &&
		positionOfMouse.Y > f.Y &&
		positionOfMouse.Y < (f.Y+f.Height) {
		return true
	} else {
		return false
	}
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

func FloatMagnitude(v1 Vector2f) float64 {
	return math.Sqrt(v1.X*v1.X + v1.Y*v1.Y)
}

func FloatNormalise(v Vector2f) Vector2f {
	m := FloatMagnitude(v)
	value := Vector2f{}
	if m > 0 {
		value.X = v.X / m
		value.Y = v.Y / m
	}
	return value
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

func AngleFromFVecToFVec(v1, v2 Vector2f) float64 {
	angleVec := Vector2f{X: v2.X - v1.X, Y: v2.Y - v1.Y}
	ang := math.Atan2(-angleVec.Y, -angleVec.X)

	return ang
}
