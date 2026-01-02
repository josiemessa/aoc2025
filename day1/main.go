package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/josiemessa/aoc2025/utils"
)

func main() {
	log.SetFlags(0)
	var debug = flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}

	currValue := 50
	result1, result2 := 0, 0

	lines := utils.ReadFileAsLines("input")
	log.Println("Starting calc value:", currValue)
	for i, line := range lines {
		log.Println("\n", line)
		n, err := parseLine(line)
		if err != nil {
			log.Fatalf("could not parse line %q (%d): %v\n", line, i, err)
		}

		currValue += n
		log.Println("Current calc value:", currValue)

		// easier to capture here as capturing with result1 double accounts
		if currValue == 0 {
			log.Println("hit")
			result2++
		} else if x := currValue / 100; x != 0 {
			// calculate if we went over +/-100 and if so how many times
			if x < 0 {
				x = x * -1
			}

			log.Println("hit", x)
			result2 += x
		}

		// if we've gone negative, excluding starting at 0, then add 1 for passing through 0.
		if currValue < 0 && currValue != n {
			log.Println("hit")
			result2++
		}

		// convert calculated value to dial value
		currValue = currValue % 100
		if currValue < 0 {
			currValue = 100 + currValue
		}
		if currValue == 0 {
			result1++
		}

		log.Println("Current dial value:", currValue)
	}

	fmt.Println("Part 1:", result1)
	fmt.Println("Part 2:", result2)
}

func parseLine(line string) (int, error) {
	if line[0] != 'R' && line[0] != 'L' {
		return 0, errors.New("line does not start with 'R' or 'L'")
	}

	n, err := strconv.Atoi(line[1:])
	if err != nil {
		return 0, err
	}

	if line[0] == 'L' {
		return n * -1, nil
	} else {
		return n, nil
	}

}

var testInput = `L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`
