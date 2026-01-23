package queue

import (
	"fmt"
	"strings"
)

type Item struct {
	Value any
}

type Queue []*Item

func (q *Queue) Enqueue(a *Item) {
	*q = append(*q, a)
}

func (q *Queue) Dequeue() *Item {
	if len(*q) == 0 {
		return nil
	}
	result := (*q)[0]
	(*q)[0] = nil
	*q = (*q)[1:]
	return result
}

func (q *Queue) CheckFront() *Item {
	if len(*q) == 0 {
		return nil
	}
	return (*q)[0]
}

func (q *Queue) ChcekBack() *Item {
	if len(*q) == 0 {
		return nil
	}
	return (*q)[len(*q)-1]
}

func (q *Queue) Length() int {
	return len(*q)
}

func (q *Queue) String() string {
	var s strings.Builder
	for _, v := range *q {
		fmt.Fprintf(&s, "%s ", v.Value)
	}
	return s.String()
}
