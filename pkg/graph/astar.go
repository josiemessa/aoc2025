package graph

// func (g *GridGraph) FloodFill(start GridCoord, f func(current GridCoord, neighbours []GridCoord)) {
// 	frontier := make(chan GridCoord)
// 	frontier <- start
// 	reached := map[GridCoord]struct{}{start: struct{}{}}

// 	for len(frontier) != 0 {
// 		current := <-frontier
// 	}
// 	close(frontier)
// }

// func (g GridGraph) Neighbours(start GridCoord) []GridCoord {
// 	var result []GridCoord
// 	for i := -1; i <= 1; i++ {
// 		// out of bounds check
// 		checkRow := start.Y + i
// 		if checkRow < 0 || checkRow >= g.NumRows {
// 			continue
// 		}

// 		for j := -1; j <= 1; j++ {
// 			if i == 0 && j == 0 {
// 				// don't check current value
// 				continue
// 			}

// 			// out of bounds check
// 			checkCol := start.X + j
// 			if checkCol < 0 || checkCol >= g.NumCols {
// 				continue
// 			}

// 			result = append(result, GridCoord{X: checkCol, Y: checkRow, Value: })

// 		}
// 	}
// 	return result
// }

// if lines[checkRow][checkCol] == '@' {
// 	result++
// }
