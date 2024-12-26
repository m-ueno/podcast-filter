package filter

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

type DescriptionIncludeRule struct {
	DesciptionInclude string
}

func (r DescriptionIncludeRule) Match(item *PodcastItem) bool {
	return strings.Contains(item.Description, r.DesciptionInclude)
}

type FilterFeedTaskSet struct {
	Tasks     []FilterFeedTask
	OutputDir string
}

func (ts FilterFeedTaskSet) Run() error {
	if err := os.MkdirAll(ts.OutputDir, 0755); err != nil && err != os.ErrExist {
		return err
	}
	for _, task := range ts.Tasks {
		if err := task.Run(ts.OutputDir); err != nil {
			return err
		}
	}
	return nil
}

type FilterFeedTask struct {
	FeedUrl        string
	RuleSet        RuleSet
	OutputFileName string
}

func (t FilterFeedTask) Run(basedir string) error {
	resp, err := http.Get(t.FeedUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(filepath.Join(basedir, t.OutputFileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return ParseFilterWrite(resp.Body, file, t.RuleSet)
}

type RuleSet []DescriptionIncludeRule

func ParseFilter(r io.Reader, ruleSet RuleSet) (*PodcastFeed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.Parse(r)
	if err != nil {
		return nil, err
	}

	var items []*PodcastItem
	for _, item := range feed.Items {
		includeItem := true
		for _, rule := range ruleSet {
			if !rule.Match(item) {
				includeItem = false
			}
		}
		if includeItem {
			items = append(items, item)
		}
	}

	feed.Items = items
	return feed, nil
}

func convertGofeedToGorillaFeed(feed *PodcastFeed) *feeds.Feed {
	newFeed := &feeds.Feed{
		Title: feed.Title,
		Link: &feeds.Link{
			Href: feed.Link,
		},
		Description: feed.Description,
		Author:      &feeds.Author{Name: feed.Author.Name},
		Image: &feeds.Image{
			Url:   feed.Image.URL,
			Title: feed.Image.Title,
		},
	}
	if feed.UpdatedParsed != nil {
		newFeed.Updated = *feed.UpdatedParsed
	}
	if feed.PublishedParsed != nil {
		newFeed.Created = *feed.PublishedParsed
	}

	for _, item := range feed.Items {
		newFeed.Add(&feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: item.Link},
			Description: item.Description,
			Id:          item.GUID,
			Created:     *item.PublishedParsed,
			Enclosure: &feeds.Enclosure{
				Url:    item.Enclosures[0].URL,
				Length: item.Enclosures[0].Length,
				Type:   item.Enclosures[0].Type,
			},
		})
	}
	return newFeed
}

func ParseFilterWrite(r io.Reader, w io.Writer, ruleSet RuleSet) error {
	feed, err := ParseFilter(r, ruleSet)
	if err != nil {
		return err
	}
	newFeed := convertGofeedToGorillaFeed(feed)

	return newFeed.WriteRss(w)
}
