/*
Serve is a very simple static file server in go
Usage:

	-p="8100": port to serve on
	-d=".":    the directory of static files to host

Navigating to http://localhost:8100 will display the index.html or directory
listing file.
*/
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	natural_deploy "github.com/Akilan1999/Natural-deploy"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Daemon Runs the daemon first to learn about nodes in the network
func Daemon() error {
	_, err := os.Stat("Task.lock")
	if err != nil {
		err := natural_deploy.CreateTaskMachine("Test-2", nil)
		if err != nil {
			return err
		}

		os.Create("Task.lock")
	}

	err = natural_deploy.RunDaemon()
	if err != nil {
		return err
	}

	return nil
}

// SelfCompileLinux Compile the current program for linux x86.
func SelfCompileLinux() error {
	cmd := exec.Command("sh", "compile.sh")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

// Deploy deploys to a particular machine called test 2
func Deploy(Name string) error {
	var task natural_deploy.Task
	task.Name = Name
	task.TaskFile = "run.sh"
	task.KillTaskFile = "exit.sh"

	var err error
	task.NodeInfo, err = natural_deploy.SearchMachine("Test-Server")
	if err != nil {
		return err
	}

	// Allocate ports
	var port natural_deploy.Ports
	port.Port = "8034"
	port.DomainName = "nd.akilan.io"
	task.ExposedPorts = append(task.ExposedPorts, &port)

	// Check if the linux binary to run is compiled
	_, err = os.Stat("Deploy-Linux")
	if err != nil {
		err = SelfCompileLinux()
		if err != nil {
			return err
		}
	}

	fmt.Println("Uploading Deploy-Linux.....")

	// Sends files for setup
	err = task.SendFile("Deploy-Linux")
	if err != nil {
		return err
	}

	fmt.Println("Uploading index.html.....")
	err = task.SendFile("index.html")
	if err != nil {
		return err
	}

	fmt.Println("Uploading style.css.....")
	err = task.SendFile("style.css")
	if err != nil {
		return err
	}

	fmt.Println("Starting to run the task.....")
	// Creates and runs the task
	err = task.CreateTask()
	if err != nil {
		return err
	}

	// Test to print out the process
	ListProcess()

	return nil
}

func ListProcess() {
	natural_deploy.PrintTasks()
}

func KillProcess(Name string) error {
	task, avaliable := natural_deploy.ViewTasks(Name)
	if !avaliable {
		return errors.New("task not found")
	}

	err := task.KillTask()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	deploy := flag.Bool("deploy", false, "Deploy the webpage to the designated machine")
	kill := flag.Bool("kill", false, "Kills the deployed task")
	daemon := flag.Bool("daemon", false, "Starts the background process to learn about nodes in the network")
	ls := flag.Bool("ls", false, "List processes")
	port := flag.String("p", "8034", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	if *deploy {
		err := Deploy("nd-docs")
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	if *kill {
		err := KillProcess("nd-docs")
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	if *daemon {
		err := Daemon()
		if err != nil {
			fmt.Println(err)
			return
		}
		for {

		}
	}

	if *ls {
		ListProcess()
		return
	}

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
