package main

import (
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
	body     string
}

func ParseRequestMessage(rawMessage string) (Request, error) {
	message := strings.Split(rawMessage, CRLF+CRLF)
	rest := strings.Split(message[0], CRLF)
	status_line := rest[0]
	status_line_arr := strings.Fields(status_line)

	if len(status_line_arr) != 3 {
		return Request{}, errors.New("Invalid Status line")
	}

	parsed_body := message[1]
	parsed_header := rest[1:]

	hd := make(map[string]string)

	for _, s := range parsed_header {
		keyValArr := strings.Split(s, ": ")
		if len(keyValArr) != 2 {
			log.Println("ERROR: Unable to parse header: ", s)
			return Request{}, errors.New("Invalid header")
		}
		hd[keyValArr[0]] = keyValArr[1]
	}

	return Request{
		method:   strings.TrimSpace(status_line_arr[0]),
		path:     strings.TrimSpace(status_line_arr[1]),
		protocol: strings.TrimSpace(status_line_arr[2]),
		header:   hd,
		body:     strings.TrimSpace(parsed_body),
	}, nil
}

func HandleEchoRequest(req Request) string {
	body := strings.TrimPrefix(req.path, "/echo/")
	header := fmt.Sprintf("Content-Type: text/plain%sContent-Length: %d%s", CRLF, len(body), CRLF)
	res := makeResponse(STATUS_LINE_OK, header, body)
	return res

}

func HandleFileRequest(req Request, dirPath string) string {
	filename := strings.TrimPrefix(req.path, "/files/")
	reqFilePath := fmt.Sprintf("%s%s", dirPath, filename)
	//	reqFilePath := fmt.Sprintf("/home/astral/Workspace/codecrafters/codecrafters-http-server-go%s", filename)

	res := ""
	fmt.Println("The protocol is: ", req.protocol)
	switch req.method {
	case "GET":
		{
			info, err := os.Stat(reqFilePath)
			if err != nil {
				log.Println("ERROR: Cant get file info", err.Error())
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", "")
				return res
			}
			data, err := os.ReadFile(reqFilePath)
			if err != nil {
				log.Println("ERROR: Failed to read file: ", err.Error())
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", "")
				return res
			}
			stringData := string(data)
			header := fmt.Sprintf("Content-Type: application/octet-stream%sContent-Length: %d%s", CRLF, info.Size(), CRLF)
			res = makeResponse(STATUS_LINE_OK, header, stringData)

		}
	case "POST":
		{
			data := req.body
			log.Println("Request body to write", data)
			err := os.WriteFile(reqFilePath, []byte(data), 0644)
			if err != nil {
				log.Println("Failed to write file")
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", "")
				return res
			}
		}
		res = makeResponse(STATUS_LINE_CREATED, "", "")

	}

	return res

}
