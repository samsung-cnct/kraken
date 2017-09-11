package cmd

import (
	"fmt"
	"testing"
)

const randConfigNameLenght = 12

func TestProviderFlagDefaultsToAWS(t *testing.T) {
	if provider != "aws" {
		t.Error("Provider should default to 'aws'")
	}
}

func TestPreRunSetsVariables(t *testing.T) {
	args := make([]string, 2)
	args[0] = fmt.Sprintf("$HOME/sandbox/%s.tmp", randStringBytesMaskImprSrc(randConfigNameLenght))
	preRunEFunc(nil, args)

	if configPath == "" || generatePath == "" {
		t.Error("Expected configPath and generatePath to be set")
	}
}

func TestAWSProviderCreatesAWSconfig(t *testing.T) {
	args := make([]string, 2)
	args[0] = fmt.Sprintf("$HOME/sandbox/%s.tmp", randStringBytesMaskImprSrc(randConfigNameLenght))

	preRunEFunc(nil, args)

	if configPath != "ansible/roles/kraken.config/files/config.yaml" {
		t.Error("Expected the ansible config path to point to config.yaml")
	}
}

func TestGKEProviderCreatesGKEConfig(t *testing.T) {
	provider = "gke"
	args := make([]string, 2)
	args[0] = fmt.Sprintf("$HOME/sandbox/%s.tmp", randStringBytesMaskImprSrc(randConfigNameLenght))

	preRunEFunc(nil, args)

	if configPath != "ansible/roles/kraken.config/files/gke-config.yaml" {
		t.Error("Expected the ansible config path to point to gke-config.yaml")
	}
}
