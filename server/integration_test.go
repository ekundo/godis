// +build integration

package server

import (
	"github.com/ekundo/godis/client"
	"log"
	"os"
	"testing"
	"time"
)

const testHost = "127.0.0.1"
const testPort = 12121

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

var srv *server
var c *client.Client

func setup() {
	srv = NewServer(false)
	srv.Start(testHost, testPort)

	c = client.New()
	err := c.Connect(testHost, testPort, time.Second)
	if err != nil {
		log.Fatal("can't connect to test server: ", err)
	}
}

func shutdown() {
	c.Close()
	srv.Stop()
}

func key(t *testing.T, key string) string {
	return t.Name() + "#" + key
}

func ttl(ttl time.Duration) *time.Duration {
	return &ttl
}
