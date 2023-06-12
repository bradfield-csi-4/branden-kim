package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// For now im going to force the in-memory representation to have to be strings for the keys and strings for the values
// in a JSON structure. Haven't really thought of a use case where the user should be able to query
// in a memory complicated manner such as a key that points to a list of values rather than just a string
type MemoryStore struct {
	data map[string]string
	fd   *os.File
}

func CreateMemoryStore() *MemoryStore {
	memory_store := new(MemoryStore)
	memory_store.data = make(map[string]string)

	// Open the file with
	// O_RDWR since we want to read and write to the file
	// O_CREATE since we want to create the file if it doesn't exist
	fd, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0777)

	if err != nil {
		fmt.Println("Opening the database file failed due to " + err.Error())
	}

	// Store the file descriptor so that we can operate on it in later calls to SET
	memory_store.fd = fd

	byteValue, _ := ioutil.ReadAll(fd)

	json.Unmarshal(byteValue, &memory_store.data)

	fmt.Println(memory_store.data)

	// for s := range memory_store.data["testing_arrays"].([]interface{}) {
	// 	fmt.Println(s)
	// }

	return memory_store
}

func (ms *MemoryStore) Get(key string) string {
	return ms.data[key]
	// switch ms.data[key].(type) {
	// case string:
	// 	return ms.data[key].(string)
	// case float64:
	// 	fmt.Println("Key value was float")
	// 	return "float64"
	// case int:
	// 	fmt.Println("Key value was int")
	// 	return "int"
	// case interface{}:
	// 	fmt.Println("Key value was interface")
	// 	return "interface"
	// default:
	// 	fmt.Println("Got an invalid return type! Json might be invalid.")
	// 	return ""
	// }
}

func (ms *MemoryStore) Set(key string, value string) {
	// Add the value to the memory store
	ms.data[key] = value

	// TODO: Implement something better or more optimized
	// Write the whole store to disk
	// Ignoring error handling

	// Need to truncate the file first to clear out the contents
	ms.fd.Truncate(0)
	// Also need to seek to the beginning or else it will write from the previous position in the file
	ms.fd.Seek(0, 0)

	write_data, _ := json.Marshal(ms.data)
	ms.fd.Write(write_data)
}
