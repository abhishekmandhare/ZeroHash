package queue

import "fmt"

type Queue[T any] struct {
	queue []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		queue: make([]T, 0),
	}
}

func (q *Queue[T]) Push(data T) {
	q.queue = append(q.queue, data)
}

func (q *Queue[T]) Pop() (T, error) {
	var result T
	if q.IsEmpty() {

		return result, fmt.Errorf("pop failed: queue empty")
	}

	result = q.queue[0]
	q.queue = q.queue[1:]
	return result, nil
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.queue) == 0
}

func (q *Queue[T]) Len() int {
	return len(q.queue)
}

func (q *Queue[T]) GetQueue() []T {
	return q.queue
}
