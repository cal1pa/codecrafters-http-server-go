package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	STATUS_LINE_OK        = "HTTP/1.1 200 OK"
	CRLF                  = "\r\n"
	STATUS_LINE_NOT_FOUND = "HTTP/1.1 404 Not Found"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	var res string
	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind port 4221")
		os.Exit(1)
	}

	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buff := make([]byte, 1024)

	n, err := conn.Read(buff)

	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes\n", n)
	fmt.Println("The message was: ", string(buff))

	msg := string(buff)

	msgArr := strings.Fields(msg)

	path := strings.TrimSpace(msgArr[1])
	if strings.HasPrefix(path, "/echo/") {
		body := strings.TrimPrefix(path, "/echo/")
		header := fmt.Sprintf("Content-Type: text/plain%sContent-Length: %d%s", CRLF, len(body), CRLF)
		res = makeResponse(STATUS_LINE_OK, header, body)
	} else if path != "/" {
		res = makeResponse(STATUS_LINE_NOT_FOUND, "", "")
	} else {
		res = makeResponse(STATUS_LINE_OK, "", "")
	}

	fmt.Println(path)
	_, err = conn.Write([]byte(res))
	if err != nil {
		fmt.Println("Failed to write response: ", err.Error())
		os.Exit(1)
	}
	conn.Close()
}

func makeResponse(statusline, header, body string) string {
	return fmt.Sprintf("%s%s%s%s%s\n", statusline, CRLF, header, CRLF, body)
}
