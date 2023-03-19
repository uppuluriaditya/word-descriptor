package csvwriter

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/uppuluriaditya/word-descriptor/constants"
	"github.com/uppuluriaditya/word-descriptor/models"
)

var Writer *CSVWriter

type CSVWriter struct {
	fileName        string
	failedFileName  string
	specialFileName string
	dataChan        chan models.NounForm
	failedRec       chan []string
	specialRec      chan []string
	doneChan        chan bool
	numOfWriteReqs  int
	totalRecords    int
}

func NewCSVWriter(fileName, failedFileName, specialFileName string) *CSVWriter {
	return &CSVWriter{
		fileName:        fileName,
		failedFileName:  failedFileName,
		dataChan:        make(chan models.NounForm),
		doneChan:        make(chan bool),
		failedRec:       make(chan []string),
		specialRec:      make(chan []string),
		specialFileName: specialFileName,
	}
}

func (c *CSVWriter) Start() {
	// create a new CSV file
	file, err := os.OpenFile(c.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	failedFile, err := os.OpenFile(c.failedFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	specialFile, err := os.OpenFile(c.specialFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	// create a CSV writer
	writer := csv.NewWriter(file)

	failedWriter := csv.NewWriter(failedFile)

	specialWriter := csv.NewWriter(specialFile)

	// start a separate goroutine to write to the CSV file
	go func() {
		for {
			select {
			case row := <-c.dataChan:
				log.Printf("Got write request for: %v", row.String())
				c.numOfWriteReqs++
				for i := 0; i < len(row.Forms); i++ {
					nounForm := row.Forms[i]
					record := []string{
						nounForm.Shabdam,
						nounForm.Vibhakthi,
						nounForm.Ekavachan,
						nounForm.Dvivachan,
						nounForm.Bahuvachan,
						nounForm.Nirdesh,
						nounForm.Lingam,
					}
					err = writer.Write(record)

					if err != nil {
						log.Printf("Error in writing %v", err)
						panic(err)
					}
					// flush any remaining buffered data
					writer.Flush()
				}
				if c.numOfWriteReqs == c.totalRecords {
					constants.DoneRecWrite <- true
				}

			case row := <-c.failedRec:
				log.Printf("Got fail request for: %v", row)
				c.numOfWriteReqs++
				err = failedWriter.Write(row)
				if err != nil {
					log.Printf("Error in writing %v", err)
					panic(err)
				}
				// flush any remaining buffered data
				failedWriter.Flush()

				if c.numOfWriteReqs == c.totalRecords {
					constants.DoneRecWrite <- true
				}

			case row := <-c.specialRec:
				log.Printf("Got special request for: %v", row)
				err = specialWriter.Write(row)
				if err != nil {
					log.Printf("Error in writing %v", err)
					panic(err)
				}
				// flush any remaining buffered data
				specialWriter.Flush()

			case num := <-constants.InputRecordCount:
				c.totalRecords = num
			case <-c.doneChan:
				// flush any remaining buffered data
				writer.Flush()
				file.Close()
				failedFile.Close()
				specialFile.Close()
				return
			}
		}
	}()
}

func (c *CSVWriter) Write(nouns models.NounForm) {
	c.dataChan <- nouns
}

func (c *CSVWriter) WriteToFailed(str []string) {
	c.failedRec <- str
}

func (c *CSVWriter) WriteToSpecialFile(str []string) {
	c.specialRec <- str
}

func (c *CSVWriter) Close() {
	// signal the writer goroutine to exit
	c.doneChan <- true
}
