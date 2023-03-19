package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/uppuluriaditya/word-descriptor/models"
)

var rwLock sync.Mutex

func csvPrinter() {
	file, err := os.Open("sample.csv")

	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	for {
		record, err := csvReader.Read()

		if err != nil {
			if err == io.EOF {
				// we have reached the end of file
				break
			}
			log.Fatalln(err)
		}

		for i := 0; i < len(record); i++ {
			if i == 0 {
				words := strings.Split(record[i], " ")
				fmt.Println(words[0])
				fmt.Println(words[1])
				continue
			}
			fmt.Println(record[i])
		}
		fmt.Println("============================")
	}
}

func Colly_test() {
	c := colly.NewCollector()
	// setting a valid User-Agent header
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	// c.OnHTML("div[id=divFullAnswer]", func(e *colly.HTMLElement) {
	// 	// e.Request.Visit(e.Attr("href"))
	// 	allRows := e.DOM.Find("div[class=row]").Find("span[class=h5]")

	// 	for _, row := range allRows.Nodes {
	// 		// if idx == 0 {
	// 		// 	continue
	// 		// }

	// 		fmt.Println(row)

	// 	}
	// })
	// c.OnHTML("div.col-3.heading", func(e *colly.HTMLElement) {
	// 	fmt.Println(e)
	// })

	c.OnHTML("span.h5.light.answer", func(e *colly.HTMLElement) {
		fmt.Println("=================")
		fmt.Println(e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://sanskritabhyas.in/%E0%A4%A8%E0%A5%80-Shabd-Roop")
}

func Write(noun *models.NounForm) error {
	rwLock.Lock()
	file, err := os.OpenFile("output.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var record []string

	for i := 0; i < len(noun.Forms); i++ {
		nounForm := noun.Forms[i]
		record = []string{
			nounForm.Shabdam,
			nounForm.Vibhakthi,
			nounForm.Ekavachan,
			nounForm.Dvivachan,
			nounForm.Bahuvachan,
			nounForm.Nirdesh,
			nounForm.Lingam,
		}
		writer.Write(record)

		if err != nil {
			return err
		}
	}
	rwLock.Unlock()
	return nil
}

func Write_to_csv(noun *models.NounForm, doneChan chan<- bool) {
	file, _ := os.OpenFile("output.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	return err
	// }
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var record []string

	for i := 0; i < len(noun.Forms); i++ {
		nounForm := noun.Forms[i]
		record = []string{
			nounForm.Shabdam,
			nounForm.Vibhakthi,
			nounForm.Ekavachan,
			nounForm.Dvivachan,
			nounForm.Bahuvachan,
			nounForm.Nirdesh,
			nounForm.Lingam,
		}
		writer.Write(record)

		// if err != nil {
		// 	return err
		// }
	}
	doneChan <- true
	// return nil

}

func Scrape(shabd, nirdesh, lingam, url string) (*models.NounForm, error) {
	// vibhakthis := []string{"प्रथमा", "सम्बोधन", "द्वितीया", "तृतीया", "चतुर्थी", "पञ्चमी", "षष्ठी", "सप्तमी"}

	// c := colly.NewCollector()
	c := colly.NewCollector(colly.MaxDepth(1), colly.DetectCharset(), colly.Async(true), colly.AllowURLRevisit())
	// c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1})
	c.SetRequestTimeout(30 * time.Second)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	var nouns models.NounForm
	c.OnHTML("div[id=divFullAnswer]", func(e *colly.HTMLElement) {
		lines := e.DOM.Find("span.h5.light.answer")

		headings := e.DOM.Find("div.col-3.heading").Find("span.h5")

		// for idx, head := range headings.Nodes {
		// 	if idx <= 3 { // remove first row
		// 		continue
		// 	}
		// 	fmt.Println(idx, head.FirstChild.Data)
		// }

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

func Scrape_shabd(shabd, nirdesh, lingam, url string, nounForm chan<- *models.NounForm) {

	vibhakthis := []string{"प्रथमा", "सम्बोधन", "द्वितीया", "तृतीया", "चतुर्थी", "पञ्चमी", "षष्ठी", "सप्तमी"}

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	var nouns models.NounForm
	c.OnHTML("div[id=divFullAnswer]", func(e *colly.HTMLElement) {
		lines := e.DOM.Find("span.h5.light.answer")

		for idx := 0; idx < 8; idx++ {
			var record models.Record
			record.Shabdam = shabd
			record.Nirdesh = nirdesh
			record.Lingam = lingam

			vidx := 3 * idx
			record.Ekavachan = strings.TrimSpace(lines.Nodes[vidx].FirstChild.Data)
			record.Dvivachan = strings.TrimSpace(lines.Nodes[vidx+1].FirstChild.Data)
			record.Bahuvachan = strings.TrimSpace(lines.Nodes[vidx+2].FirstChild.Data)
			record.Vibhakthi = vibhakthis[idx]
			nouns.Forms[idx] = record
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit(url)
	nounForm <- &nouns
}

func write_to_failed_req(str string) {
	file, _ := os.OpenFile("failed_to_scrape.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var record []string = []string{str}
	writer.Write(record)

	// if err != nil {
	// 	return err
	// }
}
