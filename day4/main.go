package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/josiemessa/aoc2025/pkg/slowgraph"
	"github.com/josiemessa/aoc2025/pkg/utils"
)

func main() {
	log.SetFlags(0)
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}

	// old()

	var result1 int
	graph := slowgraph.LinesToGridGraph(utils.ReadFileAsLines("input"), true)
	graph.FloodFill(slowgraph.GridCoord{X: 0, Y: 0}, func(current slowgraph.GridCoord, neighbours []slowgraph.GridCoord) {
		var paper int
		d := graph.GetCoordData(current)
		if d == '@' {
			for _, c := range neighbours {
				if graph.GetCoordData(c) == '@' {
					paper++
				}
			}
			if paper < 4 {
				result1++
				log.Println(current)
			}
		}
	})

	fmt.Println("Part 1:", result1)

}

func old() {
	var result1 int
	var resultString string

	lines := utils.ReadFileAsLines("input")
	for row, line := range lines {
		for col, char := range line {
			if char == '@' {
				if lookAround(lines, row, col) < 4 {
					result1++
					resultString += "x"
					continue
				}
			}
			resultString += string(char)
		}
		resultString += "\n"
	}

	fmt.Println("Part 1:", result1)
	log.Println(resultString)
}

func lookAround(lines []string, row int, col int) int {
	var result int
	for i := -1; i <= 1; i++ {
		// out of bounds check
		checkRow := row + i
		if checkRow < 0 || checkRow >= len(lines) {
			continue
		}

		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				// don't check current value
				continue
			}

			// out of bounds check
			checkCol := col + j
			if checkCol < 0 || checkCol >= len(lines[0]) {
				continue
			}

			if lines[checkRow][checkCol] == '@' {
				result++
			}
		}
	}
	return result
}
