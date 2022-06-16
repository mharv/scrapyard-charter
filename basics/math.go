package basics

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
