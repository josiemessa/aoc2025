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

func main() {
	log.SetFlags(0)
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}

	startTime := time.Now()

	lines := utils.ReadFileAsLines("input")
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
		idx, found := sort.Find(fresh.Len(), func(j int) int {
			if ingredient == fresh.starts[j] {
				return 0
			}
			if ingredient == fresh.ends[j] {
				return 0
			}
			if ingredient > fresh.starts[j] {
				if ingredient < fresh.ends[j] {
					// inside this range, we found it!
					return 0
				}
				// must be greater than this range
				return 1
			}
			return -1
		})

		if found {
			result1++
			log.Printf("Found %d in range %d-%d\n", ingredient, fresh.starts[idx], fresh.ends[idx])
		}
	}

	fmt.Printf("Part 1: %d (%s)\n", result1, time.Since(startTime).String())
}
