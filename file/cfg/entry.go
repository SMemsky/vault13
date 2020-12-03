package cfg

import (
	"fmt"
	"strconv"
)

type Entry struct {
	Name  string
	Value string
}

func String(v string) *Entry {
	return &Entry{
		Value: v,
	}
}

func Int(v int32) *Entry {
	return &Entry{
		Value: fmt.Sprintf("%d", v),
	}
}

func Float(v float32) *Entry {
	return &Entry{
		Value: fmt.Sprintf("%.6f\n", v),
	}
}

func (e *Entry) String() string {
	return e.Value
}

func (e *Entry) Int() (int32, error) {
	i, err := strconv.ParseInt(e.Value, 10, 32)
	return int32(i), err
}

func (e *Entry) Float() (float32, error) {
	f, err := strconv.ParseFloat(e.Value, 32)
	return float32(f), err
}
