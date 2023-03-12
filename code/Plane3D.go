package code

import (
	"fmt"
	"math"
)

// Plane3D represents a 3D plane
type Plane3D struct {
	A float64
	B float64
	C float64
	D float64
}

// Plane3D with supporting points
type Plane3DwSupport struct {
	Plane3D
 	SupportSize int
	SupportingPoints []Point3D
}

// computes the plane defined by a set of 3 points
func GetPlane(p1, p2, p3 Point3D) Plane3D {
	// compute the normal of the plane
	normal := GetNormal(p1, p2, p3)

	// compute the distance of the plane from the origin
	distance := -(normal.X*p1.X + normal.Y*p1.Y + normal.Z*p1.Z)

	// return the plane
	return Plane3D{normal.X, normal.Y, normal.Z, distance}
}

// received array containing 3 Point3D objects and sends back a Plane3D object through output channel
func GetPlaneC(pointsIn <-chan [3]Point3D, done <-chan bool) <-chan Plane3D {
	fmt.Println("********** GetPlaneC started **********")
	// outbound channel
	planeOut := make(chan Plane3D)
	// goroutine to compute the plane
	go func() {
		defer close(planeOut)
		defer fmt.Println("********** GetPlaneC done **********")

		for points := range pointsIn {
			// get array of points from inbound channel
			// compute the plane
			// send plane on outbound channel
			planeOut <- GetPlane(points[0], points[1], points[2])
		}

		// // until we receive a message on the done channel
		// // send a Plane3D object on the outbound channel
		// for {
		// 	select {
		// 	case <-done:
		// 		return
		// 	case points := <-pointsIn:
		// 		// get array of points from inbound channel
		// 		// compute the plane
		// 		// send plane on outbound channel
		// 		planeOut <- GetPlane(points[0], points[1], points[2])
		// 	}
		// }
	}()
	// return the outbound channel
	return planeOut
}

func GetNormal(point3D1, point3D2, point3D3 Point3D) Point3D {
	// compute the vectors v1 and v2
	v1 := Point3D{point3D2.X - point3D1.X, point3D2.Y - point3D1.Y, point3D2.Z - point3D1.Z}
	v2 := Point3D{point3D3.X - point3D1.X, point3D3.Y - point3D1.Y, point3D3.Z - point3D1.Z}

	// compute the normal vector
	normal := Point3D{v1.Y*v2.Z - v1.Z*v2.Y, v1.Z*v2.X - v1.X*v2.Z, v1.X*v2.Y - v1.Y*v2.X}

	// return the normal vector
	return normal
}

// calculate distance of a point to a plane
func (p *Plane3D) GetDistance(point *Point3D) float64 {
	return math.Abs(p.A*point.X + p.B*point.Y + p.C*point.Z + p.D) / math.Sqrt(p.A*p.A+p.B*p.B+p.C*p.C)
}

// string representation of a Plane3D
func (p Plane3D) String() string {
	return fmt.Sprintf("a=%f, b=%fy, c=%fz, d=%f", p.A, p.B, p.C, p.D)
}

// method to return an array of points that support the plane
func (p *Plane3D) GetSupportingPointss(points []Point3D, eps float64) *[]Point3D {
	// create an array of points that support the plane
	supportingPoints := make([]Point3D, 0)

	// iterate over all points
	for _, point := range points {
		// if the point is on the plane, add it to the array
		if p.GetDistance(&point) <= eps {
			supportingPoints = append(supportingPoints, point)
		}
	}

	// return the array of supporting points
	return &supportingPoints
}