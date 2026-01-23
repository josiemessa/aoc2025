package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/josiemessa/aoc2025/pkg/utils"
)

type rangeList struct {
	starts []uint64
	ends   []uint64
}

func (r *rangeList) Len() int {
	return len(r.starts)
}

func (r *rangeList) Less(i, j int) bool {
	if r.starts[i] == r.starts[j] {
		return r.ends[i] < r.ends[j]
	}
	return r.starts[i] < r.starts[j]
}

func (r *rangeList) Swap(i, j int) {
	oldStarts := r.starts[i]
	r.starts[i] = r.starts[j]
	r.starts[j] = oldStarts

	oldEnds := r.ends[i]
	r.ends[i] = r.ends[j]
	r.ends[j] = oldEnds
}

func (r *rangeList) In(ingredient uint64, j int) (int, bool) {
	if j >= r.Len() {
		return -1, false
	}
	if ingredient == r.starts[j] {
		return j, true
	}
	if ingredient == r.ends[j] {
		return j, true
	}
	if ingredient > r.starts[j] {
		if ingredient < r.ends[j] {
			// inside this range, we found it!
			return j, true
		}
	}
	return -1, false
}

func main() {
	log.SetFlags(0)
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}

	startTime := time.Now()

	lines := utils.ReadFileAsLines("test-input")
	var index int

	for i, v := range lines {
		if v == "" {
			index = i
			break
		}
	}

	fresh := &rangeList{
		starts: make([]uint64, index),
		ends:   make([]uint64, index),
	}

	for i := 0; i < index; i++ {
		split := strings.Split(lines[i], "-")
		start, err := strconv.ParseUint(split[0], 10, 64)
		if err != nil {
			log.Fatalf("could not parse line %d %q: %s\n", i, lines[i], err)
		}
		end, err := strconv.ParseUint(split[1], 10, 64)
		if err != nil {
			log.Fatalf("could not parse line %d %q: %s\n", i, lines[i], err)
		}
		fresh.starts[i] = start
		fresh.ends[i] = end
	}

	// sort by smallest range start
	sort.Sort(fresh)

	var result1 int
	for i := index + 1; i < len(lines); i++ {
		ingredient, err := strconv.ParseUint(lines[i], 10, 64)
		if err != nil {
			log.Fatalf("could not parse line %d %q: %s\n", i, lines[i], err)
		}
		var found bool
		var idx int
		for j := 0; j < fresh.Len(); j++ {
			idx, found = fresh.In(ingredient, j)
			if found {
				break
			}
		}

		if found {
			result1++
			log.Printf("Found %d in range %d-%d\n", ingredient, fresh.starts[idx], fresh.ends[idx])
			continue
		}
	}

	fmt.Printf("Part 1: %d (%s)\n", result1, time.Since(startTime).String())

	// Part 2
	startTime = time.Now()
	var result2 uint64

	for i := 0; i < fresh.Len(); i++ {
		start := fresh.starts[i]
		end := fresh.ends[i]
		// check if this range overlaps with any of the next ranges
		for j := 1; j < fresh.Len()-i; j++ {
			if _, ok := fresh.In(fresh.starts[i+j], i); ok {
				if fresh.ends[i+j] > end {
					end = fresh.ends[i+j]
				}
			} else {
				// if we can't find the next starting range in the current range, we can stop looking as they are ordered
				i += j - 1
				break
			}
		}

		result2 += end - start + 1
	}

	fmt.Printf("Part 2: %d (%s)\n", result2, time.Since(startTime).String())

}
