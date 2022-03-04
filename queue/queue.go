package queue

var (
	Q *Queue
)

type Queue struct {
	Tasks  *[]*Task
	TaskCh map[TaskStatus]chan *Task // each TaskStatus has it's own individual Task channel
}

func NewQueue() *Queue {
	chMap := make(map[TaskStatus]chan *Task) // channel map
	// populate channel map
	for _, status := range TaskStatuses {
		// create channel
		ch := make(chan *Task)
		go func() {
			for {
				<-ch
			}
		}()

		// register channel to map
		chMap[status] = ch
	}

	return &Queue{
		Tasks:  new([]*Task),
		TaskCh: chMap,
	}
}

func (q *Queue) Add(t *Task) {
	*q.Tasks = append(*q.Tasks, t)
	q.TaskCh[TaskStatusCreated] <- t
}

func (q *Queue) GetCurrentAllTasks() *[]Task {
	t := new([]Task)
	for _, task := range *q.Tasks {
		*t = append(*t, *task)
	}
	return t
}
