package main

import (
	"bufio"
	"fmt"
	"net"
)

func connectToFreeLing(host string) (c net.Conn, err error) {
	c, err = net.Dial("tcp", host)
	if err != nil {
		// TODO delete this
		fmt.Println(err)
		return
	}

	return
}

func getMorphResult(msg string, c net.Conn) (status string, err error) {
	fmt.Fprintf(c, "%s%c", msg, '\x00')
	status, err = bufio.NewReader(c).ReadString('\x00')
	if err != nil {
		// TODO delete this
		fmt.Println(err)
		return
	}

	return
}
