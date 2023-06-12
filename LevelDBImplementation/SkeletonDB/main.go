package main

import (
	"fmt"
	"os"
	"skeletondb/db"
)

// func main() {
// 	fmt.Println("Hello SkeletonDB")

// 	db, err := leveldb.OpenFile("./test", nil)

// 	if err != nil {
// 		fmt.Printf("Could not open the directory to the database!")
// 		os.Exit(1)
// 	}

// 	put_err := db.Put([]byte("testkey"), []byte("testvalue"), nil)

// 	if put_err != nil {
// 		fmt.Printf("Could not add a key value pair to the database!")
// 		os.Exit(1)
// 	}

// 	data, get_err := db.Get([]byte("testkey"), nil)

// 	if get_err != nil {
// 		fmt.Printf("Could not get the key value pair from the database!")
// 		os.Exit(1)
// 	}

// 	fmt.Println(string(data[:]))

// 	defer db.Close()
// }

func main() {
	fmt.Println("Hello from SkeletonDB")

	cwd, cwd_err := os.Getwd()

	if cwd_err != nil {
		fmt.Println(cwd_err.Error())
		os.Exit(1)
	}

	f, err := os.OpenFile(cwd+"/writeahead.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()

	db := db.CreateSkipListSkeletonDB(f)

	put_err := db.Put([]byte("testkey"), []byte("testvalue"))

	if put_err != nil {
		fmt.Println(put_err.Error())
		os.Exit(1)
	}

	// data, get_err := db.Get([]byte("testkey"))

	// if get_err != nil {
	// 	fmt.Println(get_err.Error())
	// 	os.Exit(1)
	// }

	// fmt.Println(string(data[:]))

	// has_key, has_err := db.Has([]byte("testkey"))

	// if has_err != nil {
	// 	fmt.Println(get_err.Error())
	// 	os.Exit(1)
	// }

	// if has_key {
	// 	fmt.Println("DB found the key value")
	// }

	// for i := 0; i < 100; i++ {
	// 	put_err = db.Put([]byte("testkey"+strconv.Itoa(i)), []byte("testvalue"+strconv.Itoa(i)))

	// 	if put_err != nil {
	// 		fmt.Println(put_err.Error())
	// 		os.Exit(1)
	// 	}
	// }

	// iterator, scan_err := db.RangeScan([]byte("testkey5"), []byte("testkey60"))

	// if scan_err != nil {
	// 	fmt.Println(scan_err.Error())
	// 	os.Exit(1)
	// }

	// for {
	// 	key := iterator.Key()
	// 	value := iterator.Value()

	// 	fmt.Println("key: " + string(key[:]) + " value: " + string(value[:]))

	// 	done := iterator.Next()

	// 	if !done {
	// 		break
	// 	}
	// }

	// data, get_err = db.Get([]byte("willfail"))

	// if get_err != nil {
	// 	fmt.Println(get_err.Error())
	// 	os.Exit(1)
	// }
}
