package main

import (
    "fmt"
    natural_deploy "github.com/Akilan1999/P2PRC-natural-deploy"
    "os"
)

func main() {
    // Example output from Root Node
    //{
    //    "Name": "RootNode",
    //    "MachineUsername": "",
    //    "IPV4": "215.41.125.246",
    //    "IPV6": "",
    //    "Latency": 0,
    //    "Download": 0,
    //    "Upload": 0,
    //    "ServerPort": "8088",
    //    "UDPServerPort": "",
    //    "BareMetalSSHPort": "22",
    //    "NAT": false,
    //    "EscapeImplementation": "",
    //    "ProxyServer": false,
    //    "UnSafeMode": false,
    //    "PublicKey": "",
    //    "CustomInformation": ""
    //}

    _, err := os.Stat("Regular.lock")
    if err != nil {
        // Create the root node to be added
        var rootnode natural_deploy.RootNode
        // make sure to fill in "<>" area with the appropriate information
        rootnode.IPAddress = "<IPV4>"
        rootnode.Port = "<ServerPort>"

        err := natural_deploy.CreateRegularNode("Test-1", &rootnode)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
            return
        }
        os.Create("Regular.lock")
    }

    // Start a daemon
    err = natural_deploy.RunDaemon()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
        return
    }

    for {

    }
}
