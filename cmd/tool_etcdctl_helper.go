package cmd

import (
	"fmt"
	"os"
	path2 "path"
	"path/filepath"
	"strings"
)

const (
	etcdDockerImage               string = "quay.io/coreos/etcd"
	defaultEtcdDockerVersion      string = "v3.0.17"
	defaultKrakenEtcdSnapshotFile string = "/tmp/snapshots/snapshot.db"
	defaultEtcdSSL                string = "/etc/etcd/ssl"
	envETCDCTL_API                string = "ETCDCTL_API"
	envETCDCTL_CACERT             string = "ETCDCTL_CA_FILE"
	envETCDCTL_CERT               string = "ETCDCTL_CERT_FILE"
	envETCDCTL_KEY                string = "ETCDCTL_KEY_FILE"
	envSSH_AUTH_SOCK              string = "SSH_AUTH_SOCK"
	krakenSSH                     string = "/kraken/bin/ssh_command.sh"
)

var allowedApiVersions = []int{2, 3}

// createDockerCommand assumes that the environmentVars will be of the form NAME=val and volumes of the form VOLUME:VOLUME
// this command is tailored to etcdctl
func createDockerCommand(environmentVars []string, volumes []string, etcdClientImage string, etcdClientVersion string, command string) string {
	var flatVolumes string
	var flatEnvVars string

	if volumes != nil || len(volumes) > 0 {
		flatVolumes = flattenVolumes(volumes)
	}

	if environmentVars != nil || len(environmentVars) > 0 {
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

func flattenFlags(flags []string) string {
	return strings.Join(flags, " ")
}

func flattenArgsToCommand(command string, args []string) string {
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

// etcdctlArgs returns environment variables, volume variables, and etcd flags
// this wil return empty values, and not nils.
func etcdctlArgs(apiVersion int) ([]string, []string, []string) {
	flags := []string{}
	volumes := []string{}
	environmentVars := []string{}

	if etcdEndpoints != "" {
		flags = append(flags, fmt.Sprintf("--endpoints=%s", etcdEndpoints))
	}

	volumes = append(volumes, volumeMountFmt(defaultEtcdSSL, ""))

	if snapshotSaveFile != "" {
		volumes = append(volumes, volumeMountFmt(path2.Dir(snapshotSaveFile), ""))
	}

	environmentVars = append(environmentVars, fmt.Sprintf("%s=%d", envETCDCTL_API, apiVersion))
	environmentVars = append(environmentVars, fmt.Sprintf("%s=%s", envETCDCTL_CACERT, "/etc/etcd/ssl/client-ca.pem"))
	environmentVars = append(environmentVars, fmt.Sprintf("%s=%s", envETCDCTL_CERT, "/etc/etcd/ssl/client.pem"))
	environmentVars = append(environmentVars, fmt.Sprintf("%s=%s", envETCDCTL_KEY, "/etc/etcd/ssl/client-key.pem"))

	return environmentVars, volumes, flags
}

func runEtcdCtl(args []string) (int, error) {
	var hostString []string
	var krakenSSHBase []string
	envVars, volumes, flags := etcdctlArgs(apiVersion)

	if verbosity {
		fmt.Printf("volumes passed: %v \n", volumes)
		fmt.Printf("environment variables passed: %v \n", envVars)
		fmt.Printf("flags passed: %v \n", flags)
		fmt.Printf("arguments passed: %v \n", args)

		krakenSSHBase = []string{krakenSSH, "--verbose"}
	} else {
		krakenSSHBase = []string{krakenSSH}
	}

	// here we append top level flags before subcommands
	command := flattenFlags(flags)

	// append flags belonging to subcommands
	if args != nil {
		command = flattenArgsToCommand(command, args)
	}

	dockerCommand := createDockerCommand(envVars, volumes, etcdDockerImage, etcdVersion, command)

	authSocket := os.Getenv(envSSH_AUTH_SOCK)
	sshConfigPath := fmt.Sprintf("%s/%s/ssh_config", filepath.Dir(ClusterConfigPath), getFirstClusterName())
	sshConfigPathDir := filepath.Dir(sshConfigPath)

	additionalVolumes = []string{volumeMountFmt(sshConfigPathDir, "")}
	if authSocket !="" {
		additionalEnvVars = []string{fmt.Sprintf("%s=%s", envSSH_AUTH_SOCK, authSocket)}
		additionalVolumes = append(additionalVolumes, volumeMountFmt(authSocket, ""))

	}

	if hostName == "" {
		hostString = []string{"--user", hostUsername, "--address", hostAddr, "--port", fmt.Sprintf("%d", hostPort)}
	} else {
		hostString = []string{"--hostname", hostName}
	}

	krakenSSHBase = append(krakenSSHBase, hostString...)
	sshCommand := append(krakenSSHBase, []string{"--ssh-config", sshConfigPath, dockerCommand}...)

	if verbosity {
		fmt.Printf("command passed: %v \n", sshCommand)
	}

	onFailure := func(out []byte) {
		fmt.Printf("%s \n", out)
	}

	onSuccess := func(out []byte) {
		fmt.Printf("%s \n", out)
	}

	return runKrakenLibCommandNoSpinner(sshCommand, ClusterConfigPath, onFailure, onSuccess)
}

// checkIfSnapshotSave returns true if using snapshot subCommand is found,
// the file should come up after those two subcommands as long as they are not files
func checkIfSnapshotSubcommand(subCommand string, args []string) (bool, string) {
	isSnapshot := false
	isSubcmd := false
	file := ""

	for _, k := range args {
		if !isSnapshot && k == "snapshot" {
			isSnapshot = true
			continue
		}

		if isSnapshot && !isSubcmd && k == subCommand {
			isSubcmd = true
			continue
		}

		if isSnapshot && isSubcmd && !strings.HasPrefix(k, "--") {
			file = k
			break
		}
	}

	return isSnapshot && isSubcmd, file
}
