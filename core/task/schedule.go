package task

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type dailyTime struct {
	hour int
	min  int
}

type TaskSchedule struct {
	taskBuffer      map[string]*Task //待执行任务buffer,按ticker 周期触发
	dailyTaskBuffer map[dailyTime][]*Task
}

func (s *TaskSchedule) Init() {
	log.Println("taskSchedule init")
	s.taskBuffer = make(map[string]*Task)
	s.dailyTaskBuffer = make(map[dailyTime][]*Task)
}
func (s *TaskSchedule) Start() {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for name, tk := range s.taskBuffer {
					second := tk.second
					now := time.Now()
					if second == 0 || now.Unix()%second == 0 {
						go func() {
							log.Println("执行定时任务:", name)
							tk.Run()
						}()
					}
				}
				for nextTime, taskArr := range s.dailyTaskBuffer { //
					now := time.Now()
					next := time.Date(now.Year(), now.Month(), now.Day(), nextTime.hour, nextTime.min, 0, 0, now.Location())
					if now.Unix() == next.Unix() {
						for _, task := range taskArr {
							go func() {
								log.Println("执行定时任务:", task.name)
								task.Run()
							}()
						}
					}
				}
			}
		}
	}()
}

func (s *TaskSchedule) AddTask(tk *Task, second int64) {
	name := tk.name
	tk.second = second
	s.taskBuffer[name] = tk
	log.Println("添加任务", name, strconv.Itoa(int(second))+"s")
}

func (s *TaskSchedule) AddDailyTask(tk *Task, hour, minute int) {
	name := tk.name
	daily := dailyTime{
		hour: hour,
		min:  minute,
	}
	s.dailyTaskBuffer[daily] = append(s.dailyTaskBuffer[daily], tk)
	log.Println("添加每日任务", name, hour, minute)
}

func (s *TaskSchedule) StopTask(name string) {
	for tkName, tk := range s.taskBuffer {
		if tkName == name {
			log.Println("停止定时任务", name)
			tk.Stop()
		}
	}
	for _, taskArr := range s.dailyTaskBuffer {
		for _, task := range taskArr {
			if task.name == name {
				log.Println("停止每日任务", name)
				task.Stop()
			}
		}
	}
}

func (s *TaskSchedule) DelTask(name string) {
	for tkName, _ := range s.taskBuffer {
		if tkName == name {
			log.Println("删除定时任务", name)
			delete(s.taskBuffer, name)
		}
	}
	for daily, taskArr := range s.dailyTaskBuffer {
		for i := 0; i < len(taskArr); i++ {
			if taskArr[i].name == name {
				log.Println("删除每日任务", name)
				s.dailyTaskBuffer[daily] = append(taskArr[:i], taskArr[i+1:]...)
				break
			}
		}
	}
	fmt.Println("zzzz", s.dailyTaskBuffer)
}
