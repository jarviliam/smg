package requester

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/jarviliam/smg/internal/target"
)

var (
	ErrRequestFailed = errors.New("err request failed")
)

type Requester struct {
	http.Client
}

func NewRequester(opts ...func(*Requester)) *Requester {
	transport := &http.Transport{
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 1,
		MaxConnsPerHost:     5,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	requester := &Requester{Client: http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}}
	for _, opt := range opts {
		opt(requester)
	}
	return requester
}

func (w *Requester) Fetch(t *target.Target) error {
	start := time.Now()
	request, err := http.NewRequest("GET", t.BaseURL, nil)
	if err != nil {
		return ErrRequestFailed
	}
	request.Header.Set("User-Agent", "Go/CCC-Sitemap")
	resp, err := w.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer func(url string) { fmt.Printf("Url: %s ; Time: %.2f\n", url, time.Since(start).Seconds()) }(t.BaseURL)
	t.Content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (w *Requester) Pipe(wg *sync.WaitGroup, in <-chan *target.Target, out chan<- *target.Target) {
	defer close(out)
	//Limit Requests with a sephamore to prevent timeout
	s := make(chan interface{}, 10)
	defer close(s)

	for t := range in {
		wg.Add(1)
		s <- struct{}{}
		go func(tgt *target.Target) {
			//TODO: Logger
			if err := w.Fetch(tgt); err != nil {
				log.Printf("Requester: %s\n", err.Error())
				//Remove from queue
				wg.Done()
			} else {
				out <- tgt
			}
			wg.Done()
			<-s
		}(t)
	}
}
