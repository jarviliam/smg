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
	Spider *Spider
}

//TODO: Figure Out Why Internal is not Linking
type Spider struct {
	client   *http.Client
	visted   []string
	MaxDepth int
	Robots   Robots
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
func (s *Spider) Fetch(u string) error {
	depth := 1
	parsedUrl, err := url.Parse(u)
	s.Robots.isAllowed(parsedUrl)
	if err != nil {
		return errors.New("Url is empty")
	}
	if s.MaxDepth > 0 && s.MaxDepth < depth {
		return nil
	}
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
		fmt.Errorf(err.Error())
		return
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			for _, x := range n.Attr {
				if x.Key == "href" {
					fmt.Printf("Link Found : %s \n", x.Val)
				}
			}
		}
	})
}

func isValidUrl(u url.URL) bool {
	return u.Host != ""
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
