package main

import (
	"io/ioutil"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var myIP string
var port string = "1337"

func init() {
	// Determine your machine's local network IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal("Could not get local IP")
	}
	defer conn.Close()
	myIP = conn.LocalAddr().String()[:strings.Index(conn.LocalAddr().String(), ":")]

	// Create a dummy server for all the scanner tests
	go func() {
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatal(err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			_, err = ioutil.ReadAll(conn)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}

func BenchmarkSequential(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Sequential("80")
	}
}

func BenchmarkPool10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Pool("80", 10)
	}
}

func BenchmarkPool50(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Pool("80", 50)
	}
}

func BenchmarkPool100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Pool("80", 100)
	}
}

func BenchmarkPool255(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Pool("80", 255)
	}
}

func BenchmarkSwarm(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Swarm("80")
	}
}

// TestSwarm tests the swarm scanner
func TestSwarm(t *testing.T) {
	is := is.New(t)
	results := Swarm(port)
	is.Equal([]string{myIP + ":" + port}, results) // Swarm found server
}

func TestPool(t *testing.T) {
	is := is.New(t)
	results := Pool(port, 100)
	is.Equal([]string{myIP + ":" + port}, results) // Pool found server
}

func TestSequential(t *testing.T) {
	is := is.New(t)
	results := Sequential(port)
	is.Equal([]string{myIP + ":" + port}, results) // Sequential found server
}
