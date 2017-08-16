# Kraken
Kraken is a command-line interface for [Krakenlib](https://github.com/samsung-cnct/k2).

## Getting Started
This Getting Started guide describes a Kubernetes deployment to AWS. Krakenlib currently also supports deployments to GKE, but not
by default.

### Requirements
Docker must be installed on the machine where you run Kraken, and your user must have permissions to run it.

### Installing/Fetching the Official Build
Installation on OSX can happen via Brew by:

```
brew tap 'samsung-cnct/homebrew-k2cli'
brew install k2cli
```

Otherwise, the latest official build can be found here: https://github.com/samsung-cnct/k2cli/releases. You should use the latest version unless you have a specific reason to use a different version.

### Building a Configuration File
Krakenlib uses a yaml configuration file for all aspects of the both the Kubernetes cluster and the infrastructure that is
running it. To build a generic AWS configuration file that has a large number of sensible defaults, you can run:
```
kraken generate
```
which will create a file at `${HOME}/.kraken/config.yaml`  **Note:** If a config file already exists the `generate` subcommand will fail with the message: `A Krakenlib config file already exists at <your config location> - rename, delete or move it to generate a new default Krakenlib config file`

Or you may specify a path with:
```
kraken generate ${HOME}/krakenlibconfigs/
```
which will create a file at `${HOME}/krakenlibconfigs/config.yaml`.

For a GKE configuration file, run:
```
kraken generate --provider gke
```

#### Required Configuration Changes
For an AWS cluster there are several fields that need to be set before this file can be used:
*  **Cluster name**  All Krakenlib clusters should have a unique name so their assets can be easily identified by humans in the
AWS console (no more than 13 characters). The cluster name is set in the `deployment.clusters.name` field.  This dotted notation refers to the hierarchical structure of a yaml file where cluster is a sub field of deployment. This line is towards the bottom of the file in the `deployment` section.

The following fields are in the `definitions` section of the configuration file.
In lieu of specifying all of the following, you may just put your credentials file and Krakenlib will grab the authentication specs from there.
*  **AWS access key**  Your AWS access key is required for programmatic access to AWS. The field is named
`providerConfigs.authentication.accessKey`. This can be either set to the literal value, or to an environment
variable that Krakenlib will use.
*  **AWS access secret**  This is your AWS access secret that is paired to the above access key. This field is named
`providerConfigs.authentication.accessSecret`. This can be either set to the literal value, or to an environment
variable that Krakenlib will use.
*  **AWS credentials file**  This is your AWS credentials file that is paired to the below profile. The field is named
`providerConfigs.authentication.credentialsFile`. This file and path must exist bind mounted to /root
inside the container, ie. ${HOME}/.aws/credentials.
*  **AWS credentials profile** This is the AWS credentials profile name used to select the credentials set from the credentials
file above.

When the required fields are set your configuration file is ready to go! The default file will create a production-ready
cluster with the following configuration:

Role | # | Type
--- | ---  | ---
Primary etcd cluster | 5 | t2.small
Events etcd cluster | 5 | t2.small
Master nodes | 3 | m4.large
Cluster nodes | 10 | c4.large
Special nodes | 2 | m4.large

We have chosen this configuration based on our own, and other's, publicly available research. It creates an underpowered cluster
in terms of cluster nodes, but that's an easy setting to change (see below). The important point is to ensure that the control plane is
production quality.

#### Optional Configuration Changes (More Advanced)
First time users looking to setup a simple evaluation cluster can skip this section and go directly to [Creating Your First Cluster](#creating-your-first-cluster).  

There are a great many options that you can modify to control the deployment of your Kubernetes cluster.  Here we introduce only a couple options that may be of interest to you before starting your first cluster.  

*  **Deployment Region and Availability Zones**  
In the default generated configuration file, all clusters begin their lives in the AWS Region US East. This can be modified to locate the cluster in a region or to modify the availability zones to a value that may be more suitable to individual use cases.  For reference, the [Global AWS Infrastructure](https://aws.amazon.com/about-aws/global-infrastructure/) provides a complete list of Regions and Availability Zones.  These fields are named `definitions.providerConfigs.region`, and `definitions.providerConfigs.subnet.az` respectively.  Note that there are actually three total `.subnet.az` values defined so that the cluster can be spread across multiple failure domains.  You should update all three to availability zones that are within the region selected.  

*  **Cluster Node Size and Count**  
This setting affects the type and total number of nodes in your cluster that can be used to schedule workloads on.  The default EC2 Instance type is `c4.large` and the default cluster includes 10 of these instances.  These fields are named  `definitions.nodeConfigs.defaultAwsClusterNode.providerConfig.type` for the instance type, and `deployment.clusters.nodePools.name: clusterNodes`, `.count: 10`, for the total number of worker nodes.  Keep in mind when setting the total number of nodes that these nodes will be spread across all the availability zones in your cluster.

#### To Create a Small Research or Development Cluster
To spinup a small low-resource consuming cluster, alter your config to the following:

Role | # | Type
--- | ---  | ---
Primary etcd cluster | 1 | t2.small
Events etcd cluster | 1 | t2.small
Master nodes | 1 | m4.large
Cluster nodes | 1 | c4.large
~~Special~~ ~~nodes~~ | ~~2~~ | ~~m4.large~~

Delete 'Special nodes'

yaml:

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


### Creating Your First Cluster
To create your first cluster, run the following command. (This assumes you have a configuration built as described above.)
If you have used the default config location:
```
kraken cluster up
```
Or you may specify the location of the config file:
```
kraken cluster up ${HOME}/krakenlibconfigs/config.yaml
```
This will take anywhere from five to twenty minutes depending on how AWS is feeling when you execute this command. When
complete the cluster will exist in its own VPC and will be accesible via the `tool` subcommands. The output artifacts
will be stored in the default location: `${HOME}/.kraken/<cluster name>`.

### Working with Your Cluster (using Kraken)
Kraken uses the Krakenlib image (github.com/samsung_cnct/k2) for all of its operations. The Krakenlib image ships with `kubectl` and `helm`
installed and through the `kraken tool` subcommand you can access these tools. Using the `kraken tool` subcommand helps ensure you are
using the correct version of the relevant CLI for your cluster.

`kubectl` (http://kubernetes.io/docs/user-guide/kubectl-overview/) is a CLI for working with a Kubernetes cluster. It is
used for deploying application, checking system status and more. See the linked documentation for more details.

`helm` (https://github.com/kubernetes/helm) is a CLI for deploying and packaging applications to deploy
to Kubernetes. See the linked documentation for more details.

#### Example Usage - Kraken tool kubectl

If you have specified a path for your config.yaml, then you will need to include the `--config ${HOME}/path_to_config/config.yaml` option when running the following commands. Otherwise it will assume your config lives at `${HOME}/.kraken/config.yaml`

To see all nodes in your Kubernetes cluster (and specify path to config file):

```
kraken tool kubectl --config ${HOME}/krakenlibconfigs/config.yaml get nodes
```

To see all installed applications across all namespaces:
```
kraken tool kubectl --config ${HOME}/krakenlibconfigs/config.yaml -- get pods --all-namespaces
```

#### Example usage - Kraken tool helm
To list all installed charts with default config.yaml location:
```
kraken tool helm list
```

To install the Kafka chart maintained by Samsung CNCT.
```
kraken tool helm install atlas/kafka
```

### Working with your cluster (using host installed tools)
The file that is required for both helm and kubectl to connect to, and interact with, your Kubernetes deployment is saved to your
local machine in the output directory. By default, this directory is `${HOME}/.kraken/<cluster name>/`. The filename is `admin.kubeconfig`.

#### Example usage - local kubectl
To list all nodes in your Kubernetes cluster
```
kubectl --kubeconfig=${HOME}/.kraken/<cluster name>/admin.kubeconfig --cluster=<cluster name> get nodes
```

To list all installed applications across all namespaces
```
kubectl --kubeconfig=${HOME}/.kraken/<cluster name>/admin.kubeconfig --cluster=<cluster name> get pods --all-namespaces
```

#### Example usage - local helm
Helm requires both the admin.kubeconfig file and the saved local helm state. The helm state directory is also in the output
directory

To list all installed charts
```
KUBECONFIG=${HOME}/.kraken/<cluster name>/admin.kubeconfig HELM_HOME=${HOME}/.kraken/<cluster name>/.helm helm list
```

To install the Kafka chart maintained by Samsung CNCT.
```
KUBECONFIG=${HOME}/.kraken/<cluster name>/admin.kubeconfig HELM_HOME=${HOME}/.kraken/<cluster name>/.helm helm install atlas/kafka
```

### Updating your cluster
You may update your nodepools with Kraken, specifically the Kubernetes version, the nodepool counts and instance types. To do so, please make desired changes in your configuration file, and then run Kraken's cluster update command, as described below, pointing to your configuration file.

#### Running Kraken update
You can specify different versions of Kubernetes in each nodepool. Note: this may affect the compatibility of your cluster's Krakenlib-provided services. Specify which nodepools you wish to update with a comma-separated list of the names of the nodepools. Please be patient; this process may take a while; about ten minutes per node.

- Step 1: Make appropriate changes to configuration file
- Step 2: Run
```bash
kraken cluster update ${HOME}/krakenlibconfigs/config.yaml <your,nodepools,here>
```

### Destroying the running cluster
While not something to be done in production, during development when you are done with your cluster (or with a quickstart) it's
best to clean up your resources. To destroy the running cluster from this guide, simply run:
```
kraken cluster down ${HOME}/kkrakenlibconfigs/config.yaml
```

**Note:** if you have specified an '--output' directory during the creation command, make sure you specify it here or the cluster
will still be running!

## Note on using environment variables in yaml configuration

Kraken will automatically attempt to expand all ```$VARIABLE_NAME``` strings in your configuration file. It will pass the variable and value to the Krakenlib Docker container, and mount the path (if it's a path to an existing file or folder) into the Krakenlib Docker container.
### Environment variable expansion

For example, given a variable such as ```export MY_SERVICE_ACCOUNT_KEYFILE=/Users/kraken/.ssh/keyfile.json```, the configuration:

```
deployment:
  cluster: production-cluster
  keypair:
    -
      name: key
      publickeyFile:
      privatekeyFile:
      providerConfig:
        username:
        serviceAccount: "serviceaccount@project.iam.gserviceaccount.com"
        serviceAccountKeyFile: "$MY_SERVICE_ACCOUNT_KEYFILE"

...
```

will be expanded to:

```
deployment:
  cluster: production-cluster
  keypair:
    -
      name: key
      publickeyFile:
      privatekeyFile:
      providerConfig:
        username:
        serviceAccount: "serviceaccount@project.iam.gserviceaccount.com"
        serviceAccountKeyFile: "/Users/kraken/.ssh/keyfile.json"

...
```

and the KRAKENLIB Docker container would get a ```/Users/kraken/.ssh/keyfile.json:/Users/kraken/.ssh/keyfile.json``` mount and a ```KRAKENLIB_SERVICE_ACCOUNT_KEYFILE=/Users/kraken/.ssh/keyfile.json``` environment variable

### Automatic binding of environment variables

Environment variables with a `KRAKENLIB` prefix can also bind automatically to configuration values.

For example, given that ```export KRAKENLIB_DEPLOYMENT_CLUSTER=production-cluster```, the configuration:

```
deployment:
  cluster: changeme
  keypair:
    -
      name: key
      publickeyFile:
      privatekeyFile:
      providerConfig:
        username:
        serviceAccount: "serviceaccount@project.iam.gserviceaccount.com"
        serviceAccountKeyFile: "/Users/kraken/.ssh/keyfile.json"

...
```

will be expanded to:

```
deployment:
  cluster: production-cluster
  keypair:
    -
      name: key
      publickeyFile:
      privatekeyFile:
      providerConfig:
        username:
        serviceAccount: "serviceaccount@project.iam.gserviceaccount.com"
        serviceAccountKeyFile: "/Users/kraken/.ssh/keyfile.json"

...
```

If you have further questions or needs please read through the rest of the docs and then open an issue!


## Developing features or bugfixes
We accept pull requests from all users and do not require a contributor license agreement. To simplify merging we prefer pull requests that are based on the current Master
and from a feature branch in your personal fork of the Kraken repo.

### To build
This is a go project with vendored dependencies so building is a snap.

```
git clone https://github.com/<your github account>/k2cli.git
cd k2cli
go build
```

This will create a Kraken binary that can be executed directly like so:

```
./kraken
```

## Cutting a release

* Install github-release from https://github.com/c4milo/github-release
* Create a github personal access token with repo read/write permissions and export it as GITHUB_TOKEN
* Adjust VERSION and TYPE variables in the [Makefile](Makefile) as needed or set them as command line parmaters to `make`
* Run ```make release``` or with paramaters ```make release VERSION=v0.1```
