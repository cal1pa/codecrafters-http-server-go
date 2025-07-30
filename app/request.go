package main

import (
	"bytes"
	"errors"
	"log"
	"strings"
)

type Request struct {
	method   string
	path     string
	protocol string
	header   map[string]string
	body     []byte
}

func ParseRequestMessage(rawMessage []byte) (Request, error) {
	message := bytes.Split(rawMessage, []byte(CRLF+CRLF))
	rest := bytes.Split(message[0], []byte(CRLF))
	status_line := string(rest[0])
	status_line_arr := strings.Fields(status_line)

	if len(status_line_arr) != 3 {
		return Request{}, errors.New("invalid Status line")
	}

	parsed_body := message[1]
	parsed_header := rest[1:]

	hd := make(map[string]string)

	for _, s := range parsed_header {
		keyValArr := strings.Split(string(s), ": ")
		if len(keyValArr) != 2 {
			log.Println("ERROR: Unable to parse header: ", string(s))
			return Request{}, errors.New("invalid header")
		}
		hd[keyValArr[0]] = keyValArr[1]
	}

	return Request{
		method:   strings.TrimSpace(status_line_arr[0]),
		path:     strings.TrimSpace(status_line_arr[1]),
		protocol: strings.TrimSpace(status_line_arr[2]),
		header:   hd,
		body:     bytes.TrimSpace(parsed_body),
	}, nil
}
