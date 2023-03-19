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
	inputFileName := "sample.2.csv"
	outputFileName := "output.csv"
	failedFileName := "failed.csv"
	specialFileName := "special.csv"

	numOfWorkers := 2

	startTime := time.Now()
	file, err := os.Open(inputFileName)
	var rowCount int

	constants.DoneRecWrite = make(chan bool)
	constants.InputRecordCount = make(chan int)

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
					constants.InputRecordCount <- rowCount
					break
				}
				log.Fatalln(err)
			}
			rowCount++

			//	format ईकारान्त पुंलिङ्गम्,नी,https://sanskritabhyas.in/%E0%A4%A8%E0%A5%80-Shabd-Roop
			words := strings.Split(record[0], " ")
			scrapeUrl := record[2]

			reader := &InputRecord{
				Shabd:   record[1],
				Nirdesh: words[0],
				Lingam:  words[1],
				URL:     scrapeUrl,
			}

			workerqueue.JobQueue <- reader
		}
	}()

	<-constants.DoneRecWrite
	dispatcher.Close()
	csvwriter.Writer.Close()
	timeSince := time.Since(startTime)
	log.Printf("Total of %v records written in %v seconds", rowCount, timeSince)
}
