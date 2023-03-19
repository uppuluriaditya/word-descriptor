package main

import (
	"fmt"
	"log"

	"github.com/uppuluriaditya/word-descriptor/csvwriter"
	"github.com/uppuluriaditya/word-descriptor/utils"
)

// This implements the Job interface
type InputRecord struct {
	Shabd   string
	Nirdesh string
	Lingam  string
	URL     string
}

func (r InputRecord) String() string {
	return fmt.Sprintf("%s:%s", r.Shabd, r.URL)
}

func (r InputRecord) Process() error {
	// scrape
	nouns, err := utils.Scrape(r.Shabd, r.Nirdesh, r.Lingam, r.URL)
	if err != nil {
		errStr := fmt.Sprintf("error in scraping: %v, err : %s\n", r.Shabd, err)
		log.Println(errStr)
		return err
	}

	if len(nouns.Forms) == 0 || nouns.Forms[0].Ekavachan == "" {
		err := fmt.Sprintf("No data found for %v", r.URL)
		log.Println(err)
		var str []string = []string{r.URL}
		csvwriter.Writer.WriteToFailed(str)
		return nil
	}

	if len(nouns.Forms) < 8 {
		// 8 rows were not present. store it separately for any interesting observations
		err := fmt.Sprintf("Special Record Found: %v", r.URL)
		log.Println(err)
		var str []string = []string{r.URL}
		csvwriter.Writer.WriteToSpecialFile(str)
		// Do not return
	}

	csvwriter.Writer.Write(*nouns)
	return nil
}
