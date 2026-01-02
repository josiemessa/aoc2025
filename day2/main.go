package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/josiemessa/aoc2025/utils"
)

func main() {
	log.SetFlags(0)
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}

	lines := utils.ReadFileAsLines("input")
	// var result1 int
	var result2 int

	for i, idRange := range strings.Split(lines[0], ",") {
		split := strings.Split(idRange, "-")
		f, err := strconv.Atoi(split[0])
		if err != nil {
			log.Fatalf("could not parse first id in range for id range %q (%d): %v\n", idRange, i, err)
		}
		l, err := strconv.Atoi(split[1])
		if err != nil {
			log.Fatalf("could not parse last id in range for id range %q (%d): %v\n", idRange, i, err)
		}

		// log.Println(idRange)

		// Part 1
		// for i := f; i <= l; i++ {
		// 	if s := strconv.Itoa(i); isInvalidP1(s) {
		// 		result1 += i
		// 	}
		// }
		for i := f; i <= l; i++ {
			if s := strconv.Itoa(i); isInvalidP2(s) {
				log.Println(s)
				result2 += i
			}
		}
	}

	// fmt.Println("Part 1:", result1)
	fmt.Println("Part 2:", result2)
}

func isInvalidP1(a string) bool {
	if len(a)%2 != 0 {
		return false
	}
	half := len(a) / 2
	return a[0:half] == a[half:]
}

// e.g. let a be 789789. len(a) = 6
func isInvalidP2(a string) bool {
	for i := 2; i <= len(a); i++ {
		if len(a)%i != 0 {
			// not divisible by i, ignore
			continue
		}
		// segment length, i is the number of segments
		matching := true
		sl := len(a) / i
		for j := 0; j < i-1; j++ {
			seg1 := a[sl*j : sl*(j+1)]
			seg2 := a[sl*(j+1) : sl*(j+2)]
			matching = matching && (seg1 == seg2)
			if !matching {
				break
			}
		}
		if matching {
			return true
		}
	}
	return false
}
