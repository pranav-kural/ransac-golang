package code

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
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
func (pointCloud *PointCloud) GetRandomPoint() Point3D {
	// return the point at the random index
	return pointCloud.points[rand.Intn(len(pointCloud.points))]
}

// get three random points from PointCloud
func (pointCloud *PointCloud) GetRandomPoints() (Point3D, Point3D, Point3D) {
	// get three random points
	p1 := pointCloud.GetRandomPoint()
	p2 := pointCloud.GetRandomPoint()
	p3 := pointCloud.GetRandomPoint()

	// return the three points
	return p1, p2, p3
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