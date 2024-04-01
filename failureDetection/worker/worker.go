package main

import (
	"failureDetection/heartbeat"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func connectToServer(port string, serverPort string) {
	conn, err := net.Dial("tcp", serverPort)

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer conn.Close()

	heartbeatHandler := heartbeat.NewHeartBeatHandler(conn)
	hostName, nameErr := os.Hostname()

	if nameErr != nil {
		log.Fatalln(nameErr.Error())
	}

	isAlive := "yes"

	msg := &heartbeat.HeartbeatMessage{HostName: hostName, PortNumber: port, IsAlive: isAlive}

	heartbeatHandler.Send(msg)
	fmt.Println(msg, " sent")
	// serverMsg, _ := heartbeatHandler.Receive()

	// fmt.Println("Server Message: ", serverMsg)
}

func main() {
	port := os.Args[1]
	serverPort := os.Args[2] // Server port number

	_, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		connectToServer(port, serverPort)
		time.Sleep(5 * time.Second)
	}
}
