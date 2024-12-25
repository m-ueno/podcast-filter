package filter

type Podcast struct {
    Title       string
    Description string
    Category    string
    Link        string
}

func FilterFeed(podcasts []Podcast, criteria string) []Podcast {
    var filtered []Podcast
    for _, podcast := range podcasts {
        if podcast.Category == criteria {
            filtered = append(filtered, podcast)
        }
    }
    return filtered
}