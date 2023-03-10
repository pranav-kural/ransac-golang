package ransac

import "fmt"

type Point3D struct {
	X float64
	Y float64
	Z float64
}

func (p Point3D) String() string {
	 return fmt.Sprintf("%f %f %f", p.X, p.Y, p.Z)
}

