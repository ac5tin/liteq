package proto

import (
	"liteq/queue"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Task2ProtoTask(task *queue.Task) *Task {
	t := new(Task)
	t.Id = task.ID
	t.Data = task.Data
	t.CreatedAt = timestamppb.New(task.CreationDate)
	t.Status = *TaskStatus(task.Status).Enum()
	return t
}

func ProtoTask2Task(task *Task) *queue.Task {
	t := new(queue.Task)
	t.ID = task.Id
	t.Data = task.Data
	t.CreationDate = task.CreatedAt.AsTime()
	t.Status = queue.TaskStatus(task.Status)
	return t
}
