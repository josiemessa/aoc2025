// https://www.quasilyte.dev/blog/post/pathfinding/

package graph

import "math"

type GridCoord struct {
	X int
	Y int
}

type Grid struct {
	NumCols uint
	NumRows uint

	// serializes the values of the graph
	Data []byte

	tilesPerByte uint
	shiftFactor  uint
	encode       func(rune) uint8
	decode       func(uint8) rune
}

// tilesPerByte must be 8, 4, 2 or 1
func LinesToGrid(lines []string, tilesPerByte int, encode func(rune) uint8, decode func(uint8) rune) Grid {
	g := Grid{
		NumCols:      uint(len(lines[0])),
		NumRows:      uint(len(lines)),
		encode:       encode,
		decode:       decode,
		tilesPerByte: uint(tilesPerByte),
		shiftFactor:  8 / uint(tilesPerByte),
	}

	length := math.Ceil(float64(g.NumCols*g.NumRows) / float64(tilesPerByte))
	g.Data = make([]byte, int(length))

	var i uint

	for _, line := range lines {
		for _, char := range line {
			sliceIndex := i / g.tilesPerByte
			shift := (i % g.tilesPerByte) * g.shiftFactor // position in byte to append encoded value
			g.Data[sliceIndex] |= (byte(encode(char)) & g.getMask()) << shift
			i++
		}
	}

	return g
}

func (g *Grid) GetCellTile(c GridCoord) uint8 {
	x := uint(c.X)
	y := uint(c.Y)
	if x >= g.NumCols || y >= g.NumRows {
		return 0 // Out-of-bounds access - this is important for BFS
	}
	i := y*g.NumCols + x
	byteIndex := i / uint(g.tilesPerByte)
	shift := (i % g.tilesPerByte) * g.shiftFactor
	return (g.Data[byteIndex] >> shift) & g.getMask()
}

func (g *Grid) getMask() byte {
	return byte(math.Pow(2, float64(g.shiftFactor)) - 1)
}
