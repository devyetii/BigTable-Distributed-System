package main

import (
	"fmt"
	"os"
	"sync"
)

type SafeUpdateLog struct {
	mu sync.Mutex
	file os.File
}

func (l *SafeUpdateLog) LogAddRow(row_key string, cols BigTableEntry) {
	l.mu.Lock() // Lock acquire attempt (Blocking)
	defer l.mu.Unlock()

	l.file.WriteString(fmt.Sprintf("add_row %v\n", row_key))
	for k, v := range cols {
		l.file.WriteString(fmt.Sprintf("set_cell %v %v %v\n", row_key, k, v))
	}
}

func (l *SafeUpdateLog) LogDeleteRow(row_key string) {
	l.mu.Lock() // Lock acquire attempt (Blocking)
	defer l.mu.Unlock()

	l.file.WriteString(fmt.Sprintf("delete_row %v\n", row_key))
}

func (l *SafeUpdateLog) LogSetCells(row_key string, cols BigTableEntry) {
	l.mu.Lock() // Lock acquire attempt (Blocking)
	defer l.mu.Unlock()

	for k, v := range cols {
		l.file.WriteString(fmt.Sprintf("set_cell %v %v %v\n", row_key, k, v))
	}
}

func (l *SafeUpdateLog) LogDeleteCells(row_key string, col_keys []string) {
	l.mu.Lock() // Lock acquire attempt (Blocking)
	defer l.mu.Unlock()

	for _,k := range col_keys {
		l.file.WriteString(fmt.Sprintf("delete_cell %v %v\n", row_key, k))
	}
}