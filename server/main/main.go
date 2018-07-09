package main

import (
	"flag"
	"github.com/ekundo/godis/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	host := flag.String("host", "127.0.0.1", "bind to host")
	port := flag.Int("port", 2121, "bind to port")
	wal := flag.Bool("wal", true, "enable write-ahead log")
	flag.Parse()

	srv := server.NewServer(*wal)
	srv.Start(*host, *port)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	srv.Stop()
}
