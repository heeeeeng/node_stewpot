package main

import (
	"fmt"
	"github.com/heeeeeng/node_stewpot/tasks"
	"github.com/heeeeeng/node_stewpot/types"
	"sync"
	"time"
)

type Timeline struct {
	timestamp int64

	current *TimeUnit
	next    *TimeUnit
	db      *MemDB

	mu    sync.RWMutex
	close chan struct{}
}

func newTimeline(db *MemDB) *Timeline {
	t := &Timeline{}

	t.timestamp = 0
	t.current = newTimeUnit(t.timestamp)
	t.next = newTimeUnit(t.timestamp + 1)
	t.db = db

	t.close = make(chan struct{})

	return t
}

func (tl *Timeline) Start() {
	go tl.loop()
}

func (tl *Timeline) Stop() {
	close(tl.close)
}

func (tl *Timeline) CurrentTime() int64 { return tl.timestamp }

func (tl *Timeline) NextTime() int64 { return tl.timestamp + 1 }

//func (tl *Timeline) ImportCurrentTask(task Task) {
//	tl.mu.Lock()
//	defer tl.mu.Unlock()
//
//	tl.importTask(tl.CurrentTime(), task)
//}
//
//func (tl *Timeline) ImportNextTask(task Task) {
//	tl.mu.Lock()
//	defer tl.mu.Unlock()
//
//	tl.importTask(tl.NextTime(), task)
//}

func (tl *Timeline) SendNewMsg(node *Node, msg types.Message) int64 {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	task := tasks.NewMsgProcessTask(tl.NextTime(), node, msg)
	tl.ImportTask(task.StartTime(), task)
	return task.StartTime()
}

func (tl *Timeline) ImportTask(startTime int64, task types.Task) {
	if startTime != tl.current.timestamp || startTime != tl.next.timestamp {
		// TODO
		// log out this situation, this should not happen.
		return
	}
	if startTime == tl.current.timestamp {
		tl.current.appendTask(task)
		return
	}
	if startTime == tl.next.timestamp {
		tl.next.appendTask(task)
		return
	}
	return
}

func (tl *Timeline) GetTimeUnit(t int64) *TimeUnit {
	if t == tl.current.timestamp {
		return tl.current
	} else if t == tl.next.timestamp {
		return tl.next
	} else {
		return tl.db.Get(t)
	}
}

func (tl *Timeline) loop() {
	for {
		select {
		case <-tl.close:
			return

		default:
			if tl.current == nil {
				continue
			}

			for {
				task := tl.current.nextTask()
				if task == nil {
					break
				}
				task.Process(tl)
			}

			tl.mu.Lock()
			tl.db.Insert(tl.current)
			tl.timestamp++
			tl.current = tl.next
			tl.next = newTimeUnit(tl.timestamp + 1)

			tl.mu.Unlock()

			time.Sleep(time.Millisecond)
		}
	}
}

type TimeUnit struct {
	timestamp int64
	tasks     []types.Task

	mu sync.RWMutex
}

func newTimeUnit(timestamp int64) *TimeUnit {
	return &TimeUnit{
		timestamp: timestamp,
		tasks:     make([]types.Task, 0),
	}
}

func (t *TimeUnit) appendTask(task types.Task) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fmt.Println("appendTask")
	t.tasks = append(t.tasks, task)
}

func (t *TimeUnit) nextTask() (task types.Task) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.tasks) == 0 {
		return nil
	}
	task = t.tasks[0]
	t.tasks = t.tasks[1:]
	return task
}
