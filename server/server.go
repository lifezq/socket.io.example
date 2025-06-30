package main

import (
	"fmt"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	httpServer := types.CreateServer(nil)
	io := socket.NewServer(httpServer, nil)

	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		auth := client.Handshake().Auth.(map[string]any)

		fmt.Printf("connection:%v\n", client.Handshake().Auth)

		exit := false
		room := socket.Room(auth["username"].(string))
		client.Join(room)
		client.On("event", func(datas ...any) {
		})
		client.On("disconnect", func(...any) {
			client.Leave(room)
			exit = true
		})

		go func() {
			midx := 0
			for {
				midx++
				fmt.Printf("%d.send message to room: %s\n", midx, room)
				io.Of("/", nil).In(room).Emit("message", types.NewStringBufferString(fmt.Sprintf("%d.hello world", midx)))
				time.Sleep(time.Second)
				if exit {
					return
				}
			}
		}()
	})

	httpServer.Listen("127.0.0.1:3000", nil)

	exit := make(chan struct{})
	SignalC := make(chan os.Signal)

	signal.Notify(SignalC, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range SignalC {
			switch s {
			case os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	httpServer.Close(nil)
	os.Exit(0)
}
