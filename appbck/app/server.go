package main

import (
	"log"
	"net"
)

type Server struct {
	net.Listener
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
