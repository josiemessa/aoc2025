package slowgraph

import (
	"container/heap"
	"math"
	"slices"

	"github.com/josiemessa/aoc2025/pkg/queue"
)

type GridCoord struct {
	X uint
	Y uint
}

type Heuristic interface {
	Distance(GridCoord, GridCoord) uint
	Neighbours(GridCoord) []GridCoord
}

type gridGraph struct {
	NumRows uint
	NumCols uint

	Data []rune
	Cost func(GridCoord, GridCoord) uint
	Heuristic
}

type ChessGridGraph struct {
	*gridGraph
}

type ManhattanGridGraph struct {
	*gridGraph
}

func LinesToGridGraph(lines []string, chess bool) gridGraph {
	g := gridGraph{
		NumRows: uint(len(lines)),
		NumCols: uint(len(lines[0])),
	}
	if chess {
		g.Heuristic = &ChessGridGraph{gridGraph: &g}
	} else {
		g.Heuristic = &ManhattanGridGraph{gridGraph: &g}
	}

	g.Data = make([]rune, g.NumRows*g.NumRows)

	var i int
	for _, line := range lines {
		for _, c := range line {
			g.Data[i] = c
			i++
		}
	}

	return g
}

func (g *gridGraph) GetCoordData(coord GridCoord) rune {
	i := coord.Y*g.NumCols + coord.X
	return g.Data[i]
}

// Chessboard/Chebyshev neighbours
func (g *ChessGridGraph) Neighbours(start GridCoord) []GridCoord {
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

// Axial/Manhattan neighbours
func (g *ManhattanGridGraph) Neighbours(start GridCoord) []GridCoord {
	var result []GridCoord
	candidates := [4][2]int{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0},
	}
	for _, c := range candidates {
		// out of bounds check
		checkRow := int(start.Y) + c[1]
		if checkRow < 0 || checkRow >= int(g.NumRows) {
			continue
		}
		// out of bounds check
		checkCol := int(start.X) + c[0]
		if checkCol < 0 || checkCol >= int(g.NumCols) {
			continue
		}

		result = append(result, GridCoord{X: uint(checkCol), Y: uint(checkRow)})
	}
	return result
}

func (g *ManhattanGridGraph) Distance(a GridCoord, b GridCoord) uint {
	return uint(math.Abs(float64(a.X)-float64(b.X)) + math.Abs(float64(a.Y)-float64(b.Y)))
}

// ChebyshevDistance is the distance between coords on a grid where diagonal moves are allowed
// aka chessboard distance
func (g *ChessGridGraph) Distance(a GridCoord, b GridCoord) uint {
	return uint(math.Max((math.Abs(float64(a.X) - float64(a.Y))), math.Abs(float64(a.Y)-float64(b.Y))))
}

// FloodFill finds every tile in the graph and executes f() on that
func (g *gridGraph) FloodFill(start GridCoord, f func(current GridCoord, neighbours []GridCoord)) {
	frontier := make(queue.Queue, 0)
	frontier.Enqueue(&queue.Item{Value: start})
	reached := map[GridCoord]struct{}{start: {}}

	for len(frontier) != 0 {
		current := frontier.Dequeue().Value.(GridCoord)
		// TODO: this assumes chess neighbours
		neighbours := g.Neighbours(current)
		for _, next := range neighbours {
			if _, ok := reached[next]; !ok {
				frontier.Enqueue(&queue.Item{Value: next})
				reached[next] = struct{}{}
			}
		}
		f(current, neighbours)
	}
}

// BreadthFirstSearch generates a map of which tile we came from to reach the current tile, starting at start.
// Goal provides an early exit criteria
func (g *gridGraph) BreadthFirstSearch(start GridCoord, goal GridCoord) map[GridCoord]GridCoord {
	frontier := make(queue.Queue, 0)
	frontier.Enqueue(&queue.Item{Value: start})
	cameFrom := map[GridCoord]GridCoord{start: {}}

	for len(frontier) != 0 {
		current := frontier.Dequeue().Value.(GridCoord)

		if current == goal {
			break
		}
		// TODO: this assumes chess neighbours
		for _, next := range g.Neighbours(current) {
			if _, ok := cameFrom[next]; !ok {
				frontier.Enqueue(&queue.Item{Value: next})
				cameFrom[next] = current
			}
		}
	}
	return cameFrom
}

// DijkstraSearch generates a map of the cost of reaching the current tile from the previous tiles,
// starting from start (used for Dijkstra)
// Graph must have a `Cost()` function defined on it so it can calculate the cost of travelling
// from one tile to another
func (g *gridGraph) DijkstraSearch(start GridCoord, goal GridCoord) map[GridCoord]GridCoord {
	frontier := make(queue.PriorityQueue, g.NumCols*g.NumRows)
	heap.Push(&frontier, queue.NewPriorityItem(start, 0))
	cameFrom := map[GridCoord]GridCoord{start: {}}
	costSoFar := map[GridCoord]uint{start: 0}

	for len(frontier) != 0 {
		current := heap.Pop(&frontier).(*queue.PriorityItem).GetValue().(GridCoord)

		if current == goal {
			break
		}

		// TODO: this assumes chess neighbours
		for _, next := range g.Neighbours(current) {
			newCost := costSoFar[current] + g.Cost(current, next)
			if _, ok := costSoFar[next]; !ok || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				heap.Push(&frontier, queue.NewPriorityItem(next, int(newCost)))
				cameFrom[next] = current
			}
		}
	}
	return cameFrom
}

// GreedyBestFirstSearch implements the GreedyBFS algorithm
func (g *gridGraph) GreedyBestFirstSearch(start GridCoord, goal GridCoord, h Heuristic) map[GridCoord]GridCoord {
	frontier := make(queue.PriorityQueue, g.NumCols*g.NumRows)
	heap.Push(&frontier, queue.NewPriorityItem(start, 0))
	cameFrom := map[GridCoord]GridCoord{start: {}}
	for len(frontier) != 0 {
		current := heap.Pop(&frontier).(*queue.PriorityItem).GetValue().(GridCoord)

		if current == goal {
			break
		}

		// TODO: this assumes chess neighbours
		for _, next := range g.Neighbours(current) {
			if _, ok := cameFrom[next]; !ok {
				heap.Push(&frontier, queue.NewPriorityItem(next, int(g.Distance(goal, next))))
				cameFrom[next] = current
			}
		}
	}
	return cameFrom
}

// AStar implements the A* algorithm
// Graph must have a Cost() defined to
func (g *gridGraph) AStarSearch(start GridCoord, goal GridCoord) map[GridCoord]GridCoord {
	frontier := make(queue.PriorityQueue, g.NumCols*g.NumRows)
	heap.Push(&frontier, queue.NewPriorityItem(start, 0))
	cameFrom := map[GridCoord]GridCoord{start: {}}
	costSoFar := map[GridCoord]uint{start: 0}

	for len(frontier) != 0 {
		current := heap.Pop(&frontier).(*queue.PriorityItem).GetValue().(GridCoord)

		if current == goal {
			break
		}

		// TODO: this assumes chess neighbours
		for _, next := range g.Neighbours(current) {
			newCost := costSoFar[current] + g.Cost(current, next)
			if _, ok := costSoFar[next]; !ok || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				heap.Push(&frontier, queue.NewPriorityItem(next, int(newCost+g.Distance(goal, next))))
				cameFrom[next] = current
			}
		}
	}
	return cameFrom
}

func (g *gridGraph) FindPath(start GridCoord, goal GridCoord, search map[GridCoord]GridCoord) []GridCoord {
	path := make([]GridCoord, 0)
	current := goal
	for current != start {
		path = append(path, current)
		current = search[current]
	}
	path = append(path, current)
	slices.Reverse(path)
	return path
}
