package constants

// Number of records that were read from the input file
// channel will be written by reader and read by the writer
var RecordCount chan int

// Returns whether all the writes are done or not
// channel will be written by csv writer and read by dispatcher
// this helps dispatcher know that all the records are done and it can close the workers
var WritesDone chan bool

// written by dispatcher and read by main
var DoneProgram chan bool
