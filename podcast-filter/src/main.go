package main

import (
    "fmt"
    "podcast-filter/src/filter"
)

func main() {
    // Initialize application configurations here

    // Example podcast feed (this would typically come from an external source)
    podcastFeed := []string{
        "Episode 1: Introduction to Go",
        "Episode 2: Advanced Go Techniques",
        "Episode 3: Go for Podcasting",
    }

    // Call the filtering function
    filteredFeed := filter.FilterFeed(podcastFeed)

    // Output the filtered feed
    fmt.Println("Filtered Podcast Feed:")
    for _, episode := range filteredFeed {
        fmt.Println(episode)
    }
}