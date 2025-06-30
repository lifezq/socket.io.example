package main

import (
	"time"

	"github.com/zishang520/engine.io-client-go/transports"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
	"github.com/zishang520/socket.io-client-go/socket"
)

func main() {
	opts := socket.DefaultOptions()
	opts.SetAuth(map[string]any{"username": "lisi"})
	opts.SetTransports(types.NewSet(transports.Polling, transports.WebSocket /*transports.WebTransport*/))

	manager := socket.NewManager("http://127.0.0.1:3000", opts)
	// Listening to manager events
	manager.On("error", func(errs ...any) {
		utils.Log().Warning("Manager Error: %v", errs)
	})

	manager.On("ping", func(...any) {
		utils.Log().Warning("Manager Ping")
	})

	manager.On("reconnect", func(...any) {
		utils.Log().Warning("Manager Reconnected")
	})

	manager.On("reconnect_attempt", func(...any) {
		utils.Log().Warning("Manager Reconnect Attempt")
	})

	manager.On("reconnect_error", func(errs ...any) {
		utils.Log().Warning("Manager Reconnect Error: %v", errs)
	})

	manager.On("reconnect_failed", func(errs ...any) {
		utils.Log().Warning("Manager Reconnect Failed: %v", errs)
	})

	io := manager.Socket("/", opts)
	utils.Log().Error("socket %v", io)
	io.On("connect", func(args ...any) {
		utils.Log().Warning("io iD %v", io.Id())
		utils.SetTimeout(func() {
			io.Emit("message", types.NewStringBufferString("test"))
		}, 1*time.Second)
		utils.Log().Warning("connect %v", args)
	})

	io.On("connect_error", func(args ...any) {
		utils.Log().Warning("connect_error %v", args)
	})

	io.On("disconnect", func(args ...any) {
		utils.Log().Warning("disconnect: %+v", args)
	})

	io.OnAny(func(args ...any) {
		utils.Log().Warning("OnAny: %+v", args)
	})

	io.On("message-back", func(args ...any) {
		// io.Emit("message", types.NewStringBufferString("88888"))
		utils.Log().Question("message-back: %+v", args)
	})

	select {}
}
