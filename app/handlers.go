package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func HandleEchoRequest(req Request, resHeader Header) []byte {

	var res []byte
	var body []byte
	bodyStr := strings.TrimPrefix(req.path, "/echo/")

	if resHeader.contentEncoding != "" {
		body = compressBytesTogzip([]byte(bodyStr))

	} else {

		body = []byte(bodyStr)
	}
	resHeader.contentLength = fmt.Sprintf("%d", len(body))
	resHeader.contentType = CONTENT_TYPE_PLAIN_TEXT

	res = makeResponse(STATUS_LINE_OK, createResHeader(resHeader), []byte(body))
	log.Println("HEADER: ", createResHeader(resHeader))

	return res
}

func HandleUserAgentRequest(req Request, resHeader Header) []byte {
	var res []byte
	val, ok := req.header["User-Agent"]
	if ok {
		//	header := fmt.Sprintf("Content-Type: text/plain%sContent-Length: %d%s", CRLF, len(val), CRLF)
		resHeader.contentLength = fmt.Sprintf("%d", len(val))
		resHeader.contentType = CONTENT_TYPE_PLAIN_TEXT

		res = makeResponse(STATUS_LINE_OK, createResHeader(resHeader), []byte(val))
	} else {
		res = makeResponse(STATUS_LINE_NOT_FOUND, "", nil)
	}
	log.Println("HEADER: ", createResHeader(resHeader))
	return res
}
func HandleFileRequest(req Request, dirPath string, resHeader Header) []byte {
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
			resHeader.contentLength = fmt.Sprintf("%d", info.Size())
			resHeader.contentType = CONTENT_TYPE_OCTET_STREAM
			res = makeResponse(STATUS_LINE_OK, createResHeader(resHeader), data)

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
