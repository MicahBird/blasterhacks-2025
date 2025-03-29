package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"
)

// Read from socket
func rsock() []byte {
	// DialUnix does not take a context.Context parameter. This example shows
	// how to dial a Unix socket with a Context. Note that the Context only
	// applies to the dial operation; it does not apply to the connection once
	// it has been established.

	c, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return nil
	}

	connection, err := c.Accept()

	defer c.Close()

	for {
		buf := make([]byte, 1024)
		for {
			_, err := connection.Read(buf[:])
			if err != nil {
				return nil
			}
			//println("Client got:", string())
			return buf[:]
		}
	}
}

// Write to socket
func wsock(text string) {
	// DialUnix does not take a context.Context parameter. This example shows
	// how to dial a Unix socket with a Context. Note that the Context only
	// applies to the dial operation; it does not apply to the connection once
	// it has been established.
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	d.LocalAddr = nil // if you have a local addr, add it here
	raddr := net.UnixAddr{Name: socketPath, Net: "unix"}
	conn, err := d.DialContext(ctx, "unix", raddr.String())
	if err != nil {
		log.Printf("Failed to dial (write): %v", err)
		return
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	if _, err := conn.Write([]byte(text)); err != nil {
		log.Print(err)

	}
}
