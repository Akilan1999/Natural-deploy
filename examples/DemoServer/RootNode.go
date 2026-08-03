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
		// Create the root node configuration
		rootNode, err := natural_deploy.CreateRootNode("RootNode")
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
