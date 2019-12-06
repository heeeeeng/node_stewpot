package main

import (
	"sync"
)

type Timeline struct {
	start int64
	end   int64
	line  map[int64]*TimeUnit

	mu sync.RWMutex
}

func (tl *Timeline) ImportTask(startTime int64, task Task) {
	if startTime < tl.start {
		// TODO
		// log out this situation, this should not happen.
		return
	}

	timeUnit := tl.getTimeUnit(startTime)
	if timeUnit == nil {
		timeUnit = tl.initTimeUnit(startTime)
	}
	timeUnit.appendTask(task)
	tl.setTimeUnit(timeUnit)

	if startTime > tl.end {
		tl.end = startTime
	}
}

func (tl *Timeline) getTimeUnit(t int64) *TimeUnit {
	return tl.line[t]
}

func (tl *Timeline) setTimeUnit(tu *TimeUnit) {
	tl.line[tu.timestamp] = tu
}

func (tl *Timeline) initTimeUnit(t int64) *TimeUnit {
	timeUnit := newTimeUnit(t)
	tl.line[t] = timeUnit
	return timeUnit
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
