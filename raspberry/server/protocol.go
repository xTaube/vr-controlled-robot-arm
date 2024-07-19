package server

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ResponseCode byte

const (
	OK ResponseCode = iota
)

type ErrorCode byte

const (
	UNKNOWN_COMMAND ErrorCode = iota + 10
	UNKNOWN_ERROR
)

func readFloat32(data string) float32 {
	float, _ := strconv.ParseFloat(data, 32)
	return float32(float)
}

func ParseRequestArguments(request string) (CommandIdentifier, []string) {
	arguments := strings.Split(request, "$")
	log.Printf("%s\n", arguments[0])
	command_id, _ := strconv.ParseInt(arguments[0], 10, 8)
	return CommandIdentifier(command_id), arguments[1:]
}

func SuccessResponse(code ResponseCode, additional_context string) []byte {
	return []byte(fmt.Sprintf("%d$%s", code, additional_context))
}

func ErrorResponse(code ErrorCode, err error) []byte {
	return []byte(fmt.Sprintf("%d$%s", code, err))
}