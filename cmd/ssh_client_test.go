package cmd

import (
	"testing"
)

func TestNewSSHClient(t *testing.T) {
	// ensure we get the object we want.
	username := "some_user"
	address := "dummy_address"
	port := 0
	sshClient := newSSHClient(username, address, port)

	if sshClient.Username != username {
		t.Errorf("Username not properly set, expected", username, "and got", sshClient.Username)
	}

	if sshClient.HostName != address {
		t.Errorf("Username not properly set, expected", address, "and got", sshClient.HostName)
	}

	if sshClient.Port != port {
		t.Errorf("Username not properly set, expected", port, "and got", sshClient.Port)
	}

	if sshClient.client != nil {
		t.Errorf("new client object should have been nil")
	}


}

func TestSshSession(t *testing.T) {
	// create a session that should fail. check if objects fail to create and error message sent.
	username := "some_user"
	address := "dummy_address"
	port := 0
	sshClient := newSSHClient(username, address, port)
	session, err := sshClient.sshSession()

	// client should fail to create
	if sshClient.client != nil {
		t.Errorf("ssh client unexpectedly is not nil")
	}

	// session should not be created
	if session != nil && err != nil {
		t.Errorf("ssh session and error are unexpectedly bot not nil")
	}
}

