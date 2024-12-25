package filter

import (
	"io"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

type Rule struct {
	DesciptionInclude string
}

func (r Rule) Match(description string) bool {
	return strings.Contains(description, r.DesciptionInclude)
}

type RuleSet []Rule

func ParseFilter(r io.Reader, ruleSet RuleSet) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.Parse(r)
	if err != nil {
		return nil, err
	}

	var items []*gofeed.Item
	for _, item := range feed.Items {
		includeItem := true
		for _, rule := range ruleSet {
			if !rule.Match(item.Description) {
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

func ParseFilterWrite(r io.Reader, w io.Writer, ruleSet RuleSet) error {
	feed, err := ParseFilter(r, ruleSet)
	if err != nil {
		return err
	}

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

	return newFeed.WriteRss(w)
}
