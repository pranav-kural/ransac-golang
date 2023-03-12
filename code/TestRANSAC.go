package code

import (
	"fmt"
	"os"
)

func TestRANSAC(filename string, confidence, percentageOfPointsOnPlane, eps float64) {
	fmt.Println("Test RANSAC run")
	// get the PointCloud
	pointCloud, err := readXYZ(filename)
	// if error extracting point cloud, print error and exit
	if err != nil {
		fmt.Println("Unable to get Point Cloud", err)
		os.Exit(1)
	}

	// calculate number of iterations
	numOfIterations := getNumberOfIterations(confidence, percentageOfPointsOnPlane)

	// get the dominant planes and the point cloud without the points belonging to the dominant planes
	dominantPlanes, cloud := getDominantPlanes(numOfIterations, pointCloud, eps)

	// size of points covered by dominant planes
	dominantPlanesSize := 0

	// total size of points covered by dominant planes
	for _, plane := range dominantPlanes {
		// update the size of points covered by dominant planes
		dominantPlanesSize += plane.SupportSize
	}

	// verify result
	if dominantPlanesSize + len(cloud.points) != len(pointCloud.points) {
		fmt.Println("RANSAC result is incorrect")
		os.Exit(1)
	} else {
		fmt.Println("Test RANSAC run completed")
	}
}