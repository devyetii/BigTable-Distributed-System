package main

import "strconv"

type RowKeyType int
type ColKeyType string
type ValType interface{}
type BigTableEntry map[ColKeyType]ValType
type BigTablePartition map[RowKeyType]BigTableEntry

// If you change the row key type, you only have to change the compare func
func (first RowKeyType) LowerBound(second RowKeyType) bool {
    return int(first) >= int(second)
}

func (first RowKeyType) UpperBound(second RowKeyType) bool {
    return int(first) > int(second)
}

func RowKeyFromString(val string) (RowKeyType, error) {
    i, e := strconv.Atoi(val)
    return RowKeyType(i), e
}

func MapStringsToRowKeys(in []string, f func (string) (RowKeyType, error)) ([]RowKeyType) {
	var out []RowKeyType

	for _,v := range in {
        rk, err := f(v)
        
        if (err != nil) {
            return nil
        }

		out = append(out, rk)
	}
	return out
}