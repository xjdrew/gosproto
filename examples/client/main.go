package main

import (
	"log"
	"net"
	"strconv"

	"github.com/xjdrew/gosproto"
	"github.com/xjdrew/gosproto/examples/sproto_echo"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8686")
	if err != nil {
		log.Fatal(err)
	}

	client, _ := sproto.NewService(conn, sproto_echo.Protocols)
	count := 100000
	log.Printf("test echo %d times", count)
	ping := sproto_echo.PingRequest{}

	// dispatch
	go client.Dispatch()

	// call
	for i := 0; i < count; i++ {
		v := strconv.Itoa(i)
		ping.Ping = &v
		resp, err := client.Call("echo.ping", &ping)
		if err != nil {
			log.Fatalf("ping failed:%s", err)
		}
		pong := resp.(*sproto_echo.PingResponse)
		if pong.Pong == nil || *pong.Pong != v {
			log.Fatalf("ping failed")
		}
	}
	log.Printf("test echo end")
}
