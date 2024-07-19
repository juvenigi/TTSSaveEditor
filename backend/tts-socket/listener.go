package tts_socket

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

// ListenToAppTCP listens for incoming TCP connections from Tabletop Simulator
func ListenToAppTCP(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("[TTS] Error creating listener: " + err.Error())
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
	var reader = bufio.NewReader(conn)
	// Read data from the connection until EOF
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
