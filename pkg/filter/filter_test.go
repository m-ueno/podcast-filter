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
		Name string
		Rule filter.Rule
		Want int
	}{
		{
			Name: "empty",
			Rule: filter.DescriptionIncludeRule{"奥山昌次郎"},
			Want: 0,
		},
		{
			Name: "exclude",
			Rule: filter.DescriptionExcludeRule{"神田大介"},
			Want: 878,
		},
		{
			Name: "single",
			Rule: filter.And(
				filter.DescriptionIncludeRule{"奥山晶二郎"},
				filter.DescriptionIncludeRule{"甲子園の応援"},
			),
			Want: 1,
		},
		{
			Name: "the three",
			Rule: filter.And(
				filter.DescriptionIncludeRule{"奥山"},
				filter.DescriptionIncludeRule{"伊藤"},
				filter.DescriptionIncludeRule{"神田"},
			),
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

			feed, err := filter.ParseFilter(reader, test.Rule)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if len(feed.Items) != test.Want {
				t.Errorf("expected %d items, got %d", test.Want, len(feed.Items))
			}
		})
	}
}

func TestParseFilter_badinput(t *testing.T) {
	badinput := strings.NewReader("<xml></xml>")
	_, err := filter.ParseFilter(badinput, filter.DescriptionIncludeRule{})
	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func TestParseFilterWrite(t *testing.T) {
	reader, err := NewTestReader()
	if err != nil {
		t.Fatal(err)
	}

	writer := &strings.Builder{}
	ruleSet := filter.And(
		filter.DescriptionIncludeRule{"甲子園の応援"},
		filter.DescriptionIncludeRule{"奥山晶二郎"},
	)

	err = filter.ParseFilterWrite(reader, writer, ruleSet)
	if err != nil {
		t.Fatal(err)
	}

	if !testing.Short() {
		t.Log(writer.String())
	}
}
