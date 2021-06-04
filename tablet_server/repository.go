package main

import (
	"fmt"
	"sort"
)

func (repo *Repository) keysLowerBoundComparator(key RowKeyType) func(int) bool {
    return func(i int) bool { return repo.keys[i].LowerBound(key) }
}
type Repository struct {
    data BigTablePartition
    keys []RowKeyType
    updateLogsFile *SafeUpdateLog
}

func (repo *Repository) insertKeySorted(key RowKeyType) {
    idx := sort.Search(len(keys), repo.keysLowerBoundComparator(key))
	
    // Expand capacity
    repo.keys = append(repo.keys, key)
	
    // Shift (O(n-idx))
    copy(repo.keys[idx+1:], repo.keys[idx:len(repo.keys)-1])
	
    // Insert
    repo.keys[idx] = key
}

func (repo *Repository) deleteKeySorted(key RowKeyType) {
    idx := sort.Search(len(repo.keys), repo.keysLowerBoundComparator(key))
    
    copy(repo.keys[idx:], repo.keys[idx+1:])

    repo.keys = repo.keys[:len(repo.keys)-1]
}

func (repo *Repository) checkRowExists(row_key RowKeyType) bool {
    _, ok := repo.data[row_key]
    return ok
}

func (repo *Repository) getByRange(from RowKeyType, to RowKeyType) BigTablePartition {
    if (from.UpperBound(to)) {
        return nil
    }

    entries := make(BigTablePartition)
    for i := sort.Search(len(repo.keys), repo.keysLowerBoundComparator(from)) ; i < len(repo.keys) && to.LowerBound(repo.keys[i]); i++ {
        entries[repo.keys[i]] = data[repo.keys[i]]
    }
    return entries
}

func (repo *Repository) getByKeysList(keys []RowKeyType) BigTablePartition {
    entries := make(BigTablePartition)
    for _, k := range keys {
        v, ok := repo.data[k]
        if (ok) {
            entries[k] = v
        }
    }
    return entries
}

func (repo *Repository) addRow(row_key RowKeyType, cols BigTableEntry) BigTableEntry {
    if (repo.checkRowExists(row_key)) {
        return nil
    }

	repo.data[row_key] = cols
    repo.insertKeySorted(row_key)
    fmt.Println(repo.keys)
    repo.updateLogsFile.LogAddRow(row_key, cols)
	return repo.data[row_key]
}

func (repo *Repository) setCells(row_key RowKeyType, cols BigTableEntry) BigTableEntry {
    if (!repo.checkRowExists(row_key)) {
        return nil
    }

    for k, v := range cols {
        repo.data[row_key][k] = v
    }
    repo.updateLogsFile.LogSetCells(row_key, cols)
    return repo.data[row_key]
}

func (repo *Repository) deleteCells(row_key RowKeyType, col_keys []ColKeyType) BigTableEntry {
    if (!repo.checkRowExists(row_key)) {
        return nil
    }

    for _, col := range col_keys {
        delete(repo.data[row_key], col)
    }
    repo.updateLogsFile.LogDeleteCells(row_key, col_keys)
    return repo.data[row_key]
}

func (repo *Repository) deleteRow(row_key RowKeyType) bool {
    if (!repo.checkRowExists(row_key)) {
        return false
    }

    repo.deleteKeySorted(row_key)
    delete(repo.data, row_key)
    repo.updateLogsFile.LogDeleteRow(row_key)
	return true
}