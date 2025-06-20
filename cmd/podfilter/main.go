package main

import (
	"log"

	"github.com/m-ueno/podcast-filter/pkg/filter"
)

func main() {
	taskSet := filter.NewFilterFeedTaskSet(nil, "./public")
	taskSet.AddTask(
		filter.NewFilterFeedTask(
			"https://www.omnycontent.com/d/playlist/1e3bd144-9b57-451a-93cf-ac0e00e74446/50382bb4-3af3-4250-8ddc-ac0f0033ceb5/07a1de49-67cf-4714-8581-ac1000059302/podcast.rss",
			filter.And(
				filter.DescriptionIncludeRule{Substring: "奥山"},
				filter.DescriptionIncludeRule{Substring: "伊藤"},
				filter.DescriptionIncludeRule{Substring: "神田"},
			),
			"feed-iok.rss"),
	)
	taskSet.AddTask(
		filter.NewFilterFeedTask( // ニュースの現場から
			"https://www.omnycontent.com/d/playlist/1e3bd144-9b57-451a-93cf-ac0e00e74446/50382bb4-3af3-4250-8ddc-ac0f0033ceb5/684015f9-2396-4ac4-bc1f-ac0f0033d08c/podcast.rss",
			filter.And(
				filter.DescriptionExcludeRule{Substring: "笑い飯"},
				filter.DescriptionExcludeRule{Substring: "スポーツ部"},
			),
			"feed-news.rss"),
	)

	if err := taskSet.Run(); err != nil {
		log.Fatal(err)
	}
}
