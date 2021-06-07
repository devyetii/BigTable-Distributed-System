package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type SafeUpdateLog struct {
	mu sync.Mutex
	file *os.File
}

func (l *SafeUpdateLog) GetFileForRead() *os.File {
	_, err := l.file.Seek(0,0)
	justPrintErr(err)
	return l.file
}

func (l *SafeUpdateLog) LogAddRow(row_key RowKeyType, cols BigTableEntry) {
	l.mu.Lock() // Lock acquire attempt (Blocking)
	defer l.mu.Unlock()

	l.file.WriteString(fmt.Sprintf("add_row %v\n", row_key))
	for k, v := range cols {
		l.file.WriteString(fmt.Sprintf("set_cell %v %v %v\n", row_key, k, v))
	}
}

func (l *SafeUpdateLog) LogDeleteRow(row_key RowKeyType) {
	l.mu.Lock() // Lock acquire attempt (Blocking)
	defer l.mu.Unlock()

	l.file.WriteString(fmt.Sprintf("delete_row %v\n", row_key))
}

func (l *SafeUpdateLog) LogSetCells(row_key RowKeyType, cols BigTableEntry) {
	l.mu.Lock() // Lock acquire attempt (Blocking)
	defer l.mu.Unlock()

	for k, v := range cols {
		l.file.WriteString(fmt.Sprintf("set_cell %v %v %v\n", row_key, k, v))
	}
}

func (l *SafeUpdateLog) LogDeleteCells(row_key RowKeyType, col_keys []ColKeyType) {
	l.mu.Lock() // Lock acquire attempt (Blocking)
	defer l.mu.Unlock()

	for _,k := range col_keys {
		l.file.WriteString(fmt.Sprintf("delete_cell %v %v\n", row_key, k))
	}
}

func (l *SafeUpdateLog) ClearLogs() {
	l.file.Close()
	update_logs_file, err := os.OpenFile("updates.log", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	if (err != nil) {
		log.Println(fmt.Sprintf("Error in clearing logs: %v", err))
	}
	l.file = update_logs_file
}