// test RANSAC performance

package test

import (
	"fmt"
	"time"

	"github.com/pranav-kural/ransac-golang/code"
)

func TestRANSAC() {
	/***********************************/
	/* Update Test RANSAC parameters below */

	// point cloud datasets
	pointCloudFiles := []string{ 
			"data/datasets/PointCloud1.xyz",
			"data/datasets/PointCloud2.xyz",
			"data/datasets/PointCloud3.xyz",
	}
	// number of point clouds
	numPC := len(pointCloudFiles)

	// number of tests
	n := 30

	// ransac parameters
	confidence := 0.99
	percentageOfPointsOnPlane := 0.3
	eps := 0.5

	/* End of Test RANSAC parameters */
	/***********************************/

	// print parameters
	fmt.Println("Test RANSAC run parameters:")
	fmt.Println("Point Cloud Files: ", pointCloudFiles)
	fmt.Println("Number of tests: ", n)
	fmt.Println("Confidence: ", confidence)
	fmt.Println("Percentage of points on plane: ", percentageOfPointsOnPlane)
	fmt.Println("Epsilon: ", eps)

	// variable to alternate between point cloud
	pc := 0

	// store run times for each point cloud in a slice
	runTimes := make([]float64, numPC)

	// for number of tests to perform
	for i := 0; i < n; i++ {
		// record start time
		start := time.Now()
		// run RANSAC
		code.TestRANSAC(pointCloudFiles[pc], confidence, percentageOfPointsOnPlane, eps)
		// record run time
		runTimes[pc] += time.Since(start).Seconds()
		// alternate between point cloud
		pc = (pc + 1) % numPC
	}

	// print average run times
	fmt.Println("Average run times:")
	for i := 0; i < numPC; i++ {
		fmt.Println("PointCloud", i+1, ": ", runTimes[i]/float64(n/numPC))
	}

	fmt.Println("Test completed")
}