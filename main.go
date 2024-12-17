package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var delay time.Duration

func main() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP)
		for range c {
			if delay == 0 {
				delay = time.Second * 15
			} else {
				delay = 0
			}
		}
	}()

	listener, err := net.Listen("tcp", "0.0.0.0:10000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	log.Println("handling conn")
	defer log.Println("/handling conn")
	if err := _handle(conn); err != nil {
		log.Println("failed to handle:", err)
	}
}

func _handle(conn net.Conn) error {
	wg := &sync.WaitGroup{}

	wr, err := net.Dial("tcp", os.Getenv("TARGET"))
	if err != nil {
		return err
	}
	defer wr.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(wr, conn)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		go io.Copy(conn, SlowReader{wr})
	}()

	wg.Wait()
	return nil
}

type SlowReader struct{ r io.Reader }

func (sr SlowReader) Read(b []byte) (int, error) {
	time.Sleep(delay)
	return sr.r.Read(b)
}
