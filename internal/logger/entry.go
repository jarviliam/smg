package logger

import (
	"errors"
	"time"
)

type LogEntry struct {
	at      time.Time
	message string
	error   error
}

var (
	ErrLogEntryVoidMsg = errors.New("Message cannot be empty")
)

func NewLogEntry(message string) *LogEntry {
	if message == "" {
		return nil
	}
	e := &LogEntry{
		at:      time.Now(),
		message: message,
		error:   nil,
	}
	return e
}

func (e *LogEntry) At() time.Time {
	return e.at
}

func (e *LogEntry) Message() string {
	return e.message
}

func (e *LogEntry) Error() error {
	return e.error
}

func (e *LogEntry) SetTime(t time.Time) {
	e.at = t
}

func (e *LogEntry) SetMessage(message string) error {
	if message == "" {
		return ErrLogEntryVoidMsg
	}
	e.message = message
	if e.IsError() {
		e.error = nil
	}
	return nil
}

func (e *LogEntry) SetError(error error) {
	e.error = error
	if e.message != "" {
		e.message = ""
	}
}

func (e *LogEntry) IsError() bool {
	return e.error != nil
}
