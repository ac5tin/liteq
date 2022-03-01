package queue

var (
	Q *Queue
)

type Queue struct {
	Tasks *[]Task
}

func NewQueue() *Queue {
	return &Queue{
		Tasks: new([]Task),
	}
}

func (q *Queue) Add(t Task) {
	*q.Tasks = append(*q.Tasks, t)
}
