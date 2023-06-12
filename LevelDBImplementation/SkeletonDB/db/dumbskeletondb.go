package db

import (
	"errors"
	"fmt"
	"sort"
)

type DumbSkeletonDB struct {
	data_map map[string]string
}

func CreateDumbSkeletonDB() *DumbSkeletonDB {
	db := new(DumbSkeletonDB)
	db.data_map = make(map[string]string, 0)

	return db
}

func (db *DumbSkeletonDB) Get(key []byte) (value []byte, err error) {
	// Convert byte to string because we can't make a map type with slices
	key_string := string(key[:])

	if value, found := db.data_map[key_string]; found {
		return []byte(value), nil
	} else {
		return nil, errors.New("Could not find the key in the database!")
	}
}

func (db *DumbSkeletonDB) Has(key []byte) (ret bool, err error) {
	// Convert byte to string because we can't make a map type with slices
	key_string := string(key[:])

	if _, found := db.data_map[key_string]; found {
		return true, nil
	} else {
		return false, errors.New("Could not find the key in the database!")
	}
}

func (db *DumbSkeletonDB) Put(key, value []byte) error {
	// Need to protect against empty key value byte arrays
	if len(key) == 0 {
		return errors.New("Cannot add an empty key to the database!")
	}

	if len(value) == 0 {
		return errors.New("Cannot add an empty value to the database!")
	}

	// Convert byte to string because we can't make a map type with slices
	key_string := string(key[:])
	value_string := string(value[:])

	db.data_map[key_string] = value_string

	return nil
}

func (db *DumbSkeletonDB) Delete(key []byte) error {
	found, err := db.Has(key)

	if err != nil {
		return errors.New("Cannot delete the key because it doesn't exist in the database!")
	}

	if !found {
		return errors.New("Cannot delete the key because it doesn't exist in the database!")
	}

	// Convert byte to string because we can't make a map type with slices
	key_string := string(key[:])

	delete(db.data_map, key_string)

	return nil
}

func (db *DumbSkeletonDB) RangeScan(start, limit []byte) (Iterator, error) {
	// Convert byte to string because we can't make a map type with slices
	start_string := string(start[:])
	end_string := string(limit[:])

	if start_string > end_string {
		return nil, errors.New("Cannot return an iterator when the start is greater than the limit!")
	}

	key_list := make([]string, 0)
	value_list := make([]string, 0)

	for key, _ := range db.data_map {
		key_list = append(key_list, key)
	}

	sort.Strings(key_list)

	start_index := 0
	end_index := 0

	for index, key := range key_list {
		if start_string >= key {
			start_index = index
		}

		if end_string >= key {
			end_index = index
		}
	}

	key_list = key_list[start_index : end_index+1]

	for _, key := range key_list {
		value_list = append(value_list, db.data_map[key])
	}

	return CreateDumbSkeletonDBIterator(key_list, value_list), nil
}

func (db *DumbSkeletonDB) Print() {
	for k, v := range db.data_map {
		fmt.Println("key: " + k + " value: " + v)
	}
}

// Terrible implementation since this wastes a bunch of memory, but it is one of the simplest
// especially due to the fact that you can't have pointers to keys in a map or have then
// in any sort of ordered structure
type DumbSkeletonDBIterator struct {
	index           int
	iterator_errors error
	key_array       []string
	value_array     []string
}

func CreateDumbSkeletonDBIterator(keys, values []string) *DumbSkeletonDBIterator {
	iterator := new(DumbSkeletonDBIterator)
	iterator.index = 0
	iterator.iterator_errors = nil
	iterator.key_array = keys
	iterator.value_array = values

	return iterator
}

func (iterator *DumbSkeletonDBIterator) Next() bool {
	if iterator.index == len(iterator.key_array)-1 {
		return false
	} else {
		iterator.index++
		return true
	}
}

func (iterator *DumbSkeletonDBIterator) Error() error {
	if iterator.iterator_errors == nil {
		return nil
	} else {
		return iterator.iterator_errors
	}
}

func (iterator *DumbSkeletonDBIterator) Key() []byte {
	if iterator.index == len(iterator.key_array) {
		if iterator.iterator_errors == nil {
			iterator.iterator_errors = errors.New("Cannot get the key as the iterator is exhausted!")
		} else {
			iterator.iterator_errors = fmt.Errorf("%w; Cannot get the key as the iterator is exhausted!", iterator.iterator_errors)
		}
		return nil
	} else {
		return []byte(iterator.key_array[iterator.index])
	}
}

func (iterator *DumbSkeletonDBIterator) Value() []byte {
	if iterator.index == len(iterator.key_array) {
		if iterator.iterator_errors == nil {
			iterator.iterator_errors = errors.New("Cannot get the value as the iterator is exhausted!")
		} else {
			iterator.iterator_errors = fmt.Errorf("%w; Cannot get the value as the iterator is exhausted!", iterator.iterator_errors)
		}
		return nil
	} else {
		return []byte(iterator.value_array[iterator.index])
	}
}
