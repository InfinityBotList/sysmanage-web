package main

import (
	"encoding/json"
	"sync"
	"time"
)

const logPrefix = "rlogentry:"
const markerDoneSuffix = "_marker"
const logTime = 24 * time.Hour

var inDeploy = sync.Mutex{}

type autoLogger struct {
	id    string
	Error bool
}

func (a autoLogger) Write(p []byte) (n int, err error) {
	if a.Error {
		addToLog(a.id, "ERROR: "+string(p), false)
	} else {
		addToLog(a.id, string(p), false)
	}

	return len(p), nil
}

func markLogDone(id string) {
	rdb.Set(ctx, logPrefix+id+markerDoneSuffix, "1", logTime)
}

func addToLog(id string, data string, newline bool) error {
	if newline {
		data += "\n"
	}

	currLog := rdb.Get(ctx, logPrefix+id).Val()

	if currLog == "" {
		currLog = "[]"
	}

	var logs []string

	err := json.Unmarshal([]byte(currLog), &logs)

	if err != nil {
		return err
	}

	logs = append(logs, data)

	newLog, err := json.Marshal(logs)

	if err != nil {
		return err
	}

	rdb.Set(ctx, logPrefix+id, newLog, logTime)

	return nil
}
