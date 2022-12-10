package origin

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

type Poller struct {
	routineGroup *routineGroup
	workerNum    int
}

// generate corresponding goroutine
func NewPoller(workerNum int) *Poller {
	return &Poller{
		routineGroup: newRoutineGroup(),
		workerNum:    workerNum,
	}
}

func (p *Poller) Poll(ctx context.Context) error {
	for i := 0; i < p.workerNum; i++ {
		// step01
		p.routineGroup.Run(func() {
			for {
				select {
				// job is done
				case <-ctx.Done():
					return
				default:
					// step02
					task, err := p.fetch(ctx)
					if err != nil {
						log.Println("can't get task", err.Error())
						break
					}
					// step 03
					if err := p.execute(ctx, task); err != nil {
						log.Println("execute task error:", err.Error())
					}
				}
			}
		})
	}
	p.routineGroup.Wait()
	return nil
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
	log.Printf("task %s executed", task)
	return nil
}
