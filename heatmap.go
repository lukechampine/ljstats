package main

import (
	"fmt"
	"github.com/dustin/go-heatmap"
	"github.com/dustin/go-heatmap/schemes"
	"image"
	"image/png"
	"os"
	"time"
)

// makeDataPoint is a trivial helper function for turning a set of input values into a DataPoint.
// It exists only to make it easier to change the x and y values of the heatmap.
// If this step were performed inside csvToPoints, Go would complain about unused variables.
func makeDataPoint(date, timeofday time.Time, seconds, rate, keyspiece float32, numblocks, keys int) heatmap.DataPoint {
	x := numblocks
	y := keyspiece
	return heatmap.P(float64(x), float64(y))
}

// csvToPoints reads games in csv format, extracts the desired values, and returns them as a list of DataPoints.
func csvToPoints(games [][]byte) (points []heatmap.DataPoint) {
	var (
		date, timeofday          time.Time
		seconds, rate, keyspiece float32
		numblocks, keys          int
	)
	for _, game := range games {
		date, _ = time.Parse("2006-01-02", string(game[:10]))
		timeofday, _ = time.Parse("15:04", string(game[12:16]))
		fmt.Sscanf(string(game[17:]), "%d,%f,%f,%d,%f", &numblocks, &seconds, &rate, &keys, &keyspiece)

		dp := makeDataPoint(date, timeofday, seconds, rate, keyspiece, numblocks, keys)
		points = append(points, dp)
	}
	return
}

// writeHeatmap generates a png image from the input points and writes it to a file.
// This is where you can modify the size of the output image, the radius of the dots, etc.
func writeHeatmap(points []heatmap.DataPoint, file *os.File) {
	img := heatmap.Heatmap(image.Rect(0, 0, 1024, 768), points, 10, 200, schemes.Classic)
	png.Encode(file, img)
}
