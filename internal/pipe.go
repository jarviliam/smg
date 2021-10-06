package internal

import (
	"sync"

	"github.com/jarviliam/smg/internal/target"
)

type Pipe interface {
	Pipe(*sync.WaitGroup, <-chan *target.Target, chan<- *target.Target)
}
