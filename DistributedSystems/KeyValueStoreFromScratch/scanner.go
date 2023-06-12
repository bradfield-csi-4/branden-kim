package main

import (
	"strings"
)

type Command struct {
	Operation string
	Key       string
	Value     string
}

func HandleCommand(command string) *Command {
	new_command := &Command{}
	command_split := strings.Split(command, " ")

	operation := strings.ToUpper(command_split[0])

	switch operation {
	case "SET":
		new_command.Operation = "SET"
		new_command.Key = strings.Split(command_split[1], "=")[0]
		new_command.Value = strings.Split(command_split[1], "=")[1]

	case "GET":
		new_command.Operation = "GET"
		new_command.Key = strings.Split(command_split[1], "=")[0]
		new_command.Value = ""

	default:
		new_command.Operation = "INVALID"
		new_command.Key = ""
		new_command.Value = ""
	}

	return new_command
}
