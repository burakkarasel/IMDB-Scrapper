package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// Star holds our actor's data
type Star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []Movie
}

// Movie holds movie data
type Movie struct {
	Title string
	Year  string
}

func main() {
	month := 2
	day := 19
	crawl(month, day)
}

// crawl crawls the URL we specified according to month and day and logs out the data
func crawl(month, day int) {
	// here we created new collector to visit these domains
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	// now we created a new clone to collect the informations from these domains
	infoCollector := c.Clone()

	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileUrl := e.ChildAttr("div.lister-item-image > a", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})

	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	infoCollector.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		tmpProfile := Star{}

		tmpProfile.Name = e.ChildText("h1.header > span.itemprop")
		tmpProfile.Photo = e.ChildAttr("#name-poster", "src")
		tmpProfile.JobTitle = e.ChildText("#name-job-categories > a > span.itemprop")
		tmpProfile.BirthDate = e.ChildAttr("#name-born-info time", "datetime")
		tmpProfile.Bio = strings.TrimSpace(e.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))

		e.ForEach("div.knownfor-title", func(_ int, kf *colly.HTMLElement) {
			tmpMovie := Movie{}
			tmpMovie.Title = kf.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
			tmpMovie.Year = kf.ChildText("div.knownfor-year > span.knownfor-ellipsis")
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		})

		js, err := json.MarshalIndent(tmpProfile, "", "    ")

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(js))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting:", r.URL.String())
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting profile URL:", r.URL.String())
	})

	startURL := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	c.Visit(startURL)
}
