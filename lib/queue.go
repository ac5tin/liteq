package lib

import (
	"context"
	"io"
	"liteq/queue"
	"liteq/queue/proto"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Client) GetTasks() (<-chan *queue.Task, error) {

	client := proto.NewLiteQClient(c.conn)

	// stream tasks
	stream, err := client.GetTasks(context.Background(), new(emptypb.Empty))
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
					// log.Println("EOF") // debug
					close(ch)
					return
				}
				if err != nil {
					log.Printf("GetTasks Err: %v\n", err)
					close(ch)
					return
				}
			}
			// handle nil
			{
				if t == nil {
					log.Println("GetTasks: nil")
					continue
				}
			}

			// parse task data
			task := new(queue.Task)
			task.ID = t.Id
			task.Data = t.Data
			task.CreationDate = t.CreatedAt.AsTime()
			task.Status = queue.TaskStatus(t.Status)
			// send task to channel
			ch <- task

		}
	}()

	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err.Error())
		}
		// close if not already closed
		if _, ok := <-ch; ok {
			close(ch)
		}

	}()

	return ch, nil
}
