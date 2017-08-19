package main

import (
	"flag"
	"time"
	"fmt"
	"net"
	"strconv"
	"strings"
	"os"
)

func main() {
	interval := flag.Duration("interval", 24 * time.Hour, "Specify the interval the service tries to restart the target server.")
	pollRate := flag.Duration("pollrate", 30 * time.Second, "Specify how many seconds to wait between polls after the interval duration has passed.")
	ipAddress := flag.String("ip", "", "Specify the target server IP address.")
	port := flag.Int("port", 27960, "Specify the target server port.")
	rconPassword := flag.String("rconpassword", "", "Specify the target server rcon password.")
	numChecksBeforeRestart := flag.Int("numchecks", 5, "How many times should the service check that the server is empty before restarting it.")

	flag.Parse()
	fmt.Printf("Starting the server restarter with params:\n");
	fmt.Printf("Interval: %v\n", interval)
	fmt.Printf("Poll rate: %v\n", pollRate)
	fmt.Printf("IP address: %s:%d\n", *ipAddress, *port)
	fmt.Printf("Rcon password: %s\n", *rconPassword)

	fmt.Println("Testing connection and rcon password")
	switch testConnection(ipAddress, port, rconPassword) {
	case InvalidPassword:
		fmt.Println("Invalid rcon password, exiting.")
		os.Exit(1)
		break
	case Timeout:
		fmt.Println("Destination server is unreachable, exiting.")
		os.Exit(1)
		break
	default:
		fmt.Println("Connection ok, starting the service")
	}

	var emptyCount int
	for {
		time.Sleep(*interval)
		fmt.Printf("Atleast %v has passed since last restart at %v. Polling the server to see if it's empty\n", interval, time.Now())

		emptyCount = 0
		for {
			if isEmpty(ipAddress, port) {
				emptyCount++
			} else {
				emptyCount = 0
			}

			if emptyCount >= *numChecksBeforeRestart {
				fmt.Printf("Server was empty for %v. Killing server\n", (*pollRate) * time.Duration(*numChecksBeforeRestart))
				killServer(ipAddress, port, rconPassword)
				break
			}

			time.Sleep(*pollRate)
		}
	}
}
const (
	Timeout = iota
	InvalidPassword = iota
	Ok = iota
)
func testConnection(ipAddress *string, port *int, rconPassword *string) int {
	conn, err := net.DialTimeout("udp", *ipAddress + ":" + strconv.Itoa(*port), time.Second)
	if err != nil {
		return Timeout
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	conn.SetDeadline(time.Now().Add(time.Second))
	conn.Write([]byte("\xff\xff\xff\xffrcon " + *rconPassword + " status"))
	_, err = conn.Read(buffer)
	if err != nil {
		return Timeout
	}

	response := string(buffer)
	if strings.Contains(response, "Bad rconpassword.") {
		return InvalidPassword
	}
	return Ok
}
func killServer(ipAddress *string, port *int, rconPassword *string) {
	conn, err := net.DialTimeout("udp", *ipAddress + ":" + strconv.Itoa(*port), time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second))
	conn.Write([]byte("\xff\xff\xff\xffrcon " + *rconPassword + " quit"))
}
func isEmpty(ipAddress *string, port *int) bool {
	conn, err := net.DialTimeout("udp", *ipAddress + ":" + strconv.Itoa(*port), time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	conn.SetDeadline(time.Now().Add(time.Second))
	conn.Write([]byte("\xff\xff\xff\xffgetstatus"))
	conn.Read(buffer)

	statusResponse := string(buffer)
	rows := strings.Split(statusResponse, "\n")

	return len(rows) <= 3
}
