package sql

import (
	"context"
	"fmt"
	"strings"
	"sync"

	go_sql "database/sql"

	goqueue "github.com/antonio-alexander/go-queue"
	goqueuepriority "github.com/antonio-alexander/go-queue-priority"

	"github.com/pkg/errors"
)

type sqlQueue struct {
	sync.RWMutex
	sync.WaitGroup
	context.Context
	*go_sql.DB
	started      bool
	configured   bool
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
	goqueuepriority.PriorityEnqueuer
	Owner
} {
	var config *Configuration

	s := &sqlQueue{}
	for _, parameter := range parameters {
		switch v := parameter.(type) {
		case ErrorHandlerFx:
			s.errorHandler = v
		case *Configuration:
			config = v
		}
	}
	if config != nil {
		if err := s.Configure(config); err != nil {
			panic(err)
		}
		if err := s.Initialize(); err != nil {
			panic(err)
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
	var cancel context.CancelFunc

	s.Context, cancel = context.WithCancel(context.Background())
	s.Add(1)
	started := make(chan struct{})
	go func() {
		defer s.WaitGroup.Done()

		close(started)
		select {
		case <-s.stopper:
			cancel()
		case <-s.Context.Done():
		}
	}()
	<-started
}

func (s *sqlQueue) length() (int, error) {
	var length int

	query := fmt.Sprintf("SELECT count(*) FROM %s", s.config.Table)
	row := s.QueryRow(query)
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

func (s *sqlQueue) enqueueWithPriority(items ...priorityBytes) error {
	if len(items) <= 0 {
		return nil
	}
	args := make([]interface{}, 0, len(items)*2)
	values := make([]string, 0, len(items))
	for _, priorityBytes := range items {
		args = append(args, priorityBytes.bytes, priorityBytes.priority)
		values = append(values, "(?,?)")
	}
	query := fmt.Sprintf("INSERT INTO %s (data,priority) VALUES %s;", s.config.Table, strings.Join(values, ","))
	if _, err := s.Exec(query, args...); err != nil {
		return err
	}
	return nil
}

func (s *sqlQueue) dequeue(n ...int) ([]goqueue.Bytes, error) {
	if s.config.WithPriority {
		return s.dequeueWithPriority(n...)
	}

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

func (s *sqlQueue) dequeueWithPriority(n ...int) ([]goqueue.Bytes, error) {
	var args []interface{}
	var bytes []goqueue.Bytes
	var query string

	switch {
	default:
		query = fmt.Sprintf("DELETE FROM %s ORDER BY priority DESC RETURNING data;", s.config.Table)
	case len(n) > 0:
		query = fmt.Sprintf("DELETE FROM %s ORDER BY priority DESC LIMIT ? RETURNING data;", s.config.Table)
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

func (s *sqlQueue) Configure(config *Configuration) error {
	s.Lock()
	defer s.Unlock()

	if config == nil {
		return errors.New("config is nil")
	}

	//TODO: validate configuration
	s.config = *config
	s.configured = true
	return nil
}

func (s *sqlQueue) Initialize() error {
	if s.started {
		return errors.New("started")
	}
	if !s.configured {
		return errors.New("not configured")
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		s.config.Username, s.config.Password, s.config.Hostname, s.config.Port, s.config.Database, s.config.ParseTime)
	sql, err := go_sql.Open(s.config.DriverName, dataSourceName)
	if err != nil {
		return err
	}
	s.DB = sql
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

func (s *sqlQueue) PriorityEnqueue(item interface{}, priorities ...int) bool {
	var priority int

	for _, p := range priorities {
		priority = p
	}
	bytes, err := convertSingle(item)
	if err != nil {
		s.error(err)
		return true
	}
	if err = s.enqueueWithPriority(priorityBytes{
		bytes:    bytes,
		priority: priority,
	}); err != nil {
		s.error(err)
		return true
	}
	return false
}

func (s *sqlQueue) PriorityEnqueueMultiple(items []interface{}, priorities ...int) ([]interface{}, bool) {
	priorityBytes, err := convertMultipleWithPriority(items, priorities...)
	if err != nil {
		s.error(err)
		return nil, true
	}
	if err = s.enqueueWithPriority(priorityBytes...); err != nil {
		s.error(err)
		return nil, true
	}
	return nil, false
}
