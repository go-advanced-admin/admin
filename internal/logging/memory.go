package logging

import (
	"fmt"
)

type InMemoryLogStore struct {
	logEntryMap map[string]*LogEntry
	logIDs      []string
	maxLength   uint
}

func NewInMemoryLogStore(maxLength uint) *InMemoryLogStore {
	return &InMemoryLogStore{
		logEntryMap: make(map[string]*LogEntry),
		logIDs:      make([]string, 0),
		maxLength:   maxLength,
	}
}

func (store *InMemoryLogStore) InsertLogEntry(log *LogEntry) error {
	logID := fmt.Sprint(log.ID)
	if _, exists := store.logEntryMap[logID]; exists {
		return fmt.Errorf("log entry with ID %s already exists", logID)
	}

	if len(store.logIDs) >= int(store.maxLength) {
		lastID := store.logIDs[len(store.logIDs)-1]
		delete(store.logEntryMap, lastID)
		store.logIDs = store.logIDs[:len(store.logIDs)-1]
	}

	store.logIDs = append([]string{logID}, store.logIDs...)
	store.logEntryMap[logID] = log

	return nil
}

func (store *InMemoryLogStore) GetLogEntry(id interface{}) (*LogEntry, error) {
	logID := fmt.Sprint(id)

	logEntry, exists := store.logEntryMap[logID]
	if !exists {
		return nil, nil
	}
	return logEntry, nil
}

func (store *InMemoryLogStore) GetLogEntries() ([]*LogEntry, error) {
	logEntries := make([]*LogEntry, len(store.logIDs))
	for i, logID := range store.logIDs {
		logEntries[i] = store.logEntryMap[logID]
	}
	return logEntries, nil
}
