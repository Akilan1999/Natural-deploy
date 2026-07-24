package main

import (
	"fmt"
	"github.com/Akilan1999/p2p-rendering-computation/abstractions"
	"github.com/Akilan1999/p2p-rendering-computation/client"
	"github.com/Akilan1999/p2p-rendering-computation/config"
	"github.com/Akilan1999/p2p-rendering-computation/p2p"
	"github.com/melbahja/goph/v2"
	"strconv"
	"strings"
	"time"
)

// This library assumes all nodes run using P2PRC bare
// mode to ensure all node can ssh into each of them.

type Task struct {
	Name         string
	NodeInfo     *p2p.IpAddress
	ExposedPorts []*client.ResponseMAPPort
	// This needs to be a bash script to start a task
	TaskFile string
	// This needs to be a bash script to kill a task
	KillTaskFile string
	Active       bool
}

func (task *Task) CreateTask() error {
	// Get config information of P2PRC
	Config, err := config.ConfigInit(nil, nil)
	if err != nil {
		return err
	}

	// SSH port
	SSHPort, err := strconv.Atoi(task.NodeInfo.BareMetalSSHPort)
	if err != nil {
		return err
	}

	// SSH into the node and deploy bash
	client, err := goph.New(task.NodeInfo.MachineUsername, task.NodeInfo.Ipv4, goph.WithKeyFile(Config.PrivateKeyFile, ""), goph.WithPort(uint(SSHPort)), goph.WithInsecureIgnoreHostKey())
	if err != nil {
		return err
	}

	// Defer closing the network connection.
	defer client.Close()

	client.Run("mkdir ~/p2prc-task/")

	fmt.Println(task.TaskFile)

	// Get home directory path of the remote machine
	out, err := client.Run("pwd")
	if err != nil {
		return err
	}

	path := strings.TrimSuffix(string(out), "\n")

	path = path + "/p2prc-task/" + task.TaskFile

	// upload the file to directory
	err = client.Upload(task.TaskFile, path)
	if err != nil {
		return err
	}

	// run the task
	//_, err = client.Run("cd ~/p2prc-task/ && sh " + task.TaskFile)
	//
	//if err != nil {
	//    return err
	//}

	for i, port := range task.ExposedPorts {
		// Creates port on the node running P2PRC
		task.ExposedPorts[i], err = abstractions.MapPort(port.PortNo, "", task.NodeInfo.Ipv4+":"+task.NodeInfo.ServerPort, false)
		fmt.Println(task.ExposedPorts[i])
		if err != nil {
			return err
		}
	}

	// set task active to trust
	task.Active = true

	return nil
}

func (task *Task) KillTask() error {
	// Run kill script and send feedback
	return nil
}

// RunAsP2PRCNode Runs node as P2PRC instance
func RunAsP2PRCNode() error {
	// P2PRC configuration
	abstractions.Init(nil)

	Config, err := config.ConfigInit(nil, nil)
	if err != nil {
		return err
	}

	// Changing the name of the machine
	Config.MachineName = "Test"
	Config.BareMetal = true

	err = Config.WriteConfig()
	if err != nil {
		return err
	}

	// Start P2PRC server as a background process
	abstractions.Start()

	return nil
}

// ---------------------------------------------------------------------------------------

func main() {
	fmt.Println("Starting server procedure")
	fmt.Println(".................................")
	// Start P2PRC instance
	RunAsP2PRCNode()

	fmt.Println("Starting server (Please wait for 7 to 12 seconds)")
	fmt.Println(".................................")

	// Create a 5 second delay from mapping port
	time.Sleep(5 * time.Second)

	// ------------------- Create a task -----------------
	var task Task
	task.Name = "Test"
	task.TaskFile = "test.sh"

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

}
