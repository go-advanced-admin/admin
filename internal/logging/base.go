package logging

import "time"

type LogEntry struct {
	ID          interface{}
	ActionTime  time.Time
	UserID      interface{}
	UserRepr    string
	ContentType string
	ObjectID    interface{}
	ObjectRepr  string
	ActionFlag  LogStoreLevel
	Message     string
}

func (l *LogEntry) Repr() interface{} {
	if l.ObjectRepr != "" {
		return l.ObjectRepr
	}
	return l.Message
}

type LogStore interface {
	InsertLogEntry(logEntry *LogEntry) error
	GetLogEntry(id interface{}) (*LogEntry, error)
	GetLogEntries() ([]*LogEntry, error)
}
