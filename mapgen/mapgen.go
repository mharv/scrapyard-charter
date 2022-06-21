package mapgen

import (
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/mharv/scrapyard-charter/globals"
)

func GenerateMap() [globals.ScreenWidth][globals.ScreenHeight]float64 {
	const w = globals.ScreenWidth
	const h = globals.ScreenHeight
	// generate fall off map

	// var terrain [w][h]float64
	var fallOffMap [w][h]float64
	var terrain [w][h]float64

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			x := float64(i)/float64(w)*2 - 1
			y := float64(j)/float64(h)*2 - 1

			v := math.Max(math.Abs(x), math.Abs(y))
			fallOffMap[i][j] = math.Pow(v, 3) / (math.Pow(v, 3) + math.Pow(3-3*v, 3))
		}
	}

	// generate fall off map with sides open
	l_fallOffMap := false
	r_fallOffMap := false
	u_fallOffMap := false
	d_fallOffMap := true

	if l_fallOffMap {
		l_offset := w/2 - 1
		for i := 0; i < w/4; i++ {
			for j := 0; j < h; j++ {
				fallOffMap[i][j] = fallOffMap[i+l_offset][j]
			}
		}
	}
	if r_fallOffMap {
		r_offset := w/2 - 1
		for i := w/2 + w/4; i < w; i++ {
			for j := 0; j < h; j++ {

				fallOffMap[i][j] = fallOffMap[i-r_offset][j]
			}
		}
	}
	if u_fallOffMap {
		u_offset := h/2 - 1
		for i := 0; i < w; i++ {
			for j := 0; j < h/4; j++ {

				fallOffMap[i][j] = fallOffMap[i][j+u_offset]
			}
		}
	}
	if d_fallOffMap {
		d_offset := h/2 - 1
		for i := 0; i < w; i++ {
			for j := h/2 + h/4; j < h; j++ {

				fallOffMap[i][j] = fallOffMap[i][j-d_offset]
			}
		}
	}

	// setup perlin noise gen -- probably wrong useage
	var iterations int32 = 2
	perlinNoise := perlin.NewPerlin(2, 3, iterations, int64(rand.Int()))
	scale := 0.2

	// apply fall off map to perlin noise
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {

			xOffset := (rand.Float64()*2 - 1) * 5000
			yOffset := (rand.Float64()*2 - 1) * 5000
			randomChanceToAdd := perlinNoise.Noise2D(float64(x)*scale+xOffset, float64(y)*scale+yOffset)
			// try pure random - nah, looks too boring
			// randomChanceToAdd = rand.Float64()

			// forms hard border around edge of screen
			// if x == globals.ScreenWidth-offsetForGrid || x == 0 {
			// 	randomChanceToAdd = 1
			// }
			// if y == globals.ScreenHeight || y == 0 {
			// 	randomChanceToAdd = 1
			// }

			// use fall off map to reduce the chance of scrap spawning in the middle
			// creating an island like terrain
			randomChanceToAdd += fallOffMap[x][y]

			terrain[x][y] = randomChanceToAdd

		}
	}

	// smooth out values using filter to reduce noise

	for i := 0; i < 10; i++ {
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				// if not the edge, apply 3x3 filter
				if x != globals.ScreenWidth-1 &&
					x != 0 &&
					y != globals.ScreenHeight-1 &&
					y != 0 {
					median := (sum(terrain[x-1][y-1:y+2]) + sum(terrain[x][y-1:y+2]) + sum(terrain[x+1][y-1:y+2])) / 9
					terrain[x][y] = median
				}
			}
		}
	}

	return terrain
}

func sum(numbers []float64) float64 {
	total := 0.0
	for i := range numbers {
		total += numbers[i]
	}
	return total
}
