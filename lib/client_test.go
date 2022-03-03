package lib

import (
	"context"
	"liteq/queue"
	"liteq/server"
	"log"
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

	values2 := []queue.Task{
		{
			ID:           "3",
			Data:         []byte("test3"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
		{
			ID:           "4",
			Data:         []byte("test4"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
	}

	queue.Q = queue.NewQueue()
	for _, v := range values {
		queue.Q.Add(&v)
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

	t.Run("Stream Existing", func(t *testing.T) {
		// stream existing task data
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
		log.Println("received tasks, asserting length")

		assert.Len(t, tasks, len(values))

		//t.Log(values2)
	})

	t.Run("Stream new", func(t *testing.T) {
		// stream existing task data
		// client
		t.Log("initialising client")
		c := new(Client)
		c.conn = conn
		t.Log("getting new tasks")
		ch, err := c.GetTasks()
		if err != nil {
			t.Errorf("Failed to execute client.GetTasks: %v", err)
		}

		// check tasks
		tasks := make([]queue.Task, 0)
		go func() {
			for {
				task := <-ch
				if task == nil {
					break
				}
				tasks = append(tasks, *task)
			}
		}()

		// add new tasks
		for _, v := range values2 {
			queue.Q.Add(&v)
		}

		time.Sleep(time.Second * 5)

		t.Log("received tasks, asserting length")

		assert.Len(t, tasks, len(values)+len(values2))
	})

}
