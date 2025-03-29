package main

import (
	"context"
	"log"
	"net"
	"time"
)

// Read from socket
func rsock() {
	// DialUnix does not take a context.Context parameter. This example shows
	// how to dial a Unix socket with a Context. Note that the Context only
	// applies to the dial operation; it does not apply to the connection once
	// it has been established.
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	d.LocalAddr = nil // if you have a local addr, add it here
	raddr := net.UnixAddr{Name: "/tmp/unix.sock", Net: "unix"}
	conn, err := d.DialContext(ctx, "unix", raddr.String())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	var readData []byte

	_, err = conn.Read(readData)
	if err != nil {
		return
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
	raddr := net.UnixAddr{Name: "/tmp/unix.sock", Net: "unix"}
	conn, err := d.DialContext(ctx, "unix", raddr.String())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	if _, err := conn.Write([]byte(text)); err != nil {
		log.Fatal(err)
	}
}
