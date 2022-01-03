package natural_deploy

import (
	"fmt"
	"testing"
)

// Test to test if SSH into node works
func TestSSHInfo_SSHIntoNode(t *testing.T) {
	var sshInfo SSHInfo
	var sshInfoCommon SSHInfoCommon

	// Common SSH info
	sshInfoCommon.Password = "password"
	sshInfo.SSHCommon = &sshInfoCommon
	sshInfo.Username = "master"
	sshInfo.PortNo = 41031
	sshInfo.Host = "0.0.0.0"

	// Trying to SSH into the node
	node, err := sshInfo.SSHIntoNode()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(node.Client.LocalAddr().String())

}
