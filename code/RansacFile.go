package code

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

// method to read the points from a file, and return an array containing Point3D instances
func readXYZ(filename string, args ...string) (pointsCloud PointCloud, err error) {
		// validate filename
		if (filename == "") {
				return pointsCloud, errors.New("no filename provided")
		}

		// if separator provided, use it
		if len(args) > 0 {
			POINTS_SEPARATOR = args[0]
		}

		// open the file
		file, err := os.Open(filename)
		if (err != nil) {
				return pointsCloud, errors.New("could not open file")
		}

		// if open successful, defer closing the file
		defer file.Close()

		// create a scanner to read the file
		scanner := bufio.NewScanner(file)
		// store points in an array
		points := []Point3D{}

		// read the first line of the file to get the points coordinates labels
		if scanner.Scan() {
			pointsCoordinatesLabels = scanner.Text()
		}

		// read the file line by line and for each line read, extract the Point3D object and store it in the points array
		for scanner.Scan() {
			// attempt to extract a Point3D from the line information
			point, err := getPoint3D(scanner.Text())
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

// save a file with provided filename and points data
func saveXYZ(filename string, points []Point3D, args ...string) error {
	// validate filename
	if (filename == "") {
			return errors.New("no filename provided")
	}

	fmt.Println("Saving file: " + filename)

	// if custom coordinates labels provided, use them
	if len(args) > 0 {
		pointsCoordinatesLabels = args[0]
	}

	// create & open the file if doesn't exist
	file, err := os.Create(filename)
	if (err != nil) {
			return errors.New("could not open file")
	}

	// if open successful, defer closing the file
	defer file.Close()

	// create a writer to write to the file
	writer := bufio.NewWriter(file)

	// write the header
	_, err = writer.WriteString(pointsCoordinatesLabels + "\n")
	if err != nil {
		return err
	}

	// write the points to the file
	for _, point := range points {
		// write the point to the file
		_, err := writer.WriteString(point.String() + "\n")
		if err != nil {
			fmt.Println("Error writing point to file: " + point.String() + " (" + err.Error() + ")")
			return err
		}
	}

	// flush the writer
	writer.Flush()

	return nil
}