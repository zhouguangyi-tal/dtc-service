package task

type TaskStatus int

const (
	Ready TaskStatus = iota
	Going
	End
)

type Task struct {
	name   string
	state  TaskStatus
	fn     func()
	second int64
}

func (t *Task) GetStatus() TaskStatus {
	return t.state
}

func (t *Task) CreateTask(name string, fn func()) {
	t.name = name
	t.fn = fn
	t.state = Ready
}

func (t *Task) Run() {
	if t.fn != nil {
		t.state = Going
		t.fn()
	}
}

func (t *Task) Stop() { //取消任务
	t.state = End
}
