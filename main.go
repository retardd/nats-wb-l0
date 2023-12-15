package main

import (
	"context"
	"fmt"
	"l0/config"
	"l0/internal/datastorage/cache"
	"l0/internal/nats"
	nats_server "l0/nats-server"
)

func main() {
	if config.SetEnv() == false {
		return
	}
	fmt.Println("CONFIGURE GOOD")

	_, cch := cache.InitCache(context.TODO())
	hn, _ := nats.InitConnection(cch)

	defer hn.StopProcess()

	tempPub := hn.Pub
	server := nats_server.InitApi(cch, tempPub)

	// Запуск сервера (До нажатия CTRL+C)
	server.Start()
	return
}
