package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func HandleEchoRequest(req Request, encoding string) []byte {

	var res []byte
	var body []byte
	bodyStr := strings.TrimPrefix(req.path, "/echo/")

	if encoding != "" {
		body = compressBytesTogzip([]byte(bodyStr))

	} else {

		body = []byte(bodyStr)
	}
	header := createResHeader(Header{
		contentType:     CONTENT_TYPE_PLAIN_TEXT,
		contentLength:   fmt.Sprintf("%d", len(body)),
		contentEncoding: encoding,
	})
	res = makeResponse(STATUS_LINE_OK, header, []byte(body))
	return res
}

func HandleUserAgentRequest(req Request, encoding string) []byte {
	var res []byte
	val, ok := req.header["User-Agent"]
	if ok {
		//	header := fmt.Sprintf("Content-Type: text/plain%sContent-Length: %d%s", CRLF, len(val), CRLF)
		header := createResHeader(Header{
			contentType:     CONTENT_TYPE_PLAIN_TEXT,
			contentLength:   fmt.Sprintf("%d", len(val)),
			contentEncoding: encoding,
		})
		res = makeResponse(STATUS_LINE_OK, header, []byte(val))
	} else {
		res = makeResponse(STATUS_LINE_NOT_FOUND, "", nil)
	}

	return res
}
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
