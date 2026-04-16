package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer li.Close()
	// for {
	// 	conn, err := li.Accept()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	// 	io.WriteString(conn, "Content-Type: text/plain\r\n")
	// 	io.WriteString(conn, "\r\n")
	// 	io.WriteString(conn, "Hello, World!")
	// 	conn.Close()
	// }

	for {
		conn, err := li.Accept()
		if err != nil {
			panic(err)
		}
		go handle(conn)
	}

}

func handle(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	fmt.Println("Code got here.")
}
