package error_handler

import sql "github.com/antonio-alexander/go-queue-sql"

func New(chErr chan error) sql.ErrorHandlerFx {
	return func(err error) {
		select {
		case <-chErr:
		case chErr <- err:
		}
	}
}
