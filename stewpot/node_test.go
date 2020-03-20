package stewpot

import "testing"

func TestCacheDB_Insert(t *testing.T) {
	db := newNodeCacheDB()

	db.Insert(2020)
	if db.coldData != 1000 {
		t.Fatalf("coldData not correct, expected: %d, gdt: %d", 1000, db.coldData)
	}
	if !db.Exist(2020) {
		t.Fatal("hotData not correct, 1020 should exists")
	}
	if db.Exist(2000) {
		t.Fatal("2000 should not exists")
	}
}
