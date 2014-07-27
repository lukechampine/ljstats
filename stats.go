package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

// processGame extracts and formats values from a game. Note that the "time" value is converted to seconds for easier processing.
func processGame(game []byte) (line []byte) {
	var date, timeofday, numblocks, time, rate, keys, keyspiece []byte
	_, err := fmt.Sscanf(string(game), "on %s at %s\n Played %s tetrominoes in %s (%s\nPressed %s keys (%s",
		&date, &timeofday, &numblocks, &time, &rate, &keys, &keyspiece)
	if err != nil {
		panic(err)
	}
	// special processing for values enclosed by parens
	rate = bytes.TrimRight(rate, "/min)")
	keyspiece = bytes.TrimRight(keyspiece, "/piece)")

	// convert time to seconds
	var min int
	var sec float32
	_, err = fmt.Sscanf(string(time), "%d:%f", &min, &sec)
	if err != nil {
		panic(err)
	}
	time = []byte(fmt.Sprintf("%.2f", float32(min)*60.0+sec))

	joined := bytes.Join([][]byte{date, timeofday, numblocks, time, rate, keys, keyspiece}, []byte{','})
	return append(joined, '\n')
}

// parse reads the lj-scores.txt data, removes unwanted games, and extracts the desired values, returning them in csv format/
func parse(data []byte) (lines [][]byte) {
	// split data into games
	games := bytes.Split(data, []byte("\r\n\r\n\r\n"))[1:]

	// parse each game
	for _, game := range games {
		// discard incomplete and/or non-standard games
		if !bytes.Contains(game, []byte("Cleared 40 lines")) ||
			!bytes.Contains(game, []byte("tetromino")) ||
			!bytes.Contains(game, []byte("Bag of Tetrominoes")) {
			continue
		}

		splitGame := bytes.Split(game, []byte("\r\n"))
		trimmedGame := bytes.Join(append(splitGame[4:6], splitGame[7]), []byte{'\n'})
		lines = append(lines, processGame(trimmedGame))
	}
	return
}

// normalize parses csv lines and removes those with "abnormal" values
func normalize(lines [][]byte) (normal [][]byte) {
	// scan values
	var numblocks, keys int
	var time, rate, keyspiece float32
	for _, line := range lines {
		fmt.Sscanf(string(line[17:]), "%d,%f,%f,%d,%f\n", &numblocks, &time, &rate, &keys, &keyspiece)
		// filter out values that do not lie inside reasonable ranges
		if numblocks < 125 && numblocks > 100 && time < 150.0 && time > 30.0 {
			normal = append(normal, line)
		}
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE:\n\t"+os.Args[0], "INPUT [OUTPUT]\n\nOUTPUT must be a .csv or .png file.\nIf OUTPUT is not specified, csv data is written to stdout.")
		return
	}
	// open file
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// parse args and create output file
	var outFile *os.File
	var makeHeatmap bool
	switch {
	case len(os.Args) < 3:
		outFile = os.Stdout

	case bytes.HasSuffix([]byte(os.Args[2]), []byte(".csv")):
		outFile, err = os.Create(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		makeHeatmap = false

	case bytes.HasSuffix([]byte(os.Args[2]), []byte(".png")):
		outFile, err = os.Create(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		makeHeatmap = true

	default:
		fmt.Println("invalid output file (must be a .csv or .png)")
		return
	}

	// parse data
	lines := parse(data)

	// remove outliers
	normal := normalize(lines)

	if makeHeatmap {
		// create heatmap and write to file
		writeHeatmap(csvToPoints(normal), outFile)
	} else {
		// output in CSV format
		outFile.Write([]byte("date,timeofday,numblocks,time,rate,keys,keyspiece\n"))
		for _, line := range normal {
			outFile.Write(line)
		}
	}
}
