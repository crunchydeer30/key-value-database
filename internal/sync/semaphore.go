package sync

import "sync"

type Semaphore struct {
	mu    sync.Mutex
	cond  *sync.Cond
	count int
	max   int
}

func NewSemaphore(max int) *Semaphore {
	//nolint:exhaustruct
	s := &Semaphore{
		count: 0,
		max:   max,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for s.count >= s.max {
		s.cond.Wait()
	}

	s.count++
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.count--
	s.cond.Signal()
}
