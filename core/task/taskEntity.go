package task

type TaskStatus int

const (
	Ready TaskStatus = iota
	Going
	End
)

type Task struct {
	name         string
	state        TaskStatus
	fn           func()
	TickerSecond int
}

func (t *Task) GetStatus() TaskStatus {
	return t.state
}

func (t *Task) CreateTask(name string, fn func()) {
	t.name = name
	t.fn = fn
	t.state = Ready
}
func (t *Task) CreatTickerTask(name string, fn func(), second int) { //周期任务
	t.name = name
	t.fn = fn
	t.state = Ready
	t.TickerSecond = second
}

func (t *Task) Run() {
	if t.fn != nil {
		t.state = Going
		t.fn()
	}
	if t.TickerSecond > 0 {
		t.state = Ready
	} else {
		t.state = End //非定时任务执行完结束
	}

}

func (t *Task) Stop() { //取消任务
	t.state = End
}
