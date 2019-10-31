package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// netscanner is a program that quickly scans your network for open ports

var mode string
var addr string

func init() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: netscanner [-mode] port\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "  port\n\tThe port to scan\n")
	}

	flag.StringVar(&mode, "mode", "swarm", "Set the netscanner method (seq, pool, swarm)")
	flag.StringVar(&addr, "address", "192.168.1", "Set the network addres to scan. Exclude the last number")
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	port := string(flag.Args()[0])

	switch mode {
	case "seq":
		fmt.Printf("%+q\n", Sequential(port))
	case "pool":
		fmt.Printf("%+q\n", Pool(port, 5))
	case "swarm":
		fmt.Printf("%+q\n", Swarm(port))
	}
}

// Sequential scans for the given port sequentially on a single thread
func Sequential(port string) []string {
	dialer := net.Dialer{Timeout: time.Millisecond * 100}
	output := make([]string, 0)

	for i := 0; i <= 255; i++ {
		addr := addr + "." + strconv.Itoa(i) + ":" + port
		conn, err := dialer.Dial("tcp", addr)
		if err != nil {
			continue
		}
		output = append(output, conn.RemoteAddr().String())
		conn.Close()
	}

	return output
}

// Pool creates a pool of [num] workers to scan for given port
func Pool(port string, num int) []string {
	jobs := make(chan int)
	results := make(chan string, 255)
	output := make([]string, 0)

	for i := 0; i < num; i++ {
		go worker(port, jobs, results)
	}

	for i := 0; i <= 255; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 0; i <= 255; i++ {
		select {
		case ip := <-results:
			if ip != "" {
				output = append(output, ip)
			}
		}
	}

	return output
}

func worker(port string, jobs <-chan int, results chan<- string) {
	dialer := net.Dialer{Timeout: time.Millisecond * 100}
	for job := range jobs {
		addr := addr + "." + strconv.Itoa(job) + ":" + port
		conn, err := dialer.Dial("tcp", addr)
		if err == nil {
			results <- conn.RemoteAddr().String()
			conn.Close()
		} else {
			results <- ""
		}
	}
}

// Swarm launches a goroutine for each address (256 total)
func Swarm(port string) []string {
	output := make([]string, 0)
	results := make(chan string)

	for i := 0; i <= 255; i++ {
		go func(i int, results chan<- string) {
			dialer := net.Dialer{Timeout: time.Millisecond * 100}
			addr := addr + "." + strconv.Itoa(i) + ":" + port
			conn, err := dialer.Dial("tcp", addr)
			if err == nil {
				results <- conn.RemoteAddr().String()
				conn.Close()
			} else {
				results <- ""
			}
		}(i, results)
	}

	for j := 0; j <= 255; j++ {
		select {
		case ip := <-results:
			if ip != "" {
				output = append(output, ip)
			}
		}
	}
	return output
}
