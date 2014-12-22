/*
  say what?

  test chat client in GOlang
*/
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func readStdin(c chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		stdin, _ := reader.ReadString('\n')
		c <- stdin
	}
}

func readUDP(addr string) {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	listen, _ := net.ListenUDP("udp", udpAddr)
	var msg string
	defer listen.Close()
	var buf []byte
	buf = make([]byte, 1500)
	for {
		n, _, err := listen.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		if n > 0 {
			msg = strings.TrimSpace(string(buf[0:n]))
			fmt.Println(msg)
		}
	}
}

func main() {
	running := true
	var conn net.Conn
	var err error
	var connarg string
	stdin := make(chan string)
	go readStdin(stdin)
	go readUDP("0.0.0.0:5235")
	for running {
		select {
		case t := <-stdin:
			switch strings.Split(strings.TrimSpace(t), " ")[0] {
			case "/logout":
				running = false
			case "/connect":
				if conn != nil {
					conn.Close()
				}
				connarg = strings.Split(t, " ")[1]
				conn, err = net.Dial("udp4", connarg)
				if err != nil {
					panic(err)
				}
				fmt.Println("Connected!")
			default:
				if conn != nil {
					fmt.Println("Sending:", t)
					fmt.Fprintf(conn, t)
				}
			}
		default:
		}
	}
	if conn != nil {
		conn.Close()
	}
}
