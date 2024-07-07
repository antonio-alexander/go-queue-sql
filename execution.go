package sql

import (
	"encoding"

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

func convertMultipleWithPriority(items []interface{}, priorities ...int) ([]priorityBytes, error) {
	bytes := make([]priorityBytes, 0, len(items))
	if len(priorities) != len(items) {
		priority := DefaultPriority
		if len(priorities) > 0 {
			priority = priorities[0]
		}
		bytes = append(bytes, priorityBytes{
			priority: priority,
		})
	}
	for i, item := range items {
		b, err := convertSingle(item)
		if err != nil {
			return nil, err
		}
		bytes[i].bytes = b
	}
	return bytes, nil
}
