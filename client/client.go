package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func readServer(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server")
			return
		}

		fmt.Print(line)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:4242")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	go readServer(conn)

	stdin := bufio.NewReader(os.Stdin)

	for {
		text, err := stdin.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Fprint(conn, text)
	}
}
