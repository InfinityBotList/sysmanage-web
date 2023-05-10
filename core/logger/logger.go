package logger

import (
	"time"
)

const logTime = 8 * time.Hour
const logTreshold = 100000

type LogEntry struct {
	LastUpdate time.Time
	LastLog    []string
	IsDone     bool
}

type LogEntryMap map[string]LogEntry

func (l LogEntryMap) Get(id string) LogEntry {
	entry, ok := l[id]

	if !ok {
		return LogEntry{
			LastLog: []string{},
		}
	}

	if time.Since(entry.LastUpdate) > logTime {
		delete(l, id)
		return LogEntry{}
	}

	if len(entry.LastLog) > logTreshold {
		delete(l, id) // Prevent memory leaks
		return LogEntry{}
	}

	return entry
}

func (l LogEntryMap) Set(id string, entry LogEntry) {
	l[id] = entry
}

func (l LogEntryMap) Add(id string, data string, newline bool) {
	if newline {
		data += "\n"
	}

	currLog := l.Get(id)

	currLog.LastUpdate = time.Now()
	currLog.LastLog = append(currLog.LastLog, data)

	l.Set(id, currLog)
}

func (l LogEntryMap) MarkDone(id string) {
	entry := l.Get(id)

	entry.IsDone = true

	l.Set(id, entry)
}

var LogMap = LogEntryMap{}

type AutoLogger struct {
	ID      string
	Error   bool
	Newline bool
}

func (a AutoLogger) Write(p []byte) (n int, err error) {
	if a.Error {
		LogMap.Add(a.ID, "ERROR: "+string(p), a.Newline)
	} else {
		LogMap.Add(a.ID, string(p), a.Newline)
	}

	return len(p), nil
}
