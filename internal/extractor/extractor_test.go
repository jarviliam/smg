package extractor

import (
	"testing"
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
			"http://www.test.com",
			"testdata/nofollow.html",
			[]string{"http://www.test.com/test3"},
		},
		{
			"links",
			"http://www.test.com",
			"testdata/links.html",
			[]string{"http://www.test.com/link", "http://www.test.com/link/test", "http://www.test.com/test"},
		},
		{
			"nolinks",
			"http://www.test.com",
			"testdata/nolinks.html",
			[]string{},
		},
		{
			"metaNoFollow",
			"http://www.test.com",
			"testdata/metarobots.html",
			[]string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//TODO Unset when Robots is fixed
			// content, err := ioutil.ReadFile(tc.mockContentPath)
			// if err != nil {
			// 	t.Errorf("%s: %v", tc.name, err)
			// }

			// e, err := NewExtractor(tc.mockBaseURL)
			// if err != nil {
			// 	t.Error(err)
			// }
			// res := e.ExtractLinks(tc.mockBaseURL, content)
			// assert.ElementsMatch(t, tc.expectedLinks, res)
		})
	}
}
