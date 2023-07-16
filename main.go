package main

import (
	"alekseikormski.com/server-status-monitoring/app/core"
	appserver "alekseikormski.com/server-status-monitoring/app/server"
	appWebsocker "alekseikormski.com/server-status-monitoring/app/web-socket"
	"embed"
	"log"
)

var (
	//go:embed front-end/build
	Res embed.FS
)

func main() {
	config := core.NewConfig()
	update := make(chan []*core.Application)

	go core.Start(update, config)

	wss := appWebsocker.NewWebSocket(update)

	server := appserver.NewServer(config.Port, wss, Res)
	if err := server.Start(server.RegisterRoutes()); err != nil {
		log.Fatalf("Problem with server: %w", err)
		return
	}
}
