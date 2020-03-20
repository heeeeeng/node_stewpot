package stewpot

import (
	"github.com/heeeeeng/node_stewpot/types"
	"sync"
)

type Record struct {
	TaskType  int   `json:"task_type"`
	Timestamp int64 `json:"timestamp"`
}

type MemDB struct {
	latestTime int64
	timeline   map[int64]*TimeUnit

	messages map[int64]*types.Message

	mu sync.RWMutex
}

func newMemDB() *MemDB {
	db := &MemDB{}
	db.timeline = make(map[int64]*TimeUnit)

	return db
}

func (db *MemDB) GetTimeUnit(t int64) *TimeUnit {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.timeline[t]
}

func (db *MemDB) InsertTimeUnit(t *TimeUnit) {
	if t == nil {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	//fmt.Println(fmt.Sprintf("insert time unit, time: %d, tasks len: %d", t.timestamp, len(t.tasks)))
	db.timeline[t.timestamp] = t
}

func (db *MemDB) GetMessage(msgID int64) *types.Message {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.messages[msgID]
}

func (db *MemDB) InsertMessage(msg *types.Message) {
	if msg == nil {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.messages[msg.ID] = msg
}
