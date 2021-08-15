package smg

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type SMG struct {
	TargetHost string
	Spider     *Spider
	Options    interface{}
}

func NewSMG() *SMG {
	s := &SMG{}
	return s
}

//TODO -> Make this into Starting Point
func (s *SMG) Run(target string) error {
	if target == "" {
		return ErrNoUrl
	}
	parsedUrl, err := url.Parse(target)
	if err != nil {
		return err
	}
	//TODO Refactor
	s.TargetHost = parsedUrl.Host
	s.Spider.TargetHost = parsedUrl.Host
	s.Spider.Fetch(target, 0)
	return nil
}

//TODO: Figure Out Why Internal is not Linking
type Spider struct {
	client     *http.Client
	visted     []string
	MaxDepth   int
	TargetHost string
	//Robots   Robots
	//robots   map[string]string
}

type LinkElement struct {
	Attributes map[string]string
}

func NewSpider() *Spider {
	jar, _ := cookiejar.New(nil)
	s := &Spider{
		visted: make([]string, 10),
		client: &http.Client{
			Jar: jar,
		},
	}
	return s
}

func (s *Spider) Fetch(u string, depth int) error {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return errors.New("Url is empty")
	}
	err = s.isValidUrl(u, parsedUrl, depth)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	//	s.Robots.isAllowed(parsedUrl)
	if s.hasVisited(u) {
		return nil
	}
	s.visted = append(s.visted, parsedUrl.Path)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ctx := NewContext()
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

var (
	ErrNoUrl    = errors.New("No Url Recieved")
	ErrMaxDepth = errors.New("Max Depth")
	ErrBlocked  = errors.New("Regex Blocked")
)

//TODO: Robots Check
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
	for _, v := range s.visted {
		if v == u {
			visited = true
			break
		}
	}
	return visited
}

func (s *Spider) GetRobots(u *url.URL) {
	// robot, ok := s.robots[u.Host]
	// if !ok {

	// }
}
