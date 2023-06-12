package writeaheadlog

type LogEntry struct {
	key_length   int64  // 8 bytes
	key          []byte // variable bytes
	value_length int64  // 8 bytes
	value        []byte // variable bytes
}

// Format for the WriteAhead Logs.
// Only need to store write operations within the write ahead log one after another.
// Furthermore, "delete" operations can be added by adding a log with a value_length of 0 bytes.
type WriteAheadLog struct {
	data []LogEntry
}

func InitializeWriteAheadLog() *WriteAheadLog {
	log := new(WriteAheadLog)
	log.data = make([]LogEntry, 0)

	return log
}

func (l *WriteAheadLog) AddEntry(key []byte, value []byte) error {
	key_length := len(key)
	value_length := len(value)

	log_entry := LogEntry{
		key_length:   int64(key_length),
		key:          key,
		value_length: int64(value_length),
		value:        value,
	}
	l.data = append(l.data, log_entry)

	return nil
}
