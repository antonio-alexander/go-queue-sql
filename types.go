package sql

const (
	ErrUnsupportedTypef string = "unsupported type: %T\n"
)

type Owner interface {
	SetErrorHandler(ErrorHandlerFx)
	Configure(config *Configuration) error
	Initialize() error
	Shutdown()
}

type ErrorHandlerFx func(error)
