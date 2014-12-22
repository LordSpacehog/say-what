package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	udpMulticastAddress := &net.UDPAddr{
		IP:   net.ParseIP("239.255.23.42"),
		Port: 5235,
	}
	conn, err := net.ListenMulticastUDP("udp", nil, udpMulticastAddress)
	if err != nil {
		log.Panic(err)
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", "10.21.12.3:5235")
	if err != nil {
		log.Panic(err)
	}

	w := UDPWriter{
		writer: conn,
		addr:   remoteAddr,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go copyStream(os.Stdin, w, &wg)

	go copyStream(conn, os.Stdout, &wg)

	wg.Wait()
}

func copyStream(r io.Reader, w io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	b := make([]byte, 1500)
	var done bool
	for done != true {
		n, err := r.Read(b)
		if err == io.EOF {
			done = true
		} else if err != nil {
			log.Panic(err)
		}
		for n > 0 {
			m, err := w.Write(b[:n])
			if err != nil {
				log.Panic(err)
			}
			n -= m
			b = b[m:]
		}
	}
}

type UDPWriter struct {
	writer *net.UDPConn
	addr   *net.UDPAddr
}

func (uw UDPWriter) Write(b []byte) (int, error) {
	log.Print("Writing to UDP: %v", b)
	return uw.writer.WriteToUDP(b, uw.addr)
}
