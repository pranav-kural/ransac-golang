package code

import (
	"math/rand"
)

// structure of a point cloud
type PointCloud struct {
	// store the points
	points []Point3D
}

// get a random point from PointCloud
func (pointCloud *PointCloud) RandomPointGenerator(done <-chan bool) <-chan Point3D {
	dprint("********** RandomPointGenerator started **********")
	// outbound channel
	pointOut := make(chan Point3D)
	// goroutine to generate random points
	go func() {
		defer close(pointOut)
		defer dprint("********** RandomPointGenerator done **********")
		for {
			// until we receive a message on the done channel
			// send a random point on the outbound channel
			select {
			case <-done:
				return
			default:
				// check if there are points in the point cloud
				if len(pointCloud.points) == 0 {
					return
				}
				pointOut <- pointCloud.points[rand.Intn(len(pointCloud.points)+1) % (len(pointCloud.points)-1)]
			}
		}
	}()
	// return the outbound channel
	return pointOut
}

// get three random points from PointCloud
func (pointCloud *PointCloud) GetRandomPoints(done <-chan bool) <-chan [3]Point3D {
	dprint("********** GetRandomPoints started **********")
	// outbound channel
	pointsOut := make(chan [3]Point3D)
	// channel to receive random points
	chanRPG := pointCloud.RandomPointGenerator(done)
	// goroutine to generate random points
	go func() {
		defer close(pointsOut)
		defer dprint("********** GetRandomPoints done **********")
		for {
			// until we receive a message on the done channel
			// send an array containing points on the outbound channel
			select {
			case <-done:
				return
			default:
				pointsOut <- [3]Point3D{<-chanRPG, <-chanRPG, <-chanRPG}
			}
		}
	}()
	// return the outbound channel
	return pointsOut
}

// receive arrays containing 3 Point3D through incoming channel and resend the array on the outbound channel until N arrays of Point3D
func (pointCloud *PointCloud) TakeN(n int) <-chan [3]Point3D {
	dprint("********** TakeN started **********")
	// outbound channel (buffered since size already known)
	arrOut := make(chan [3]Point3D, n)
	// done channel
	done := make(chan bool)
	// inbound channel
	arrIn := pointCloud.GetRandomPoints(done)
	// goroutine to send the array N times
	go func() {
		defer close(arrOut)
		defer close(done)
		// for n times
		for i := 0; i < n; i++ {
			// get array from inbound channel
			arr := <-arrIn
			// send array on outbound channel
			arrOut <- arr
		}
		// stop upstream channels
		//done <- true
		dprint("********** TakeN done **********")
	}()
	// return the outbound channel
	return arrOut
}

// method to return an array of points that support the plane
func (pointCloud *PointCloud) GetSupportingPoints(plane Plane3D, eps float64) *[]Point3D {
	// create an array of points that support the plane
	supportingPoints := make([]Point3D, 0)

	// iterate over all points
	for _, point := range pointCloud.points {
		// if the point is on the plane, add it to the array
		if plane.GetDistance(&point) <= eps {
			supportingPoints = append(supportingPoints, point)
		}
	}

	// return the array of supporting points
	return &supportingPoints
}

// method that receives Plane3D instance from inbound channel
// returns Plane3DwSupport instance containing plane and the supporting points
func (pointCloud *PointCloud) GetSupportingPointsC(planeIn <-chan Plane3D, eps float64) <-chan Plane3DwSupport {
	dprint("********** GetSupportingPointsC started **********")
	// outbound channel
	planeOut := make(chan Plane3DwSupport)
	// goroutine to get supporting points
	go func() {
		defer close(planeOut)
		defer dprint("********** GetSupportingPointsC done **********")
		// until we receive a plane from the inbound channel
		for plane := range planeIn {
			// create an array of points that support the plane
			supportingPoints := make([]Point3D, 0)

			// iterate over all points
			for _, point := range pointCloud.points {
				// if the point is on the plane, add it to the array
				if plane.GetDistance(&point) <= eps {
					supportingPoints = append(supportingPoints, point)
				}
			}
			// send the plane with the supporting points on the outbound channel
			planeOut <- Plane3DwSupport{
				Plane3D: plane,
				SupportingPoints: supportingPoints,
				SupportSize: len(supportingPoints),
			}
		}
	}()
	// return the outbound channel
	return planeOut
}



// creates a new slice of points in which all points
// belonging to the plane have been removed
func (pointsCloud *PointCloud) RemovePlane(plane *Plane3D, eps float64) PointCloud {
	// create a new slice of points
	newPoints := []Point3D{}

	// iterate through the points and add them to the new slice if they don't belong to the plane
	for _, point := range pointsCloud.points {
		// if the point is not on the plane, add it to the new slice
		if plane.GetDistance(&point) > eps {
			newPoints = append(newPoints, point)
		}
	}

	// return the new slice
	return PointCloud{ newPoints }
}