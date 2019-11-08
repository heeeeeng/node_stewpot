package main

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Dumper struct {
	StartTime time.Time
	SpeedRate int
}

func newDumper(startTime time.Time, rate int) Dumper {
	return Dumper{
		StartTime: startTime,
		SpeedRate: rate,
	}
}

func (d *Dumper) Log(pattern string, data ...interface{}) {
	pattern = "<time: %f> " + pattern
	duration := time.Since(d.StartTime) * time.Duration(d.SpeedRate)

	data = append([]interface{}{duration.Seconds()}, data...)
	logrus.Infof(pattern, data...)
}
