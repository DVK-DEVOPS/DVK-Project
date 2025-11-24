package scraper

import (
	"strings"

	"DVK-Project/db"
	"DVK-Project/logging"

	"github.com/gocolly/colly"
)

func Run(url string, repo *db.PageRepository) error {
	c := colly.NewCollector(
		colly.UserAgent("SearchIndexerBot/1.0"),
		colly.MaxDepth(2),
	)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		title := e.DOM.Find("title").Text()
		content := strings.TrimSpace(e.DOM.Find("body").Text())

		logging.Log.Info().
			Str("event", "page_scraped").
			Str("url", url).
			Msg("scraped")

		err := repo.InsertScrapedPage(title, url, content, "en")
		if err != nil {
			logging.Log.Error().
				Err(err).
				Str("url", url).
				Msg("failed to save scraped page")
		}
	})

	return c.Visit(url)
}
