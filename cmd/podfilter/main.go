package main

import (
	"log"
	"net/http"
	"os"

	"github.com/m-ueno/podcast-filter/pkg/filter"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	feedUrl := "https://www.omnycontent.com/d/playlist/1e3bd144-9b57-451a-93cf-ac0e00e74446/50382bb4-3af3-4250-8ddc-ac0f0033ceb5/07a1de49-67cf-4714-8581-ac1000059302/podcast.rss"

	resp, err := http.Get(feedUrl)
	must(err)
	defer resp.Body.Close()

	reader := resp.Body
	writer := os.Stdout
	rules := filter.RuleSet{
		{DesciptionInclude: "奥山"},
		{DesciptionInclude: "伊藤"},
		{DesciptionInclude: "神田"},
	}

	err = filter.ParseFilterWrite(reader, writer, rules)
	must(err)
}
