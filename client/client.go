package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	msg := "Hello from the client!"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error writing: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Message sent successfully!")
}
