package tts_socket

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func ListenToTabletopApp(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("[TTS] Error creating listener:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[TTS] Listening on %s\n", addr)

	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[TTS] Error accepting connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

// load the full json file into memory then log it
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a buffer to store received data
	var buffer bytes.Buffer
	reader := bufio.NewReader(conn)

	for {
		// Read data from the connection
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("[TTS] Error reading from connection:", err)
				return
			}
			// If EOF, append the last part of data and break
			buffer.Write(data)
			break
		}

		// Append data to the buffer
		buffer.Write(data)
	}

	// Print the received message once
	fmt.Printf("[TTS] Received message:\n%s\n", buffer.String())
}
