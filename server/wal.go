package server

import (
	"bufio"
	"github.com/ekundo/godis/resp"
	"io"
	"log"
	"os"
	"sync"
)

type wal struct {
	in        chan *resp.Message
	out       chan *resp.Message
	file      *os.File
	writer    *bufio.Writer
	stop      *chan bool
	stopGroup *sync.WaitGroup
}

func newWal(stop *chan bool, stopGroup *sync.WaitGroup) *wal {
	file, err := os.OpenFile("cache.wal", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("can't open WAL file for writing: ", err)
	}

	writer := bufio.NewWriter(file)
	in := make(chan *resp.Message, 1000)
	out := make(chan *resp.Message, 1000)
	wal := &wal{in: in, out: out, file: file, writer: writer, stop: stop, stopGroup: stopGroup}

	go wal.readLoop()
	go wal.writeLoop()

	return wal
}

func (wal *wal) write(msg *resp.Message) {
	wal.in <- msg
}

func (wal *wal) read() (*resp.Message, bool) {
	msg, ok := <-wal.out
	return msg, ok
}

func (wal *wal) readLoop() {
	reader := resp.NewReader(bufio.NewReader(wal.file))
	msg := &resp.Message{}
	for {
		select {
		case <-*wal.stop:
			return
		default:
		}
		parsed, err := msg.Parse(reader)
		if err != nil {
			if err == io.EOF {
				close(wal.out)
				return
			}
			log.Panic("can't read WAL:", err)
		}

		if parsed {
			wal.out <- msg
			msg = &resp.Message{}
		}
	}
}

func (wal *wal) writeLoop() {
	defer func() {
		log.Println("flushing WAL")
		_ = wal.writer.Flush()
		_ = wal.file.Sync()
		_ = wal.file.Close()
	}()

	for {
		select {
		case msg := <-wal.in:
			wal.writeMsg(msg)
		case <-*wal.stop:
			return
		default:
		}
	}
}

func (wal *wal) writeMsg(msg *resp.Message) {
	_, err := wal.writer.Write(msg.Ser())
	if err != nil {
		log.Panic("can't write to WAL:", err)
	}
}
