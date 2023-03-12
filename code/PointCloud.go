package code

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

// structure of a point cloud
type PointCloud struct {
	// store the points
	points []Point3D
}

// default separator used to separate the coordinates of a point in point cloud data file
var POINTS_SEPARATOR = " "
// store the default points coordinate labels
var pointsCoordinatesLabels = "x y z"

// Given a string containing information of three point coordinates, returns a Point3D
func getPoint3D(pointsData string) (Point3D, error) {
	// check if pointsData is valid
	if pointsData == "" {
		return Point3D{}, errors.New("pointsData invalid, is empty")
	}

	// get data for each point separated by the separator and trim spaces
	pointData := strings.Fields(pointsData)
	
	// check we do have a 3d point
	if len(pointData) != 3 {
		return Point3D{}, errors.New("invalid number of points provided in pointsData: " + pointsData)
	}

	// parse points (might throw exception)
	x, err := strconv.ParseFloat(pointData[0], 64)
	if err != nil {
		return Point3D{}, err
	}
	y, err := strconv.ParseFloat(pointData[1], 64)
	if err != nil {
		return Point3D{}, err
	}
	z, err := strconv.ParseFloat(pointData[2], 64)
	if err != nil {
		return Point3D{}, err
	}

	return Point3D{x, y, z}, nil
}


// get a random point from PointCloud
func (pointCloud *PointCloud) RandomPointGenerator(done <-chan bool) <-chan Point3D {
	fmt.Println("********** RandomPointGenerator started **********")
	// outbound channel
	pointOut := make(chan Point3D)
	// goroutine to generate random points
	go func() {
		defer close(pointOut)
		defer fmt.Println("********** RandomPointGenerator done **********")
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
				pointOut <- pointCloud.points[rand.Intn(len(pointCloud.points)+1)]
			}
		}
	}()
	// return the outbound channel
	return pointOut
}

// get three random points from PointCloud
func (pointCloud *PointCloud) GetRandomPoints(done <-chan bool) <-chan [3]Point3D {
	fmt.Println("********** GetRandomPoints started **********")
	// outbound channel
	pointsOut := make(chan [3]Point3D)
	// channel to receive random points
	chanRPG := pointCloud.RandomPointGenerator(done)
	// goroutine to generate random points
	go func() {
		defer close(pointsOut)
		defer fmt.Println("********** GetRandomPoints done **********")
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
	fmt.Println("********** TakeN started **********")
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
			fmt.Println("********** TakeN done **********")
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
func (pointCloud *PointCloud) GetSupportingPointsC(planeIn *<-chan Plane3D, eps *float64, done <-chan bool) <-chan Plane3DwSupport {
	fmt.Println("********** GetSupportingPointsC started **********")
	// outbound channel
	planeOut := make(chan Plane3DwSupport)
	// goroutine to get supporting points
	go func() {
		defer close(planeOut)
		defer fmt.Println("********** GetSupportingPointsC done **********")
		for {
			// until we receive a message on the done channel
			// send an array containing points on the outbound channel
			select {
			case <-done:
				return
			case plane := <-*planeIn:
				// receive plane from inbound channel
				// create an array of points that support the plane
				supportingPoints := make([]Point3D, 0)
				// synchoronize checking of supporting points
				var wg sync.WaitGroup

				// iterate over all points
				for _, point := range pointCloud.points {
					wg.Add(1)
					// goroutine to check if current point is on the plane
					go func(point *Point3D, supportingPoints *[]Point3D) {
						// if the point is on the plane, add it to the array
						if plane.GetDistance(point) <= *eps {
							*supportingPoints = append(*supportingPoints, *point)
						}
						wg.Done()
					}(&point, &supportingPoints)
				}
				// wait until all goroutines have finished (all points have been checked)
				wg.Wait()
				// send Plane3DwSupport instance on outbound channel
				planeOut <- Plane3DwSupport{
					Plane3D: plane,
					SupportingPoints: supportingPoints,
					SupportSize: len(supportingPoints),
				}
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

	// waitgroup to synchronize the goroutines
	var wg sync.WaitGroup

	// iterate through the points and add them to the new slice if they don't belong to the plane
	for _, point := range pointsCloud.points {
		// add 1 to the waitgroup
		wg.Add(1)
		// goroutine to check if the point is not on the plane
		go func() {
			// if the point is not on the plane, add it to the new slice
			if plane.GetDistance(&point) > eps {
				newPoints = append(newPoints, point)
			}
			wg.Done()
		}()
	}
	// wait until all goroutines have finished
	wg.Wait()
	// return the new slice
	return PointCloud{ newPoints }
}