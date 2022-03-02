package server

import (
	"liteq/queue"
	"liteq/queue/proto"

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

	for _, t := range *queue.Q.Tasks {
		stream.Send(&proto.Task{
			Id:        t.ID,
			Data:      t.Data,
			CreatedAt: timestamppb.New(t.CreationDate),
			Status:    *proto.Task_TaskStatus(t.Status).Enum(),
		})
	}

	return nil
}
