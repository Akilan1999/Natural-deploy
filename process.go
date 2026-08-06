package natural_deploy

import (
	"bufio"
	"encoding/gob"
	"errors"
	"github.com/Akilan1999/p2p-rendering-computation/abstractions"
	"github.com/Akilan1999/p2p-rendering-computation/client"
	"github.com/Akilan1999/p2p-rendering-computation/config"
	"github.com/Akilan1999/p2p-rendering-computation/p2p"
	"github.com/kr/pretty"
	"github.com/melbahja/goph/v2"
	"net/http"
	"os"
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
var TaskPath = "./task.store"

func init() {
	// read from tasks previous tasks running from disk
	read()
	// Ping all tasks
	PingAllTasks()
}

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

	client.Run("mkdir ~/p2prc-task/" + task.Name)

	// Get home directory path of the remote machine
	out, err := client.Run("pwd")
	if err != nil {
		return err
	}

	path := strings.TrimSuffix(string(out), "\n")

	path = path + "/p2prc-task/" + task.Name + "/" + task.TaskFile

	// upload the file to directory
	err = client.Upload(task.TaskFile, path)
	if err != nil {
		return err
	}

	out, err = client.Run("cd ~/p2prc-task/" + task.Name + "/ && sh " + task.TaskFile)

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

func (task *Task) SendFile(file string) error {
	// SSH into the intended node
	client, err := task.MakeConnection()
	if err != nil {
		return err
	}

	// Defer closing the network connection.
	defer client.Close()

	client.Run("mkdir ~/p2prc-task/")

	client.Run("mkdir ~/p2prc-task/" + task.Name)

	// Get home directory path of the remote machine
	out, err := client.Run("pwd")
	if err != nil {
		return err
	}

	path := strings.TrimSuffix(string(out), "\n")

	path = path + "/p2prc-task/" + task.Name + "/" + file

	// upload the file to directory
	err = client.Upload(file, path)
	if err != nil {
		return err
	}

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

	path = path + "/p2prc-task/" + task.Name + "/" + task.KillTaskFile

	// upload the kill file to directory
	err = client.Upload(task.KillTaskFile, path)
	if err != nil {
		return err
	}

	// Run to kill the process
	out, err = client.Run("cd ~/p2prc-task/" + task.Name + "/ && sh " + task.KillTaskFile)

	if err != nil {
		return err
	}

	task.Comment = "Server killed"
	task.Active = false

	// Remove the task folder
	out, err = client.Run("cd ~/p2prc-task/ && rm -rf " + task.Name)

	if err != nil {
		return err
	}

	// unregister the task
	task.UnregisterTask()

	return nil
}

// RegisterTask tracks the task through memory
func (task *Task) RegisterTask() error {
	Tasks[task.Name] = task
	// Write to disk
	err := write()
	if err != nil {
		return err
	}
	return nil
}

// UnregisterTask Remove task from the map
func (task *Task) UnregisterTask() error {
	delete(Tasks, task.Name)
	// Write to disk
	err := write()
	if err != nil {
		return err
	}
	return nil
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

func PingAllTasks() {
	for _, t := range Tasks {
		t.PingProgress()
	}
}

func PrintTasks() {
	pretty.Println(Tasks)
}

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

// Read file for the tasks
func read() error {
	file, err := os.Open(TaskPath)
	if err != nil {
		_, err := os.Create(TaskPath)
		if err != nil {
			return errors.New("unable to create store file")
		}
	}

	fileReader := bufio.NewReader(file)
	fileDecoder := gob.NewDecoder(fileReader)
	if err := fileDecoder.Decode(&Tasks); err != nil {
		return errors.New("unable to decode")
	}
	file.Close()
	return nil
}

// Write to file for the tasks
func write() error {
	file, err := os.Create(TaskPath)
	if err != nil {
		return errors.New("unable to write file")
	}
	defer file.Close()
	fileWriter := bufio.NewWriter(file)
	fileEncoder := gob.NewEncoder(fileWriter)

	if err := fileEncoder.Encode(Tasks); err != nil {
		return errors.New("unable to encode")
	}

	fileWriter.Flush()
	return nil
}
