package queue

var (
	Q *Queue
)

type Queue struct {
	Tasks  *[]*Task
	TaskCh chan *Task
}

func NewQueue() *Queue {
	ch := make(chan *Task)
	go func() {
		for {
			<-ch
		}
	}()
	return &Queue{
		Tasks:  new([]*Task),
		TaskCh: ch,
	}
}

func (q *Queue) Add(t *Task) {
	*q.Tasks = append(*q.Tasks, t)
	q.TaskCh <- t
}

func (q *Queue) GetCurrentAllTasks() *[]Task {
	t := new([]Task)
	for _, task := range *q.Tasks {
		*t = append(*t, *task)
	}
	return t
}
