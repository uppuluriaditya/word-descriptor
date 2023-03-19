package workerqueue

type Job interface {
	Process() error
	String() string
}

var JobQueue chan Job
