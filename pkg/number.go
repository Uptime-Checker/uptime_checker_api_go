package pkg

import (
	"math/rand"
	"time"
)

// RandomNumber takes minimum and maximum range and return a random value
func RandomNumber(min, max int) int {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	return random.Intn(max-min) + min
}
