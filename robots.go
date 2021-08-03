package smg

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

type Robots struct {
	IgnoreRobots bool
	HasFetched   bool
	AllowAll     bool
	DisallowAll  bool
}
type RobotGroup struct {
	Name  string
	Rules []*Rule
}
type Rule struct {
	isAllowed bool
	path      string
	reg       *regexp.Regexp
}

func NewRobots(ignore bool) *Robots {
	r := &Robots{
		IgnoreRobots: ignore,
		HasFetched:   false,
	}
	return r
}

func (r *Robots) ParseFromResponse(res *http.Response) {
	if res == nil {

	}
	buff, err := ioutil.ReadAll(res.Body)
	if err != nil {

	}

	statusCode := res.StatusCode
	switch {
	case statusCode >= 200 && statusCode < 300:
		parseRobots(buff)
		break
	case statusCode >= 400 && statusCode < 500:
		r.AllowAll = true
		r.DisallowAll = false
		break

	case statusCode >= 500 && statusCode < 600:
		r.DisallowAll = true
		r.AllowAll = false
	}
}

//TODO: Parse Robots.txt -> Tokens
//TODO: Tokens -> UserAgentGroups -> Rulesc:W

func parseRobots(body []byte) {
	lines, err := intoLines(body)
}

func intoLines(body []byte) ([]string, error) {
	var lines []string
	newLine := "/\r\n|\r|\n/"
	for {
          tk := body.
	}
	return lines, nil
}