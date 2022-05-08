package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueueSuccess(t *testing.T) {
	q := NewQueue[int]()

	q.Push(1)
	q.Push(2)

	ret, err := q.Pop()

	require.NoError(t, err)
	require.Equal(t, 1,ret)
}
