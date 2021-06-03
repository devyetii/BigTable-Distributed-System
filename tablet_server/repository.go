package main

import (
	"sort"
	"strings"
)

type BigTableEntry map[string]string
type BigTablePartition map[string]BigTableEntry

var data BigTablePartition = BigTablePartition{
	"a.person.com" : BigTableEntry{
		"name" : "ebrahim",
		"age" : "22",
	},
    "b.person.com" : BigTableEntry{
        "name" : "Farha",
        "age" : "23",
    },
    "b.persona.com" : BigTableEntry{
        "name" : "Mahmoud",
        "age" : "21",
    },
};

var keys []string = []string{
    "a.person.com",
    "b.person.com",
    "b.persona.com",
}

func insertKeySorted(key string) {
    idx := sort.SearchStrings(keys, key)
	
    // Expand capacity
    keys = append(keys, "")
	
    // Shift (O(n-idx))
    copy(keys[idx+1:], keys[idx:len(keys)-1])
	
    // Insert
    keys[idx] = key

}

func getByRange(from string, to string) BigTablePartition {
    if (strings.Compare(from, to) > 0) {
        return nil
    }

    i, f := sort.SearchStrings(keys, from), sort.SearchStrings(keys, to)
    if  f == len(keys) {
        f = len(keys) - 1
    }

    entries := make(BigTablePartition)
    for ; i <= f; i++ {
        entries[keys[i]] = data[keys[i]]
    }
    return entries
}

func getByKeysList(keys []string) BigTablePartition {
    entries := make(BigTablePartition)
    for _, k := range keys {
        v, ok := data[k]
        if (ok) {
            entries[k] = v
        }
    }
    return entries
}

func addRow(row_key string, cols BigTableEntry) BigTableEntry {
	_, ok := data[row_key]
    if (ok) {
        return nil
    }

	data[row_key] = cols
    insertKeySorted(row_key)
	return data[row_key]
}

func setCells(row_key string, cols BigTableEntry) BigTableEntry {
    _, ok := data[row_key]
    if (!ok) {
        return nil
    }

    for k, v := range cols {
        data[row_key][k] = v
    }
    return data[row_key]
}

func deleteCells(row_key string, col_keys []string) BigTableEntry {
    _, ok := data[row_key]
    if (!ok) {
        return nil
    }

    for _, col := range col_keys {
        delete(data[row_key], col)
    }
    return data[row_key]
}

func deleteRow(row_key string) bool {
    _, ok := data[row_key]
    if (!ok) {
        return false
    }

	delete(data, row_key)
	return true
}