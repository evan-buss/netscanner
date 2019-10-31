package main

import (
	"flag"
	"fmt"

	"github.com/pkg/profile"
)

// netscanner is a program that quickly scans your network for open ports
// Goals:
//  Scan all local IP for specific port
//  Scan all local IP for port range
//    Offer multiple methods (sequential, 5 threads *old method*, thread pool, goroutine mania)
// While we do this, I want to learn more about profiling.

var mode string

func init() {
	flag.StringVar(&mode, "mode", "seq", "Set the netscanner method (seq, pool, swarm)")
}

func main() {
	defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()

	flag.Parse()

	fmt.Println(mode)
}
