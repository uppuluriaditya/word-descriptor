package constants

// Number of records that were read from the input file
// channel will be written by reader and read by the writer
var InputRecordCount chan int

// written by csvWriter and read by main
var DoneRecWrite chan bool
