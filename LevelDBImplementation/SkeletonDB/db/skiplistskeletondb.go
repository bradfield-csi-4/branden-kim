package db

import (
	"errors"
	"os"

	skiplist "skeletondb/utils"
	writeaheadlog "skeletondb/writeaheadlog"
)

type SkipListSkeletonDB struct {
	skip_list       *skiplist.SkipList
	write_ahead_log *writeaheadlog.WriteAheadLog
	log_file        *os.File
}

func CreateSkipListSkeletonDB(log_file *os.File) *SkipListSkeletonDB {
	db := new(SkipListSkeletonDB)
	db.skip_list = skiplist.InitializeSkipList()
	db.write_ahead_log = writeaheadlog.InitializeWriteAheadLog()
	db.log_file = log_file

	return db
}

func (db *SkipListSkeletonDB) Get(key []byte) (value []byte, err error) {
	// Convert byte to string
	key_string := string(key[:])

	node, search_err := db.skip_list.SearchNode(key_string)

	if search_err != nil {
		return nil, search_err
	} else {
		return []byte(node.Value), nil
	}
}

func (db *SkipListSkeletonDB) Has(key []byte) (ret bool, err error) {
	// Convert byte to string
	key_string := string(key[:])

	_, search_err := db.skip_list.SearchNode(key_string)

	if search_err != nil {
		return true, nil
	} else {
		return false, search_err
	}
}

func (db *SkipListSkeletonDB) Put(key, value []byte) error {
	// Need to protect against empty key value byte arrays
	if len(key) == 0 {
		return errors.New("Cannot add an empty key to the database!")
	}

	if len(value) == 0 {
		return errors.New("Cannot add an empty value to the database!")
	}

	// Convert byte to string
	key_string := string(key[:])
	value_string := string(value[:])

	add_err := db.skip_list.AddNode(key_string, value_string)

	if add_err != nil {
		return add_err
	} else {

		// If the write is successful to the SkipList, add it to the Write Ahead Log as well
		write_log_err := db.write_ahead_log.AddEntry(key, value)
		if write_log_err != nil {

			// Don't want to keep this entry if the write to the log failed
			db.Delete(key)
			return write_log_err
		}

		return nil
	}
}

func (db *SkipListSkeletonDB) Delete(key []byte) error {
	// Convert byte to string
	key_string := string(key[:])

	delete_err := db.skip_list.DeleteNode(key_string)

	if delete_err != nil {
		return delete_err
	} else {
		return nil
	}
}

func (db *SkipListSkeletonDB) RangeScan(start, limit []byte) (Iterator, error) {
	// Convert byte to string
	start_string := string(start[:])
	end_string := string(limit[:])

	if start_string > end_string {
		return nil, errors.New("Cannot return an iterator when the start is greater than the limit!")
	}

	start_node, err := db.skip_list.SearchNode(start_string)

	if err != nil {
		return nil, err
	}

	end_node, err := db.skip_list.SearchNode(end_string)

	if err != nil {
		return nil, err
	}

	return CreateSkipListDBIterator(start_node, end_node), nil
}

// For the SkipList implementation of the DB, the iterator should just have a reference to the node at level 1 that
// it should start with and a reference to the node at level 1 that it ends with
type SkipListDBIterator struct {
	current         *skiplist.SkipNode
	end             *skiplist.SkipNode
	iterator_errors error
}

func CreateSkipListDBIterator(start, end *skiplist.SkipNode) *SkipListDBIterator {
	iterator := new(SkipListDBIterator)
	iterator.current = start
	iterator.end = end

	return iterator
}

func (iterator *SkipListDBIterator) Next() bool {
	if iterator.current == iterator.end {
		return false
	} else {
		iterator.current = iterator.current.Level_List[0]
		return true
	}
}

func (iterator *SkipListDBIterator) Error() error {
	if iterator.iterator_errors == nil {
		return nil
	} else {
		return iterator.iterator_errors
	}
}

func (iterator *SkipListDBIterator) Key() []byte {
	return []byte(iterator.current.Key)
}

func (iterator *SkipListDBIterator) Value() []byte {
	return []byte(iterator.current.Value)
}
