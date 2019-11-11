package main

import "sync"

type Timeline struct {
	t    int64
	line []TimeUnit

	mu sync.RWMutex
}

func (tl *Timeline) ImportTask(t int64, task Task) {

}

type TimeUnit struct {
	tasks []Task
}
