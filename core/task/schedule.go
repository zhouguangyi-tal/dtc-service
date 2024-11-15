package task

import (
	"log"
	"strconv"
	"time"
)

type TaskSchedule struct {
	taskBuffer map[int]map[string]*Task //待执行任务buffer,按ticker 周期触发
}

func (s *TaskSchedule) Init() {
	log.Println("taskSchedule init")
	s.taskBuffer = make(map[int]map[string]*Task)
}
func (s *TaskSchedule) Start() {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for second, tasks := range s.taskBuffer {
					if second == 0 { //立即任务
						for name, tk := range tasks {
							go func() {
								log.Println("执行任务:", name)
								tk.Run()
								delete(tasks, name)
							}()
						}
					} else {
						now := time.Now().Unix()
						if now%int64(second) == 0 {
							for name, tk := range tasks {
								if tk.GetStatus() == Ready {
									go func() {
										log.Println("执行定时任务:", name)
										tk.Run()
									}()
								}
							}
						}
					}
				}
			}
		}
	}()
}

func (s *TaskSchedule) AddTask(name string, tk *Task) {
	secondTime := tk.TickerSecond
	if s.taskBuffer[secondTime] == nil {
		s.taskBuffer[secondTime] = make(map[string]*Task)
	}
	s.taskBuffer[secondTime][name] = tk

	msg := "添加任务"
	if secondTime > 0 {
		msg = "添加定时任务"
	}
	msg += name
	log.Println(msg, strconv.Itoa(secondTime)+"s")
}

func (s *TaskSchedule) StopTask(name string) {
	for _, tasks := range s.taskBuffer {
		for tkName, tk := range tasks {
			if tkName == name {
				log.Println("停止定时任务", name)
				tk.Stop()
			}

		}
	}
}

func (s *TaskSchedule) DelTask(name string) {
	for _, tasks := range s.taskBuffer {
		for tkName, _ := range tasks {
			if tkName == name {
				log.Println("删除定时任务", name)
				delete(tasks, name)
			}

		}
	}
}
