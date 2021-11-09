package main

import (
	"os"
	"os/signal"
	"syscall"
)

var serv *server

func init() {
	serv = newServer()
}

func main() {
	go serv.openSocket()
	serv.start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit

	serv.stop()
}
