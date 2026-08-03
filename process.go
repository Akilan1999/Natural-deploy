package natural_deploy

import (
	"github.com/Akilan1999/p2p-rendering-computation/abstractions"
	"github.com/Akilan1999/p2p-rendering-computation/client"
	"github.com/Akilan1999/p2p-rendering-computation/config"
	"github.com/Akilan1999/p2p-rendering-computation/p2p"
	"github.com/kr/pretty"
	"github.com/melbahja/goph/v2"
	"net/http"
	"strconv"
	"strings"
)

// This library assumes all nodes run using P2PRC bare
// mode to ensure all node can ssh into each of them.

type Task struct {
	Name         string
	NodeInfo     *p2p.IpAddress
	ExposedPorts []*Ports
	// This needs to be a bash script to start a task
	TaskFile string
	// This needs to be a bash script to kill a task
	KillTaskFile string
	Comment      string
	Active       bool
}

type Ports struct {
	Port       string
	DomainName string
	Response   *client.ResponseMAPPort
}

var Tasks = make(map[string]*Task)

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

	out, err = client.Run("cd ~/p2prc-task/ && sh " + task.TaskFile)

	if err != nil {
		return err
	}

	for i, port := range task.ExposedPorts {
		// Creates port on the node running P2PRC
		task.ExposedPorts[i].Response, err = abstractions.MapPort(port.Port, port.DomainName, task.NodeInfo.Ipv4+":"+task.NodeInfo.ServerPort, false)
		if err != nil {
			return err
		}
	}

	// set task active to trust
	task.Active = true

	// Append information to the task tracker
	//Tasks.Tasks = append(Tasks.Tasks, task)

	// register the task
	task.RegisterTask()

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

	// unregister the task
	task.UnregisterTask()

	return nil
}

// RegisterTask tracks the task through memory
func (task *Task) RegisterTask() {
	Tasks[task.Name] = task
}

// UnregisterTask Remove task from the map
func (task *Task) UnregisterTask() {
	delete(Tasks, task.Name)
}

// PrintTask Prints a particular task
func (task *Task) PrintTask() {
	pretty.Println(task)
}

// ViewTasks Search task based on the name provided.
func ViewTasks(name string) (*Task, bool) {
	task, ok := Tasks[name]
	return task, ok
}

func PrintTasks() {
	pretty.Println(Tasks)
}

// PingProgress Tracker future task
// PingProgress Checks if the process is active or not
// if a single port is not running the task is considered
// as inactive.
func (task *Task) PingProgress() bool {
	for _, port := range task.ExposedPorts {
		resp, err := http.Get("http://" + port.Response.EntireAddress)
		if err != nil || resp.StatusCode != 200 {
			task.Active = false
			task.Comment = "Address " + port.Response.EntireAddress + " is not active, which belongs to task machine's local port " + port.Port
			task.RegisterTask()
			return false
		}
	}
	return true
}
