package main

import (
	"fmt"
	natural_deploy "github.com/Akilan1999/Natural-deploy"
	"os"
	"time"
)

func main() {
	//_, err := os.Stat("Task.lock")
	//if err != nil {
	// If you want to add a custom root node
	//var rootnode natural_deploy.RootNode
	// make sure to fill in "<>" area with the appropriate information
	//rootnode.IPAddress = "<IPV4>"
	//rootnode.Port = "<ServerPort>"
	//natural_deploy.CreateTaskMachine("Test-2", &rootnode)

	// or

	// Test reasons
	//err := natural_deploy.CreateRegularNode("Test-2", nil)
	//if err != nil {
	//    fmt.Println(err)
	//    os.Exit(1)
	//    return
	//}
	//
	//// os.Create("Task.lock")
	////}
	//
	//err = natural_deploy.RunDaemon()
	//if err != nil {
	//    fmt.Println(err)
	//    os.Exit(1)
	//    return
	//}

	// Create a 5 second delay from mapping port
	// time.Sleep(10 * time.Second)

	// ------------------- Create a task -----------------
	var task natural_deploy.Task
	task.Name = "Test"
	task.TaskFile = "test.sh"
	task.KillTaskFile = "kill.sh"

	var err error
	task.NodeInfo, err = natural_deploy.SearchMachine("Test-Server")
	if err != nil {
		fmt.Println(err)
		return
	}

	for err != nil {
		task.NodeInfo, err = natural_deploy.SearchMachine("Test-Server")
		if err != nil {
			fmt.Println(err)
			// os.Exit(1)
		}
	}

	fmt.Println(task.NodeInfo)

	// Allocate ports
	var port natural_deploy.Ports
	port.Port = "8000"
	task.ExposedPorts = append(task.ExposedPorts, &port)

	err = task.CreateTask()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// Adds tasks to the tracked list
	task.RegisterTask()

	natural_deploy.PrintTasks()

	//fmt.Println(task.PingProgress())

	// Create a 10 second delay from mapping port
	time.Sleep(7 * time.Second)

	fmt.Println(task.PingProgress())

	natural_deploy.PrintTasks()

	err = task.KillTask()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	time.Sleep(4 * time.Second)

	fmt.Println(task.PingProgress())

	natural_deploy.PrintTasks()

}
