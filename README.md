# k2cli
k2cli is a command-line interface for K2.

Auto-generated reference documentation can be found [here](docs/k2cli.md).

## Getting Started
This Getting Started guide describes a Kubernetes deployment to AWS. K2 currently also supports deployments to GKE, but not
by default.

### Requirements
Docker must be installed on the machine where you run k2cli, and your user must have permissions to run it.

### Installing/Fetching the official build
Installation on OSX can happen via Brew by:

```
brew tap 'samsung-cnct/homebrew-k2cli'
brew install k2cli
```

Otherwise, the latest official build can be found here: https://github.com/samsung-cnct/k2cli/releases. You should use the latest version unless you have a specific reason to use a different version.

### Building a configuration file
K2 uses a yaml configuration file for all aspects of the both the Kubernetes cluster and the infrastructure that is
running it. To build a generic AWS configuration file that has a large number of sensible defaults, you can run:
```
./k2cli generate
```
which will create a file at `${HOME}/.kraken/config.yaml`

Or you may specify a path with:
```
./k2cli generate ${HOME}/k2configs/
```
which will create a file at `${HOME}/k2configs/config.yaml`.

For a GKE configuration file, run:
```
./k2cli generate --provider gke
```

For an AWS cluster there are several fields that need to be set before this file can
be used:
*  **Cluster name**  All K2 clusters should have a unique name so their assets can be easily identified by humans in the
AWS console (no more than 13 characters). The cluster name is set in the `deployment.clusters.name` field.  This dotted notation refers to the hierarchical
structure of a yaml file where cluster is a sub field of deployment. This line is towards the bottom of the file in the `deployment` section.

The following fields are in the `definitions` section of the configuration file. 
In lieu of specifying all of the following, you may just put your credentials file and K2 will grab the authentication specs from there.
*  **AWS access key**  Your AWS access key is required for programmatic access to AWS. The field is named
`providerConfig.authentication.accessKey`. This can be either set to the literal value, or to an environment
variable that K2 will use.
*  **AWS access secret**  This is your AWS access secret that is paired to the above access key. This field is named
`providerConfig.authentication.accessSecret`. This can be either set to the literal value, or to an environment
variable that K2 will use.
*  **AWS credentials file**  This is your AWS credentials file that is paired to the below profile. The field is named
`providerConfig.authentication.credentialsFile`. This file and path must exist bind mounted to /root
inside the container, ie. ${HOME}/.aws/credentials.
*  **AWS credentials profile** This is the AWS credentials profile name used to select the credentials set from the credentials
file above.

When the required fields are set your configuration file is ready to go! The default file will create a production-ready
cluster with the following configuration:

Role | # | Type
--- | ---  | ---
Primary etcd cluster | 5 | t2.small
Events etcd cluster | 5 | t2.small
Master nodes | 3 | m3.medium
Cluster nodes | 3 | c4.large
Special nodes | 2 | m3.medium

We have chosen this configuration based on our own, and other's, publicly available research. It creates an underpowered cluster
in terms of cluster nodes, but that's an easy setting to change. The important point is to ensure that the control plane is
production quality.

### Creating your first cluster
To create your first cluster, run the following command. (This assumes you have a configuration built as described above.)
If you have used the default config location:
```
./k2cli cluster up
```
Or you may specify the location of the config file:
```
./k2cli cluster up ${HOME}/k2configs/config.yaml
```
This will take anywhere from five to twenty minutes depending on how AWS is feeling when you execute this command. When
complete the cluster will exist in its own VPC and will be accesible via the `tool` subcommands. The output artifacts
will be stored in the default location: `${HOME}/.kraken/<cluster name>`.

### Working with your cluster (using k2cli)
k2cli uses the K2 image (github.com/samsung_cnct/k2) for all of its operations. The K2 image ships with `kubectl` and `helm`
installed and through the `k2cli tool` subcommand you can access these tools. Using the `k2cli tool` subcommand helps ensure you are
using the correct version of the relevant CLI for your cluster.

`kubectl` (http://kubernetes.io/docs/user-guide/kubectl-overview/) is a CLI for working with a Kubernetes cluster. It is
used for deploying application, checking system status and more. See the linked documentation for more details.

`helm` (https://github.com/kubernetes/helm) is a CLI for deploying and packaging applications to deploy
to Kubernetes. See the linked documentation for more details.

#### Example usage - k2cli tool kubectl

If you have specified a path for your config.yaml, then you will need to include the `--config ${HOME}/path_to_config/config.yaml` option when running the following commands. Otherwise it will assume your config lives at `${HOME}/.kraken/config.yaml`

To see all nodes in your Kubernetes cluster (and specify path to config file):

```
./k2cli tool kubectl --config ${HOME}/k2configs/config.yaml get nodes
```

To see all installed applications across all namespaces:
```
./k2cli tool kubectl -- get pods --all-namespaces
```

#### Example usage - k2cli tool helm
To list all installed charts with default config.yaml location:
```
./k2cli tool helm list
```

To install the Kafka chart maintained by Samsung CNCT.
```
./k2cli tool helm install atlas/kafka
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

### Destroying the running cluster
While not something to be done in production, during development when you are done with your cluster (or with a quickstart) it's
best to clean up your resources. To destroy the running cluster from this guide, simply run:
```
./k2cli cluster down ${HOME}/k2configs/config.yaml
```

**Note:** if you have specified an '--output' directory during the creation command, make sure you specify it here or the cluster
will still be running!

## Note on using environment variables in yaml configuration

k2cli will automatically attempt to expand all ```$VARIABLE_NAME``` strings in your configuration file. It will pass the variable and value to the K2 Docker container, and mount the path (if it's a path to an existing file or folder) into the K2 Docker container.
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

and the K2 Docker container would get a ```/Users/kraken/.ssh/keyfile.json:/Users/kraken/.ssh/keyfile.json``` mount and a ```K2_SERVICE_ACCOUNT_KEYFILE=/Users/kraken/.ssh/keyfile.json``` environment variable

### Automatic binding of environment variables

Environment variables with a `K2` prefix can also bind automatically to configuration values.

For example, given that ```export K2_DEPLOYMENT_CLUSTER=production-cluster```, the configuration:

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
and from a feature branch in your personal fork of the k2cli repo.

### To build
This is a go project with vendored dependencies so building is a snap.

```
git clone https://github.com/<your github account>/k2cli.git
cd k2cli
go build
```

This will create a k2cli binary that can be executed directly like so:

```
./k2cli
```

## Cutting a release

* Install github-release from https://github.com/c4milo/github-release
* Create a github personal access token with repo read/write permissions and export it as GITHUB_TOKEN
* Adjust VERSION and TYPE variables in the [Makefile](Makefile) as needed
* Run ```make release```
