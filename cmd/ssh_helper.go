package cmd

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

func runSSHCommand(command string, session *ssh.Session) (int, error) {
	if session == nil {
		return -1, fmt.Errorf("cannot pass commands to an empty session")
	}
	// connect session pipes stdin, stdout, stderr
	stdin, err := session.StdinPipe()
	if err != nil {
		return -1, fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return -1,fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return -1,fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	if verbosity {
		fmt.Printf("running command: %s \n", command)
	}

	err = session.Run(command)
	if err != nil{
		return -1, err
	}

	return 0, err
}