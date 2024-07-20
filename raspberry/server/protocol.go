package server

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ResponseCode byte

const (
	RESPONSE_OK ResponseCode = iota
)

type ErrorCode byte

const (
	RESPONSE_UNKNOWN_COMMAND ErrorCode = iota + 10
	RESPONSE_UNKNOWN_ERROR
	RESPONSE_STREAM_ERROR
	RESPONSE_INVALID_PARAMETERS_NUMBER
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

type Response interface {
	Parse() []byte
}

type BaseResponse struct {
	Code ResponseCode
}

func (r *BaseResponse) Parse() []byte {
	return []byte(fmt.Sprintf("%d", r.Code))
}

type ResponseWithFloat32Arguments struct {
	Code ResponseCode
	Args []float32
}

func (r *ResponseWithFloat32Arguments) Parse() []byte {
	response := fmt.Sprintf("%d", r.Code)
	for _, arg := range r.Args {
		response += fmt.Sprintf("$%f", arg)
	}
	return []byte(response)
}

type ResponseWithStringArguments struct {
	Code ResponseCode
	Args []string
}

func (r *ResponseWithStringArguments) Parse() []byte {
	response := fmt.Sprintf("%d", r.Code)
	for _, arg := range r.Args {
		response += fmt.Sprintf("$%s", arg)
	}
	return []byte(response)
}

type ErrorResponse struct {
	Code ErrorCode
	Err  error
}

func (er *ErrorResponse) Parse() []byte {
	return []byte(fmt.Sprintf("%d$%s", er.Code, er.Err))
}
