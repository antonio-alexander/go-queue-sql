package example

import (
	"encoding/json"

	goqueue "github.com/antonio-alexander/go-queue"
)

type Data struct {
	Int    int    `json:"int"`
	String string `json:"string"`
}

func (d *Data) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Data) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, d)
}

func Dequeue(queue goqueue.Dequeuer, n ...int) (*Data, bool) {
	item, overflow := queue.Dequeue()
	if overflow {
		return nil, true
	}
	switch v := item.(type) {
	default:
		return nil, true
	case *Data:
		return v, false
	case []byte:
		data := &Data{}
		if err := data.UnmarshalBinary(v); err != nil {
			return nil, true
		}
		return data, false
	}
}
