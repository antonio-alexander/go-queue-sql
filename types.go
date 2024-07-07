package sql

import goqueue "github.com/antonio-alexander/go-queue"

const (
	ErrUnsupportedTypef string = "unsupported type: %T\n"
	DefaultPriority     int    = 0
)

type priorityBytes struct {
	bytes    goqueue.Bytes
	priority int
}

type Owner interface {
	SetErrorHandler(ErrorHandlerFx)
	Configure(config *Configuration) error
	Initialize() error
	Shutdown()
}

type ErrorHandlerFx func(error)
