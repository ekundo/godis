package server

import (
	"fmt"
	"github.com/ekundo/godis/resp"
	"github.com/ekundo/godis/shared"
	"log"
	"net"
	"sync"
	"time"
)

type server struct {
	wal        *wal
	controller *controller
	stop       *chan bool
	stopGroup  *sync.WaitGroup
}

func NewServer(enableWal bool) *server {
	stop := make(chan bool)
	stopGroup := &sync.WaitGroup{}
	var wal *wal
	if enableWal {
		wal = newWal(&stop, stopGroup)
	}
	return &server{controller: newController(wal), wal: wal, stop: &stop, stopGroup: stopGroup}
}

func (srv *server) Start(host string, port int) {
	srv.applyWal()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err) // can't resolve address
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err) // can't bind to port
	}
	log.Println("listening on", listener.Addr())

	srv.stopGroup.Add(1)
	go srv.accept(listener)
}

func (srv *server) Stop() {
	close(*srv.stop)
	srv.stopGroup.Wait()
}

func (srv *server) applyWal() {
	if srv.wal == nil {
		return
	}
	log.Println("reading write-ahead log")
	for req, ok := srv.wal.read(); ok; req, ok = srv.wal.read() {
		_, _ = srv.controller.processRequest(req, false)
	}
}

func (srv *server) accept(listener *net.TCPListener) {
	defer srv.stopGroup.Done()
	defer listener.Close()

	for {
		select {
		case <-*srv.stop:
			log.Println("stopping listening on", listener.Addr())
			return
		default:
		}
		listener.SetDeadline(time.Now().Add(time.Second))
		conn, err := listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		srv.stopGroup.Add(1)
		go srv.handleConn(conn)
	}
}

func (srv *server) handleConn(conn *net.TCPConn) {
	defer conn.Close()
	defer srv.stopGroup.Done()

	reader := resp.NewReader(conn)
	msg := &resp.Message{}
	for {
		select {
		case <-*srv.stop:
			log.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}

		conn.SetDeadline(time.Now().Add(time.Second))
		parsed, err := msg.Parse(reader)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			srv.processError(conn, err)
			return
		}

		if parsed {
			srv.processRequest(conn, msg)
			msg = &resp.Message{}
		}
	}
}

func (srv *server) processError(c net.Conn, err error) {
	terr, ok := err.(shared.TypedError)
	var errType string
	if ok {
		errType = terr.ErrorType()
	} else {
		errType = "ERR"
	}
	res := &resp.Message{Element: &resp.Error{Kind: []byte(errType), Data: []byte(err.Error())}}
	fmt.Fprint(c, string(res.Ser()))
}

func (srv *server) processRequest(c net.Conn, req *resp.Message) {
	res, err := srv.controller.processRequest(req, srv.wal != nil)
	if err != nil {
		srv.processError(c, err)
		return
	}
	if res != nil {
		fmt.Fprint(c, string(res.Ser()))
	}
}
