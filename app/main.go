package main

import (
	"fmt"
	"net"
	"os"
)

const (
	STATUS_LINE_OK = "HTTP/1.1 200 OK"
	CRLF           = "\r\n"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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

	res := fmt.Sprintf("%s%s%s", STATUS_LINE_OK, CRLF, CRLF)

	_, err = conn.Write([]byte(res))
	if err != nil {
		fmt.Println("Failed to write response: ", err.Error())
		os.Exit(1)
	}
	conn.Close()
}
