package server

import (
	"context"
	"liteq/queue"
	"liteq/queue/proto"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// QueueService
func (*server) GetTasks(empty *emptypb.Empty, stream proto.LiteQ_GetTasksServer) error {

	/* Receive stream message
	 taskType, err := stream.Recv()
	log.Println("GetTasks:", taskType)
	// error handling
	{
		if err == io.EOF {
			log.Println("EOF")
			return nil
		}
		if err != nil {
			return err
		}
	}
	*/

	// stream all existing tasks
	for _, t := range *queue.Q.Tasks {
		if err := stream.Send(&proto.Task{
			Id:        t.ID,
			Data:      t.Data,
			CreatedAt: timestamppb.New(t.CreationDate),
			Status:    *proto.TaskStatus(t.Status).Enum(),
		}); err != nil {
			log.Println(err.Error())
		}
	}

	// keep streaming on new tasks

	go func() {
		for {
			t := <-queue.Q.TaskCh
			if err := stream.Send(&proto.Task{
				Id:        t.ID,
				Data:      t.Data,
				CreatedAt: timestamppb.New(t.CreationDate),
				Status:    *proto.TaskStatus(t.Status).Enum(),
			}); err != nil {
				log.Printf("(Failed to send new task through stream) Err: %s\n", err.Error())
			}
		}
	}()

	go func() {
		<-stream.Context().Done()
		log.Println("GetTasks: stream closed")
		if err := stream.Context().Err(); err != nil {
			log.Printf("GetTasks: %s\n", err.Error())
		}
	}()

	return nil

	/*
		for {
			select {
			case <-stream.Context().Done():
				if err := stream.Context().Err(); err != nil {
					log.Printf("GetTasks: %s\n", err.Error())
				}
				return nil
			case t := <-queue.Q.TaskCh:
				if err := stream.Send(&proto.Task{
					Id:        t.ID,
					Data:      t.Data,
					CreatedAt: timestamppb.New(t.CreationDate),
					Status:    *proto.Task_TaskStatus(t.Status).Enum(),
				}); err != nil {
					log.Printf("(Failed to send new task through stream) Err: %s\n", err.Error())
				}
			}
		}
	*/
}

func (*server) TaskStatusUpdate(ctx context.Context, req *proto.TaskStatusUpdateRequest) (*emptypb.Empty, error) {
	// get task from id
	for _, q := range *queue.Q.Tasks {
		if q.ID == req.Id {
			q.Status = queue.TaskStatus(req.Status)
			break
		}
	}
	return new(emptypb.Empty), nil
}
