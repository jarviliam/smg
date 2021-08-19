package spider

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/jarviliam/smg/internal/context"
	"github.com/jarviliam/smg/internal/logger"
	"github.com/jarviliam/smg/internal/robots"
)

var (
	ErrNoUrl    = errors.New("No Url Recieved")
	ErrMaxDepth = errors.New("Max Depth")
	ErrBlocked  = errors.New("Regex Blocked")
)

type Spider struct {
	client     *http.Client
	Logger     *logger.Logger
	visited    []LinkEntry
	MaxDepth   int
	TargetHost string
	Robots     *robots.Robots
}

type LinkEntry struct {
	url   string
	depth int
}

type LinkElement struct {
	Attributes map[string]string
}

func NewSpider() *Spider {
	jar, _ := cookiejar.New(nil)
	s := &Spider{
		visited: make([]LinkEntry, 10),
		client: &http.Client{
			Jar: jar,
		},
		Robots: robots.NewRobots(true),
	}
	return s
}

func (s *Spider) Fetch(u string, depth int) error {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return errors.New("Url is empty")
	}
	fmt.Printf("To Visit Url %s \n", u)
	err = s.isValidUrl(u, parsedUrl, depth)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if s.hasVisited(parsedUrl.Path) {
		s.replaceIfNeeded(parsedUrl.Path, depth)
		fmt.Printf("Has Visited Url %s \n", parsedUrl.Path)
		return nil
	}
	if parsedUrl.Host != s.TargetHost && parsedUrl.Host != "" {
		fmt.Printf("URL HOST : %s - TARGET : %s \n", parsedUrl.Host, s.TargetHost)
		return nil
	}
	//fmt.Printf("URL HOST : %s - TARGET : %s \n", parsedUrl.Host, s.TargetHost)
	var e LinkEntry
	if parsedUrl.IsAbs() {
		e := &LinkEntry{
			url:   parsedUrl.Path,
			depth: depth,
		}
	} else {
		//TODO: Absolute the URL
	}
	s.visited = append(s.visited, *e)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.Logger.LogInfo("Visting : %s At Depth: %d \n", u, depth)
	ctx := context.NewContext()
	//Request Struct
	request := &Request{
		URL:     req.URL,
		Headers: &req.Header,
		Depth:   depth,
		Ctx:     ctx,
		spider:  s,
	}
	res, err := s.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil
	}
	response := &Response{
		Code: res.StatusCode,
		Ctx:  ctx,
		Body: body,
	}
	s.parse(request, response)
	return nil
}

func absoluteUrl() {

}

func (s *Spider) parse(req *Request, res *Response) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(res.Body))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		for _, n := range sel.Nodes {
			for _, x := range n.Attr {
				if x.Key == "href" {
					fmt.Printf("Link Found : %s \n", x.Val)
					s.Fetch(x.Val, req.Depth+1)
				}
			}
		}
	})
}

func (s *Spider) isValidUrl(plainUrl string, parsedUrl *url.URL, depth int) error {
	if plainUrl == "" {
		return ErrNoUrl
	}
	if s.MaxDepth > 0 && s.MaxDepth < depth {
		return ErrMaxDepth
	}
	if parsedUrl.Hostname() != s.TargetHost {
		//TODO: Make SMG struct -> Run
		return nil
		//return ErrBlocked
	}
	return nil
}
func (s *Spider) hasVisited(u string) bool {
	visited := false
	for _, v := range s.visited {
		if v.url == u {
			visited = true
			break
		}
	}
	return visited
}

func (s *Spider) replaceIfNeeded(u string, d int) {
	for _, v := range s.visited {
		if v.url == u && v.depth > d {
			v.depth = d
			break
		}
	}
}

func (s *Spider) CheckRobots(u *url.URL) bool {
	allowed := s.Robots.IsAllowed(u.Path)
	return allowed
}
