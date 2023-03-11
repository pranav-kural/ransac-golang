package ransac

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// structure of a point cloud
type PointCloud struct {
	// store the points
	points []Point3D
}

// default separator used to separate the coordinates of a point in point cloud data file
var POINTS_SEPARATOR = "\\s+"
// store the default points coordinate labels
var POINTS_COORDINATES_LABELS = "x y z"

// method to read the points from a file, and return an array containing Point3D instances
func readXYZ(filename string, args ...string) (pointsCloud PointCloud, err error) {
		// validate filename
		if (filename == "") {
				return pointsCloud, errors.New("No filename provided")
		}

		// if separator provided, use it
		if len(args) > 0 {
			POINTS_SEPARATOR = args[0]
		}

		// open the file
		file, err := os.Open(filename)
		if (err != nil) {
				return pointsCloud, errors.New("Could not open file")
		}

		// if open successful, defer closing the file
		defer file.Close()

		// create a scanner to read the file
		scanner := bufio.NewScanner(file)
		// temporary variable to store the point
		point := Point3D{}
		// store points in an array
		points := []Point3D{}

		// read the file line by line and for each line read, extract the Point3D object and store it in the points array
		for scanner.Scan() {
			// attempt to extract a Point3D from the line information
			point, err = getPoint3D(scanner.Text())
			// if error, return the error and stop reading the file
			if err != nil {
				return pointsCloud, err
			}
			// append the point to the points array
			points = append(points, point)
		}

		// return the pointsCloud
		return PointCloud{ points }, nil
}

// Given a string containing information of three point coordinates, returns a Point3D
func getPoint3D(pointsData string) (Point3D, error) {
	// check if pointsData is valid
	if pointsData == "" {
		return Point3D{}, errors.New("pointsData invalid, is empty")
	}

	// get data for each point
	pointData := strings.Split(pointsData, POINTS_SEPARATOR)

	// check we do have a 3d point
	if len(pointData) != 3 {
		return Point3D{}, errors.New("Invalid number of points provided in pointsData")
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

// save a file with provided filename and points data
func saveXYZ(filename string, points []Point3D, args ...string) error {
	// validate filename
	if (filename == "") {
			return errors.New("No filename provided")
	}

	// if custom coordinates labels provided, use them
	if len(args) > 0 {
		POINTS_COORDINATES_LABELS = args[0]
	}

	// create & open the file if doesn't exist
	file, err := os.Create(filename)
	if (err != nil) {
			return errors.New("Could not open file")
	}

	// if open successful, defer closing the file
	defer file.Close()

	// create a writer to write to the file
	writer := bufio.NewWriter(file)

	// write the header
	_, err = writer.WriteString(POINTS_COORDINATES_LABELS + "\n")
	if err != nil {
		return err
	}

	// write the points to the file
	for _, point := range points {
		// write the point to the file
		_, err := writer.WriteString(point.String() + "\n")
		if err != nil {
			return err
		}
	}

	return nil
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