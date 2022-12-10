package origin

import "sync"

type routineGroup struct {
	waitGroup sync.WaitGroup
}

func newRoutineGroup() *routineGroup {
	return new(routineGroup)
}

// launch a goroutine with WaitGroup
func (g *routineGroup) Run(fn func()) {
	g.waitGroup.Add(1)
	go func() {
		defer g.waitGroup.Done()
		fn()
	}()
}

// Wait for the WaitGroup is 0
func (g *routineGroup) Wait() {
	g.waitGroup.Wait()
}
