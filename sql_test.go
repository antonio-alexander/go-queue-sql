package sql_test

import (
	"os"
	"strings"
	"testing"
	"time"

	sql "github.com/antonio-alexander/go-queue-sql"

	goqueue "github.com/antonio-alexander/go-queue"
	infinite_tests "github.com/antonio-alexander/go-queue/infinite/tests"
	goqueue_tests "github.com/antonio-alexander/go-queue/tests"

	_ "github.com/go-sql-driver/mysql"
)

var configuration = new(sql.Configuration)

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configuration.Default()
	configuration.FromEnvs(envs)
}

func newQueue(t *testing.T, configuration *sql.Configuration) interface {
	goqueue.Dequeuer
	goqueue.Enqueuer
	goqueue.Owner
	goqueue.Peeker
	goqueue.Length
	sql.Owner
} {
	errorHandlerFx := sql.ErrorHandlerFx(func(err error) { t.Log(err) })
	sql := sql.New(configuration, errorHandlerFx)
	//KIM: if we don't flush because sql is persistent
	// all the tests will fail because of data from a
	// previous test
	sql.Flush()
	return sql
}

func testSqlQueue(t *testing.T) {
	timeout, rate := time.Second, time.Second
	t.Run("Test Dequeue", goqueue_tests.TestDequeue(t, rate, timeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Dequeue Multiple", goqueue_tests.TestDequeueMultiple(t, rate, timeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Flush", goqueue_tests.TestFlush(t, rate, timeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Enqueue", infinite_tests.TestEnqueue(t, rate, timeout, func() interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Enqueue Multiple", infinite_tests.TestEnqueueMultiple(t, rate, timeout, func() interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Peek", goqueue_tests.TestPeek(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Peek From Head", goqueue_tests.TestPeekFromHead(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Length", goqueue_tests.TestLength(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Queue (goqueue)", goqueue_tests.TestQueue(t, rate, timeout, func(i int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return newQueue(t, configuration)
	}))
	t.Run("Test Asynchronous", goqueue_tests.TestAsync(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return newQueue(t, configuration)
	}))
	//TODO: test error handler
}

func TestSqlQueue(t *testing.T) {
	testSqlQueue(t)
}
