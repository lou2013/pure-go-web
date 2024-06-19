package main

import (
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(8)
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for i := 0; i < 500; i++ {
		go ListenForConnection(listener)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGBUS)
	<-sigs
}

func ListenForConnection(listener net.Listener) error {
	for {
		conn, err := listener.Accept()

		if err != nil {
			return err
		}

		go func() {
			HandleConnection(conn)
			conn.Close()
		}()
	}
}

func HandleConnection(conn net.Conn) {

	incomingMessage := make([]byte, 1024)

	_, err := conn.Read(incomingMessage)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}
		panic(err)
	}

	httpResponse := `HTTP/1.1 200 OK
					Date: Mon, 27 Jul 2009 12:28:53 GMT
					Server: Apache/2.2.14 (Win32)
					Last-Modified: Wed, 22 Jul 2009 19:15:56 GMT
					Content-Type: text/html

					<html>
					<body>
					<h1> Hello, World<h1/>
					</body>
					</html>`
	if _, err := conn.Write([]byte(httpResponse)); err != nil {
		panic(err)
	}
}
