package sql

import (
	"context"
	"fmt"
	"strings"
	"sync"

	go_sql "database/sql"

	goqueue "github.com/antonio-alexander/go-queue"

	"github.com/pkg/errors"
)

type sqlQueue struct {
	sync.WaitGroup
	context.Context
	*go_sql.DB
	started      bool
	config       Configuration
	errorHandler ErrorHandlerFx
	stopper      chan struct{}
}

func New(parameters ...interface{}) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Peeker
	goqueue.Length
	Owner
	ErrorHandler
} {
	s := &sqlQueue{}
	for _, parameter := range parameters {
		switch v := parameter.(type) {
		case *Configuration:
			if err := s.Initialize(v); err != nil {
				panic(err)
			}
		}
	}
	return s
}

func (s *sqlQueue) error(err error) {
	if s.errorHandler == nil {
		return
	}
	s.errorHandler(err)
}

func (s *sqlQueue) launchContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Add(1)
	go func() {
		defer s.WaitGroup.Done()

		select {
		case <-s.stopper:
			cancel()
		case <-ctx.Done():
		}
	}()
	s.Context = ctx
}

func (s *sqlQueue) length() (int, error) {
	query := fmt.Sprintf("SELECT count(*) FROM %s", s.config.Table)
	row := s.QueryRow(query)
	length := 0
	if err := row.Scan(&length); err != nil {
		return -1, err
	}
	return length, nil
}

func (s *sqlQueue) enqueue(items ...goqueue.Bytes) error {
	if len(items) <= 0 {
		return nil
	}
	args := make([]interface{}, 0, len(items))
	values := []string{}
	for _, bytes := range items {
		args = append(args, bytes)
		values = append(values, "(?)")
	}
	query := fmt.Sprintf("INSERT INTO %s (data) VALUES %s;", s.config.Table, strings.Join(values, ","))
	if _, err := s.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func (s *sqlQueue) dequeue(n ...int) ([]goqueue.Bytes, error) {
	var args []interface{}
	var bytes []goqueue.Bytes
	var query string

	switch {
	default:
		query = fmt.Sprintf("DELETE FROM %s ORDER BY id ASC RETURNING data;", s.config.Table)
	case len(n) > 0:
		query = fmt.Sprintf("DELETE FROM %s ORDER BY id ASC LIMIT ? RETURNING data;", s.config.Table)
		args = append(args, n[0])
	}
	rows, err := s.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		b := []byte{}
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		bytes = append(bytes, b)
	}
	return bytes, nil
}

func (s *sqlQueue) peek(n ...int) ([]goqueue.Bytes, error) {
	var args []interface{}
	var bytes []goqueue.Bytes
	var query string

	switch {
	default:
		query = fmt.Sprintf("SELECT data FROM %s ORDER BY id ASC;", s.config.Table)
	case len(n) > 0:
		query = fmt.Sprintf("SELECT data FROM %s ORDER BY id ASC LIMIT ?;", s.config.Table)
		args = append(args, n[0])
	}
	rows, err := s.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		b := []byte{}
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		bytes = append(bytes, b)
	}
	return bytes, nil
}

func (s *sqlQueue) createTable(config *Configuration) error {
	if !config.CreateTable {
		return nil
	}
	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id SERIAL,
			data BLOB
		) ENGINE = InnoDB;`, config.Table)
	if _, err := s.Query(query); err != nil {
		return err
	}
	return nil
}

func (s *sqlQueue) QueryRow(query string, args ...interface{}) *go_sql.Row {
	if s.config.QueryTimeout <= 0 {
		return s.DB.QueryRow(query, args...)
	}
	ctx, cancel := context.WithTimeout(s.Context, s.config.QueryTimeout)
	defer cancel()
	return s.QueryRowContext(ctx, query, args...)
}

func (s *sqlQueue) Query(query string, args ...interface{}) (*go_sql.Rows, error) {
	if s.config.QueryTimeout <= 0 {
		return s.DB.Query(query, args...)
	}
	ctx, cancel := context.WithTimeout(s.Context, s.config.QueryTimeout)
	defer cancel()
	return s.DB.QueryContext(ctx, query, args...)
}

func (s *sqlQueue) Exec(query string, args ...interface{}) (go_sql.Result, error) {
	if s.config.QueryTimeout <= 0 {
		return s.DB.Exec(query, args...)
	}
	ctx, cancel := context.WithTimeout(s.Context, s.config.QueryTimeout)
	defer cancel()
	return s.DB.ExecContext(ctx, query, args...)
}

func (s *sqlQueue) Initialize(config *Configuration) error {
	if s.started {
		return errors.New("started")
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		config.Username, config.Password, config.Hostname, config.Port, config.Database, config.ParseTime)
	sql, err := go_sql.Open(config.DriverName, dataSourceName)
	if err != nil {
		return err
	}
	if err = sql.Ping(); err != nil {
		return err
	}
	if err := s.createTable(config); err != nil {
		return err
	}
	s.DB = sql
	s.config.Table = config.Table
	s.stopper = make(chan struct{})
	s.launchContext()
	s.started = true
	return nil
}

func (s *sqlQueue) SetErrorHandler(errorHandler ErrorHandlerFx) {
	s.errorHandler = errorHandler
}

func (s *sqlQueue) Shutdown() {
	if !s.started {
		return
	}
	close(s.stopper)
	s.Wait()
	if err := s.DB.Close(); err != nil {
		s.error(err)
	}
	s.started = false
	s.errorHandler = nil
}

func (s *sqlQueue) Close() (items []interface{}) {
	s.Shutdown()
	return nil
}

func (s *sqlQueue) Peek() (items []interface{}) {
	bytes, err := s.peek()
	if err != nil {
		s.error(err)
		return nil
	}
	for _, bytes := range bytes {
		items = append(items, bytes)
	}
	return items
}

func (s *sqlQueue) PeekHead() (item interface{}, underflow bool) {
	bytes, err := s.peek(1)
	if err != nil {
		s.error(err)
		return nil, true
	}
	if len(bytes) <= 0 {
		return nil, true
	}
	return bytes[0], false
}

func (s *sqlQueue) PeekFromHead(n int) (items []interface{}) {
	bytes, err := s.peek(n)
	if err != nil {
		s.error(err)
		return nil
	}
	for _, bytes := range bytes {
		items = append(items, bytes)
	}
	return items
}

func (s *sqlQueue) Dequeue() (item interface{}, underflow bool) {
	bytes, err := s.dequeue(1)
	if err != nil {
		s.error(err)
		return nil, true
	}
	if len(bytes) <= 0 {
		return nil, true
	}
	return []byte(bytes[0]), false
}

func (s *sqlQueue) DequeueMultiple(n int) []interface{} {
	bytes, err := s.dequeue(n)
	if err != nil {
		s.error(err)
		return nil
	}
	if len(bytes) <= 0 {
		return nil
	}
	items := make([]interface{}, 0, len(bytes))
	for _, b := range bytes {
		items = append(items, []byte(b))
	}
	return items
}

func (s *sqlQueue) Flush() []interface{} {
	bytes, err := s.dequeue()
	if err != nil {
		s.error(err)
		return nil
	}
	if len(bytes) <= 0 {
		return nil
	}
	items := make([]interface{}, 0, len(bytes))
	for _, b := range bytes {
		items = append(items, []byte(b))
	}
	return items
}

func (s *sqlQueue) Enqueue(item interface{}) bool {
	bytes, err := convertSingle(item)
	if err != nil {
		s.error(err)
		return true
	}
	if err = s.enqueue(bytes); err != nil {
		s.error(err)
		return true
	}
	return false
}

func (s *sqlQueue) EnqueueMultiple(items []interface{}) ([]interface{}, bool) {
	bytes, err := convertMultiple(items)
	if err != nil {
		s.error(err)
		return nil, true
	}
	if err = s.enqueue(bytes...); err != nil {
		s.error(err)
		return nil, true
	}
	return nil, false
}

func (s *sqlQueue) Length() int {
	length, err := s.length()
	if err != nil {
		s.error(err)
		return -1
	}
	return length
}
