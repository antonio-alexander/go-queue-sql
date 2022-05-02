package sql_test

import (
	"os"
	"strings"
	"testing"

	sql "github.com/antonio-alexander/go-queue-sql"

	goqueue "github.com/antonio-alexander/go-queue"
	infinite_tests "github.com/antonio-alexander/go-queue/infinite/tests"
	goqueue_tests "github.com/antonio-alexander/go-queue/tests"

	_ "github.com/go-sql-driver/mysql"
)

var configuration = &sql.Configuration{
	Hostname:     sql.DefaultDatabaseHost,
	Port:         sql.DefaultDatabasePort,
	Database:     sql.DefaultDatabaseName,
	Table:        sql.DefaultDatabaseTable,
	DriverName:   sql.DefaultDriverName,
	Username:     sql.DefaultUsername,
	Password:     sql.DefaultPassword,
	ParseTime:    sql.DefaultParseTime,
	CreateTable:  sql.DefaultCreateTable,
	QueryTimeout: sql.DefaultQueryTimeout,
}

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configuration = sql.ConfigFromEnv(envs)
}

func new(configuration *sql.Configuration) interface {
	goqueue.Dequeuer
	goqueue.Enqueuer
	goqueue.Owner
	goqueue.Peeker
	goqueue.Length
	sql.Owner
} {
	sql := sql.New(configuration)
	//KIM: if we don't flush because sql is persistent
	// all the tests will fail because of data from a
	// previous test
	sql.Flush()
	return sql
}

func TestNew(t *testing.T) {
	infinite_tests.New(t, func(size int) interface {
		goqueue.Owner
		goqueue.Length
	} {
		return new(configuration)
	})
}

func TestDequeue(t *testing.T) {
	goqueue_tests.Dequeue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return new(configuration)
	})
}

func TestDequeueMultiple(t *testing.T) {
	goqueue_tests.DequeueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return new(configuration)
	})
}

func TestFlush(t *testing.T) {
	goqueue_tests.Flush(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return new(configuration)
	})
}

func TestEnqueue(t *testing.T) {
	infinite_tests.Enqueue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Length
	} {
		return new(configuration)
	})
}

func TestEnqueueMultiple(t *testing.T) {
	infinite_tests.EnqueueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Length
	} {
		return new(configuration)
	})
}

func TestPeek(t *testing.T) {
	goqueue_tests.Peek(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return new(configuration)
	})
}

func TestPeekFromHead(t *testing.T) {
	goqueue_tests.PeekFromHead(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return new(configuration)
	})
}

func TestLength(t *testing.T) {
	goqueue_tests.Length(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return new(configuration)
	})
}

func TestQueue(t *testing.T) {
	goqueue_tests.Queue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return new(configuration)
	})
	infinite_tests.Queue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return new(configuration)
	})
}

func TestAsync(t *testing.T) {
	goqueue_tests.Async(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return new(configuration)
	})
}

//TODO: test error handler
