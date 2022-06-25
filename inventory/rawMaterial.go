package inventory

import (
	"math/rand"
	"time"
)

type RawMaterial struct {
	min, max, amount int
}

func (r *RawMaterial) RandomiseAmount() {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	r.amount = rnd.Intn(r.max-r.min) + r.min
}

func (r *RawMaterial) SetMinAndMax(min, max int) {
	r.min = min
	r.max = max
	r.RandomiseAmount()
}

func (r *RawMaterial) GetAmount() int {
	return r.amount
}
