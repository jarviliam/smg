package prioritizer

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/jarviliam/smg/internal"
	"github.com/jarviliam/smg/internal/target"
)

type Prioritizer struct {
	fns []PriorityFN
	def float32
}

var (
	DEFAULT_PRIORITY = 1.0
)

type PriorityFN func(doc *goquery.Document) float32

var _ internal.Pipe = (*Prioritizer)(nil)

func NewPrioritizer(pf ...PriorityFN) *Prioritizer {
	p := &Prioritizer{
		def: float32(DEFAULT_PRIORITY),
		fns: make([]PriorityFN, 0),
	}
	for _, f := range pf {
		p.fns = append(p.fns, f)
	}
	return p
}

func (p *Prioritizer) SetDefault(value float32) {
	p.def = value
}

func (p *Prioritizer) AddFN(check PriorityFN) {
	p.fns = append(p.fns, check)
}
func (p *Prioritizer) SetPriority(t *target.Target) {
	priority := p.def
	if len(t.Content) > 0 {
		doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(t.Content))
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, v := range p.fns {
			priority = v(doc)
		}
	}
	t.Priority = priority
}

func (p *Prioritizer) Pipe(wg *sync.WaitGroup, in <-chan *target.Target, out chan<- *target.Target) {
	defer close(out)
	for t := range in {
		wg.Add(1)
		go func(t *target.Target) {
			p.SetPriority(t)
			wg.Done()
			out <- t
		}(t)
	}
}
