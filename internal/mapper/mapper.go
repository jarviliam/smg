package mapper

import (
	"fmt"
	"os"
	"sync"

	"github.com/jarviliam/smg/internal"
	"github.com/jarviliam/smg/internal/target"
)

var _ internal.Pipe = (*Mapper)(nil)

type Mapper struct {
	MaxUrls    int
	mu         *sync.RWMutex
	list       []string
	path       string
	currentUrl int
	curr       *os.File
}

//TODO: Options
//Makes a New Mapper
func NewMapper() *Mapper {
	m := &Mapper{MaxUrls: 50000,
		list: make([]string, 0), mu: &sync.RWMutex{}, path: "./"}
	m.appendSM()
	return m
}

func (m *Mapper) appendSM() {
	m.list = append(m.list, fmt.Sprintf("sitemap%d.xml", len(m.list)+1))
}

func (m *Mapper) Write(entry *target.Target) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentUrl++
	if m.currentUrl >= m.MaxUrls {
		m.appendSM()
		m.CreateFile()
		m.currentUrl = 0
	}
	out := fmt.Sprintf("<url><loc>%s</loc><priority>%.1f</priority></url>\n", entry.BaseURL, entry.Priority)
	_, err := m.curr.WriteString(out)
	if err != nil {
		fmt.Printf("Error %v", err)
	}
}

//Creates a new SM file. Rotates Writing
func (m *Mapper) CreateFile() {
	file, err := os.Create(fmt.Sprintf("%s%s", m.path, m.list[len(m.list)-1]))
	if err != nil {
		panic(err)
	}
	if m.curr != nil {
		m.curr.Close()
	}
	m.curr = file
}

func (m *Mapper) Pipe(wg *sync.WaitGroup, in <-chan *target.Target, out chan<- *target.Target) {
	defer close(out)
	m.CreateFile()
	for t := range in {
		//If It has been fetched (has a priority) then write
		if len(t.Content) != 0 {
			m.Write(t)
		}
		out <- t
	}
}
