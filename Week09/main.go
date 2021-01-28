package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", ":1234", "default listen port")
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("running tcp server error: %v", err)
	}
	log.Printf("tcp server is listening on %s\n", port)

	c := make(chan os.Signal, 1)
	exit := make(chan struct{})

	// 注册监听这两个信号
	// syscall.SIGINT: ctrl + c
	// syscall.SIGTERM: 结束程序
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println(<-c)
		exit <- struct{}{}
	}()

	go server(listener)

	<-exit
	listener.Close()
}

func server(listener net.Listener) {
	done := make(chan []byte)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept tcp connection error: %v", err)
		}
		defer listener.Close()

		go readConn(conn, done)
		go writeConn(conn, done)
	}
}

func readConn(conn net.Conn, done chan []byte) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatalf("read error:%v", err)
		}
		done <- line
	}
}

func writeConn(conn net.Conn, done chan []byte) {
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	for {
		receive := <-done
		_, err := writer.Write(receive)
		if err != nil {
			log.Fatalf("write error: %v", err)
		}
		writer.Write([]byte("\n"))
		writer.Flush()
	}
}
