// Connect to FreeLing server and get results
package main

import (
	"bufio"
	"fmt"
	"net"
)

// Connect to FreeLing
func connectToFreeLing(host string) (c net.Conn, err error) {
	c, err = net.Dial("tcp", host)

	return
}

// Get results from FreeLing
func getMorphResult(msg string, c net.Conn) (status string, err error) {
	fmt.Fprintf(c, "%s%c", msg, '\x00')
	status, err = bufio.NewReader(c).ReadString('\x00')

	return
}
