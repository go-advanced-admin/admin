package logging

import (
	"fmt"
	"github.com/google/uuid"
)

type InMemoryLogStore struct {
	logEntryMap map[uuid.UUID]*LogEntry
	logIDs      []uuid.UUID
	maxLength   uint
}

func NewInMemoryLogStore(maxLength uint) *InMemoryLogStore {
	return &InMemoryLogStore{
		logEntryMap: make(map[uuid.UUID]*LogEntry),
		logIDs:      make([]uuid.UUID, 0),
		maxLength:   maxLength,
	}
}

func (store *InMemoryLogStore) InsertLogEntry(log *LogEntry) error {
	logID, ok := log.ID.(uuid.UUID)
	if !ok {
		return fmt.Errorf("invalid log ID type: %T", log.ID)
	}
	if _, exists := store.logEntryMap[logID]; exists {
		return fmt.Errorf("log entry with ID %s already exists", logID)
	}

	if len(store.logIDs) >= int(store.maxLength) {
		lastID := store.logIDs[len(store.logIDs)-1]
		delete(store.logEntryMap, lastID)
		store.logIDs = store.logIDs[:len(store.logIDs)-1]
	}

	store.logIDs = append([]uuid.UUID{logID}, store.logIDs...)
	store.logEntryMap[logID] = log

	println(store.GetLogEntries())

	return nil
}

func (store *InMemoryLogStore) GetLogEntry(id interface{}) (*LogEntry, error) {
	logID, ok := id.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("invalid log ID type: %T", id)
	}

	logEntry, exists := store.logEntryMap[logID]
	if !exists {
		return nil, fmt.Errorf("log entry with ID %s not found", logID)
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
