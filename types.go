package sql

type bytes []byte

type Owner interface {
	Close()
	Start(config *Configuration) (err error)
	Stop()
}
