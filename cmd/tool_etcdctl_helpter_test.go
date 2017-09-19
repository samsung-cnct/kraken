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

func TestEtcdctlArguments (t *testing.T) {

	for _, api := range allowedApiVersions {
		etcdctlArugmentTester(api, t)
	}
}

func etcdctlArugmentTester(version int, t *testing.T)  {
	var envVars []string
	var volumes []string
	var flags []string
	var certFile string
	var caFile string
	var keyFile string

	switch version {
	case 2:
		// check that krakenCerts are true, and no cluster name
		etcdUseKrakenCerts = true
		etcdClusterName = ""
		envVars, volumes, flags = etcdctlArguments(version)
		//envVar tests
		if len(envVars) != 1 {
			t.Error("Expected that only", 1, "environment variable found")
		}

		if envVars[0] != envVarApiVersion2 {
			t.Error("Expected to find environment variable", envVarApiVersion2, "but found", envVars[0])
		}

		// volumes
		if len(volumes) != 1 {
			t.Error("Expected that only", 1, "volume variable found")
		}

		if volumes[0] != "/etc/etcd/ssl:/etc/etcd/ssl" {
			t.Error("Expected to find volume", "/etc/etcd/ssl/:/etc/etcd/ssl", "but found", volumes[0])
		}

		// flags
		if len(flags) != 3 {
			t.Error("Expected only", 3, "flags, found", flags)
		}

		certFile = "--cert-file /etc/etcd/ssl/client.pem"
		caFile = "--ca-file /etc/etcd/ssl/client-ca.pem"
		keyFile = "--key-file /etc/etcd/ssl/client-key.pem"

		if flags[0] != certFile {
			t.Error("Expected to find the flag for --cert-file to be ", certFile, "but found ",flags[0])
		}

		if flags[1] != caFile {
			t.Error("Expected to find the flag for --ca-file to be ", caFile, "but found ",flags[1])
		}

		if flags[2] != keyFile {
			t.Error("Expected to find the flag for --key-file to be ", keyFile, "but found ",flags[2])
		}

		// check that krakenCerts are true, and a cluster name
		etcdUseKrakenCerts = true
		etcdClusterName = "etcdEvents"
		envVars, volumes, flags = etcdctlArguments(version)
		//envVar tests
		if len(envVars) != 1 {
			t.Error("Expected that only", 1, "environment variable found")
		}

		if envVars[0] != envVarApiVersion2 {
			t.Error("Expected to find environment variable", envVarApiVersion2, "but found", envVars[0])
		}

		// volumes
		if len(volumes) != 2 {
			t.Error("Expected only", 2, "volume variable, found", volumes)
		}

		if volumes[0] != "/etc/etcdEvents/ssl:/etc/etcdEvents/ssl" && volumes[0] != "/etc/etcd/etcdEvents:/etc/etcd/etcdEvents" {
			t.Error("Expected to find a volume variable", "/etc/etcdEvents/ssl:/etc/etcdEvents/ssl", "or", "/etc/etcd/etcdEvents:/etc/etcd/etcdEvents", "but found", volumes[0])
		}

		if volumes[1] != "/etc/etcd/etcdEvents:/etc/etcd/etcdEvents" && volumes[1] != "/etc/etcdEvents/ssl:/etc/etcdEvents/ssl" {
			t.Error("Expected to find a volume variable", "/etc/etcdEvents/ssl:/etc/etcdEvents/ssl", "or", "/etc/etcd/etcdEvents:/etc/etcd/etcdEvents", "but found", volumes[1])
		}

		// flags
		if len(flags) != 3 {
			t.Error("Expected that only", 3, "volume variable found")
		}

		certFile = "--cert-file /etc/etcdEvents/ssl/client-certfile.pem"
		caFile = "--ca-file /etc/etcdEvents/ssl/client-ca.pem"
		keyFile = "--key-file /etc/etcd/etcdEvents/client-keyfile.pem"

		if flags[0] != certFile {
			t.Error("Expected to find the flag for --cert-file to be '", certFile, "' but found '",flags[0],"'")
		}

		if flags[1] != caFile {
			t.Error("Expected to find the flag for --ca-file to be '", caFile, " but found '",flags[1],"'")
		}

		if flags[2] != keyFile {
			t.Error("Expected to find the flag for --key-file to be '", keyFile, "' but found '",flags[2],"'")
		}

	case 3:
		// check that krakenCerts are true, and no cluster name
		etcdUseKrakenCerts = true
		etcdClusterName = ""
		envVars, volumes, flags = etcdctlArguments(version)
		//envVar tests
		if len(envVars) != 1 {
			t.Error("Expected that only", 1, "environment variable found")
		}

		if envVars[0] != envVarApiVersion3 {
			t.Error("Expected to find environment variable", envVarApiVersion3, "but found", envVars[0])
		}

		// volumes
		if len(volumes) != 1 {
			t.Error("Expected that only", 1, "volume variable found")
		}

		if volumes[0] != "/etc/etcd/ssl:/etc/etcd/ssl" {
			t.Error("Expected to find environment variable", "/etc/etcd/ssl/:/etc/etcd/ssl", "but found", volumes[0])
		}

		// flags
		if len(flags) != 3 {
			t.Error("Expected only", 3, "flags, found", flags)
		}

		certFile = "--cert /etc/etcd/ssl/client.pem"
		caFile = "--cacert /etc/etcd/ssl/client-ca.pem"
		keyFile = "--key /etc/etcd/ssl/client-key.pem"

		if flags[0] != certFile {
			t.Error("Expected to find the flag for --cert to be ", certFile, "but found ",flags[0])
		}

		if flags[1] != caFile {
			t.Error("Expected to find the flag for --cacert to be ", caFile, "but found ",flags[1])
		}

		if flags[2] != keyFile {
			t.Error("Expected to find the flag for --key to be ", keyFile, "but found ",flags[2])
		}

		// check that krakenCerts are true, and a cluster name
		etcdUseKrakenCerts = true
		etcdClusterName = "etcdEvents"
		envVars, volumes, flags = etcdctlArguments(version)
		//envVar tests
		if len(envVars) != 1 {
			t.Error("Expected that only", 1, "environment variable found")
		}

		if envVars[0] != envVarApiVersion3 {
			t.Error("Expected to find environment variable", envVarApiVersion3, "but found", envVars[0])
		}

		// volumes
		if len(volumes) != 2 {
			t.Error("Expected only", 2, "volume variable, found", volumes)
		}

		if volumes[0] != "/etc/etcdEvents/ssl:/etc/etcdEvents/ssl" && volumes[0] != "/etc/etcd/etcdEvents:/etc/etcd/etcdEvents" {
			t.Error("Expected to find a volume variable", "/etc/etcdEvents/ssl:/etc/etcdEvents/ssl", "or", "/etc/etcd/etcdEvents:/etc/etcd/etcdEvents", "but found", volumes[0])
		}

		if volumes[1] != "/etc/etcd/etcdEvents:/etc/etcd/etcdEvents" && volumes[1] != "/etc/etcdEvents/ssl:/etc/etcdEvents/ssl" {
			t.Error("Expected to find a volume variable", "/etc/etcdEvents/ssl:/etc/etcdEvents/ssl", "or", "/etc/etcd/etcdEvents:/etc/etcd/etcdEvents", "but found", volumes[1])
		}

		// flags
		if len(flags) != 3 {
			t.Error("Expected that only", 3, "volume variable found")
		}

		certFile = "--cert /etc/etcdEvents/ssl/client-certfile.pem"
		caFile = "--cacert /etc/etcdEvents/ssl/client-ca.pem"
		keyFile = "--key /etc/etcd/etcdEvents/client-keyfile.pem"

		if flags[0] != certFile {
			t.Error("Expected to find the flag for --cert to be '", certFile, "' but found '",flags[0],"'")
		}

		if flags[1] != caFile {
			t.Error("Expected to find the flag for --cacert to be '", caFile, " but found '",flags[1],"'")
		}

		if flags[2] != keyFile {
			t.Error("Expected to find the flag for --key to be '", keyFile, "' but found '",flags[2],"'")
		}
	}
}


func TestCreateDockerCommand(t *testing.T) {
	envVars := []string{"a"}
	volumes := []string{"v1:v1", "v2:v2"}
	command := "backup file.txt"

	expectedDockerCommand := fmt.Sprintf("docker run -e a -v v1:v1 -v v2:v2 %s:%s /usr/local/bin/etcdctl %s", etcdDockerImage, defaultEtcdDockerVersion, command)

	result := createDockerCommand(envVars, volumes, etcdDockerImage, defaultEtcdDockerVersion, command)

	if expectedDockerCommand != result {
		t.Errorf("Expected command to be; ", expectedDockerCommand, "; but found it to be instead:", result)
	}

}