package extractor

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/jarviliam/smg/internal"
	"github.com/jarviliam/smg/internal/target"
	"github.com/temoto/robotstxt"
	"golang.org/x/net/html"
)

var (
	ErrInvalidLink      = errors.New("invalid link")
	ErrInvalidHost      = errors.New("invalid host")
	ErrNotSameHost      = errors.New("not same host")
	ErrNoFollow         = errors.New("no follow link")
	ErrRobotsDisallowed = errors.New("blocked by robots")
)
var _ internal.Pipe = (*Extractor)(nil)

type Extractor struct {
	checkers   []func()
	host       string
	robots     *robotstxt.RobotsData
	ignoreLink []*regexp.Regexp
	ua         string
}

func (e *Extractor) SetUA(ua string) {
	e.ua = ua
}

func getHost(link string) (string, error) {
	urlStruct, err := url.Parse(link)
	if err != nil {
		return "", ErrInvalidHost
	} else if len(urlStruct.Host) == 0 {
		return "", ErrInvalidHost
	}
	return urlStruct.Host, nil
}

func NewExtractor(originURL string, regex ...*regexp.Regexp) (*Extractor, error) {
	host, err := getHost(originURL)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(originURL + "/robots.txt")
	if err != nil {
		return nil, err
	}
	data, err := robotstxt.FromResponse(resp)
	resp.Body.Close()
	e := &Extractor{checkers: make([]func(), 0), host: host, ignoreLink: make([]*regexp.Regexp, 0), robots: data}
	e.ignoreLink = append(e.ignoreLink, regex...)
	return e, nil
}
func (e *Extractor) isSameHost(link string) bool {
	host, err := getHost(link)
	if err != nil {
		return false
	}

	if host != e.host {
		return false
	}
	return true
}

//Gets Link from an "a" node
func (e *Extractor) getLink(n *html.Node, baseURL string) (string, error) {
	var out string
	for _, x := range n.Attr {
		if x.Key == "href" {
			//Clean Link to Absolute URL
			link, err := cleanLink(baseURL, x.Val)
			if err != nil {
				return "", err
			}
			//Ignore if Host is not same
			if !e.isSameHost(link) {
				return "", ErrNotSameHost
			}
			out = link
		} else if x.Key == "rel" && x.Val == "nofollow" {
			return "", ErrNoFollow
		}
	}
	if out == "" {
		return "", ErrInvalidLink
	}
	return out, nil
}

func (e *Extractor) validRegexLink(link string) bool {
	for _, reg := range e.ignoreLink {
		if reg.MatchString(link) {
			return false
		}
	}
	return true
}

func (e *Extractor) checkRobots(link string) error {
	u, err := url.Parse(link)
	if err != nil {
		return err
	}
	if !e.robots.TestAgent(u.EscapedPath(), "*") {
		return ErrRobotsDisallowed
	}
	return nil
}

func (e *Extractor) ExtractLinks(baseURL string, content []byte) []string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(content))
	if err != nil {
		fmt.Println(err.Error())
		return []string{}
	}
	//Meta No Robots
	if !isMetaRobotsFriendly(doc) {
		return []string{}
	}
	unique := make(map[string]bool)
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		for _, n := range sel.Nodes {
			link, err := e.getLink(n, baseURL)
			if err != nil {
				//TODO
				//LOG ERROR
				//fmt.Println("Skip1")
				continue
			}
			if err = e.checkRobots(link); err != nil {
				//Log Here
				//fmt.Println("Skip")
				continue
			}
			if !e.validRegexLink(link) {
				//fmt.Println("Skip2")
				continue
			}
			unique[link] = true
		}
	})
	//Make Unique URLs
	linksOut := make([]string, 0)
	for k := range unique {
		linksOut = append(linksOut, k)
	}
	return linksOut
}

//Extractor Provides new Targets to the pipeline.
//When the target in the IN pipe gets extracted, it gets discarded as it has already been mapped
func (e *Extractor) Pipe(wg *sync.WaitGroup, in <-chan *target.Target, out chan<- *target.Target) {
	defer close(out)
	for tget := range in {
		//Add For Go  Func
		wg.Add(1)
		go func(tgt *target.Target) {
			links := e.ExtractLinks(tgt.BaseURL, tgt.Content)
			for _, link := range links {
				wg.Add(1)
				go func(l string) { out <- target.NewTarget(l) }(link)
			}
			//Does not pass on origin (Already Recorded and Written)
			wg.Done()
			//Remove for goroutine
			wg.Done()
		}(tget)
	}
}
