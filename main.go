package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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

	flag.StringVar(&mode, "mode", "swarm", "Set the netscanner method (seq, pool, swarm, single)")
	flag.StringVar(&addr, "address", "192.168.1", "Set the network addres to scan. Exclude the last number")
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	data := string(flag.Args()[0])

	switch mode {
	case "seq":
		fmt.Printf("%+q\n", Sequential(data))
	case "pool":
		fmt.Printf("%+q\n", Pool(data, 5))
	case "swarm":
		fmt.Printf("%+q\n", Swarm(data))
	case "single":
		fmt.Printf("%+q\n", FullScan(data))
	default:
		flag.Usage()
	}
}

// FullScan scans a given IP address on all ports
func FullScan(ip string) []string {

	output := make([]string, 0)
	jobs := make(chan int, 65_535)
	results := make(chan string)

	// Spawn 100 workers
	for i := 0; i < 800; i++ {
		go func(ip string, jobs <-chan int, results chan<- string) {
			for port := range jobs {
				conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
				if err == nil {
					results <- conn.RemoteAddr().String()
					conn.Close()
				} else {
					if !strings.Contains(err.Error(), "connection refused") {
						log.Println(err)
					}
					results <- ""
				}
			}
		}(ip, jobs, results)
	}

	for i := 0; i <= 65_535; i++ {
		jobs <- i
	}
	close(jobs)

	for j := 0; j <= 65_535; j++ {
		if j%1000 == 0 {
			fmt.Println(j)
		}
		select {
		case ip := <-results:
			if ip != "" {
				output = append(output, ip)
			}
		}
	}
	return output
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
	dialer := net.Dialer{Timeout: time.Millisecond * 100}
	output := make([]string, 0)
	results := make(chan string)

	for i := 0; i <= 255; i++ {
		go func(i int, results chan<- string) {
			conn, err := dialer.Dial("tcp", addr+"."+strconv.Itoa(i)+":"+port)
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
