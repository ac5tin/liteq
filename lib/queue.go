package lib

import (
	"context"
	"io"
	"liteq/queue"
	"liteq/queue/proto"
	"log"
)

func (c *Client) GetTasks(status queue.TaskStatus) (<-chan *queue.Task, error) {

	client := proto.NewLiteQClient(c.conn)

	// stream tasks
	req := new(proto.GetTaskRequest)
	req.Status = proto.TaskStatus(status)
	stream, err := client.GetTasks(context.Background(), req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *queue.Task)
	ctx := stream.Context()

	go func() {
		for {
			t, err := stream.Recv()
			// error handling
			{

				if err == io.EOF {
					// EOF end of stream
					// log.Println("[client] EOF") // debug
					close(ch)
					return
				}

				if err != nil {
					log.Printf("[client] Failed to get tasks, Err: %v\n", err) // debug
					close(ch)
					return
				}
			}
			// handle nil
			{
				if t == nil {
					// log.Println("[client] GetTasks: nil") // debug
					continue
				}
			}

			// parse task data
			task := proto.ProtoTask2Task(t)
			// send task to channel
			ch <- task

		}
	}()

	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Printf("[client] %s\n", err.Error()) // debug
		}
		// close if not already closed
		if _, ok := <-ch; ok {
			//log.Println("[client] ctx end, closing channel") // debug
			close(ch)
		}

	}()

	return ch, nil
}

func (c *Client) UpdateTask(id string, status queue.TaskStatus) error {
	client := proto.NewLiteQClient(c.conn)
	// prepare request
	request := new(proto.TaskStatusUpdateRequest)
	request.Id = id
	request.Status = proto.TaskStatus(status)
	// send request
	if _, err := client.TaskStatusUpdate(context.Background(), request); err != nil {
		return err
	}

	return nil
}

func (c *Client) AddTask(t *queue.Task) error {
	client := proto.NewLiteQClient(c.conn)
	// prepare request
	request := proto.Task2ProtoTask(t)
	// send request
	if _, err := client.AddTask(context.Background(), request); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetCurrentTasks(status queue.TaskStatus) (*[]*queue.Task, error) {
	client := proto.NewLiteQClient(c.conn)
	// prepare request
	request := new(proto.GetTaskRequest)
	request.Status = proto.TaskStatus(status)
	resp, err := client.GetCurrentTasks(context.Background(), request)
	if err != nil {
		return nil, err
	}

	// return
	tasks := make([]*queue.Task, 0)
	for _, r := range resp.Tasks {
		t := proto.ProtoTask2Task(r)
		tasks = append(tasks, t)
	}

	return &tasks, nil
}
