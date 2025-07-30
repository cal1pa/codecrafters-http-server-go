package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
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

// func HandleEchoRequest(req Request) []byte {
// 	body := []byte(strings.TrimPrefix(req.path, "/echo/"))
// 	header := fmt.Sprintf("Content-Type: text/plain%sContent-Length: %d%s", CRLF, len(body), CRLF)
// 	res := makeResponse(STATUS_LINE_OK, header, body)
// 	return res

// }

func HandleFileRequest(req Request, dirPath string) []byte {
	filename := strings.TrimPrefix(req.path, "/files/")
	reqFilePath := fmt.Sprintf("%s%s", dirPath, filename)
	//	reqFilePath := fmt.Sprintf("/home/astral/Workspace/codecrafters/codecrafters-http-server-go%s", filename)

	var res []byte
	switch req.method {
	case "GET":
		{
			info, err := os.Stat(reqFilePath)
			if err != nil {
				log.Println("ERROR: Cant get file info", err.Error())
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", nil)
				return res
			}
			data, err := os.ReadFile(reqFilePath)
			if err != nil {
				log.Println("ERROR: Failed to read file: ", err.Error())
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", nil)
				return res
			}
			header := fmt.Sprintf("Content-Type: application/octet-stream%sContent-Length: %d%s", CRLF, info.Size(), CRLF)
			res = makeResponse(STATUS_LINE_OK, header, data)

		}
	case "POST":
		{
			data := req.body
			log.Println("Request body to write", string(data))
			err := os.WriteFile(reqFilePath, data, 0644)
			if err != nil {
				log.Println("Failed to write file")
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", nil)
				return res
			}
		}
		res = makeResponse(STATUS_LINE_CREATED, "", nil)

	}

	return res

}
