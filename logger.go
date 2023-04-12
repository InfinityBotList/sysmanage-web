package main

import (
	"encoding/json"
	"time"
)

const logPrefix = "smdl:"
const logTime = 8 * time.Hour

type LogEntry struct {
	LastUpdate time.Time
	LastLog    []string
	Valid      bool
	IsDone     bool
}

type LogEntryMap map[string]LogEntry

func (l LogEntryMap) Get(id string) LogEntry {
	entry, ok := l[id]

	if !ok {
		// Check redis
		currLog := rdb.Get(ctx, logPrefix+id).Val()

		if currLog == "" {
			return LogEntry{
				Valid:   true,
				LastLog: []string{},
			}
		}

		var logs *LogEntry

		err := json.Unmarshal([]byte(currLog), &logs)

		if err != nil {
			return LogEntry{
				Valid:   true,
				LastLog: []string{err.Error() + " " + currLog},
			}
		}

		// Has to be done as server restarted
		logs.IsDone = true

		l.Set(id, *logs)

		return *logs
	}

	if time.Since(entry.LastUpdate) > logTime {
		delete(l, id)
		return LogEntry{}
	}

	if !entry.Valid {
		delete(l, id)
		return LogEntry{
			LastLog: []string{"Log entry was invalid, please try again."},
		}
	}

	return entry
}

func (l LogEntryMap) Set(id string, entry LogEntry) {
	l[id] = entry
}

// Persist will persist the current state of the log entry to redis
// overwriting the old one
func (l LogEntryMap) Persist(id string) error {
	// First get the entry itself
	entry := l[id]

	// Load in redis
	newLog, err := json.Marshal(entry)

	if err != nil {
		return err
	}

	err = rdb.Set(ctx, logPrefix+id, newLog, logTime).Err()

	if err != nil {
		return err
	}

	return nil
}

func (l LogEntryMap) Add(id string, data string, newline bool) {
	if newline {
		data += "\n"
	}

	currLog := l.Get(id)

	if !currLog.Valid {
		currLog = LogEntry{
			LastUpdate: time.Now(),
			LastLog: []string{
				data,
			},
		}

		l.Set(id, currLog)
		return
	}

	currLog.LastLog = append(currLog.LastLog, data)

	l.Set(id, currLog)
}

func (l LogEntryMap) MarkDone(id string) {
	entry := l.Get(id)

	entry.IsDone = true

	l.Set(id, entry)
}

var logMap = LogEntryMap{}

type autoLogger struct {
	id      string
	Error   bool
	Newline bool
}

func (a autoLogger) Write(p []byte) (n int, err error) {
	if a.Error {
		logMap.Add(a.id, "ERROR: "+string(p), a.Newline)
	} else {
		logMap.Add(a.id, string(p), a.Newline)
	}

	return len(p), nil
}
