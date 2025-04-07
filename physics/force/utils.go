package force

import (
	"math"

	"github.com/alexanderi96/go-space-engine/core/constants"
)

func CalculateOrbitalVelocity(mass float64, distance float64) float64 {
	// Formula: v = sqrt(G * M / r)
	velocity := math.Sqrt(constants.G*mass/distance) / 2
	return velocity
}
