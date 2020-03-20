package stewpot

type Recorder struct {
	db *MemDB
}

func newRecorder(db *MemDB) Recorder {
	return Recorder{
		db: db,
	}
}

func (r *Recorder) Record(pattern string, data ...interface{}) {
	//pattern = "<time: %f> " + pattern
	//duration := time.Since(r.StartTime) * time.Duration(r.SpeedRate)
	//
	//data = append([]interface{}{duration.Seconds()}, data...)
	//logrus.Infof(pattern, data...)
}
