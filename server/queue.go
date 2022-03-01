package server

import (
	"io"
	"liteq/queue"
	"liteq/queue/proto"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// QueueService
func (*server) GetTasks(stream proto.LiteQ_GetTasksServer) error {
	for {
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

		for _, t := range *queue.Q.Tasks {
			stream.Send(&proto.Task{
				Id:        t.ID,
				Data:      t.Data,
				CreatedAt: timestamppb.New(t.CreationDate),
				Status:    *proto.Task_TaskStatus(t.Status).Enum(),
			})
		}

	}
}
