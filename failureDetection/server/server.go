package main

import (
	"failureDetection/heartbeat"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func handleWorker(heartbeatHandler *heartbeat.HeartBeatHandler, workerTimeMap map[string]time.Time) {
	defer heartbeatHandler.Close()
	// worker should connect, disconnect

	workerMsg, _ := heartbeatHandler.Receive()
	workerHostName := workerMsg.GetHostName()
	workerPortNum := workerMsg.GetPortNumber()
	key := workerHostName + ":" + workerPortNum
	fmt.Println("Got a message from ", key)
	isAlive := workerMsg.GetIsAlive()

	if isAlive != "" {

		currTime := time.Now().Format("2006-01-02 15:04:05")
		currTimeFormatted, formatErr := time.Parse("2006-01-02 15:04:05", currTime)
		if formatErr != nil {
			log.Fatalln(formatErr)
		}
		fmt.Println("Current Time Formatted", currTimeFormatted)
		// successMsg := &heartbeat.HeartbeatMessage{Status: "Success"}
		if value, ok := workerTimeMap[key]; ok {
			diff := currTimeFormatted.Sub(value)
			fmt.Println("Time diff", diff)
			// fmt.Println("Time in map", value)
			if diff > 15*time.Second {
				// failureMsg := &heartbeat.HeartbeatMessage{Status: "Failure"}
				// heartbeatHandler.Send(failureMsg)
				fmt.Println("Failure. Reinitalize yourself as new node")
			} else {
				fmt.Println("Success")
				workerTimeMap[key] = currTimeFormatted
			}
		} else {
			workerTimeMap[key] = currTimeFormatted
			// heartbeatHandler.Send(successMsg)
		}

	}

}

func checkWorkerValidity(workerTimeMap map[string]time.Time) {

	for {
		currTime := time.Now().Format("2006-01-02 15:04:05")
		currTimeFormatted, formatErr := time.Parse("2006-01-02 15:04:05", currTime)

		if formatErr != nil {
			log.Fatalln(formatErr)
		}

		for key, value := range workerTimeMap {

			diff := currTimeFormatted.Sub(value)
			// fmt.Println("Time diff", diff)
			if diff > 15*time.Second {
				fmt.Println(key, "is Off", "Reinitalize yourself as new node")
			} else {
				fmt.Println(key, "is On")
			}

		}
		time.Sleep(10 * time.Second)
	}

}
func main() {

	workerTimeMap := make(map[string]time.Time)
	listner, err := net.Listen("tcp", ":"+os.Args[1])

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	go checkWorkerValidity(workerTimeMap)
	// Keep on running the server infinitely
	for {

		fmt.Println("Started an infinte loop")
		if conn, err := listner.Accept(); err == nil {
			heartbeatHandler := heartbeat.NewHeartBeatHandler(conn)

			handleWorker(heartbeatHandler, workerTimeMap)
		}

	}
}
