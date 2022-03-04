package sql_test

import (
	"os"
	"strings"
	"testing"

	goqueue "github.com/antonio-alexander/go-queue"
	sql "github.com/antonio-alexander/go-queue-sql"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

var configuration *sql.Configuration

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configuration = sql.ConfigFromEnv(envs)
}

func initQueue() (interface {
	goqueue.Dequeuer
	goqueue.Enqueuer
	//REVIEW: do we want to implement events?
	// goqueue.Event
	goqueue.Info
	sql.Owner
}, error) {
	s := sql.New()
	if err := s.Start(configuration); err != nil {
		return nil, err
	}
	return s, nil
}

func TestEnqueue(t *testing.T) {
	s, err := initQueue()
	assert.Nil(t, err)
	bytes := []byte("Hello, World!")
	overflow := s.Enqueue(bytes)
	assert.False(t, overflow)
	length := s.Length()
	assert.GreaterOrEqual(t, length, 1)
	item, underflow := s.Dequeue()
	assert.False(t, underflow)
	bytes, ok := item.([]byte)
	assert.True(t, ok)
	assert.Equal(t, []byte("Hello, World!"), bytes)
	s.Close()
}
