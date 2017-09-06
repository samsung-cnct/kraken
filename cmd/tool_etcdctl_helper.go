package cmd

import (
	"fmt"
	path2 "path"
	"strings"

	"github.com/spf13/cobra"
)

const (
	etcdDockerImage          string = "quay.io/coreos/etcd"
	defaultEtcdDockerVersion string = "v3.2.5"
	envVarETCDCTL_API        string = "ETCDCTL_API"
)
var allowedApiVersions = []int { 2, 3 }

func preRunEtcdctl(cmd *cobra.Command, args []string) error{
	if !isApiVersionSupported(apiVersion) {
		return fmt.Errorf("The api-version %d is not supported.", apiVersion)
	}
	return nil
}

// createDockerCommand assumes that the environmentVars will be of the form NAME=val and volumes of the form VOLUME:VOLUME
// this command is tailored to etcdctl
func createDockerCommand(environmentVars []string, volumes []string, etcdClientImage string, etcdClientVersion string, command string) string {
	var flatVolumes string
	var flatEnvVars string

	if volumes != nil || len(volumes) >0 {
		flatVolumes = flattenVolumes(volumes)
	}

	if environmentVars != nil || len(environmentVars) >0 {
		flatEnvVars = flattenEnvVars(environmentVars)
	}

	return fmt.Sprintf("docker run %s %s %s:%s /usr/local/bin/etcdctl %s", flatEnvVars, flatVolumes, etcdClientImage, etcdClientVersion, command)
}

func flattenVolumes(volumes []string) string {
	return fmt.Sprintf("-v %s", strings.Join(volumes, " -v "))
}

func flattenEnvVars(volumes []string) string {
	return fmt.Sprintf("-e %s", strings.Join(volumes, " -e "))
}

func flattenFlags(flags [] string ) string {
	return strings.Join(flags, " ")
}

func runEtcdSSHDockerCommand(baseCommand string, flags []string, args []string, addVolumes []string ) (int, error) {
	sshClient := newSSHClient(hostUsername, hostAddr, hostPort)
	session, err := sshClient.sshSession()
	if err != nil {
		return -1, err
	}
	defer sshClient.Close(session)

	envVars, volumes, filteredFlags := etcdctlArguments(apiVersion)

	volumes = append(volumes, addVolumes...)
	flags = append(flags, filteredFlags...)

	if verbosity {
		fmt.Printf("volumes passed: %v \n", volumes)
		fmt.Printf("environment variables passed: %v \n", envVars)
		fmt.Printf("flags passed: %v \n", flags)
		fmt.Printf("arguments passed: %v \n", args)

	}

	// here we append top level flags before subcommands
	command := fmt.Sprintf("%s %s", flattenFlags(flags), baseCommand)

	// append flags belonging to subcommands
	if args != nil {
		command = flattenArgsToCommand(command, args)
	}

	dockerCommand := createDockerCommand(envVars, volumes, etcdDockerImage, etcdVersion, command)

	return runSSHCommand(dockerCommand, session)
}

func flattenArgsToCommand(command string, args []string) string{
	return fmt.Sprintf("%s %s", command, strings.Join(args, " "))
}

func isApiVersionSupported(version int) bool {
	for _, v := range allowedApiVersions {
		if v == version {
			return true
		}
	}

	return false
}

// etcdctlArguments returns environment variables, volume variables, and etcd flags
// this wil return empty values, and not nils.
func etcdctlArguments(apiVersion int) ([] string, []string, []string) {
	volumesMap := map[string]string{}
	flags := []string{}
	volumes := []string{}
	environmentVars := []string{}

	if etcdClusterName == "" {
		etcdCertFile = "/etc/etcd/ssl/client.pem"
		etcdCaFile = "/etc/etcd/ssl/client-ca.pem"
		etcdKeyFile = "/etc/etcd/ssl/client-key.pem"
	}



	if etcdUseKrakenCerts && etcdClusterName != "" {
		etcdCertFile = fmt.Sprintf("/etc/%s/ssl/client-certfile.pem", etcdClusterName)
		etcdCaFile = fmt.Sprintf("/etc/%s/ssl/client-ca.pem", etcdClusterName)
		etcdKeyFile = fmt.Sprintf("/etc/etcd/%s/client-keyfile.pem", etcdClusterName)
	}

	if etcdEndpoints != "" {
		flags = append(flags, fmt.Sprintf("--endpoints=%s", etcdEndpoints ))
	}

	switch apiVersion {
	case 2:
		if etcdCertFile != "" {
			flags = append(flags, fmt.Sprintf("--cert-file %s", etcdCertFile ))
			addToVolumeMapIfKeyNotExist(volumesMap, path2.Dir(etcdCertFile))
		}

		if etcdCaFile != "" {
			flags = append(flags, fmt.Sprintf("--ca-file %s", etcdCaFile ))
			addToVolumeMapIfKeyNotExist(volumesMap, path2.Dir(etcdCaFile))
		}

		if etcdKeyFile != "" {
			flags = append(flags, fmt.Sprintf("--key-file %s", etcdKeyFile ))
			addToVolumeMapIfKeyNotExist(volumesMap, path2.Dir(etcdKeyFile))
		}
	case 3:
		if etcdCertFile != "" {
			flags = append(flags, fmt.Sprintf("--cert %s", etcdCertFile ))
			addToVolumeMapIfKeyNotExist(volumesMap, path2.Dir(etcdCertFile))
		}

		if etcdCaFile != "" {
			flags = append(flags, fmt.Sprintf("--cacert %s", etcdCaFile ))
			addToVolumeMapIfKeyNotExist(volumesMap, path2.Dir(etcdCaFile))
		}

		if etcdKeyFile != "" {
			flags = append(flags, fmt.Sprintf("--key %s", etcdKeyFile ))
			addToVolumeMapIfKeyNotExist(volumesMap, path2.Dir(etcdKeyFile))
		}
	}

	for _, v := range volumesMap {
		volumes = append(volumes, v)
	}

	environmentVars = append(environmentVars, fmt.Sprintf("%s=%d", envVarETCDCTL_API, apiVersion))

	return environmentVars, volumes, flags
}

func addToVolumeMapIfKeyNotExist(m map[string]string, key string) {
	if _, ok := m[key]; !ok {
		m[key] = fmt.Sprintf("%s:%s", key, key)
	}
}
