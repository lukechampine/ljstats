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

func makeDataPoint(date, timeofday time.Time, seconds, rate, keyspiece float32, numblocks, keys int) heatmap.DataPoint {
	x := numblocks
	y := keyspiece
	return heatmap.P(float64(x), float64(y))
}

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

func writeHeatmap(points []heatmap.DataPoint, file *os.File) {
	img := heatmap.Heatmap(image.Rect(0, 0, 1024, 768), points, 10, 200, schemes.Classic)
	png.Encode(file, img)
}
