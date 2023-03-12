package code

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// default number of dominant planes to be identified
const DEFAULT_NUM_OF_DOMINANT_PLANES int = 3
// additional message are printed if DEBUG is true
var DEBUG bool = false

// method to compute the number of iterations needed for RANSAC
func getNumberOfIterations(confidence float64, perctangeOfPointsOnPlane float64) int {
		// The number of iterations is computed as follows:
		// n = log(1 - confidence) / log(1 - (percentageOfPointsOnPlane)^3)
		// where n is the number of iterations, confidence is the probability that
		// at least one of the iterations will find a good model, and
		// percentageOfPointsOnPlane is the percentage of points that are on the
		// plane.
		return int(math.Log(1 - confidence) / math.Log(1 - math.Pow(perctangeOfPointsOnPlane, 3)))
}

func DominantPlaneIdentifier(numOfIterations int, pointCloud PointCloud, eps float64) Plane3DwSupport {
	
	// receive array containing random points for numOfIterations
	randomPointsChan := pointCloud.TakeN(numOfIterations)

	// get plane
	plane := GetPlaneC(randomPointsChan)

	// get supporting points
	supportingPoints := pointCloud.GetSupportingPointsC(plane, eps)

	// get best plane
	bestPlane := fanIn(supportingPoints)
	
	return <-bestPlane
}

// fanIn method receives Plane3DwSupport instances from inbound channel and sends back the best plane with the most supporting points on the outbound channel
func fanIn(supportingPointsIn <-chan Plane3DwSupport) <-chan Plane3DwSupport {
	// outbound channel
	bestPlaneOut := make(chan Plane3DwSupport)
	// store best support
	bestSupport := 0
	// best plane
	var bestPlane Plane3DwSupport
	// goroutine to find the best plane
	go func() {
		defer close(bestPlaneOut)
		// until we receive plane on the inbound channel
		for plane := range supportingPointsIn {
			// get plane from inbound channel
			// if the plane has more supporting points than the current best plane
			if plane.SupportSize > bestSupport {
				// update the best support
				bestSupport = plane.SupportSize
				// update the best plane
				bestPlane = plane
			}
		}
		// send the best plane on the outbound channel
		bestPlaneOut <- bestPlane
	}()
	// return the outbound channel
	return bestPlaneOut
}

// method to retrieve given number of dominant planes from the point cloud
// returns an array containing dominant planes and the point cloud without the points belonging to the dominant planes
func getDominantPlanes(numOfIterations int, pointCloud PointCloud, eps float64, numOfDominantPlanes ...int) ([]Plane3DwSupport, PointCloud) {
		// store the dominant planes
		dominantPlanes := []Plane3DwSupport{}
		// if the number of dominant planes is not specified, set it to the default value
		if len(numOfDominantPlanes) == 0 {
				numOfDominantPlanes = []int{DEFAULT_NUM_OF_DOMINANT_PLANES}
		}
		// store the point cloud
		cloud := pointCloud
		// iterate for the number of dominant planes to be identified
		for i := 0; i < numOfDominantPlanes[0]; i++ {
				// identify the dominant plane from given point cloud
				dominantPlane := DominantPlaneIdentifier(numOfIterations, cloud, eps)
				// append the dominant plane to the array of dominant planes
				dominantPlanes = append(dominantPlanes, dominantPlane)
				// remove the points on the dominant plane from the point cloud
				cloud = cloud.RemovePlane(&dominantPlane.Plane3D, eps)
		}

		// return the array of dominant planes and the points cloud without the points belonging to the dominant planes
		return dominantPlanes, cloud
}

// method to get the output filename
func getOutputFilename(filename string) (file string) {
	// remove substring '.xyz' from filename if it exists, and dd output path
	// remove "/data/datasets/" from filename if it exists
	file = strings.Replace(filename, "data/datasets/", "", -1)
	file = "data/output/" + strings.Replace(file, ".xyz", "", -1) + "_p"
	return
}

func RANSAC(filename string, confidence, percentageOfPointsOnPlane, eps float64) {
	fmt.Println("Initiating RANSAC")
	// get the PointCloud
	pointCloud, err := readXYZ(filename)
	// if error extracting point cloud, print error and exit
	if err != nil {
		fmt.Println("Unable to get Point Cloud", err)
		os.Exit(1)
	}

	fmt.Println("Point Cloud extracted successfully")
	fmt.Println("Point Cloud size: ", len(pointCloud.points))

	// calculate number of iterations
	numOfIterations := getNumberOfIterations(confidence, percentageOfPointsOnPlane)
	fmt.Println("Number of iterations: ", numOfIterations)

	// get the dominant planes and the point cloud without the points belonging to the dominant planes
	dominantPlanes, cloud := getDominantPlanes(numOfIterations, pointCloud, eps)

	fmt.Println("RANSAC completed")
	fmt.Println("Number of dominant planes: ", len(dominantPlanes))

	// size of points covered by dominant planes
	dominantPlanesSize := 0

	// get the output filename
	filename = getOutputFilename(filename)

	// save each dominant plane to a file
	for i, plane := range dominantPlanes {
		err := saveXYZ(filename + strconv.Itoa(i+1) + ".xyz", plane.SupportingPoints)
		// if error saving dominant plane, print error and exit
		if err != nil {
			fmt.Println("Unable to save dominant plane", err)
			os.Exit(1)
		}
		// print size of each dominant plane
		fmt.Printf("Dominant plane %d size: %d points \n", i+1, plane.SupportSize)
		// update the size of points covered by dominant planes
		dominantPlanesSize += plane.SupportSize
	}

	fmt.Println("Dominant planes saved successfully")

	// save the point cloud without the points belonging to the dominant planes to a file
	saveXYZ(filename + "0.xyz", cloud.points)

	fmt.Println("Point cloud without dominant planes saved successfully")

	fmt.Println("Total number of points covered by dominant planes: ", dominantPlanesSize)
	fmt.Println("Total number of points not covered by dominant planes: ", len(cloud.points))
	fmt.Println("Total number of points: ", len(pointCloud.points))

	fmt.Println("Program completed successfully :)")
}

// method to print messages when DEBUG mode is on
func dprint(s string) {
	if DEBUG {
		fmt.Println(s)
	}
}