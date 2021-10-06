package storage

import (
	"sync"

	"github.com/jarviliam/smg/internal"
	"github.com/jarviliam/smg/internal/target"
)

var _ internal.Pipe = (*Storage)(nil)

type Storage struct {
	hasVisited map[string]bool
	mu         *sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		mu:         &sync.Mutex{},
		hasVisited: make(map[string]bool),
	}
}

//Checks to see if url has been stored.
func (s *Storage) hasSeen(url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hasVisited[url]; ok {
		return true
	}
	s.hasVisited[url] = true
	return false
}

func (e *Storage) Pipe(wg *sync.WaitGroup, in <-chan *target.Target, out chan<- *target.Target) {
	defer close(out)
	for t := range in {
		//Remove from pipeline if it has been checked
		if e.hasSeen(t.BaseURL) {
			wg.Done()
		} else {
			out <- t
		}
	}
}
