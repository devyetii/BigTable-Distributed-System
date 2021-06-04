package main

import (
	"sort"
	"strings"
)

type BigTableEntry map[string]string
type BigTablePartition map[string]BigTableEntry

type Repository struct {
    data BigTablePartition
    keys []string
    updateLogsFile *SafeUpdateLog
}

func (repo *Repository) insertKeySorted(key string) {
    idx := sort.SearchStrings(repo.keys, key)
	
    // Expand capacity
    repo.keys = append(repo.keys, "")
	
    // Shift (O(n-idx))
    copy(repo.keys[idx+1:], repo.keys[idx:len(repo.keys)-1])
	
    // Insert
    repo.keys[idx] = key
}

func (repo *Repository) getByRange(from string, to string) BigTablePartition {
    if (strings.Compare(from, to) > 0) {
        return nil
    }

    i, f := sort.SearchStrings(repo.keys, from), sort.SearchStrings(repo.keys, to)
    if  f == len(repo.keys) {
        f = len(repo.keys) - 1
    }

    entries := make(BigTablePartition)
    for ; i <= f; i++ {
        entries[repo.keys[i]] = repo.data[repo.keys[i]]
    }
    return entries
}

func (repo *Repository) getByKeysList(keys []string) BigTablePartition {
    entries := make(BigTablePartition)
    for _, k := range keys {
        v, ok := repo.data[k]
        if (ok) {
            entries[k] = v
        }
    }
    return entries
}

func (repo *Repository) addRow(row_key string, cols BigTableEntry) BigTableEntry {
	_, ok := repo.data[row_key]
    if (ok) {
        return nil
    }

	repo.data[row_key] = cols
    repo.insertKeySorted(row_key)
    repo.updateLogsFile.LogAddRow(row_key, cols)
	return repo.data[row_key]
}

func (repo *Repository) setCells(row_key string, cols BigTableEntry) BigTableEntry {
    _, ok := repo.data[row_key]
    if (!ok) {
        return nil
    }

    for k, v := range cols {
        repo.data[row_key][k] = v
    }
    repo.updateLogsFile.LogSetCells(row_key, cols)
    return repo.data[row_key]
}

func (repo *Repository) deleteCells(row_key string, col_keys []string) BigTableEntry {
    _, ok := repo.data[row_key]
    if (!ok) {
        return nil
    }

    for _, col := range col_keys {
        delete(repo.data[row_key], col)
    }
    repo.updateLogsFile.LogDeleteCells(row_key, col_keys)
    return repo.data[row_key]
}

func (repo *Repository) deleteRow(row_key string) bool {
    _, ok := repo.data[row_key]
    if (!ok) {
        return false
    }

	delete(repo.data, row_key)
    repo.updateLogsFile.LogDeleteRow(row_key)
	return true
}