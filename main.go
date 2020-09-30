package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	listenerNetType = "tcp"
	listenerAddress = ":8080"
	HttpMethodGet   = "GET"
)

func main() {
	listener := createListener()
	defer listener.Close()

	handleIncomingConnections(listener)
}

func createListener() net.Listener {
	listener, err := net.Listen(listenerNetType, listenerAddress)

	if err != nil {
		log.Panic(err)
	}

	return listener
}

func handleIncomingConnections(listener net.Listener) {
	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
		}

		go handleConnection(connection)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	i := 0
	for scanner.Scan() {
		if i == 0 {
			mux(conn, scanner.Text())
		}

		i++
	}
}

func mux(conn net.Conn, ln string) {
	method := strings.Fields(ln)[0]
	uri := strings.Fields(ln)[1]

	if method == HttpMethodGet && uri == "/" {
		index(conn)
	} else if method == HttpMethodGet && uri == "/about" {
		about(conn)
	} else {
		notFound(conn)
	}
}

func index(conn net.Conn) {
	body := `<html><head><title>Index</title></head><body>Index <br> <a href="/about">About</a></body></html>`

	sendHttpResponse(conn, body)
}

func about(conn net.Conn) {
	body := `<html><head><title>About</title></head><body>Index <br> <a href="/">Back to index</a></body></html>`

	sendHttpResponse(conn, body)
}

func notFound(conn net.Conn) {
	body := `<html><head><title>Page Not Found</title></head><body><h1>404</h1> <br> <a href="/">Go to Index</a></body></html>`

	sendHttpResponse(conn, body)
}

func sendHttpResponse(conn net.Conn, body string) {
	fmt.Fprintln(conn, "HTTP/1.1 200 OK")
	fmt.Fprintln(conn, "Content-Length: ", len(body))
	fmt.Fprintln(conn, "Content-Type: text/html")
	fmt.Fprintln(conn, "")
	fmt.Fprintln(conn, body)
}