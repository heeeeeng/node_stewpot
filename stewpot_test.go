package main

import (
	"fmt"
	"github.com/heeeeeng/node_stewpot/types"
	"math/rand"
	"testing"
	"time"
)

func TestStewpot_MultiSimulate(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	stewpot := NewStewpot()
	stewpot.InitNetwork(200, 8, 4, 3, 100*types.MB)
	stewpot.Start()

	conf := SimConfig{
		IterNum:   20,
		MsgSize:   256 * types.KB,
		NodeNum:   100,
		Bandwidth: 100 * types.MB,
		MaxIn:     8,
		MaxOut:    4,
	}
	avgTime := stewpot.MultiSimulate(conf)
	fmt.Println("average time usage: ", avgTime)
}
