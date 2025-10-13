package util

import "sync"

type Semaphore struct {
	maxConcurrency int
	currentCount   int
	mutex          sync.Mutex
	cond           *sync.Cond
}

func NewSemaphore(max int) *Semaphore {
	sem := &Semaphore{maxConcurrency: max}
	sem.cond = sync.NewCond(&sem.mutex)

	return sem
}

func (s *Semaphore) Acquire() {
	s.mutex.Lock() // protect currentCount
	for s.currentCount >= s.maxConcurrency {
		s.cond.Wait() // unlocks mutex before blocking, re-locks before resuming
	}
	s.currentCount++
	s.mutex.Unlock()
}

func (s *Semaphore) TryAcquire() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.currentCount >= s.maxConcurrency {
		return false
	}

	s.currentCount++
	return true
}

func (s *Semaphore) Release() {
	s.mutex.Lock()
	s.currentCount--
	s.mutex.Unlock()
	s.cond.Signal() // wake one waiting goroutine
}
