package main

import (
	"io"
	"log"
	"net"

	"github.com/xjdrew/gosproto"
	"github.com/xjdrew/gosproto/examples/sproto_echo"

	"net/http"
	_ "net/http/pprof"
)

type Echo int

var echo Echo = 0

func (e *Echo) Ping(req *sproto_echo.PingRequest, resp *sproto_echo.PingResponse) {
	resp.Pong = req.Ping
	*e++
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	server, _ := sproto.NewService(conn, sproto_echo.Protocols)
	server.Register(&echo)
	err := server.Dispatch()
	if err == io.EOF {
		log.Printf("client(%v) closed", conn.RemoteAddr())
	} else {
		log.Printf("client(%v) failed:%s", conn.RemoteAddr(), err)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8686")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}
