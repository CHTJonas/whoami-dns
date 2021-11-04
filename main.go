package main

import (
	"os"
	"os/signal"
	"syscall"

	dnstap "github.com/dnstap/golang-dnstap"
)

var serv *server

func init() {
	serv = newServer()
}

func parseCapture() {
	input, err := dnstap.NewFrameStreamSockInputFromPath("/var/lib/knot/dnstap.sock")
	if err != nil {
		panic(err)
	}

	output := dnstap.NewTextOutput(serv, dnstap.TextFormat)

	go output.RunOutputLoop()
	input.ReadInto(output.GetOutputChannel())
}

func main() {
	serv.start()
	go parseCapture()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit

	serv.stop()
}
