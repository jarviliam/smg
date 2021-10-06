package extractor

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	RegexAnchors = regexp.MustCompile(`(#.*)$`)
)

func isMetaRobotsFriendly(page *goquery.Document) bool {
	friendly := true
	page.Find("meta[name='robots']").Each(func(i int, sel *goquery.Selection) {
		for _, n := range sel.Nodes {
			for _, attr := range n.Attr {
				if attr.Val == "nofollow" {
					friendly = false
					break
				}
			}
			if !friendly {
				break
			}
		}
	})

	return friendly
}

func validScheme(in string) bool {
	return in == "https"
}

//Helper function that returns an absolute link
func cleanLink(base, link string) (string, error) {
	link = RegexAnchors.ReplaceAllString(link, "")
	if len(link) == 0 {
		return "", ErrInvalidLink
	}
	linkURL, err := url.Parse(link)
	if err != nil {
		return "", ErrInvalidLink
	}

	if validScheme(linkURL.Scheme) {
		return link, nil
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return "", ErrInvalidLink
	} else if len(baseURL.Host) == 0 {
		return "", ErrInvalidLink
	}

	if link[0] == '/' || link[len(link)-1] == '/' {
		link = strings.Trim(link, "/")
	}

	return strings.Join([]string{baseURL.Scheme, "://", baseURL.Host, "/", link}, ""), nil
}
