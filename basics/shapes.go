package basics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawCircle(dst *ebiten.Image, px, py, ox, oy, r float64, clr color.Color) {
	x1 := int(px + ox)
	y1 := int(py + oy)
	x, y, dx, dy := int(r-1), 0, 1, 1
	err := dx - int(r*2)

	for x > y {
		dst.Set(x1+x, y1+y, clr)
		dst.Set(x1+y, y1+x, clr)
		dst.Set(x1-y, y1+x, clr)
		dst.Set(x1-x, y1+y, clr)
		dst.Set(x1-x, y1-y, clr)
		dst.Set(x1-y, y1-x, clr)
		dst.Set(x1+y, y1-x, clr)
		dst.Set(x1+x, y1-y, clr)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - int(r*2)
		}
	}
}
