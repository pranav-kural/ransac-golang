package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pranav-kural/ransac-golang/code"
)

// method to parse command line arguments
func parseArguments(filename string, confidence string, percentageOfPointsOnPlane string, eps string) (string, float64, float64, float64, error) {
	// validate filename
	if filename == "" {
		return "", 0, 0, 0, fmt.Errorf("filename cannot be empty")
	}

	// parse confidence
	conf, err := strconv.ParseFloat(confidence, 64)
	if err != nil {
		return "", 0, 0, 0, err
	}

	// parse percentage of points on plane
	per, err := strconv.ParseFloat(percentageOfPointsOnPlane, 64)
	if err != nil {
		return "", 0, 0, 0, err
	}

	// parse epsilon
	e, err := strconv.ParseFloat(eps, 64)
	if err != nil {
		return "", 0, 0, 0, err
	}

	return filename, conf, per, e, nil
}

func main() {
	// main program must be supplied with 4 command line arguments
	if len(os.Args) != 5 {
		fmt.Println("Invalid number of arguments: ", len(os.Args))
		fmt.Println("Usage: ransac <input file> <confidence> <percentage of points on plane> <eps>")
		os.Exit(1)
	}

	// parse arguments
	filename, confidence, percentageOfPointsOnPlane, eps, err := parseArguments(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	// if error parsing arguments, print error and exit
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Parsing arguments completed successfully")
	fmt.Println("Filename: ", filename)
	fmt.Println("Confidence: ", confidence)
	fmt.Println("Epsilon: ", eps)

	// run RANSAC algorithm
	code.RANSAC(filename, confidence, percentageOfPointsOnPlane, eps)
}