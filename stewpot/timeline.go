package stewpot

import (
	"fmt"
	"github.com/heeeeeng/node_stewpot/tasks"
	"github.com/heeeeeng/node_stewpot/types"
	"sync"
)

type Timeline struct {
	timestamp int64

	current  *TimeUnit
	next     *TimeUnit
	nextChan chan *TimeUnit

	db *MemDB

	protocol       types.Protocol
	importCallback func(types.Task)

	mu    sync.RWMutex
	close chan struct{}
}

func newTimeline(db *MemDB, protocol types.Protocol, callback func(task types.Task)) *Timeline {
	t := &Timeline{}

	t.timestamp = 0
	t.current = newTimeUnit(t.timestamp)
	t.next = nil
	//t.next = newTimeUnit(t.timestamp + 1)
	t.nextChan = make(chan *TimeUnit)
	t.db = db
	t.protocol = protocol
	t.importCallback = callback

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

func (tl *Timeline) Protocol() types.Protocol { return tl.protocol }

func (tl *Timeline) SendNewMsg(node types.Node, msg types.Message) int64 {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	task := tasks.NewMsgProcessTask(tl.NextTime(), node, msg)
	tl.ImportTask(task.StartTime(), task)
	return task.StartTime()
}

func (tl *Timeline) ImportTask(startTime int64, task types.Task) {
	if tl.importCallback != nil {
		go tl.importCallback(task)
	}

	if startTime != tl.timestamp && startTime != tl.timestamp+1 {
		fmt.Println(fmt.Sprintf("appendTask not curr or next, task start time: %d, curr: %d", task.StartTime(), tl.current.timestamp))
		return
	}
	if startTime == tl.timestamp {
		tl.current.appendTask(task)
		return
	}
	if startTime == tl.timestamp+1 {
		if tl.next == nil {
			tl.next = newTimeUnit(tl.timestamp + 1)
			tl.next.appendTask(task)
			go func() { tl.nextChan <- tl.next }()
		} else {
			tl.next.appendTask(task)
		}
		return
	}
	return
}

func (tl *Timeline) GetTimeUnit(t int64) *TimeUnit {
	return tl.db.GetTimeUnit(t)
}

func (tl *Timeline) loop() {
	for {
		select {
		case <-tl.close:
			return

		case tl.current = <-tl.nextChan:
			if tl.current == nil {
				continue
			}
			tl.next = nil

			//fmt.Println("start processing timestamp: ", tl.current.timestamp)

			for {
				task := tl.current.nextTask()
				if task == nil {
					break
				}
				task.Process(tl)
			}

			tl.mu.Lock()
			tl.db.InsertTimeUnit(tl.current)
			tl.timestamp++

			tl.mu.Unlock()
		}
	}
}

type TimeUnit struct {
	index     int
	timestamp int64
	tasks     []types.Task

	mu sync.RWMutex
}

func newTimeUnit(timestamp int64) *TimeUnit {
	return &TimeUnit{
		index:     0,
		timestamp: timestamp,
		tasks:     make([]types.Task, 0),
	}
}

func (t *TimeUnit) appendTask(task types.Task) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tasks = append(t.tasks, task)
}

func (t *TimeUnit) nextTask() (task types.Task) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.tasks) == 0 {
		return nil
	}
	if t.index >= len(t.tasks) {
		return nil
	}
	task = t.tasks[t.index]
	t.index++
	return task
}
