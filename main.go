package natural_deploy

import (
	"fmt"
	"github.com/Akilan1999/p2p-rendering-computation/abstractions"
	"github.com/Akilan1999/p2p-rendering-computation/client"
	"github.com/Akilan1999/p2p-rendering-computation/config"
	"github.com/Akilan1999/p2p-rendering-computation/p2p"
	"github.com/melbahja/goph/v2"
	"strconv"
	"strings"
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
	Comment      string
	Active       bool
}

type TaskTracker struct {
	Tasks []*Task
}

var Tasks *TaskTracker

func (task *Task) MakeConnection() (*goph.Client, error) {
	// Get config information of P2PRC
	Config, err := config.ConfigInit(nil, nil)
	if err != nil {
		return nil, err
	}

	// SSH port
	SSHPort, err := strconv.Atoi(task.NodeInfo.BareMetalSSHPort)
	if err != nil {
		return nil, err
	}

	// SSH into the node and deploy bash
	client, err := goph.New(task.NodeInfo.MachineUsername, task.NodeInfo.Ipv4, goph.WithKeyFile(Config.PrivateKeyFile, ""), goph.WithPort(uint(SSHPort)), goph.WithInsecureIgnoreHostKey())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (task *Task) CreateTask() error {

	// SSH into the intended node
	client, err := task.MakeConnection()
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

	fmt.Println("Uploading the file")

	// upload the file to directory
	err = client.Upload(task.TaskFile, path)
	if err != nil {
		return err
	}

	out, err = client.Run("cd ~/p2prc-task/ && sh " + task.TaskFile)

	fmt.Println(string(out))

	if err != nil {
		return err
	}

	for i, port := range task.ExposedPorts {
		// Creates port on the node running P2PRC
		task.ExposedPorts[i], err = abstractions.MapPort(port.PortNo, "", task.NodeInfo.Ipv4+":"+task.NodeInfo.ServerPort, false)
		fmt.Println(task.ExposedPorts[i].PortNo)
		if err != nil {
			return err
		}
	}

	// set task active to trust
	task.Active = true

	// Append information to the task tracker
	//Tasks.Tasks = append(Tasks.Tasks, task)

	return nil
}

func (task *Task) KillTask() error {
	// SSH into the intended node
	client, err := task.MakeConnection()
	if err != nil {
		return err
	}

	// Defer closing the network connection.
	defer client.Close()

	// Get home directory path of the remote machine
	out, err := client.Run("pwd")
	if err != nil {
		return err
	}

	path := strings.TrimSuffix(string(out), "\n")

	path = path + "/p2prc-task/" + task.KillTaskFile

	// upload the kill file to directory
	err = client.Upload(task.KillTaskFile, path)
	if err != nil {
		return err
	}

	// Run the kill file
	out, err = client.Run("cd ~/p2prc-task/ && sh " + task.KillTaskFile)

	if err != nil {
		return err
	}

	task.Comment = "Server killed"
	task.Active = false

	return nil
}

// Tracker future task
// PingProgress Checks if the process is active or not
//func (task *Task) PingProgress() bool {
//
//	for _, port := range task.ExposedPorts {
//		l := time.Duration(2 * time.Second) // 2sec
//		// Ping checked 3 times
//		for i := 0; i < 3; i++ {
//			sTime := time.Now()
//			resp, err := http.Get(port.EntireAddress)
//			fTime := time.Now()
//			if err != nil || resp.StatusCode != 200 {
//				return false
//			}
//			if fTime.Sub(sTime) < l {
//				l = fTime.Sub(sTime)
//			}
//			resp.Body.Close()
//		}
//	}
//
//	return true
//}

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
