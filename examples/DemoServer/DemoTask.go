package main

import (
	"fmt"
	natural_deploy "github.com/Akilan1999/P2PRC-natural-deploy"
	"github.com/Akilan1999/p2p-rendering-computation/client"
	"github.com/Akilan1999/p2p-rendering-computation/p2p"
	"os"
	"time"
)

func main() {
	_, err := os.Stat("Task.lock")
	if err != nil {
		// If you want to add a custom root node
		var rootnode natural_deploy.RootNode
		// make sure to fill in "<>" area with the appropriate information
		rootnode.IPAddress = "<IPV4>"
		rootnode.Port = "<ServerPort>"
		natural_deploy.CreateTaskMachine("Test-2", &rootnode)

		// or

		//natural_deploy.CreateTaskMachine("Test-2", nil)
		os.Create("Task.lock")
	}

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
		os.Exit(1)
		return
	}

	fmt.Println("IP table read")

	// Searching for the information of a particular
	// node information (In this example Test-2 node)
	var machineInfo p2p.IpAddress
	for i, _ := range table.IpAddress {
		if table.IpAddress[i].Name == "Test-1" {
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
		os.Exit(1)
		return
	}

	// Create a 10 second delay from mapping port
	time.Sleep(10 * time.Second)

	err = task.KillTask()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

}
