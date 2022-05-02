package sql

import (
	"encoding"
	"strconv"

	goqueue "github.com/antonio-alexander/go-queue"

	"github.com/pkg/errors"
)

func convertSingle(item interface{}) (goqueue.Bytes, error) {
	switch v := item.(type) {
	default:
		return nil, errors.Errorf(ErrUnsupportedTypef, v)
	case struct{}:
		return goqueue.Bytes("{}"), nil
	case goqueue.Bytes:
		return v, nil
	case []byte:
		return v, nil
	case encoding.BinaryMarshaler:
		bytes, err := v.MarshalBinary()
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}
}

func convertMultiple(items []interface{}) ([]goqueue.Bytes, error) {
	var bytes []goqueue.Bytes
	for _, item := range items {
		b, err := convertSingle(item)
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, b)
	}
	return bytes, nil
}

//ConfigFromEnv can be used to generate a configuration pointer
// from a list of environments, it'll set the default configuraton
// as well
func ConfigFromEnv(envs map[string]string) *Configuration {
	c := &Configuration{
		Hostname:     DefaultDatabaseHost,
		Port:         DefaultDatabasePort,
		Database:     DefaultDatabaseName,
		Table:        DefaultDatabaseTable,
		DriverName:   DefaultDriverName,
		Username:     DefaultUsername,
		Password:     DefaultPassword,
		ParseTime:    DefaultParseTime,
		CreateTable:  DefaultCreateTable,
		QueryTimeout: DefaultQueryTimeout,
	}
	if hostname, ok := envs["HOSTNAME"]; ok {
		c.Hostname = hostname
	}
	if port, ok := envs["PORT"]; ok {
		c.Port = port
	}
	if username, ok := envs["USERNAME"]; ok {
		c.Username = username
	}
	if password, ok := envs["PASSWORD"]; ok {
		c.Password = password
	}
	if database, ok := envs["DATABASE"]; ok {
		c.Database = database
	}
	if driver, ok := envs["DRIVER"]; ok {
		c.DriverName = driver
	}
	if s, ok := envs["CREATE_TABLE"]; ok {
		if createTable, err := strconv.ParseBool(s); err != nil {
			c.CreateTable = createTable
		}
	}
	return c
}
