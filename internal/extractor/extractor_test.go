package extractor

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExtractor(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}

func TestExtractLinks(t *testing.T) {
	testCases := []struct {
		name            string
		mockBaseURL     string
		mockContentPath string
		expectedLinks   []string
	}{
		{
			"nofollow",
			"http://www.basic-one.com",
			"testdata/nofollow.html",
			[]string{"http://www.basic-one.com/test3"},
		},
		{
			"links",
			"http://www.basic-one.com",
			"testdata/links.html",
			[]string{"http://www.basic-one.com/link", "http://www.basic-one.com/link/test", "http://www.basic-one.com/test"},
		},
		{
			"nolinks",
			"http://www.basic-one.com",
			"testdata/nolinks.html",
			[]string{},
		},
		{
			"metaNoFollow",
			"http://www.basic-one.com",
			"testdata/metarobots.html",
			[]string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := ioutil.ReadFile(tc.mockContentPath)
			if err != nil {
				t.Errorf("%s: %v", tc.name, err)
			}

			e, err := NewExtractor(tc.mockBaseURL)
			if err != nil {
				t.Error(err)
			}
			res := e.ExtractLinks(tc.mockBaseURL, content)
			assert.ElementsMatch(t, tc.expectedLinks, res)
		})
	}
}
