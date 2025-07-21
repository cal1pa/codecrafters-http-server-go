package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	STATUS_LINE_OK        = "HTTP/1.1 200 OK"
	STATUS_LINE_NOT_FOUND = "HTTP/1.1 404 Not Found"
	CRLF                  = "\r\n"
	TEST_REQUEST          = "GET " + "/testPath " + "HTTP/1.1" + CRLF + "Host: localhost:4221" + CRLF + "User-Agent: foobar/1.2.3" + CRLF + CRLF
)

type Server struct {
	net.Listener
}

type Request struct {
	method   string
	path     string
	protocol string
	header   []string
	body     string
}

func CreateServer() *Server {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to build port 4221")
		os.Exit(1)
	}
	s := Server{
		l,
	}
	return &s
}

func AcceptReqest(server *Server) (string, net.Conn) {
	conn, err := server.Accept()
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

	return string(buff), conn
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	s := CreateServer()
	defer s.Close()

	rawMessage, conn := AcceptReqest(s)

	req, err := ParseRequestMessage(rawMessage)

	if err != nil {
		fmt.Println("Could not parse request message:", err.Error())
	}

	var res string

	switch {
	case req.path == "/user-agent":
		{
			userAgentValue := ""
			for _, s := range req.header {
				if value, ok := strings.CutPrefix(s, "User-Agent:"); ok {
					userAgentValue = strings.TrimSpace(value)
					header := fmt.Sprintf("Content-Type: text/plain%sContent-Length: %d%s", CRLF, len(userAgentValue), CRLF)
					res = makeResponse(STATUS_LINE_OK, header, userAgentValue)
					break
				}
			}
		}
	case strings.HasPrefix(req.path, "/echo/"):
		{
			body := strings.TrimPrefix(req.path, "/echo/")
			header := fmt.Sprintf("Content-Type: text/plain%sContent-Length: %d%s", CRLF, len(body), CRLF)
			res = makeResponse(STATUS_LINE_OK, header, body)
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
		fmt.Println("Failed to write response: ", err.Error())
		os.Exit(1)
	}
	conn.Close()

	fmt.Println(res)

}

func makeResponse(statusline, header, body string) string {
	return fmt.Sprintf("%s%s%s%s%s\n", statusline, CRLF, header, CRLF, body)
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

	return Request{
		method:   status_line_arr[0],
		path:     status_line_arr[1],
		protocol: status_line_arr[2],
		header:   parsed_header,
		body:     parsed_body,
	}, nil
}

func TestingParse() {
	req, _ := ParseRequestMessage(TEST_REQUEST)

	fmt.Println("Method: " + req.method)
	fmt.Println("Path: " + req.path)
	fmt.Println("Protocol : " + req.protocol)
	fmt.Print("Header: ")
	fmt.Println(req.header)
	fmt.Println("Body: " + req.body)

}
