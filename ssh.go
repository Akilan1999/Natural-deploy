package natural_deploy

import (
	"errors"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

// SSHInfo consisting of information required for a SSH
type SSHInfo struct {
	SSHCommon      *SSHInfoCommon
	PortNo         uint
	Username       string
	Host           string
	ConnectionType string // If the connection is TCP or UDP (Currently we only support TCP)
}

// SSHInfoCommon Common information required to do SSH
type SSHInfoCommon struct {
	PrivateKey string // Path to private key
	PublicKey  string // Path to the public key
	Passphrase string // Passphrase
	Password   string // Plain text password

}

// SSHIntoNode SSH into node based on the information on the other struct
func (sshInfo *SSHInfo) SSHIntoNode() (*goph.Client, error) {
	// Declaring variables of type auth and error
	var auth goph.Auth
	var err error

	// Check if the private key or password is provided
	// Private key takes priority to password
	if sshInfo.SSHCommon.PrivateKey != "" {
		// Start new ssh connection with private key.
		auth, err = goph.Key(sshInfo.SSHCommon.PrivateKey, sshInfo.SSHCommon.Passphrase)
		if err != nil {
			return nil, err
		}
	} else if sshInfo.SSHCommon.Password != "" { // If a plaintext password is provided
		auth = goph.Password(sshInfo.SSHCommon.Password)
	} else {
		return nil, errors.New("no Private or password not provided")
	}

	client, err := New(sshInfo.Username, sshInfo.Host, sshInfo.PortNo, auth)
	if err != nil {
		return nil, err
	}

	// Defer closing the network connection.
	//defer client.Close()

	return client, err
}

// New Modified function to of the library Goph to ensure the user can provide an SSH
// port number during the connection.
// New starts a new ssh connection, the host public key must be in known hosts.
func New(user string, addr string, port uint, auth goph.Auth) (c *goph.Client, err error) {
	_, err = goph.DefaultKnownHosts()

	if err != nil {
		return
	}

	c, err = goph.NewConn(&goph.Config{
		User:     user,
		Addr:     addr,
		Port:     port,
		Auth:     auth,
		Timeout:  goph.DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	return
}
