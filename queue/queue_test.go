package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	values := []Task{
		{
			ID:           "1",
			Data:         []byte("test"),
			Status:       TaskStatusCreated,
			CreationDate: time.Now(),
		},
		{
			ID:           "2",
			Data:         []byte("test2"),
			Status:       TaskStatusCreated,
			CreationDate: time.Now(),
		},
	}

	values2 := []Task{
		{
			ID:           "3",
			Data:         []byte("test3"),
			Status:       TaskStatusCreated,
			CreationDate: time.Now(),
		},
		{
			ID:           "4",
			Data:         []byte("test4"),
			Status:       TaskStatusCreated,
			CreationDate: time.Now(),
		},
	}

	Q = NewQueue()
	for _, v := range values {
		value := v
		Q.Add(&value)
	}

	// assert length
	assert.Equal(t, len(values), len(*Q.Tasks))

	for _, v := range values2 {
		value := v
		Q.Add(&value)
	}

	// assert length
	assert.Equal(t, len(values)+len(values2), len(*Q.Tasks))
}
