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

// store dominant planes
type DominantPlane struct {
	Plane3D
	points []Point3D
	supportSize int
}

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

// method to identify the most dominant plane in the point cloud
func getDominantPlane(numOfIterations int, pointCloud PointCloud, eps float64) DominantPlane {
		// iterate for the number of iterations
		// - get a plane by selecting 3 random points from the point cloud
		// - get the support of the plane
		// - if the support is greater than the best support, update the best support
		// return the plane with the best support
		
		// temporarily store plane
		plane3D := Plane3D{}
		// store the best dominant plane (containing the plane and the points that support it)
		bestSupport := DominantPlane{}

		// iterate for the number of iterations, to identify plane with the best support
		for i := 0; i < numOfIterations; i++ {
				// select 3 random points from the point cloud and obtain a plane
				plane3D = GetPlane(pointCloud.GetRandomPoints())
				// get supporting points of the current plane
				supportPoints := plane3D.GetSupportingPoints(pointCloud.points, eps)
				// check if the support is greater than the best support
				if len(*supportPoints) > bestSupport.supportSize {
						// update the best support dominant plane
						bestSupport = DominantPlane{plane3D, *supportPoints, len(*supportPoints)}
				}
		}
		// return plane the dominant plane with the best support
		return bestSupport
}

// method to retrieve given number of dominant planes from the point cloud
// returns an array containing dominant planes and the point cloud without the points belonging to the dominant planes
func getDominantPlanes(numOfIterations int, pointCloud PointCloud, eps float64, numOfDominantPlanes ...int) ([]DominantPlane, PointCloud) {
	fmt.Println("Identifying dominant planes...")
		// store the dominant planes
		dominantPlanes := []DominantPlane{}
		// if the number of dominant planes is not specified, set it to the default value
		if len(numOfDominantPlanes) == 0 {
				numOfDominantPlanes = []int{DEFAULT_NUM_OF_DOMINANT_PLANES}
		}
		// store the point cloud
		cloud := pointCloud
		// iterate for the number of dominant planes to be identified
		for i := 0; i < numOfDominantPlanes[0]; i++ {
				// identify the dominant plane from given point cloud
				dominantPlane := getDominantPlane(numOfIterations, cloud, eps)
				// append the dominant plane to the array of dominant planes
				dominantPlanes = append(dominantPlanes, dominantPlane)
				// remove the points on the dominant plane from the point cloud
				cloud = cloud.RemovePlane(&dominantPlane.Plane3D, eps)
				saveXYZ("data/output/points_removed_" + strconv.Itoa(i) + ".xyz", cloud.points)
		}

		fmt.Println("Dominant planes identified successfully")

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
		err := saveXYZ(filename + strconv.Itoa(i+1) + ".xyz", plane.points)
		// if error saving dominant plane, print error and exit
		if err != nil {
			fmt.Println("Unable to save dominant plane", err)
			os.Exit(1)
		}
		// print size of each dominant plane
		fmt.Printf("Dominant plane %d size: %d points \n", i+1, plane.supportSize)
		// update the size of points covered by dominant planes
		dominantPlanesSize += plane.supportSize
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