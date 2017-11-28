package cmd

import (
	"fmt"
	"testing"
)

const (
	envVarApiVersion2 string = "ETCDCTL_API=2"
	envVarApiVersion3 string = "ETCDCTL_API=3"
)

func TestIsApiVersionSupported(t *testing.T) {
	for _, v := range allowedApiVersions {
		if !isApiVersionSupported(v) {
			t.Error("Expected that (", v, ") is a supported version")
		}
	}

	unsupported := []int{0, 1, 4, 20}
	for _, v := range unsupported {
		if isApiVersionSupported(v) {
			t.Error("Expected that (", v, ") is an unsupported version")
		}
	}
}

func TestEtcdctlArguments(t *testing.T) {
	var envVars []string
	var volumes []string
	var flags []string

	// Version = 2 Case
	envVars, volumes, flags = etcdctlArgs(2)
	if len(envVars) != 4 {
		t.Error("Expected that", 4, "environment variables be found")
	}

	if envVars[0] != envVarApiVersion2 {
		t.Error("Expected to find environment variable", envVarApiVersion2, "but found", envVars[0])
	}

	if len(volumes) != 1 {
		t.Error("Expected that only", 1, "volume variable be found")
	}

	if volumes[0] != "/etc/etcd/ssl:/etc/etcd/ssl" {
		t.Error("Expected to find volume", "/etc/etcd/ssl:/etc/etcd/ssl", "but found", volumes[0])
	}

	if len(flags) != 0 {
		t.Error("Expected only", 1, "flag, but found", flags)
	}

	// Version 3 case
	snapshotSaveFile = "/path/to/my_file.db"
	mount := "/path/to:/path/to"

	envVars, volumes, flags = etcdctlArgs(3)

	if envVars[0] != envVarApiVersion3 {
		t.Error("Expected to find environment variable", envVarApiVersion2, "but found", envVars[0])
	}

	if len(volumes) != 2 {
		t.Error("Expected that", 2, "volume variable be found")
	}

	if volumes[1] != mount {
		t.Error("Expected to find volume", mount, "but found", volumes[1])
	}

}

func TestCreateDockerCommand(t *testing.T) {
	envVars := []string{"a"}
	volumes := []string{"v1:v1", "v2:v2"}
	command := "backup file.txt"

	expectedDockerCommand := fmt.Sprintf("docker run -e a -v v1:v1 -v v2:v2 %s:%s /usr/local/bin/etcdctl %s", etcdDockerImage, defaultEtcdDockerVersion, command)

	result := createDockerCommand(envVars, volumes, etcdDockerImage, defaultEtcdDockerVersion, command)

	if expectedDockerCommand != result {
		t.Error("Expected command to be: ", expectedDockerCommand, "; but found it to be instead:", result)
	}
}

func TestCheckIfSnapshotSubcommand(t *testing.T) {
	ok, file := checkIfSnapshotSubcommand("save", []string{"snapshot", "blah", "blah"})
	if ok {
		t.Error("Expected not to be a snapshot save but found: ", ok)
	}

	ok, file = checkIfSnapshotSubcommand("save", []string{"snapshot", "blah", "restore"})
	if ok {
		t.Error("Expected not to be a snapshot save but found: ", ok)
	}

	ok, file = checkIfSnapshotSubcommand("save", []string{"snapshot", "blah", "save"})
	if ok && file != "" {
		t.Error("Expected snapshot save's file to be empty but found it to be", file)
	}

	ok, file = checkIfSnapshotSubcommand("restore", []string{"snapshot", "blah", "restore"})
	if !ok {
		t.Error("Expected not to be a snapshot restore but found: ", ok)
	}

	if file != "" {
		t.Error("Expected snapshot save's file to be empty but found it to be", file)
	}

	ok, file = checkIfSnapshotSubcommand("status", []string{"snapshot", "blah", "status", "/tmp/path/file2"})
	if !ok {
		t.Error("Expected not to be a snapshot status but found: ", ok)
	}

	if file != "/tmp/path/file2" {
		t.Error("Expected snapshot save's file to be", "/tmp/path/file2", "but found it to be", file)
	}
}
