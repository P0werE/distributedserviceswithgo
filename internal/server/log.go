package server

import (
	"errors"
	"fmt"
	"sync"
)

var ErrOffsetNotFound = errors.New("offset not found")

type (
	Log struct {
		mu      *sync.RWMutex
		records []Record
	}

	Record struct {
		Value  []byte `json:"value"`
		Offset uint64 `json:"offset"`
	}
)

func NewLog() *Log {
	return &Log{mu: new(sync.RWMutex), records: make([]Record, 0, 100)}
}

func (l *Log) Append(record Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	newRecord := record

	newRecord.Value = make([]byte, len(record.Value))
	copy(newRecord.Value, record.Value)

	newRecord.Offset = uint64(len(l.records))
	l.records = append(l.records, newRecord)

	return newRecord.Offset, nil
}

func (l *Log) Read(offset uint64) (Record, error) {
	record := Record{}

	if offset >= uint64(len(l.records)) {
		return record, fmt.Errorf("%d %w", offset, ErrOffsetNotFound)
	}

	l.mu.RLock()
	defer l.mu.RUnlock()
	record.Offset = offset
	v := l.records[offset].Value
	record.Value = make([]byte, len(v))
	copy(record.Value, v)

	return record, nil
}
