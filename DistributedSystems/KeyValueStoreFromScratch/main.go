package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Setting up the signal handler for SIGINT
	sigs := make(chan os.Signal, 4)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("Received " + sig.String() + "! Gracefully exiting...")
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	memory_store := CreateMemoryStore()

	for {
		fmt.Print("> ")
		str, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		// To remove the newline that gets appended at the end
		str = str[:len(str)-1]

		scanner_state := HandleCommand(str)

		switch scanner_state.Operation {
		case "SET":
			memory_store.Set(scanner_state.Key, scanner_state.Value)
			fmt.Println("Successfully added key to store")
		case "GET":
			result := memory_store.Get(scanner_state.Key)
			fmt.Println(result)
		case "INVALID":
			fmt.Println("Invalid command supplied!")
		}
	}
}
