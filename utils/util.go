package utils

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/uppuluriaditya/word-descriptor/models"
)

func Scrape(shabd, nirdesh, lingam, url string) (*models.NounForm, error) {

	// c := colly.NewCollector()
	c := colly.NewCollector(colly.MaxDepth(1), colly.DetectCharset(), colly.Async(true), colly.AllowURLRevisit())
	// c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1})
	c.SetRequestTimeout(30 * time.Second)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	var nouns models.NounForm
	c.OnHTML("div[id=divFullAnswer]", func(e *colly.HTMLElement) {
		lines := e.DOM.Find("span.h5.light.answer")
		headings := e.DOM.Find("div.col-3.heading").Find("span.h5")

		for idx := 0; idx < len(headings.Nodes)-4; idx++ {
			var record models.Record
			record.Shabdam = shabd
			record.Nirdesh = nirdesh
			record.Lingam = lingam

			vidx := 3 * idx
			vibhakthi := 4 + idx
			record.Ekavachan = strings.TrimSpace(lines.Nodes[vidx].FirstChild.Data)
			record.Dvivachan = strings.TrimSpace(lines.Nodes[vidx+1].FirstChild.Data)
			record.Bahuvachan = strings.TrimSpace(lines.Nodes[vidx+2].FirstChild.Data)
			record.Vibhakthi = strings.TrimSpace(headings.Nodes[vibhakthi].FirstChild.Data)
			nouns.Forms = append(nouns.Forms, record)
		}
	})

	var scrapeErr error

	c.OnError(func(r *colly.Response, err error) {
		if err != nil {
			scrapeErr = err
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit(url)
	c.Wait()

	return &nouns, scrapeErr
}
