package cmd

import (
	"os"
	"strconv"

	"crypto/tls"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/tlsconfig"
	"path/filepath"
)

// DockerAPIVersion defines the docker api version used.
var DockerAPIVersion = client.DefaultVersion

// DockerClientConfig provides a simple encapsulation of parameters to construct the Docker API client
type DockerClientConfig struct {
	DockerHost       string
	DockerAPIVersion string
	TLSEnabled       bool
	TLSVerify        bool
	TLSCACertificate string
	TLSCertificate   string
	TLSKey           string
}

// GetDefaultHost produces either the environment-provided host, or a sensible default.
func (conf *DockerClientConfig) GetDefaultHost() string {
	env := os.Getenv("DOCKER_HOST")
	if env == "" {
		return client.DefaultDockerHost
	}
	return env
}

// GetDefaultTLSVerify indicates whether TLS is enabled by the current environment.
func (conf *DockerClientConfig) GetDefaultTLSVerify() bool {
	env := os.Getenv("DOCKER_TLS_VERIFY")
	if (env == "") || (env == "0") {
		return false
	}
	return true
}

// GetDefaultDockerAPIVersion produces either the environment-provided Docker API version, or a sensible default.
func (conf *DockerClientConfig) GetDefaultDockerAPIVersion() string {
	env := os.Getenv("DOCKER_API_VERSION")
	if env == "" {
		return client.DefaultVersion
	}
	return env
}

// GetDefaultTLSCertificatePath produces either the environment-provided path to TLS certificates, or a sensible default.
func (conf *DockerClientConfig) GetDefaultTLSCertificatePath() string {
	env := os.Getenv("DOCKER_CERT_PATH")
	if env == "" {
		return os.ExpandEnv("${HOME}/.docker/")
	}
	return env
}

// GetDefaultTLSCACertificate produces the path to the environment-configured CA certificate for TLS verification.
func (conf *DockerClientConfig) GetDefaultTLSCACertificate() string {
	return filepath.Join(conf.GetDefaultTLSCertificatePath(), "ca.pem")
}

// GetDefaultTLSCertificate produces the path to the environment-configured TLS certificate.
func (conf *DockerClientConfig) GetDefaultTLSCertificate() string {
	return filepath.Join(conf.GetDefaultTLSCertificatePath(), "cert.pem")
}

// GetDefaultTLSKey produces the path to the environment-configured TLS key.
func (conf *DockerClientConfig) GetDefaultTLSKey() string {
	return filepath.Join(conf.GetDefaultTLSCertificatePath(), "key.pem")
}

// Was this config derived solely by OS environment?
// The properties of `conf` will be overridden only by command line args,
// otherwise they're given the values of the associated Default methods.
func (conf *DockerClientConfig) isInheritedFromEnvironment() bool {

	// This is structured like a table-test, to improve readability. Huzzah!
	compare := map[string][]string{
		"version": {conf.DockerAPIVersion, conf.GetDefaultDockerAPIVersion()},
		"host":    {conf.DockerHost, conf.GetDefaultHost()},
		"verify":  {strconv.FormatBool(conf.TLSVerify), strconv.FormatBool(conf.GetDefaultTLSVerify())},
		"cacert":  {conf.TLSCACertificate, conf.GetDefaultTLSCACertificate()},
		"cert":    {conf.TLSCertificate, conf.GetDefaultTLSCertificate()},
		"key":     {conf.TLSKey, conf.GetDefaultTLSKey()},
	}

	for _, val := range compare {
		if val[0] != val[1] {
			return false
		}
	}

	return true

}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsExist(err)
}

func (conf *DockerClientConfig) isTLSActivated() bool {
	return (conf.TLSEnabled || conf.TLSVerify) && fileExists(conf.TLSCACertificate) && fileExists(conf.TLSCertificate) && fileExists(conf.TLSKey)
}

func (conf *DockerClientConfig) createTLSConfig() (*tls.Config, error) {
	return tlsconfig.Client(tlsconfig.Options{
		CAFile:             conf.TLSCACertificate,
		CertFile:           conf.TLSCertificate,
		KeyFile:            conf.TLSKey,
		InsecureSkipVerify: !(conf.TLSVerify),
	})
}
