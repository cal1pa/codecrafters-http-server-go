package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strings"
)

const (
	STATUS_LINE_OK          = "HTTP/1.1 200 OK"
	STATUS_LINE_NOT_FOUND   = "HTTP/1.1 404 Not Found"
	STATUS_LINE_CREATED     = "HTTP/1.1 201 Created"
	CRLF                    = "\r\n"
	CONTENT_TYPE_PLAIN_TEXT = "text/plain"
)

type Header struct {
	contentType     string
	contentLength   string
	contentEncoding string
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	var dirPath string
	dirCmd := flag.NewFlagSet("directory", flag.ExitOnError)
	if len(os.Args) > 2 {
		dirCmd.Parse(os.Args[2:])
		dirPath = dirCmd.Arg(0)

	}
	s := CreateServer()
	defer s.Close()

	for {
		conn, err := s.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn, dirPath)

	}
}

func handleConnection(conn net.Conn, dirPath string) {
	defer conn.Close()
	var res string
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		return
	}
	message := string(buff[:n])

	log.Printf("INFO: Read %d bytes\n", n)
	log.Println("INFO: The message was: ", message)

	req, err := ParseRequestMessage(message)

	log.Println("INFO: The path was: ", req.path)
	log.Println("INFO: The method was: ", req.method)

	if err != nil {
		log.Println("ERROR: Could not parse request message:", err.Error())
	}
	resEncoding, ok := req.header["Accept-Encoding"]
	accptedEncodings := strings.Split(resEncoding, ", ")
	acceptedEncoding := ""
	if ok && slices.Contains(accptedEncodings, "gzip") {
		acceptedEncoding = "gzip"
	}

	switch {
	case req.path == "/user-agent":
		{
			val, ok := req.header["User-Agent"]
			if ok {
				//	header := fmt.Sprintf("Content-Type: text/plain%sContent-Length: %d%s", CRLF, len(val), CRLF)
				header := createResHeader(Header{
					contentType:     CONTENT_TYPE_PLAIN_TEXT,
					contentLength:   fmt.Sprintf("%d", len(val)),
					contentEncoding: acceptedEncoding,
				})
				res = makeResponse(STATUS_LINE_OK, header, val)
				break
			} else {
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", "")
			}
		}
	case strings.HasPrefix(req.path, "/echo/"):
		{
			body := strings.TrimPrefix(req.path, "/echo/")
			header := createResHeader(Header{
				contentType:     CONTENT_TYPE_PLAIN_TEXT,
				contentLength:   fmt.Sprintf("%d", len(body)),
				contentEncoding: acceptedEncoding,
			})
			res = makeResponse(STATUS_LINE_OK, header, body)
		}
	case strings.HasPrefix(req.path, "/files/"):
		{
			fmt.Println()
			fmt.Println("===============================================")
			fmt.Println("Attempt to handle file")
			res = HandleFileRequest(req, dirPath)
		}
	case req.path == "/":
		{
			res = makeResponse(STATUS_LINE_OK, "", "")
		}
	default:
		res = makeResponse(STATUS_LINE_NOT_FOUND, "", "")
	}

	_, err = conn.Write([]byte(res))
	if err != nil {
		log.Println("ERROR: Failed to write response: ", err.Error())
	}
	fmt.Println(res)
}

func makeResponse(statusline, header, body string) string {
	return fmt.Sprintf("%s%s%s%s%s\n", statusline, CRLF, header, CRLF, body)
}

func createResHeader(header Header) string {
	val := ""
	if header.contentType != "" {
		val += fmt.Sprintf("Content-Type:%s%s", header.contentType, CRLF)
	}

	if header.contentLength != "" {
		val += fmt.Sprintf("Content-Length:%s%s", header.contentLength, CRLF)
	}

	if header.contentEncoding != "" {
		val += fmt.Sprintf("Content-Encoding:%s%s", header.contentEncoding, CRLF)
	}
	return val
}
