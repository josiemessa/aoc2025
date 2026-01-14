package slowgraph

type GridCoord struct {
	X uint
	Y uint
}

type GridGraph struct {
	NumRows uint
	NumCols uint

	Data []rune
}

func LinesToGridGraph(lines []string) GridGraph {
	g := GridGraph{
		NumRows: uint(len(lines)),
		NumCols: uint(len(lines[0])),
	}

	g.Data = make([]rune, g.NumRows*g.NumRows)

	var i int
	for _, line := range lines {
		for _, c := range line {
			g.Data[i] = c
		}
	}

	return g
}

func (g *GridGraph) GetCoordData(coord GridCoord) rune {
	i := coord.Y*g.NumCols + coord.X
	return g.Data[i]
}

func (g *GridGraph) FloodFill(start GridCoord, f func(current GridCoord, neighbours []GridCoord)) {
	frontier := make(chan GridCoord)
	frontier <- start
	reached := map[GridCoord]struct{}{start: {}}

	for len(frontier) != 0 {
		current := <-frontier
		neighbours := g.Neighbours(current)
		for _, next := range neighbours {
			if _, ok := reached[next]; !ok {
				frontier <- next
				reached[next] = struct{}{}
				f(current, neighbours)
			}
		}
	}
	close(frontier)
}

func (g *GridGraph) Neighbours(start GridCoord) []GridCoord {
	var result []GridCoord
	for i := -1; i <= 1; i++ {
		// out of bounds check
		checkRow := int(start.Y) + i
		if checkRow < 0 || checkRow >= int(g.NumRows) {
			continue
		}

		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				// don't check current value
				continue
			}

			// out of bounds check
			checkCol := int(start.X) + j
			if checkCol < 0 || checkCol >= int(g.NumCols) {
				continue
			}

			result = append(result, GridCoord{X: uint(checkCol), Y: uint(checkRow)})
		}
	}
	return result
}
