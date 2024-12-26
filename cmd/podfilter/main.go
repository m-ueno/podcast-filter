package main

import (
	"log"

	"github.com/m-ueno/podcast-filter/pkg/filter"
)

func main() {
	taskSet := filter.NewFilterFeedTaskSet([]filter.FilterFeedTask{
		{
			FeedUrl: "https://www.omnycontent.com/d/playlist/1e3bd144-9b57-451a-93cf-ac0e00e74446/50382bb4-3af3-4250-8ddc-ac0f0033ceb5/07a1de49-67cf-4714-8581-ac1000059302/podcast.rss",
			RuleSet: filter.RuleSet{
				{DescriptionInclude: "奥山"},
				{DescriptionInclude: "伊藤"},
				{DescriptionInclude: "神田"},
			},
			OutputFileName: "feed-iok.rss",
		},
	}, "./public")

	if err := taskSet.Run(); err != nil {
		log.Fatal(err)
	}
}
