# Natural Deploy

<h1 align="center">
  <br>
  <a href=""><img src="https://github.com/user-attachments/assets/9d3d2a0a-f24e-485d-8993-77ecf644bb9b" alt="p2prc" width="400"></a>
  <br>
</h1>

The aim of "Natural deploy" is to be a library which builds on top of [P2PRC](https://github.com/Akilan1999/p2p-rendering-computation/) to make it as easy as running a binary to deploy your custom tasks to your servers. This will gives you the flexibility to truly customise to build your own orchestration layer to run your programs across your cluster of servers which can be situated anywhere (i.e even behind NAT).

## How to install
The only requirement if you are compiling from source is to have the Go compiler.
We will show this example by demonstrating a simple example which consists of 3 computers:

### General instruction to run
It's important to use the P2PRC master branch for this library.
```go
go get github.com/Akilan1999/p2p-rendering-computation@master
```

### The first computer (Root node):
Which is has a static public address. This is required to traverse traffic for your computers
behind NAT.
- Create a file called root.go
- Copy the following below into root.go
```go
package main

import (
    "fmt"
    natural_deploy "github.com/Akilan1999/Natural-deploy"
    "github.com/Akilan1999/p2p-rendering-computation/client"
    "os"
)

func main() {

    // Checking if the lock file exists
    // if not configures the root node
    _, err := os.Stat("Root.lock")
    if err != nil {
        // Create the root node configuration. Function parameter (Node name
        // and port number)
        rootNode, err := natural_deploy.CreateRootNode("RootNode","3244")
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        fmt.Println("----- Root Node information -------")
        client.PrettyPrint(rootNode)
        fmt.Println("-----------------------------------")
        fmt.Println("Daemon starting .....")

        os.Create("Root.lock")
    }

    err = natural_deploy.RunDaemon()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
        return
    }
    fmt.Println("Daemon started .....")

    for {

    }
}
```
- Compile root.go
```bash
  go build root.go
  # Run the binary as a background process
  ./root &
```

- Store the output
```json
    {
       "Name": "RootNode",
       "MachineUsername": "",
       "IPV4": "<public ipv4 address>", // The public address of your root server
       "IPV6": "",
       "Latency": 0,
       "Download": 0,
       "Upload": 0,
       "ServerPort": "3244", // The default port of the root node
       "UDPServerPort": "",
       "BareMetalSSHPort": "22",
       "NAT": false,
       "EscapeImplementation": "",
       "ProxyServer": false,
       "UnSafeMode": false,
       "PublicKey": "",
       "CustomInformation": ""
    }
```

### The 2nd computer (The server behind NAT)
This could for instance be an old laptop which you would want to bring back to life.

- Create the file server.go
- Add the following to the server.go file. We are going call this server "Test-Server".
```go
  package main

import (
	"fmt"
	natural_deploy "github.com/Akilan1999/Natural-deploy"
	"os"
)

func main() {
	_, err := os.Stat("Regular.lock")
	if err != nil {
		// Create the root node to be added
		var rootnode natural_deploy.RootNode
		// make sure to fill in "<>" area with the appropriate information
		rootnode.IPAddress = "<public ipv4 address>"
		rootnode.Port = "<ServerPort>" // (i.e 8088 by default)

		err := natural_deploy.CreateRegularNode("Test-Server", &rootnode)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		os.Create("Regular.lock")
	}

	// Start the daemon
	err = natural_deploy.RunDaemon()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	for {

	}
}
```

- Compile server.go
  ``` bash
    go build server.go
    # Run the binary as a background process
    ./server &
  ```

### The 3rd Computer (Task machine)
This can be your current laptop which can remotely deploy tasks to your self hosted server.
We are going to run a simple python server on the remote machine called "Test-Server". We will
start the task ,which will deploy the python server from our local machine to the "Test-Server"
machine ,which will then start the server and will internally call P2PRC function calls to provide
an static ip and port number from the Root node. Therefore escaping NAT which allows you to self
host.
- Create the deploy script which we will call "test.sh" (The python server runs locally in 8000).
```bash
  nohup python3 -m http.server 8000 > server.log 2>&1 &
```
- Create a script to kill the deployed python server. We will call this "kill.sh".
```bash
  kill -9 $(lsof -t -i:8000)
```
- Create the file task.go
- Add the following to task.go
```go
      package main

    import (
        "fmt"
        natural_deploy "github.com/Akilan1999/Natural-deploy"
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
            rootnode.IPAddress = "<public ipv4 address>"
            rootnode.Port = "<ServerPort>"
            // Name of the task machine
            natural_deploy.CreateTaskMachine("Task", &rootnode)
        }

        // Create a 5 second delay for mapping port
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

        // Print the task out
        task.PrintTask()

        // Below is to kill the task
        err = task.KillTask()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
            return
        }

    }

```
- Compile server.go
  ``` bash
     go build task.go
     # Run the binary as a background process
     ./task &
  ```


## Documentation and setup
- https://nd.akilan.io
- [Motivation](https://nd.akilan.io/#org44f09d9)
- [Getting Started](https://nd.akilan.io/#orgdcf5c56)



<!-- ## The Solution:
<img width="778" height="834" alt="Proposal" src="https://github.com/user-attachments/assets/4ce16823-3be2-4b89-924c-94331e9a627e" /> -->


