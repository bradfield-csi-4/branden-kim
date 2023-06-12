package db

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestSkipListDBPut(t *testing.T) {
	cwd, cwd_err := os.Getwd()

	if cwd_err != nil {
		fmt.Println(cwd_err.Error())
		os.Exit(1)
	}

	f, err := os.OpenFile(cwd+"/testwriteahead.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db := CreateSkipListSkeletonDB(f)

	put_err := db.Put([]byte("Some Key"), []byte("Some Value"))

	if put_err != nil {
		t.Errorf("Expected put to operate successfully but it failed due to %s", put_err.Error())
	}
}

func BenchmarkSkipListGet(b *testing.B) {
	cwd, cwd_err := os.Getwd()

	if cwd_err != nil {
		fmt.Println(cwd_err.Error())
		os.Exit(1)
	}

	f, err := os.OpenFile(cwd+"/testwriteahead.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db := CreateSkipListSkeletonDB(f)

	key_prefix := "testkey"
	value_prefix := "testvalue"

	for i := 0; i < 1000; i++ {
		put_err := db.Put([]byte(key_prefix+strconv.Itoa(i)), []byte(value_prefix+strconv.Itoa(i)))

		if put_err != nil {
			b.Log(put_err.Error())
			b.Fail()
			break
		}
	}

	key_bytes := []byte("testkey0")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, get_err := db.Get(key_bytes)

		if get_err != nil {
			b.Log(get_err.Error())
			b.Fail()
			break
		}
	}

}

func BenchmarkSkipListHas(b *testing.B) {
	cwd, cwd_err := os.Getwd()

	if cwd_err != nil {
		fmt.Println(cwd_err.Error())
		os.Exit(1)
	}

	f, err := os.OpenFile(cwd+"/testwriteahead.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db := CreateSkipListSkeletonDB(f)

	key_prefix := []byte("testkey")
	value_prefix := []byte("testvalue")

	db.Put(key_prefix, value_prefix)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, has_err := db.Has(key_prefix)

		if has_err != nil {
			b.Log(has_err.Error())
			b.Fail()
			break
		}
	}
}

func BenchmarkSkipListPut(b *testing.B) {
	cwd, cwd_err := os.Getwd()

	if cwd_err != nil {
		fmt.Println(cwd_err.Error())
		os.Exit(1)
	}

	f, err := os.OpenFile(cwd+"/testwriteahead.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db := CreateSkipListSkeletonDB(f)

	key_prefix := "testkey"
	value_prefix := "testvalue"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		put_err := db.Put([]byte(key_prefix+strconv.Itoa(i)), []byte(value_prefix+strconv.Itoa(i)))

		if put_err != nil {
			b.Log(put_err.Error())
			b.Fail()
			break
		}
	}
}

// TODO: Figure out how to benchmark "Deletes"
func BenchmarkSkipListDelete(b *testing.B) {
	cwd, cwd_err := os.Getwd()

	if cwd_err != nil {
		fmt.Println(cwd_err.Error())
		os.Exit(1)
	}

	f, err := os.OpenFile(cwd+"/testwriteahead.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db := CreateSkipListSkeletonDB(f)

	for i := 0; i < b.N; i++ {
		put_err := db.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)))

		if put_err != nil {
			b.Log(put_err.Error())
			b.Fail()
			break
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		delete_err := db.Delete([]byte(strconv.Itoa(i)))

		if delete_err != nil {
			b.Log(delete_err.Error())
			b.Fail()
			break
		}
	}
}

func BenchmarkSkipListRangeScanTimeToGetIterator(b *testing.B) {
	cwd, cwd_err := os.Getwd()

	if cwd_err != nil {
		fmt.Println(cwd_err.Error())
		os.Exit(1)
	}

	f, err := os.OpenFile(cwd+"/testwriteahead.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db := CreateSkipListSkeletonDB(f)

	for i := 0; i < 10000000; i++ {
		put_err := db.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)))

		if put_err != nil {
			b.Log(put_err.Error())
			b.Fail()
			break
		}
	}

	var table = []struct {
		scan_size int
	}{
		{scan_size: 0},
		{scan_size: 99},
		{scan_size: 9999},
		{scan_size: 999999},
	}

	b.ResetTimer()

	for _, v := range table {
		b.Run(fmt.Sprintf("scan_size_%d", v.scan_size+1), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := db.RangeScan([]byte(strconv.Itoa(0)), []byte(strconv.Itoa(v.scan_size)))

				if err != nil {
					b.Log(err.Error())
					b.Fail()
					break
				}
			}
		})
	}
}
