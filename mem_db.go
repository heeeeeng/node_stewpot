package main

type Record struct {
	TaskType  int   `json:"task_type"`
	Timestamp int64 `json:"timestamp"`
}

type MemDB struct {
	latestTime int64
	timeline   map[int64]*TimeUnit
}

func newMemDB() *MemDB {
	db := &MemDB{}
	db.timeline = make(map[int64]*TimeUnit)

	return db
}

func (db *MemDB) Get(t int64) *TimeUnit {
	return db.timeline[t]
}

func (db *MemDB) Insert(t *TimeUnit) {
	db.timeline[t.timestamp] = t
}