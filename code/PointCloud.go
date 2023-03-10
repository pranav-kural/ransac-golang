package ransac

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

// separator used to separate the coordinates of a point in point cloud data file
var POINTS_SEPARATOR = " "

// method to read the points from a file, and return an array containing Point3D instances
func readXYZ(filename string) (points []Point3D, err error) {
		// validate filename
		if (filename == "") {
				return points, errors.New("No filename provided")
		}

		// open the file
		file, err := os.Open(filename)
		if (err != nil) {
				return points, errors.New("Could not open file")
		}

		// if open successful, defer closing the file
		defer file.Close()

		// create a scanner to read the file
		scanner := bufio.NewScanner(file)
		// temporary variable to store the point
		point := Point3D{}

		// read the file line by line and for each line read, extract the Point3D object and store it in the points array
		for scanner.Scan() {
			// attempt to extract a Point3D from the line information
			point, err = getPoint3D(scanner.Text())
			// if error, return the error and stop reading the file
			if err != nil {
				return points, err
			}
			// append the point to the points array
			points = append(points, point)
		}

		// return the points
		return points, nil
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
func saveXYZ(filename string, points []Point3D) error {
	// validate filename
	if (filename == "") {
			return errors.New("No filename provided")
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