package spider

import (
	"fmt"
	"sync"
	"time"

	"github.com/jarviliam/smg/internal"
	"github.com/jarviliam/smg/internal/target"
)

type Spider struct {
	urls chan *target.Target
}

func NewSpider() *Spider {
	return &Spider{
		urls: make(chan *target.Target, 10),
	}
}

//Cycle the entries
func (c *Spider) pipeEnd(in <-chan *target.Target) {
	for t := range in {
		c.urls <- t
	}
}

func (c *Spider) Run(entry *target.Target, pipes ...internal.Pipe) error {
	start := time.Now()
	wg := sync.WaitGroup{}
	in := c.urls
	for _, pipe := range pipes {
		out := make(chan *target.Target)
		go pipe.Pipe(&wg, in, out)
		in = out
	}
	go c.pipeEnd(in)
	wg.Add(1)
	c.urls <- entry
	wg.Wait()
	fmt.Printf("Finished SMG. Elapsed: %0.2f", time.Since(start).Seconds())
	close(c.urls)

	return nil
}
