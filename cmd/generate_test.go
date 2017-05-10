package cmd

import (
	"testing"
)

func TestProviderFlagDefaultsToAWS(t *testing.T) {
	if provider != "aws" {
		t.Error("Provider should default to 'aws'")
	}
}

func TestPreRunSetsVariables(t *testing.T) {
	args := make([]string, 2)
	args[0] = "$HOME/sandbox/fake-config.tmp"
	preRunEFunc(nil, args)

	if configPath == "" || generatePath == "" {
		t.Error("Expected configPath and generatePath to be set")
	}
}

func TestAWSProviderCreatesAWSconfig(t *testing.T) {
	preRunEFunc(nil, nil)
	if configPath != "ansible/roles/kraken.config/files/config.yaml " {
		t.Error("Expected generated config to be config.yaml")
	}
}

func TestGKEProviderCreatesGKEConfig(t *testing.T) {
	provider = "gke"
	preRunEFunc(nil, nil)
	if configPath != "ansible/roles/kraken.config/files/gke-config.yaml " {
		t.Error("Expected generated config to be gke-config.yaml")
	}
}
