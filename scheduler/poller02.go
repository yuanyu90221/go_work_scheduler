package scheduler

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Poller struct {
	routineGroup *routineGroup
	workerNum    int

	sync.Mutex
	ready  chan struct{}
	metric *metric
}

func NewPoller(workerNum int) *Poller {
	return &Poller{
		routineGroup: newRoutineGroup(),
		workerNum:    workerNum,
		ready:        make(chan struct{}, 1),
		metric:       newMetric(),
	}
}

func (p *Poller) schedule() {
	p.Lock()
	defer p.Unlock()
	if int(p.metric.BusyWorkers()) >= p.workerNum {
		return
	}
	select {
	case p.ready <- struct{}{}:
	default:
	}
}

func (p *Poller) Poll(ctx context.Context) error {
	// scheduler
	for {
		// step 01
		p.schedule()
		select {
		case <-p.ready:
		case <-ctx.Done():
			return nil
		}
	LOOP:
		for {
			select {
			// when finish
			case <-ctx.Done():
				break
			default:
				// step 02
				task, err := p.fetch(ctx)
				if err != nil {
					log.Println("fetch task error:", err.Error())
					break
				}
				p.metric.IncBusyWorker()
				// step 03
				p.routineGroup.Run(func() {
					if err := p.execute(ctx, task); err != nil {
						log.Println("execute task error:", err.Error())
					}
				})
				break LOOP
			}
		}
	}
}

func (p *Poller) fetch(ctx context.Context) (string, error) {
	// connect database or other service
	time.Sleep(1000 * time.Millisecond)
	result := rand.Intn(2)
	taskId := time.Now()
	if result%2 == 1 {
		return taskId.String(), nil
	}
	return "foobar", errors.New("no task")
}

func (p *Poller) execute(ctx context.Context, task string) error {
	defer func() {
		p.metric.DecBusyWorker()
		p.schedule()
	}()
	log.Printf("task %s executed", task)
	return nil
}
