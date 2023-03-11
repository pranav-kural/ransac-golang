# Planning

## Main Program

1. Read XYZ file, and create array/slice of Point3D objects
2. Create variable of type Point3DwSupport (named 'bestSupport')
3. Compute number of iterations required
4. Create and start a pipeline to find dominant plane, which stops after the required number of iterations
5. Once pipeline is finished, save a file with points of dominant plane, and a file containing all points of points cloud without the points of the dominant plane

## Components

In order:

1.  Random point generator
    1. Select random point from given points cloud
    2. Output channel transmits instance of Point3D
2.  Triplet of points generator
    1. Reads using input channel 3 points from Random Point Generator
    2. Output channel sends array of Point3D (containing those 3 points)
3.  TakeN
    1. Input channel received an array of Point3D (containing 3 points)
    2. Output channel sends array of Point3D (containing those 3 points)
    3. Repeats until 'N' arrays are received
       1. Output channel can be buffered with size N
4.  Plane estimator
    1. Input channel reads an array of Point3D (containing 3 points)
    2. Output channel transmits a Plane3D (computed using those 3 points)
5.  Supporting Points finder (multiplexed)
    1. Input channel reads Plane3D
    2. Counts number of points which support the plane
    3. Output channel transmits Point3DwSupport instance (containing plane parameters and the number of supporting points)
6.  Fan in
    1. Input channel reads Point3DwSupport instance (containing plane parameters and the number of supporting points)
    2. Combines received instances of into one
    3. From multiple input channels to one output channel
7.  Dominant plane identifier (end)
    1. Received Plane3DwSupport instances and keeps in memory the plane with the best support![Alt text](https://cwstatic.nyc3.digitaloceanspaces.com/2021/06/16/images/2/CrossTower%20Exchange%20Trading%20View.png)

# Research

## Pipelines

- A pipeline is a series of stages connected by channels
  - each stage is a group of goroutines running the same function
- In each stage, the goroutines
  - receive values from upstream via inbound channels
  - perform some function on that data, usually producing new values
  - send values downstream via outbound channels
- The first stage is sometimes called the source or producer; the last stage, the sink or consumer

### Fan-out, fan-in

Multiple functions can read from the same channel until that channel is closed; this is called **fan-out**. This provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O.

A function can read from multiple inputs and proceed _until_ all are closed by multiplexing the input channels onto a single channel that’s closed when all the inputs are closed. This is called **fan-in**.

To implement fan-in, we introduce a merge function:

The merge function converts a list of channels to a single channel by starting a goroutine for each inbound channel that copies the values to the sole outbound channel. Once all the output goroutines have been started, merge starts one more goroutine to close the outbound channel after all sends on that channel are done.

```go
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

When the number of values to be sent is known at channel creation time, a buffer can simplify the code

### Guidelines for pipeline construction:

- stages close their outbound channels when all the send operations are done
- stages keep receiving values from inbound channels until those channels are closed or the senders are unblocked
- Pipelines unblock senders either by ensuring there’s enough buffer for all the values that are sent or by explicitly signalling senders when the receiver may abandon the channel
- When using buffered channels, you could read from the channel even after the channel has been closed by the upstream, as long as there are pending elements to be read in that outbound channel
- Use done channel pattern to signal cancellation

## Testing

### Performance

Initial performance test for sequential RANSAC implementation:

```
Test RANSAC run parameters:
Point Cloud Files:  [data/datasets/PointCloud1.xyz data/datasets/PointCloud2.xyz data/datasets/PointCloud3.xyz]
Number of tests:  30
Confidence:  0.99
Percentage of points on plane:  0.3
Epsilon:  0.5

Average run times:
PointCloud 1 : 0.143937529
PointCloud 2: 0.1748237916
PointCloud 3 : 0.17469329179999998
```

## Backlog

1. Using concurrency when reading file and creating points
