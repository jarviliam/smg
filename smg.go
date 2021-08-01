package smg

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Spider struct {
	client   *http.Client
	visted   []string
	MaxDepth int
}

type Response struct {
	Code    int
	Body    []byte
	Ctx     *Context
	Request *http.Request
}
type Request struct {
	URL     *url.URL
	Headers *http.Header
	Ctx     *Context
	spider  *Spider
}
type Context struct {
	contextMap map[string]string
	lock       *sync.Mutex
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
func NewContext() *Context {
	c := &Context{
		contextMap: make(map[string]string),
		lock:       &sync.Mutex{},
	}
	return c
}

func (s *Spider) Fetch(url string) error {
	depth := 1
	if !isValidUrl(url) {
		return errors.New("Url is empty")
	}
	if s.MaxDepth > 0 && s.MaxDepth < depth {
		return nil
	}
	if s.hasVisited(url) {
		return nil
	}
	s.visted = append(s.visted, url)
	req, err := http.NewRequest("GET", url, nil)
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
		fmt.Println(err)
		return
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			fmt.Println("NODE -----")
			for _, x := range n.Attr {
				if x.Key == "href" {
					fmt.Println(x.Val)
				}
			}
			fmt.Println("-----")
		}
	})
}

func isValidUrl(u string) bool {
	return u != ""
}
func (s *Spider) hasVisited(u string) bool {
        visited := false
  for _ ,v := s.visted {
    if v == u {
      visited = true
      break
    }
  }
        return visited
}

func (c *Context) Get(key string) string {
	if v, ok := c.contextMap[key]; ok {
		return v
	}
	return ""
}

func (c *Context) Put(key, value string, ovrride bool) error {
	if _, ok := c.contextMap[key]; ok && ovrride {
		c.contextMap[key] = value
	} else if ok && !ovrride {
		return errors.New("Key already exists")
	} else {
		c.contextMap[key] = value
	}
	return nil
}
