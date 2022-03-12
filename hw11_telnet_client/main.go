package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout duration")
	flag.Parse()

	host := flag.Arg(0)
	port := flag.Arg(1)

	client := NewTelnetClient(net.JoinHostPort(host, port), *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connection was closed")
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)

	go func() {
		defer cancel()

		err := client.Receive()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer cancel()

		err := client.Send()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
