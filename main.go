package main

import (
	"context"
	// poller01 "go_work_scheduler/origin"
	poller02 "go_work_scheduler/scheduler"
)

func main() {
	// producer01 := poller01.NewPoller(5)
	// producer01.Poll(context.Background())
	producer02 := poller02.NewPoller(5)
	producer02.Poll(context.Background())
}
