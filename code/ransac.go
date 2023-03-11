package ransac

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

// store the point cloud read from file
var pointCloud []Point3D
// default number of dominant planes to be identified
const DEFAULT_NUM_OF_DOMINANT_PLANES int = 3

// method to compute the number of iterations needed for RANSAC
func getNumberOfIterations(confidence float64, perctangeOfPointsOnPlane float64) float64 {
		// The number of iterations is computed as follows:
		// n = log(1 - confidence) / log(1 - (percentageOfPointsOnPlane)^3)
		// where n is the number of iterations, confidence is the probability that
		// at least one of the iterations will find a good model, and
		// percentageOfPointsOnPlane is the percentage of points that are on the
		// plane.
		return math.Log(1 - confidence) / math.Log(1 - math.Pow(perctangeOfPointsOnPlane, 3))
}

// method to identify the most dominant plane in the point cloud
func getDominantPlane(numOfIterations int, pointCloud PointCloud, eps float64) Plane3D {
		// iterate for the number of iterations
		// - get a plane by selecting 3 random points from the point cloud
		// - get the support of the plane
		// - if the support is greater than the best support, update the best support
		// return the plane with the best support
		
		// store the best plane
		plane3D := Plane3D{}
		// store the best plane's number of inliers
		bestSupport := Plane3DwSupport{}

		// iterate for the number of iterations, to identify plane with the best support
		for i := 0; i < numOfIterations; i++ {
				// select 3 random points from the point cloud and obtain a plane
				plane3D = GetPlane(pointCloud.GetRandomPoints())
				// get support of the current plane
				support := plane3D.GetSupport(pointCloud.points, eps)
				// check if the support is greater than the best support
				if support.SupportSize > bestSupport.SupportSize {
						// update the best support plane
						bestSupport = support
				}
		}
		// return plane with the best support
		return bestSupport.Plane3D
}

// method to retrieve given number of dominant planes from the point cloud
// returns an array containing dominant planes and the point cloud without the points belonging to the dominant planes
func getDominantPlanes(numOfIterations int, pointCloud PointCloud, eps float64, numOfDominantPlanes ...int) ([]Plane3D, PointCloud) {
		// store the dominant planes
		dominantPlanes := []Plane3D{}
		// if the number of dominant planes is not specified, set it to the default value
		if len(numOfDominantPlanes) == 0 {
				numOfDominantPlanes[0] = DEFAULT_NUM_OF_DOMINANT_PLANES
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
				cloud = cloud.RemovePlane(&dominantPlane, eps)
		}
		// return the array of dominant planes and the points cloud without the points belonging to the dominant planes
		return dominantPlanes, cloud
}

func main() {
	// main program must be supplied with 4 command line arguments
	if len(os.Args) != 4 {
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

	// calculate number of iterations
	numOfIterations := getNumberOfIterations(confidence, percentageOfPointsOnPlane)
	// get the PointCloud
	pointCloud, err := readXYZ(filename)
	// if error extracting point cloud, print error and exit
	if err != nil {
		fmt.Println("Unable to get Point Cloud", err)
		os.Exit(1)
	}

}

// method to parse command line arguments
func parseArguments(filename string, confidence string, percentageOfPointsOnPlane string, eps string) (string, float64, float64, float64, error) {
	// validate filename
	if filename == "" {
		return "", 0, 0, 0, fmt.Errorf("Filename cannot be empty")
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
