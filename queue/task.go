package queue

import "time"

type TaskStatus uint8

const (
	TaskStatusCreated    TaskStatus = iota // just created, waiting to be processed
	TaskStatusProcessing                   // processing by a worker
	TaskStatusFailed                       // process failed
	TaskStatusDone                         // process completed
)

var (
	TaskStatuses []TaskStatus = []TaskStatus{
		TaskStatusCreated,
		TaskStatusProcessing,
		TaskStatusFailed,
		TaskStatusDone,
	}
)

func (s TaskStatus) Uint32() uint32 {
	return uint32(s)
}

func (s TaskStatus) String() string {
	switch s {
	case TaskStatusCreated:
		return "created"
	case TaskStatusProcessing:
		return "processing"
	case TaskStatusFailed:
		return "failed"
	case TaskStatusDone:
		return "done"
	default:
		return "unknown"
	}
}

type Task struct {
	ID           string
	Data         []byte
	CreationDate time.Time
	Status       TaskStatus
}
