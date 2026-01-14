package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadFileAsLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open file", err)
	}
	defer f.Close()
	fmt.Println(f.Name())

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scanner error", err)
	}
	return lines
}

func ReadFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open file", err)
	}
	defer f.Close()
	fmt.Println(f.Name())
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal("read all error:", err)
	}
	return b
}

func SliceAtoi(input []string) []int {
	result := make([]int, len(input))
	for i, s := range input {
		x, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Could not parse line %d: %q", i, s)
		}
		result[i] = x
	}
	return result
}
