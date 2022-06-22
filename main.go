package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Closing...")
}
