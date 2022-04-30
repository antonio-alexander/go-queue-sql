package sql_test

import (
	"os"
	"strings"
	"testing"

	sql "github.com/antonio-alexander/go-queue-sql"
	example "github.com/antonio-alexander/go-queue-sql/internal/example"

	"github.com/stretchr/testify/assert"

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

func TestDequeue(t *testing.T) {
	s := sql.New()
	err := s.Initialize(configuration)
	assert.Nil(t, err)
	data := &example.Data{
		Int:    12345,
		String: "12345",
	}
	overflow := s.Enqueue(data)
	assert.False(t, overflow)
	dataDequeued, underflow := example.Dequeue(s)
	assert.False(t, underflow)
	assert.Equal(t, data, dataDequeued)
	s.Shutdown()
}
