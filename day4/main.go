package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/josiemessa/aoc2025/pkg/utils"
)

func main() {
	log.SetFlags(0)
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}

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

type cellType int

const (
	cellEmpty cellType = iota // 0
	cellPaper                 // 1
)

func encode(c rune) uint8 {
	if c == '@' {
		return uint8(cellPaper)
	}
	return uint8(cellEmpty)
}

func decode(c uint8) rune {
	if c == uint8(cellEmpty) {
		return '.'
	}
	return '@'
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
