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
		{
			ID:           "5",
			Data:         []byte("test5"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
		{
			ID:           "6",
			Data:         []byte("test6"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
	}

	values3 := []queue.Task{
		{
			ID:           "7",
			Data:         []byte("test7"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
		{
			ID:           "8",
			Data:         []byte("test8"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
		{
			ID:           "9",
			Data:         []byte("test9"),
			Status:       queue.TaskStatusCreated,
			CreationDate: time.Now(),
		},
	}

	queue.Q = queue.NewQueue()
	for _, v := range values {
		value := v
		queue.Q.Add(&value)
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
	defer c.Close()
	c.conn = conn
	t.Log("getting task")
	ch, err := c.GetTasks(queue.TaskStatusCreated)
	if err != nil {
		t.Errorf("Failed to execute client.GetTasks: %v", err)
	}

	tasks := make([]queue.Task, 0)
	go func() {
		for {
			task := <-ch
			if task == nil {
				break
			}
			tasks = append(tasks, *task)
			t.Logf("Received (new) task %s\n", task.ID)
		}
	}()

	time.Sleep(time.Second * 5)
	assert.Equal(t, len(values), len(tasks))

	for _, v := range values2 {
		value := v
		queue.Q.Add(&value)
	}

	time.Sleep(time.Second * 5)
	assert.Equal(t, len(values)+len(values2), len(tasks))
	assert.Equal(t, len(*queue.Q.Tasks), len(tasks))

	// update
	ch, err = c.GetTasks(queue.TaskStatusDone)
	if err != nil {
		t.Error(err.Error())
	}

	updatedTask := new(queue.Task)
	go func() {
		for {
			task := <-ch
			if task == nil {
				break
			}
			updatedTask = task
			t.Logf("Received (done) task %s\n", task.ID)
		}
	}()

	time.Sleep(time.Second * 2)

	if err := c.UpdateTask("3", queue.TaskStatusDone); err != nil {
		t.Error(err.Error())
	}
	if err := c.UpdateTask("2", queue.TaskStatusDone); err != nil {
		t.Error(err.Error())
	}
	time.Sleep(time.Second * 2)

	assert.NotNil(t, updatedTask)

	assert.Equal(t, "2", updatedTask.ID)

	// adding stuff again
	for _, v := range values3 {
		value := v
		queue.Q.Add(&value)
	}
	time.Sleep(time.Second * 5)
	assert.Equal(t, len(values)+len(values2)+len(values3), len(tasks))

	// connection
	assert.Equal(t, "READY", c.conn.GetState().String())
}
