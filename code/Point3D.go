package ransac

import (
	"fmt"
	"math"
)

// Point3D represents a 3D point
type Point3D struct {
	X float64
	Y float64
	Z float64
}

// computes the distance between points p1 and p2
func (p1 *Point3D) GetDistance(p2 *Point3D) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
}

// string representation of a Point3D
func (p Point3D) String() string {
	 return fmt.Sprintf("%f %f %f", p.X, p.Y, p.Z)
}