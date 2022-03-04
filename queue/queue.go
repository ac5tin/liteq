package queue

import (
	"errors"

	uf "github.com/ac5tin/usefulgo"
)

var (
	Q *Queue
)

type QueueChan struct {
	ID     string
	Chan   *chan *Task
	Status TaskStatus
}

type Queue struct {
	Tasks  *[]*Task
	TaskCh *map[TaskStatus]map[string]QueueChan // each TaskStatus has it's own individual Task channel
}

func NewQueueChan(status TaskStatus) *QueueChan {
	c := make(chan *Task)
	uid := uf.GenUUIDV4()
	qch := QueueChan{
		ID:     uid,
		Chan:   &c,
		Status: status,
	}
	return &qch
}

func NewQueue() *Queue {
	chMap := make(map[TaskStatus]map[string]QueueChan) // channel map
	// populate channel map
	for _, status := range TaskStatuses {
		// create map
		chMap[status] = make(map[string]QueueChan)
	}

	return &Queue{
		Tasks:  new([]*Task),
		TaskCh: &chMap,
	}
}

// Add New Task
func (q *Queue) Add(t *Task) {
	t.Status = TaskStatusCreated
	*q.Tasks = append(*q.Tasks, t)

	for _, qc := range (*q.TaskCh)[TaskStatusCreated] {
		*qc.Chan <- t
	}
}

func (q *Queue) UpdateTask(id string, status TaskStatus) error {
	found := false
	for _, task := range *q.Tasks {
		if task.ID == id {
			task.Status = status
			found = true
			// stream
			for _, qc := range (*q.TaskCh)[status] {
				*qc.Chan <- task
			}

			break
		}
	}

	if !found {
		return errors.New("failed to update task, unable to find ID")
	}

	return nil
}

func (q *Queue) RegisterStatusChan(status TaskStatus) *QueueChan {
	qc := NewQueueChan(status)
	(*q.TaskCh)[status][qc.ID] = *qc

	return qc
}

func (q *Queue) UnregisterStatusChan(c *QueueChan) {
	close(*c.Chan)
	delete((*q.TaskCh)[c.Status], c.ID)
}

func (q *Queue) GetCurrentAllTasks() *[]Task {
	t := new([]Task)
	for _, task := range *q.Tasks {
		*t = append(*t, *task)
	}
	return t
}
