package filter

import (
	"io"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

type Rule interface {
	Match(item *PodcastItem) bool
}

type DescriptionIncludeRule struct {
	Substring string
}

func (r DescriptionIncludeRule) Match(item *PodcastItem) bool {
	return strings.Contains(item.Description, r.Substring)
}

type AndRule []Rule

// Match implement rule
func (r AndRule) Match(item *PodcastItem) bool {
	for _, rule := range r {
		if !rule.Match(item) {
			return false
		}
	}
	return true
}

func And(rules ...Rule) AndRule {
	return AndRule(rules)
}

func ParseFilter(r io.Reader, rule Rule) (*PodcastFeed, error) {
	feed, err := gofeed.NewParser().Parse(r)
	if err != nil {
		return nil, err
	}

	var items []*PodcastItem
	for _, item := range feed.Items {
		if rule.Match(item) {
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

func ParseFilterWrite(r io.Reader, w io.Writer, rule Rule) error {
	feed, err := ParseFilter(r, rule)
	if err != nil {
		return err
	}
	newFeed := convertGofeedToGorillaFeed(feed)

	return newFeed.WriteRss(w)
}
