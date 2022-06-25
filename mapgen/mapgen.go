package mapgen

import (
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/mharv/scrapyard-charter/globals"
)

const w = globals.ScreenWidth
const h = globals.ScreenHeight

func GenerateMap(l_open, r_open, u_open, d_open bool) [w][h]float64 {

	// generate fall off map and return terrain map
	var fallOffMap [w][h]float64
	var terrain [w][h]float64

	fallOffMap = createSquareFallOffMap(fallOffMap)

	// generate fall off map with sides open
	fallOffMap = openFallOffMapSide(fallOffMap, l_open, r_open, u_open, d_open)

	terrain = applyPerlinNoise(terrain, fallOffMap)
	// smooth out values using filter to reduce noise
	terrain = applyMedianFilterNTime(terrain, 10)

	return terrain
}

func sum(numbers []float64) float64 {
	total := 0.0
	for i := range numbers {
		total += numbers[i]
	}
	return total
}

func applyPerlinNoise(terrain, fallOffMap [w][h]float64) [w][h]float64 {
	// setup perlin noise gen -- probably wrong useage
	var iterations int32 = 2
	perlinNoise := perlin.NewPerlin(2, 3, iterations, int64(globals.GetPlayerData().GetWorldSeed()))
	scale := 0.2

	rand.Seed(int64(globals.GetPlayerData().GetWorldSeed()))

	// apply fall off map to perlin noise
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			// try pure random - nah, looks too boring
			// randomChanceToAdd = rand.Float64()

			xOffset := (rand.Float64()*2 - 1) * 5000
			yOffset := (rand.Float64()*2 - 1) * 5000
			randomChanceToAdd := perlinNoise.Noise2D(float64(x)*scale+xOffset, float64(y)*scale+yOffset)

			// use fall off map to reduce the chance of scrap spawning in the middle
			// creating an island like terrain
			randomChanceToAdd += fallOffMap[x][y]
			terrain[x][y] = randomChanceToAdd
		}
	}
	return terrain
}

func createSquareFallOffMap(fallOffMap [w][h]float64) [w][h]float64 {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			x := float64(i)/float64(w)*2 - 1
			y := float64(j)/float64(h)*2 - 1

			v := math.Max(math.Abs(x), math.Abs(y))
			fallOffMap[i][j] = math.Pow(v, 3) / (math.Pow(v, 3) + math.Pow(3-3*v, 3))
		}
	}
	return fallOffMap
}

func applyMedianFilterNTime(terrain [w][h]float64, n int) [w][h]float64 {
	// 3x3 variant
	for i := 0; i < n; i++ {
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				// if not the edge, apply 3x3 filter
				if x != w-1 &&
					x != 0 &&
					y != h-1 &&
					y != 0 {
					median := (sum(terrain[x-1][y-1:y+2]) + sum(terrain[x][y-1:y+2]) + sum(terrain[x+1][y-1:y+2])) / 9
					terrain[x][y] = median
				}
			}
		}
	}
	return terrain
}

func openFallOffMapSide(fallOffMap [w][h]float64, l_fallOffMap, r_fallOffMap, u_fallOffMap, d_fallOffMap bool) [w][h]float64 {

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
	return fallOffMap
}
