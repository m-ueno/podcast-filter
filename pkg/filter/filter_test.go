package filter_test

import (
	"os"
	"strings"
	"testing"

	"github.com/m-ueno/podcast-filter/pkg/filter"
)

var rssData string

func NewTestReader() (*strings.Reader, error) {
	if rssData == "" {
		testFile := "./testdata/podcast.rss"
		body, err := os.ReadFile(testFile)
		if err != nil {
			return nil, err
		}
		rssData = string(body)
	}
	return strings.NewReader(rssData), nil
}

func TestParseFilter(t *testing.T) {
	tests := []struct {
		Name    string
		RuleSet filter.RuleSet
		Want    int
	}{
		{
			Name: "empty",
			RuleSet: filter.RuleSet{
				{DesciptionInclude: "奥山昌次郎"},
			},
			Want: 0,
		},
		{
			Name: "single",
			RuleSet: filter.RuleSet{
				{DesciptionInclude: "奥山晶二郎"},
				{DesciptionInclude: "甲子園の応援"},
			},
			Want: 1,
		},
		{
			Name: "the three",
			RuleSet: filter.RuleSet{
				{DesciptionInclude: "奥山"},
				{DesciptionInclude: "伊藤"},
				{DesciptionInclude: "神田"},
			},
			Want: 191,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			reader, err := NewTestReader()
			if err != nil {
				t.Error(err)
			}

			feed, err := filter.ParseFilter(reader, test.RuleSet)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if len(feed.Items) != test.Want {
				t.Errorf("expected %d items, got %d", test.Want, len(feed.Items))
			}
		})
	}
}

func TestParseFilterWrite(t *testing.T) {
	reader, err := NewTestReader()
	if err != nil {
		t.Fatal(err)
	}

	writer := &strings.Builder{}
	ruleSet := filter.RuleSet{
		{DesciptionInclude: "甲子園の応援"},
		{DesciptionInclude: "奥山晶二郎"},
	}

	err = filter.ParseFilterWrite(reader, writer, ruleSet)
	if err != nil {
		t.Fatal(err)
	}

	if !testing.Short() {
		t.Log(writer.String())
	}
}
