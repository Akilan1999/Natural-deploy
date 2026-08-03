package natural_deploy

import (
	"errors"
	"github.com/Akilan1999/p2p-rendering-computation/abstractions"
	"github.com/Akilan1999/p2p-rendering-computation/config"
	"github.com/Akilan1999/p2p-rendering-computation/config/generate"
	"github.com/Akilan1999/p2p-rendering-computation/p2p"
)

// This file consists of the wrappers for P2PRC
// In this implementation the root only traverses traffic
// The proxy information assumes to hold true

func CommonSetup(MachineName string) error {
	// P2PRC configuration
	abstractions.Init(nil)

	Config, err := config.ConfigInit(nil, nil)
	if err != nil {
		return err
	}

	// Changing the name of the machine
	Config.MachineName = MachineName

	err = Config.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func CreateRootNode(MachineName string) (*p2p.IpAddress, error) {
	err := CommonSetup(MachineName)
	if err != nil {
		return nil, err
	}

	Config, err := config.ConfigInit(nil, nil)
	if err != nil {
		return nil, err
	}

	// Root nodes cannot be behind NAT
	Config.BehindNAT = false

	err = Config.WriteConfig()
	if err != nil {
		return nil, err
	}

	// ---------------- creating a root node ------------------
	var rootnode p2p.IpAddress
	var rootnodes p2p.IpAddresses

	rootnode.Ipv4, err = p2p.CurrentPublicIP()
	if err != nil {
		return nil, err
	}

	rootnode.ServerPort = Config.ServerPort
	rootnode.Name = Config.MachineName
	rootnode.NAT = false
	if Config.ProxyPort != "" {
		rootnode.ProxyServer = true
	}
	rootnode.EscapeImplementation = ""
	// Assumes SSH runs in port 22
	if Config.BareMetal {
		rootnode.BareMetalSSHPort = "22"
	}

	rootnodes.IpAddress = append(rootnodes.IpAddress, rootnode)

	// ----------------------------------------------------

	// Creates the root node entry
	generate.GenerateIPTableFile(rootnodes.IpAddress)

	return &rootnode, nil
}

type RootNode struct {
	IPAddress string
	Port      string
}

// CreateRegularNode This is for the default server to allow servers to connect through
// to run processes
func CreateRegularNode(MachineName string, node *RootNode) error {
	err := CommonSetup(MachineName)
	if err != nil {
		return err
	}

	if node != nil {
		err = abstractions.AddRootNode(node.IPAddress, node.Port)
		if err != nil {
			return err
		}
	}

	// Changing config settings
	Config, err := config.ConfigInit(nil, nil)
	if err != nil {
		return err
	}

	// To allow other nodes in network to SSH into this machine using
	// their private key
	Config.BareMetal = true
	Config.BareMetal = true

	err = Config.WriteConfig()
	if err != nil {
		return err
	}

	return nil

}

func CreateTaskMachine(MachineName string, node *RootNode) error {
	err := CreateRegularNode(MachineName, node)
	if err != nil {
		return err
	}

	if node != nil {
		err = abstractions.AddRootNode(node.IPAddress, node.Port)
		if err != nil {
			return err
		}
	}

	// Changing config settings
	Config, err := config.ConfigInit(nil, nil)
	if err != nil {
		return err
	}

	// To allow other nodes in network to SSH into this machine using
	// their private key
	Config.BareMetal = false
	Config.BareMetal = false

	err = Config.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// SearchMachine search for a machine on your p2p network
func SearchMachine(MachineName string) (*p2p.IpAddress, error) {
	table, err := p2p.ReadIpTable()
	if err != nil {
		return nil, err
	}

	for i, _ := range table.IpAddress {
		if table.IpAddress[i].Name == MachineName {
			return &table.IpAddress[i], nil
		}
	}
	return nil, errors.New("machine name not found")
}

func RunDaemon() error {
	_, err := abstractions.Start()
	if err != nil {
		return err
	}
	return nil
}
