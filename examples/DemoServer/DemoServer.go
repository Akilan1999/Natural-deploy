package main

import (
	"fmt"
	natural_deploy "github.com/Akilan1999/P2PRC-natural-deploy"
	"github.com/Akilan1999/p2p-rendering-computation/client"
	"github.com/Akilan1999/p2p-rendering-computation/p2p"
	"time"
)

func main() {
	fmt.Println("Starting server procedure")
	fmt.Println(".................................")
	// Start P2PRC instance
	//natural_deploy.RunAsP2PRCNode()

	fmt.Println("Starting server (Please wait for 7 to 12 seconds)")
	fmt.Println(".................................")

	// Create a 5 second delay from mapping port
	time.Sleep(5 * time.Second)

	// ------------------- Create a task -----------------
	var task natural_deploy.Task
	task.Name = "Test"
	task.TaskFile = "test.sh"
	task.KillTaskFile = "kill.sh"

	// Get all node information
	table, err := p2p.ReadIpTable()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("IP table read")

	// Searching for the information of a particular
	// node information
	var machineInfo p2p.IpAddress
	for i, _ := range table.IpAddress {
		if table.IpAddress[i].Name == "Test" {
			machineInfo = table.IpAddress[i]
		}
	}

	fmt.Println(machineInfo)

	// Allocate ports
	var port client.ResponseMAPPort
	port.PortNo = "8000"
	task.ExposedPorts = append(task.ExposedPorts, &port)

	task.NodeInfo = &machineInfo
	// ---------------------------------------------------------

	err = task.CreateTask()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server running for 10 seconds")

	// Create a 10 second delay from mapping port
	time.Sleep(10 * time.Second)

	err = task.KillTask()
	if err != nil {
		fmt.Println(err)
		return
	}

}
