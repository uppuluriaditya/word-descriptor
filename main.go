package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/uppuluriaditya/word-descriptor/constants"
	"github.com/uppuluriaditya/word-descriptor/csvwriter"
	"github.com/uppuluriaditya/word-descriptor/workerqueue"
)

func main() {
	// forms, err := utils.Scrape("exampl", "nire", "pumlingam", "https://sanskritabhyas.in/%E0%A4%A8%E0%A5%80-Shabd-Roop")
	// fmt.Println(forms.Forms)
	// return
	inputFileName := "sample.2.csv"
	outputFileName := "output.csv"
	failedFileName := "failed.csv"
	specialFileName := "special.csv"

	numOfWorkers := 2

	startTime := time.Now()
	file, err := os.Open(inputFileName)
	var rowCount int

	constants.DoneProgram = make(chan bool)
	constants.RecordCount = make(chan int)

	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	csvwriter.Writer = csvwriter.NewCSVWriter(outputFileName, failedFileName, specialFileName)
	csvwriter.Writer.Start()

	dispatcher := workerqueue.NewDispatcher(numOfWorkers)
	dispatcher.Run()

	go func() {
		for {
			record, err := csvReader.Read()

			if err != nil {
				if err == io.EOF {
					// we have reached the end of file
					constants.RecordCount <- rowCount
					break
				}
				log.Fatalln(err)
			}
			rowCount++

			//	format ईकारान्त पुंलिङ्गम्,नी,https://sanskritabhyas.in/%E0%A4%A8%E0%A5%80-Shabd-Roop

			words := strings.Split(record[0], " ")

			scrapeUrl := record[2]

			reader := &ReadRequest{
				Shabd:   record[1],
				Nirdesh: words[0],
				Lingam:  words[1],
				URL:     scrapeUrl,
			}

			workerqueue.JobQueue <- reader
		}
	}()

	<-constants.DoneProgram
	dispatcher.Close()
	csvwriter.Writer.Close()
	timeSince := time.Since(startTime)
	log.Println("Execution time:", timeSince)
}
