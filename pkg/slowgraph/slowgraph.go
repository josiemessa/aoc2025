package slowgraph

import (
	"container/heap"
	"math"
	"slices"

	"github.com/josiemessa/aoc2025/pkg/queue"
)

type Coord struct {
	X uint
	Y uint
}

type GridMover interface {
	Distance(Coord, Coord) uint
	Neighbours(gc Coord, cols uint, rows uint) []Coord
}

type GridGraph struct {
	NumRows uint
	NumCols uint

	Data  []rune
	Cost  func(Coord, Coord) uint
	Mover GridMover
}

type Chess struct {
	GridGraph
}

type Manhattan struct {
	GridGraph
}

// NewGraph parses a slice of strings, assuming that each character represents a new tile on the grid
// Graphmover specifies whether the movement on this grid is manhattan (no diagonals) or chess (diagonals)
func NewGraph(g GridMover, lines []string, costFunc func(Coord, Coord) uint) GridGraph {
	grid := GridGraph{
		NumRows: uint(len(lines)),
		NumCols: uint(len(lines[0])),
		Mover:   g,
		Cost:    costFunc,
	}

	grid.Data = make([]rune, grid.NumRows*grid.NumRows)

	var i int
	for _, line := range lines {
		for _, c := range line {
			grid.Data[i] = c
			i++
		}
	}

	return grid
}

func (g *GridGraph) GetCoordData(coord Coord) rune {
	i := coord.Y*g.NumCols + coord.X
	return g.Data[i]
}

// Chessboard/Chebyshev neighbours
func (g *Chess) Neighbours(start Coord, numCols, numRows uint) []Coord {
	var result []Coord
	for i := -1; i <= 1; i++ {
		// out of bounds check
		checkRow := int(start.Y) + i
		if checkRow < 0 || checkRow >= int(numRows) {
			continue
		}

		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				// don't check current value
				continue
			}

			// out of bounds check
			checkCol := int(start.X) + j
			if checkCol < 0 || checkCol >= int(numCols) {
				continue
			}

			result = append(result, Coord{X: uint(checkCol), Y: uint(checkRow)})
		}
	}
	return result
}

// Axial/Manhattan neighbours
func (g *Manhattan) Neighbours(start Coord, numCols, numRows uint) []Coord {
	var result []Coord
	candidates := [4][2]int{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0},
	}
	for _, c := range candidates {
		// out of bounds check
		checkRow := int(start.Y) + c[1]
		if checkRow < 0 || checkRow >= int(numRows) {
			continue
		}
		// out of bounds check
		checkCol := int(start.X) + c[0]
		if checkCol < 0 || checkCol >= int(numCols) {
			continue
		}

		result = append(result, Coord{X: uint(checkCol), Y: uint(checkRow)})
	}
	return result
}

func (g *Manhattan) Distance(a Coord, b Coord) uint {
	return uint(math.Abs(float64(a.X)-float64(b.X)) + math.Abs(float64(a.Y)-float64(b.Y)))
}

// ChebyshevDistance is the distance between coords on a grid where diagonal moves are allowed
// aka chessboard distance
func (g *Chess) Distance(a Coord, b Coord) uint {
	return uint(math.Max((math.Abs(float64(a.X) - float64(a.Y))), math.Abs(float64(a.Y)-float64(b.Y))))
}

// FloodFill finds every tile in the graph and executes f() on that
func (g *GridGraph) FloodFill(start Coord, f func(current Coord, neighbours []Coord)) {
	frontier := make(queue.Queue, 0)
	frontier.Enqueue(&queue.Item{Value: start})
	reached := map[Coord]struct{}{start: {}}

	for len(frontier) != 0 {
		current := frontier.Dequeue().Value.(Coord)
		// TODO: this assumes chess neighbours
		neighbours := g.Mover.Neighbours(current, g.NumCols, g.NumRows)
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
func (g *GridGraph) BreadthFirstSearch(start Coord, goal Coord) map[Coord]Coord {
	frontier := make(queue.Queue, 0)
	frontier.Enqueue(&queue.Item{Value: start})
	cameFrom := map[Coord]Coord{start: {}}

	for len(frontier) != 0 {
		current := frontier.Dequeue().Value.(Coord)

		if current == goal {
			break
		}
		// TODO: this assumes chess neighbours
		for _, next := range g.Mover.Neighbours(current, g.NumCols, g.NumRows) {
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
func (g *GridGraph) DijkstraSearch(start Coord, goal Coord) map[Coord]Coord {
	frontier := make(queue.PriorityQueue, g.NumCols*g.NumRows)
	heap.Push(&frontier, queue.NewPriorityItem(start, 0))
	cameFrom := map[Coord]Coord{start: {}}
	costSoFar := map[Coord]uint{start: 0}

	for len(frontier) != 0 {
		current := heap.Pop(&frontier).(*queue.PriorityItem).GetValue().(Coord)

		if current == goal {
			break
		}

		// TODO: this assumes chess neighbours
		for _, next := range g.Mover.Neighbours(current, g.NumCols, g.NumRows) {
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
func (g *GridGraph) GreedyBestFirstSearch(start Coord, goal Coord) map[Coord]Coord {
	frontier := make(queue.PriorityQueue, g.NumCols*g.NumRows)
	heap.Push(&frontier, queue.NewPriorityItem(start, 0))
	cameFrom := map[Coord]Coord{start: {}}
	for len(frontier) != 0 {
		current := heap.Pop(&frontier).(*queue.PriorityItem).GetValue().(Coord)

		if current == goal {
			break
		}

		// TODO: this assumes chess neighbours
		for _, next := range g.Mover.Neighbours(current, g.NumCols, g.NumRows) {
			if _, ok := cameFrom[next]; !ok {
				heap.Push(&frontier, queue.NewPriorityItem(next, int(g.Mover.Distance(goal, next))))
				cameFrom[next] = current
			}
		}
	}
	return cameFrom
}

// AStar implements the A* algorithm
// Graph must have a Cost() defined to
func (g *GridGraph) AStarSearch(start Coord, goal Coord) map[Coord]Coord {
	frontier := make(queue.PriorityQueue, g.NumCols*g.NumRows)
	heap.Push(&frontier, queue.NewPriorityItem(start, 0))
	cameFrom := map[Coord]Coord{start: {}}
	costSoFar := map[Coord]uint{start: 0}

	for len(frontier) != 0 {
		current := heap.Pop(&frontier).(*queue.PriorityItem).GetValue().(Coord)

		if current == goal {
			break
		}

		// TODO: this assumes chess neighbours
		for _, next := range g.Mover.Neighbours(current, g.NumCols, g.NumRows) {
			newCost := costSoFar[current] + g.Cost(current, next)
			if _, ok := costSoFar[next]; !ok || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				heap.Push(&frontier, queue.NewPriorityItem(next, int(newCost+g.Mover.Distance(goal, next))))
				cameFrom[next] = current
			}
		}
	}
	return cameFrom
}

func (g *GridGraph) FindPath(start Coord, goal Coord, search map[Coord]Coord) []Coord {
	path := make([]Coord, 0)
	current := goal
	for current != start {
		path = append(path, current)
		current = search[current]
	}
	path = append(path, current)
	slices.Reverse(path)
	return path
}
