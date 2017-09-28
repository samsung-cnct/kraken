package cmd



import (
	"strings"
	"testing"
	"bytes"
)

const(
	awsClusterName string = "dummyAWScluster"
	gkeClusterName string = "dummyGKEcluster"
	awsClusterKubeVersion string = "v1.7.6"
	gkeClusterKubeVersion string = "v1.7.5"
)

func readInConfig(array []byte) error {
	clusterConfig.SetConfigType("yaml")
	clusterConfig.AutomaticEnv()
	clusterConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	clusterConfig.SetEnvPrefix("krakenlib")


	return clusterConfig.ReadConfig(bytes.NewBuffer(array))
}


func TestGetFirstClusterName(t *testing.T){
	// AWS TEST
	if err := readInConfig(AWS_CONFIG_BYTE_ARRAY); err != nil {
		t.Error("For aws an error occured reading the config file: ",err)
	}

	clusterName := getFirstClusterName()

	if clusterName != awsClusterName {
		t.Error("For aws expected to find cluster name to be", awsClusterName, "but instead found it to be", clusterName)
	}

	// GKE TEST
	if err := readInConfig(GKE_CONFIG_BYTE_ARRAY); err != nil {
		t.Error("For gke an error occured reading the config file: ",err)
	}

	clusterName = getFirstClusterName()

	if clusterName != gkeClusterName {
		t.Error("For aws expected to find cluster name to be", gkeClusterName, "but instead found it to be", clusterName)
	}
}


func TestGetFirstClusterKubernetesVersion(t *testing.T) {
	// AWS TEST
	if err := readInConfig(AWS_CONFIG_BYTE_ARRAY); err != nil {
		t.Error("For aws an error occured reading the config file: ",err)
	}

	version, err := getFirstClusterKubernetesVersion("aws")
	if err != nil {
		t.Error("For aws, expected to find kubernetes version",awsClusterKubeVersion, "but found error: ", err)
	}

	if version != awsClusterKubeVersion {
		t.Error("For aws, expected to find kubernetes version", awsClusterKubeVersion, "but found instead", version)
	}

	// GKE TEST
	if err := readInConfig(GKE_CONFIG_BYTE_ARRAY); err != nil {
		t.Error("For aws an error occured reading the config file: ",err)
	}

	version, err = getFirstClusterKubernetesVersion("gke")
	if err != nil {
		t.Error("For aws, expected to find kubernetes version",gkeClusterKubeVersion, "but found error: ", err)
	}

	if version != gkeClusterKubeVersion {
		t.Error("For aws, expected to find kubernetes version", gkeClusterKubeVersion, "but found instead", version)
	}

}

func TestGetProviderFromConfig(t *testing.T) {
	// AWS TEST
	if err := readInConfig(AWS_CONFIG_BYTE_ARRAY); err != nil {
		t.Error("For aws an error occured reading the config file: ",err)
	}

	provider, err := getProviderFromConfig()
	if err != nil {
		t.Error("For aws, expected to find provider to be","aws", "but found error: ", err)
	}

	if provider != "aws" {
		t.Error("For aws, expected to find provider to be", "aws", "but found instead", provider)
	}

	// GKE TEST
	if err := readInConfig(GKE_CONFIG_BYTE_ARRAY); err != nil {
		t.Error("For aws an error occured reading the config file: ",err)
	}

	provider, err = getProviderFromConfig()
	if err != nil {
		t.Error("For gke, expected to find provider to be","gke", "but found error: ", err)
	}

	if provider != "gke" {
		t.Error("For gke, expected to find provider to be", "gke", "but found instead", provider)
	}

}
