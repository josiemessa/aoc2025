package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinesToGraph8tiles(t *testing.T) {
	lines := []string{
		"..........",
		"@@@@@@@@@@",
		".@.@.@.@.@",
	}

	encode := func(c rune) uint8 {
		if c == '.' {
			return 0
		}
		return 1
	}

	decode := func(c uint8) rune {
		if c == 0 {
			return '.'
		}
		return '@'
	}

	grid := LinesToGrid(lines, 8, encode, decode)

	require.EqualValues(t, 3, grid.NumRows)
	require.EqualValues(t, 10, grid.NumCols)
	require.Equal(t, 4, len(grid.Data)) // 30 values, 8 values per byte => 4 bytes

	expected := []byte{
		0b00000000, // [0,0]-[7,0] (RTL)
		0b11111100, // [8,0]-[5,1]
		0b10101111, // [6,1]-[3,2]
		0b00101010, // [4,2]-[9,2] (ignore last two bits)
	}

	for i, actual := range grid.Data {
		require.Equal(t, expected[i], actual,
			"byte %d did not match:\nexpected: %b\nactual:   %b\n", i, expected[i], actual)
	}
}

func TestLinesToGraph4tiles(t *testing.T) {
	lines := []string{
		"..........",
		"@@@@@@@@@@",
		"++++++++++",
		"----------",
		".@+-.@+-.@",
	}

	encode := func(c rune) uint8 {
		if c == '.' {
			return 0
		}
		if c == '@' {
			return 1
		}
		if c == '+' {
			return 2
		}
		return 3
	}

	decode := func(c uint8) rune {
		if c == 0 {
			return '.'
		}
		if c == 1 {

			return '@'
		}
		if c == 2 {
			return '+'
		}
		return '-'
	}

	grid := LinesToGrid(lines, 4, encode, decode)

	require.EqualValues(t, 5, grid.NumRows)
	require.EqualValues(t, 10, grid.NumCols)
	require.Equal(t, 13, len(grid.Data)) // 50 values, 4 values per byte => 13 bytes

	expected := []byte{
		0b00000000, // [0,0]-[0,3]
		0b00000000, // [0,4]-[0,7]
		0b01010000, // [0,8]-[1,1]
		0b01010101,
		0b01010101,
		0b10101010, // [2,0]-[2,3]
		0b10101010,
		0b11111010, // [2,8]-[3,1]
		0b11111111,
		0b11111111,
		0b11100100, // [4,0]-[4,3]
		0b11100100,
		0b00000100,
	}

	for i, actual := range grid.Data {
		require.Equal(t, expected[i], actual,
			"byte %d did not match:\nexpected: %b\nactual:   %b\n", i, expected[i], actual)
	}
}

func TestGetCellTile8Tiles(t *testing.T) {
	lines := []string{
		"..........",
		"@@@@@@@@@@",
		".@.@.@.@.@",
	}

	encode := func(c rune) uint8 {
		if c == '.' {
			return 0
		}
		return 1
	}

	decode := func(c uint8) rune {
		if c == 0 {
			return '.'
		}
		return '@'
	}

	grid := LinesToGrid(lines, 8, encode, decode)

	require.Equal(t, uint8(0b0), grid.GetCellTile(GridCoord{0, 0}), "coord [0, 0]")
	require.Equal(t, uint8(0b1), grid.GetCellTile(GridCoord{9, 2}), "coord [9, 2]")
	require.Equal(t, uint8(0b1), grid.GetCellTile(GridCoord{7, 1}), "coord [7, 1]")
	require.Equal(t, uint8(0b0), grid.GetCellTile(GridCoord{4, 2}), "coord [4, 2]")
	require.Equal(t, uint8(0b0), grid.GetCellTile(GridCoord{100, 100}), "coord [100,100] (out of bounds)")
}

func TestGetCellTile4Tiles(t *testing.T) {
	lines := []string{
		"..........",
		"@@@@@@@@@@",
		"++++++++++",
		"----------",
		".@+-.@+-.@",
	}

	encode := func(c rune) uint8 {
		if c == '.' {
			return 0
		}
		if c == '@' {
			return 1
		}
		if c == '+' {
			return 2
		}
		return 3
	}

	decode := func(c uint8) rune {
		if c == 0 {
			return '.'
		}
		if c == 1 {

			return '@'
		}
		if c == 2 {
			return '+'
		}
		return '-'
	}

	grid := LinesToGrid(lines, 4, encode, decode)

	require.Equal(t, uint8(0b00), grid.GetCellTile(GridCoord{0, 0}), "coord [0, 0]")
	require.Equal(t, uint8(0b01), grid.GetCellTile(GridCoord{9, 4}), "coord [9, 4]")
	require.Equal(t, uint8(0b11), grid.GetCellTile(GridCoord{7, 3}), "coord [7, 3]")
	require.Equal(t, uint8(0b10), grid.GetCellTile(GridCoord{4, 2}), "coord [4, 2]")
	require.Equal(t, uint8(0b00), grid.GetCellTile(GridCoord{100, 100}), "coord [100,100] (out of bounds)")
}
