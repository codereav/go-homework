package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Таймаут подключения")
	flag.Parse()

	host := flag.Arg(0)
	port := flag.Arg(1)

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	go receive(ctx, stop, client)
	go send(ctx, stop, client)

	<-ctx.Done()
}

func send(ctx context.Context, stop context.CancelFunc, c TelnetClient) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.Send(); err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
					stop()
				}
			}
		}
	}
}

func receive(ctx context.Context, stop func(), c TelnetClient) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.Receive(); err != nil {
				os.Stderr.Write([]byte(err.Error()))
				if errors.Is(err, io.EOF) {
					fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
					stop()
				}
			}
		}
	}
}
