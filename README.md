![logo](logo.png)

# kraken

[![pipeline status](https://git.cnct.io/common-tools/samsung-cnct_kraken/badges/master/pipeline.svg)](https://git.cnct.io/common/samsung-cnct_kraken/commits/master)

This document will help you get started deploying a high-availability Kubernetes cluster to AWS using kraken, a command-line interface for [kraken-lib](https://github.com/samsung-cnct/kraken-lib). kraken currently also supports deployments to GKE (see Building a Configuration File below).

## Prerequisites
Docker must be installed on the machine where you run kraken and your user must have permissions to run it.

**AWS Credentials:**
If deploying to AWS, the AWS User profile you wish to deploy under must have a policy attached with full access granted to:
* AmazonEC2FullAccess
* IAMFullAccess
* AmazonRoute53FullAccess

The User must also have "Programmatic access enabled" which will create an access key ID and secret access key which is required for kraken default `config.yaml`.

## Installing/Fetching the Official Build
You can install the official build on OS X via Brew by:

```
brew tap 'samsung-cnct/homebrew-k2cli'
brew install kraken
```

Otherwise, you can find the latest official build [here](https://github.com/samsung-cnct/kraken/releases). Use the latest version, unless you have a specific reason for using a different one.

## Building a Configuration File
kraken-lib uses a YAML configuration file for all aspects of the Kubernetes cluster and the infrastructure
running it. To build a generic AWS configuration file with a large number of sensible defaults, you can run:
```
kraken generate
```
This will create a file at `${HOME}/.kraken/config.yaml`

**Note:** If a config file already exists, the `generate` subcommand will fail with the message: `A kraken config file already exists at <your config location> - rename, delete or move it to generate a new default kraken config file`

Or you can specify a path with:
```
kraken generate ${HOME}/krakenlibconfigs/
```
This will create a file at `${HOME}/krakenlibconfigs/config.yaml`.

**For a GKE configuration file, run:**
```
kraken generate --provider gke
```

### Required configuration changes
For an AWS cluster, you need to set several fields before using the config file file you created, as listed below.
*  **Cluster name**  All kraken clusters should have a unique name so their assets can be easily identified by humans in the AWS console (no more than 32 characters). Set the cluster name in the `deployment.clusters.name` field. This dotted notation refers to the hierarchical structure of a YAML file where the cluster is a sub field of deployment. Find this line near the bottom of the file in the `deployment` section. GKE cluster names must use lower case letters.

The following fields are in the `definitions` section of the configuration file.
In lieu of specifying all of the following, you can simply put your credentials in the AWS credentials file from where kraken will access them.
*  **AWS access key**: required for programmatic access to AWS. The field is named
`providerConfigs.authentication.accessKey` that you can set to the literal value or to an environment
variable that kraken will use.
*  **AWS access secret**:  paired to the above access key. This field is named
`providerConfigs.authentication.accessSecret` that you can set to the literal value or to an environment
variable that kraken will use.
*  **AWS credentials file**: paired with the below profile. The field is named
`providerConfigs.authentication.credentialsFile`. This file and path must exist bind-mounted to /root
inside the container, (${HOME}/.aws/credentials).
*  **AWS credentials profile**: used to select the credentials set from the credentials
file above.

When you've set the required fields, your configuration file is ready to go. The default file will create a production-ready cluster with the following configuration:

Role | # | Type
--- | ---  | ---
Primary etcd cluster | 5 | t2.small
Events etcd cluster | 5 | t2.small
Master nodes | 3 | m4.large
Cluster nodes | 10 | c4.large
Special nodes | 2 | m4.large

We have chosen this configuration based on our own and others' publicly available research. It creates an underpowered cluster
in terms of cluster nodes, which is an easy setting to change (see below). The important point is to ensure the control plane is production quality.

### Optional configuration changes (more advanced)
First-time users looking to set up a simple evaluation cluster can skip this section and go directly to [Creating Your First Cluster](#creating-your-first-cluster).  

You can modify many options to control the deployment of your Kubernetes cluster. Here we focus on a couple that may be of interest before starting your first cluster. For reference, here is the [full set of kraken configuration options](https://samsung-cnct.github.io/kraken-lib/).

*  **Deployment Region and Availability Zones**  
In the default-generated configuration file, all clusters begin their lives in the AWS Region us-east-1. You can move the default region and  modify the availability zones, if needed. For reference, the [Global AWS Infrastructure](https://aws.amazon.com/about-aws/global-infrastructure/) provides a complete list of regions and availability zones. These fields are named `definitions.providerConfigs.region`, and `definitions.providerConfigs.subnet.az` respectively. Note three total `.subnet.az` values are defined, so the cluster can be spread across multiple failure domains. Be sure to update all three to availability zones within your selected region.  

*  **Cluster Node Size and Count**  
This setting defines the type and total number of cluster nodes on which you can schedule workloads. The default EC2 instance type is `c4.large`, and the default cluster includes 10 of these instances. These fields are named  `definitions.nodeConfigs.defaultAwsClusterNode.providerConfig.type` for the instance type and `deployment.clusters.nodePools.name: clusterNodes`, `.count: 10`, for the total number of worker nodes. When setting the total number of nodes, keep mind they will be spread across all of your cluster's availability zones.

### How to create a small research or development cluster
To create a small, low resource-consuming cluster, alter your config to the following:

Role | # | Type
--- | ---  | ---
Primary etcd cluster | 1 | t2.small
Events etcd cluster | 1 | t2.small
Master nodes | 1 | m4.large
Cluster nodes | 1 | c4.large
~~Special~~ ~~nodes~~ | ~~2~~ | ~~m4.large~~

Delete 'Special nodes'

YAML:

```deployment:
  clusters:
    - name:
...
      nodePools:
        - name: etcd
          count: 1
...
    - name: etcdEvents
          count: 1
...
        - name: master
          count: 1
...
        - name: clusterNodes
          count: 1
```


## Creating Your First Cluster
Assuming you have a configuration built (as described above), you're ready to create your first cluster. Run the following command.
If you have used the default config location:
```
kraken cluster up
```
Or you can specify the location of the config file:
```
kraken cluster up --config ${HOME}/krakenlibconfigs/config.yaml
```
This will take anywhere from 5 to 20 minutes, depending on AWS performance when you execute this command. When
complete, the cluster exists in its own VPC and is accessible via the `tool` subcommands. The output artifacts
are stored in the default location: `${HOME}/.kraken/<cluster name>`.

## Working with Your Cluster (Using kraken)
For all of its operations, kraken uses the [kraken-lib image](https://quay.io/samsung_cnct/kraken-lib) that ships with the installed `kubectl` and `helm`. You can access these tools through the `kraken tool` subcommand. Using this subcommand helps ensure you're using the correct version of the relevant CLI for your cluster.

`kubectl` (http://kubernetes.io/docs/user-guide/kubectl-overview/), a CLI for working with a Kubernetes cluster, is
used for deploying applications, checking system status and more. See the linked documentation for more details.

`Helm` (https://github.com/kubernetes/helm) is a CLI for packaging and deploying applications to Kubernetes. See the linked documentation for more details.

### Example usage - kraken tool kubectl

If you've specified a path for your config.yaml, you need to include the `--config ${HOME}/path_to_config/config.yaml` option when running the following commands. Otherwise, it assumes your config lives at `${HOME}/.kraken/config.yaml`

To see all nodes in your cluster (and specify the path to the config file):

```
kraken tool kubectl --config ${HOME}/krakenlibconfigs/config.yaml get nodes
```

To see all installed applications across all namespaces:
```
kraken tool kubectl --config ${HOME}/krakenlibconfigs/config.yaml -- get pods --all-namespaces
```

### Example usage - kraken tool Helm
To list all installed charts with the default config.yaml location:
```
kraken tool helm list
```

To install the Kafka chart maintained by Samsung CNCT:
```
kraken tool helm install atlas/kafka
```

## Working with Your Cluster (Using Host-Installed Tools)
Your local machine's output directory stores the file needed by Helm and kubectl for connecting to and interacting with your Kubernetes deployment. By default, this directory is `${HOME}/.kraken/<cluster name>/`. The filename is `admin.kubeconfig`.

### Example usage - local kubectl
To list all nodes in your Kubernetes cluster:
```
kubectl --kubeconfig=${HOME}/.kraken/<cluster name>/admin.kubeconfig --cluster=<cluster name> get nodes
```

To list all installed applications across all namespaces:
```
kubectl --kubeconfig=${HOME}/.kraken/<cluster name>/admin.kubeconfig --cluster=<cluster name> get pods --all-namespaces
```

### Example usage - local Helm
Helm requires the admin.kubeconfig file and the saved local Helm state. The Helm state directory is also in the output
directory.

To list all installed charts:
```
KUBECONFIG=${HOME}/.kraken/<cluster name>/admin.kubeconfig HELM_HOME=${HOME}/.kraken/<cluster name>/.helm helm list
```

To install the Kafka chart maintained by Samsung CNCT:
```
KUBECONFIG=${HOME}/.kraken/<cluster name>/admin.kubeconfig HELM_HOME=${HOME}/.kraken/<cluster name>/.helm helm install atlas/kafka
```

## Updating your Cluster
With kraken, you can update all aspects of your node pools including count, Kubernetes version, instance type and more. To do so, please make desired changes in your configuration file, and then run kraken's cluster update command, as described below, pointing to your configuration file.

### Running kraken update
You can specify different versions of Kubernetes in each node pool. Note: this may affect the compatibility of your cluster's kraken-provided services. Specify which node pools you want to update with a comma-separated list of their names. This process takes approximately 10 minutes per node.
You may also add or remove entire nodepools.

- Step 1: Make appropriate changes to configuration file
- Step 2: Run
```
kraken cluster update --config ${HOME}/krakenlibconfigs/config.yaml --update-nodepools <your,nodepools,here>
```

Similarly,
```
--add-nodepools
```
will add new nodepools specified in configuration file, and
```
--rm-nodepools
```
will remove nodepools removed from your configuration file.

## Destroying the Running Cluster
When you're done with your cluster or with a quickstart, we recommend cleaning up your resources by destroying the running cluster. From this guide, simply run:
```
kraken cluster down ${HOME}/krakenConfigs/config.yaml
```

**Note:** If you have specified an '--output' directory during the creation command, be sure to specify it here or the cluster will still be running!

# Using Environment Variables in YAML Configuration

kraken will automatically attempt to expand all ```$VARIABLE_NAME``` strings in your configuration file. It will pass the variable and value to the kraken-lib Docker container and mount the path (if it's a path to an existing file or folder) into the kraken-lib Docker container.

## Environment Variable Expansion

For example, given a variable such as ```export MY_PRIVATE_KEY_FILE=/Users/kraken/.ssh/id_rsa```, the configuration:

```
definitions:
  ...
  keyPairs:
   - &defaultKeyPair
      name: defaultKeyPair
      kind: keyPair
      publickeyFile: "$HOME/.ssh/id_rsa.pub"
      privatekeyFile: "$MY_PRIVATE_KEY_FILE"
...
```

will be expanded to:

```
definitions:
  ...
  keyPairs:
   - &defaultKeyPair
      name: defaultKeyPair
      kind: keyPair
      publickeyFile: "$HOME/.ssh/id_rsa.pub"
      privatekeyFile: "/Users/kraken/.ssh/id_rsa"

...
```

and the KRAKENLIB Docker container will get a ```/Users/kraken/.ssh/id_rsa:/Users/kraken/.ssh/id_rsa``` mount and a ```KRAKENLIB_PRIVATE_KEY_FILE=/Users/kraken/.ssh/id_rsa``` environment variable.

If you have further questions or needs, please read through the rest of the documentation and then open an issue.

# Contributing Features, Bug Fixes and More
We welcome all types of contributions from the community and and don't require a contributor license agreement. To simplify merging, we prefer pull requests based on a feature branch in your personal fork that's based off the current master of the kraken repo. For more details, please refer to our [kraken-lib Contributing](https://github.com/samsung-cnct/kraken-lib/blob/master/CONTRIBUTING.md) document.

## To build
This is a go project with vendored dependencies, so building is a snap.

```
git clone https://github.com/<your github account>/kraken.git
cd kraken
go build
```

This will create a kraken binary that can be executed directly like so:

```
./kraken
```

### Asset changes
Assets  are stored in the `/data` directory of this project's directory. Any file changes only get implemented if you
follow the steps below:

* Run `go-bindata data/`
* Move the generated `bindata.go` file to the /cmd directory
* Change package from `main` to `cmd`
* Commit the changes

We plan to automate this process in the future.

## Cutting a release

Please speak to a member of the kraken team in #kraken Slack (link below) if you need a release cut.

# Additional Resources
Here are some additional resources you might find useful:

* #kraken Slack on [slack.k8s.io](http://slack.k8s.io/)
* [kraken-lib issue tracker](https://github.com/samsung-cnct/kraken-lib/issues)
* [kraken-tools](https://github.com/samsung-cnct/kraken-tools)
* [kraken codebase](https://github.com/samsung-cnct/kraken)

# Maintainer

This document is maintained by Patrick Christopher (@coffeepac) at Samsung SDS.
