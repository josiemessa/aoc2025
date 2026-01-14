package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/josiemessa/aoc2025/pkg/utils"
)

func main() {
	log.SetFlags(0)
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}

	var result2 uint64

	lines := utils.ReadFileAsLines("input")
	for _, line := range lines {
		log.Printf("\n%v\n", line)

		// result1 += part1(line)
		result2 += part2(line)
	}

	// fmt.Println("Part 1:", result1)
	fmt.Println("Part 2:", result2)
}

func part1(line string) int {
	max10s, maxIndex := 0, 0
	// search up until the penultimate rune for the largest number
	for j := range len(line) - 1 {
		x := int(line[j] - 48)
		// want strictly greater than as if there's a repeated max character, want to use that for units rather than 10s
		if x > max10s {
			max10s = x
			maxIndex = j
		}
	}

	maxUnits := 0
	for j := maxIndex + 1; j < len(line); j++ {
		x := int(line[j] - 48)
		if x > maxUnits {
			maxUnits = x
		}
	}

	thisJoltage := (max10s * 10) + maxUnits
	log.Printf("%d [%d]", thisJoltage, maxIndex)
	return thisJoltage
}

func part2(line string) uint64 {
	nums := make([]int, len(line))
	for j := range line {
		nums[j] = int(line[j] - 48)
	}
	log.Println("Line:", line)
	log.Println("Nums:", nums)

	var joltage [12]int
	var indices [13]int
	for k := range 12 {
		for j := indices[k]; j < len(nums)-(11-k); j++ {
			if nums[j] > joltage[k] {
				joltage[k] = nums[j]
				// set the starting index for the next iteration
				indices[k+1] = j + 1
			}
		}
	}

	log.Println("Jolt:", joltage)
	log.Println("Indx:", indices)

	var result float64
	for k, v := range joltage {
		result += float64(v) * math.Pow(10, float64(11-k))
		log.Printf("%f +", result)
	}
	return uint64(result)
}
