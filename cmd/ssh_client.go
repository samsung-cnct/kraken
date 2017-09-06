package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SSHClient struct {
	Username string
	HostName string
	Port     int

	client *ssh.Client
}

func newSSHClient(username string, hostname string, port int) *SSHClient {
	return &SSHClient{username, hostname, port, nil}
}

func (s *SSHClient) sshAgent() (ssh.AuthMethod, error) {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers), nil
	}else {
		return nil, err
	}
}

func (s *SSHClient) sshConfig(username string, sshAuthMethod ssh.AuthMethod) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			sshAuthMethod,
		},
	}
}

// sshSession should create a session that uses a user's ssh-agent socket. can be extended to include
// keys if required. Please defer close the session when completed.
func (s *SSHClient) sshSession() (*ssh.Session, error) {
	var err error
	var session *ssh.Session

	if s.client == nil {
		authAgent, err := s.sshAgent()
		if err != nil {
			return nil, err
		}

		sshConfig := s.sshConfig(s.Username, authAgent)

		s.client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.HostName, s.Port), sshConfig)
		if err != nil {
			return nil, fmt.Errorf("Failed to dial: %s", err)
		}
	}

	if s.client != nil {
		session, err = s.client.NewSession()
		if err != nil {
			return nil, fmt.Errorf("Failed to create session: %s", err)
		}

	} else {
		err = fmt.Errorf("could not create ssh client")
	}

	return session, err
}

func (s *SSHClient) Close(session *ssh.Session) {
	if err := session.Close(); err != nil {
		log.Fatal(err)
	}
}