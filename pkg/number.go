package pkg

import (
	"math/rand"
	"time"

	"golang.org/x/exp/constraints"
)

const Hundred = 100.0

type Number interface {
	constraints.Integer | constraints.Float
}

// RandomNumber takes minimum and maximum range and return a random value
func RandomNumber(min, max int) int {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	return random.Intn(max-min) + min
}

func Abs[T Number](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

func CentsToDollars(cents int) float64 {
	return float64(cents) / Hundred
}
