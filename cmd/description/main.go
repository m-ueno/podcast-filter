package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
)

var inputFile = flag.String("input", "", "input rss file path")

func init() {
	flag.Parse()
}

func main() {
	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	feed, err := gofeed.NewParser().Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	htmlPattern := regexp.MustCompile(`<[^>]*>`)
	const maxBytes = 2000

	for _, item := range feed.Items {
		var runes []rune
		if len(item.Description) > maxBytes {
			runes = []rune(item.Description[:maxBytes])
		} else {
			runes = []rune(item.Description)
		}
		desc := strings.ReplaceAll(string(runes), "\n", " ")
		desc = htmlPattern.ReplaceAllString(desc, "")

		fmt.Println(item.Published + "//" + item.Title + "//" + desc)
		count += 1
	}
	fmt.Printf("count=%d\n", count)
}
