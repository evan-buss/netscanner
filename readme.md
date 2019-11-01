# netscanner

A simple CLI to scan computers on a network. Scan all computers for a specific port or scan a single computer on all ports.

## Installation

`go get github.com/evan-buss/netscanner`

## Usage

```
Usage: netscanner [-mode] port
  -address string
        Set the network addres to scan. Exclude the last number (default "192.168.1")
  -mode string
        Set the netscanner method (seq, pool, swarm, single) (default "swarm")
  port
        The port to scan
```

The modes are designed to show different work sharing techniques. The only two you should care about are `swarm` and `single`.

## Examples

Scan local network for computers with an open windows share (port 445)

`netscanner -mode swarm 445`

Scan remote network for computer with open HTTPS port (443). Note: leave off the last number and dot

`netscanner -mode swarm -address 145.28.29`

Scan a specific IP address for any open ports

`netscanner -mode single 192.168.1.9`
