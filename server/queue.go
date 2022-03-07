package server

import (
	"context"
	"liteq/queue"
	"liteq/queue/proto"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

// QueueService
func (*server) GetTasks(req *proto.GetTaskRequest, stream proto.LiteQ_GetTasksServer) error {

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

	reqTaskStatus := queue.TaskStatus(req.Status)

	// stream all existing tasks
	for _, t := range *queue.Q.Tasks {
		if t.Status == reqTaskStatus {
			// log.Printf("[server] streaming existing task %s to [%s]\n", t.ID, reqTaskStatus) // debug
			if err := stream.Send(proto.Task2ProtoTask(t)); err != nil {
				log.Printf("[server] Failed to stream task, %s\n", err.Error())
				return err
			}
		}

	}

	// keep streaming on new tasks

	qc := queue.Q.RegisterStatusChan(reqTaskStatus)
	ch := *qc.Chan
	defer queue.Q.UnregisterStatusChan(qc)
taskChanLoop:
	for {
		select {
		case <-stream.Context().Done():
			// log.Println("[server] GetTasks: stream closed") // debug
			if err := stream.Context().Err(); err != nil {
				log.Printf("[server] GetTasks: %s\n", err.Error())
				return err
			}
			break taskChanLoop

		case t := <-ch:
			// log.Printf("[server] streaming (updated/new) task %s to [%s]\n", t.ID, reqTaskStatus) // debug
			if err := stream.Send(proto.Task2ProtoTask(t)); err != nil {
				log.Printf("[server] (Failed to send new task through stream) Err: %s\n", err.Error())
				return err
			}
		}

	}
	return nil
}

func (*server) TaskStatusUpdate(ctx context.Context, req *proto.TaskStatusUpdateRequest) (*emptypb.Empty, error) {
	// get task from id
	if err := queue.Q.UpdateTask(req.Id, queue.TaskStatus(req.Status)); err != nil {
		log.Printf("[server] %s\n", err.Error())
		return new(emptypb.Empty), err
	}
	return new(emptypb.Empty), nil
}

func (*server) AddTask(ctx context.Context, in *proto.Task) (*emptypb.Empty, error) {
	t := proto.ProtoTask2Task(in)
	queue.Q.Add(t)
	return new(emptypb.Empty), nil
}

func (*server) GetCurrentTasks(ctx context.Context, in *proto.GetTaskRequest) (*proto.TaskListResponse, error) {
	currentTasks := queue.Q.GetCurrentTasks(queue.TaskStatus(in.Status))

	tasks := make([]*proto.Task, 0)

	for _, t := range *currentTasks {
		task := proto.Task2ProtoTask(t)
		tasks = append(tasks, task)
	}

	resp := new(proto.TaskListResponse)
	resp.Tasks = tasks
	return resp, nil
}
