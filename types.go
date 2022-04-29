package sql

type bytes []byte

type Owner interface {
	Initialize(config *Configuration) error
	Shutdown() error
}
