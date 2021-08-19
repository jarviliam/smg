package smg

import (
	"errors"
	"net/url"

	"github.com/jarviliam/smg/pkg/spider"
	"github.com/jarviliam/smg/pkg/writer"
)

var (
	ErrNoUrl = errors.New("No Url Recieved")
)

type SMG struct {
	TargetHost string
	Spider     *spider.Spider
	Writer     *writer.Writer
	Options    *SMGOptions
}
type SMGOptions struct {
	StickToHost bool
}

func NewSMG() *SMG {
	s := &SMG{
		Spider: spider.NewSpider(),
		Options: &SMGOptions{
			StickToHost: true,
		},
	}
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
