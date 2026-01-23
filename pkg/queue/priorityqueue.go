package queue

// https://pkg.go.dev/container/heap#example-package-PriorityQueue
import (
	"container/heap"
)

// An PriorityItem is something we manage in a priority queue.
type PriorityItem struct {
	value    any // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*PriorityItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*PriorityItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// TODO could move this somewhere else
// Update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *PriorityItem, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func NewPriorityItem(value any, priority int) *PriorityItem {
	return &PriorityItem{
		value:    value,
		priority: priority,
		index:    -1,
	}
}

func (i *PriorityItem) GetValue() any {
	return i.value
}
