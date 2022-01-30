package natural_deploy

import (
	"fmt"
	"github.com/melbahja/goph"
	"testing"
)

//---------------------------------------------------------------
// NOTE: All the SSH test cases will not work on your machine
// because they are hard coded with docker containers I
// am testing locally.
//---------------------------------------------------------------
// Starts an SSH connection in the test machine
func StartSSHConnection() (*goph.Client, error) {
	var sshInfo SSHInfo
	var sshInfoCommon SSHInfoCommon

	// Common SSH info
	sshInfoCommon.Password = "password"
	sshInfo.SSHCommon = &sshInfoCommon
	sshInfo.Username = "master"
	sshInfo.PortNo = 34621
	sshInfo.Host = "0.0.0.0"

	// Trying to SSH into the node
	node, err := sshInfo.SSHIntoNode()
	if err != nil {
		return nil, err
	}

	return node, nil
}

// Test case to ensure that the binary can run
// on a remote machine.
func TestCopyAndRunBinary(t *testing.T) {
	// Creates an SSH connection to the test machine
	connection, err := StartSSHConnection()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// Runs the hello world binary on another machine
	Output, err := CopyAndRunBinary("TestBinary/RunWebServer/", "main", connection, "")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(Output)
}
