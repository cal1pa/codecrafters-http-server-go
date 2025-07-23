package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	STATUS_LINE_OK        = "HTTP/1.1 200 OK"
	STATUS_LINE_NOT_FOUND = "HTTP/1.1 404 Not Found"
	CRLF                  = "\r\n"
)

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
		log.Fatal("Critical: Failed to build port 4221")
	}
	s := Server{
		l,
	}
	return &s
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
	log.Printf("INFO: Read %d bytes\n", n)

	message := string(buff[:n])

	log.Println("INFO: The message was: ", message)

	req, err := ParseRequestMessage(message)

	if err != nil {
		log.Println("ERROR: Could not parse request message:", err.Error())
	}
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
	case strings.HasPrefix(req.path, "/files/"):
		{
			filename := strings.TrimPrefix(req.path, "/files/")
			reqFilePath := fmt.Sprintf("%s%s", dirPath, filename)
			info, err := os.Stat(reqFilePath)
			if err != nil {
				log.Println("ERROR: Cant get file info", err.Error())
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", "")
				break

			}
			data, err := os.ReadFile(reqFilePath)
			if err != nil {
				log.Println("ERROR: Failed to read file: ", err.Error())
				res = makeResponse(STATUS_LINE_NOT_FOUND, "", "")
				break
			}

			stringData := string(data)
			header := fmt.Sprintf("Content-Type: application/octet-stream%sContent-Length: %d%s", CRLF, info.Size(), CRLF)
			res = makeResponse(STATUS_LINE_OK, header, stringData)
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

func makeResponse(statusline, header, body string) string {
	return fmt.Sprintf("%s%s%s%s%s\n", statusline, CRLF, header, CRLF, body)
}
