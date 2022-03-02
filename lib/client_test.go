package lib

import (
	"context"
	"liteq/queue"
	"liteq/server"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient(t *testing.T) {
	values := []queue.Task{
		{
			ID:           "1",
			Data:         []byte("test"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
		{
			ID:           "2",
			Data:         []byte("test2"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
	}
	queue.Q = queue.NewQueue()
	for _, v := range values {
		queue.Q.Add(v)
	}
	// start server
	t.Log("starting server")
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(server.Dialer()))
	if err != nil {
		t.Errorf("DialContext: %v", err)
	}
	defer conn.Close()
	t.Log("server started")
	// client
	t.Log("initialising client")
	c := new(Client)
	c.conn = conn
	t.Log("getting task")
	ch, err := c.GetTasks()
	if err != nil {
		t.Errorf("Failed to execute client.GetTasks: %v", err)
	}

	// check tasks
	tasks := make([]queue.Task, 0)
	for {
		task := <-ch
		if task == nil {
			break
		}
		tasks = append(tasks, *task)
	}

	assert.Len(t, tasks, len(values))

	if err := c.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}

}
