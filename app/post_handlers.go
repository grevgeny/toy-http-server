package main

import "net"

func handlePost(conn net.Conn, request *Request) {
	handleUnknown(conn)
}
