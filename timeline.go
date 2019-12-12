package main

type Timeline struct {
	timestamp int64

	current *TimeUnit
	next    *TimeUnit
	db      *MemDB

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

func (tl *Timeline) ImportTask(startTime int64, task Task) {
	if startTime < tl.timestamp || startTime > tl.timestamp+1 {
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

func (tl *Timeline) getTimeUnit(t int64) *TimeUnit {
	if t == tl.current.timestamp {
		return tl.current
	} else if t == tl.next.timestamp {
		return tl.next
	} else {
		return tl.db.Get(t)
	}
}

//func (tl *Timeline) setTimeUnit(tu *TimeUnit) {
//	tl.line[tu.timestamp] = tu
//}

//func (tl *Timeline) initTimeUnit(t int64) *TimeUnit {
//	timeUnit := newTimeUnit(t)
//	tl.line[t] = timeUnit
//	return timeUnit
//}

func (tl *Timeline) loop() {
	for {
		select {
		case <-tl.close:
			return

		default:
			if tl.current == nil {
				continue
			}

			for _, task := range tl.current.tasks {
				task.Process(tl)
			}
			tl.current = tl.next
			tl.next = newTimeUnit(tl.timestamp)
			tl.timestamp = tl.NextTime()
		}
	}
}

type TimeUnit struct {
	timestamp int64
	tasks     []Task
}

func newTimeUnit(timestamp int64) *TimeUnit {
	return &TimeUnit{
		timestamp: timestamp,
		tasks:     make([]Task, 0),
	}
}

func (t *TimeUnit) appendTask(task Task) {
	t.tasks = append(t.tasks, task)
}
