# k2cli
CLI for K2

Generated [docs](docs/k2cli.md)

## Getting Started
This getting started guide will assume a deployment to AWS.  K2 currently also supports deploys to GKE but that is not
the default.

### Fetching the official build
The latest official build can be found here:  https://github.com/samsung-cnct/k2cli/releases  If not doing development, 
you should use the latest version.

### Building a configuration file
K2 uses a yaml configuration file for all aspects of the both the kubernetes cluster and the infrastructure that is
running it.  To build a generic aws config that has a large number of sensible defaults, you can run: 
```
./k2cli generate ~/k2configs/
```
which will create a file at ~/k2configs/config.yaml.  There are three fields that need to be set before this file can
be used.
1. cluster name.  All k2 clusters should have a unique name so their assets can be easily identified by humans in the
aws console.  The cluster name is set at field 'deployment.cluster'.  This dotted notation refers to the hierarchical 
structure of a yaml file where cluster is a sub field of deployment.  This is also the third line in the generated file.
2. aws access key.  This is your aws access key used for programatic access to AWS.  The field is named
'deployment.providerConfig.authentication.accessKey'.  This can be either set to the actual value or a environment
variable and k2 will perform the replacement.
3. aws secret key.  This is your aws secret key that is paired to the above access key.  This field is named
'deployment.providerConfig.authentication.secretKey'.  This can be either set to the actual value or a environment
variable and k2 will perform the replacement.
node:  although the field 'deployment.providerConfig.authentication.credentialsFile' is present it is not fully functional
see:  https://github.com/samsung-cnct/k2/issues/128

When those fields are set your configuration file is ready to go!  The default file will create a production-ready
cluster with the following configuration:
Primary etcd cluster: 5 t2.small
Events etcd cluster: 5 t2.small
Master nodes: 3 m3.medium
Cluster nodes: 3 c4.large
Special nodes: 2 m3.medium

We have arrived at these numbers based on our own and other publically available research.  This is an underpowered cluster
from a number of cluster nodes view but that's an easy setting to change.  We wanted to ensure that the control plane was
production quality.

### Creating your first cluster
To create your first cluster, run the following command (assuming you followed the previous section on configuration building)
```
./k2cli cluster up ~/k2configs/config.yaml
```
This will take anywhere from five to twenty minutes depending on how AWS is feeling when you execute this command.  When 
complete the above cluster will exist in its own VPC and will be accesible via the 'tool' subcommands.  The output artifacts
will be stored in the default location of ~/.kraken/<cluster name>.  

### Working with your cluster (using k2cli)
k2cli uses the k2 image (github.com/samsung_cnct/k2) for all of its operations.  The k2 image ships with kubectl and helm 
installed and through the tool subcommand you can access these tools.  Using the k2cli tool subcommand helps ensure you are
using the correct version of the relevant cli for your cluster.

kubectl (http://kubernetes.io/docs/user-guide/kubectl-overview/) is a cli for working with a kubernetes cluster.  It is
used for deploying application, checking system status and more.  See the linked documentation for more details.

helm (http://kubernetes.io/docs/user-guide/kubectl-overview/) is a cli for deploying and packaging applications to deploy 
to kubernetes.  See the linked documentation for more details.

#### Example usage - k2cli tool kubectl
To see all nodes in your kuberentes cluster
```
./k2cli tool kubectl --config ~/k2configs/config.yaml get nodes 
```

To see all installed application accross all namespaces
```
./k2cli tool kubectl --config ~/k2configs/config.yaml -- get pods --all-namespaces
```

#### Example usage - k2cli tool helm
To list all installed charts
```
./k2cli tool helm --config ~/k2configs/config.yaml list
```

To install the samsung_cnct maintainted kafka chart
```
./k2cli tool helm --config ~/k2configs/config.yaml install atlas/kafka
```

### Working with your cluster (using host installed tools)
The file that is required for both helm and kubectl to connect to and interact with your kubernetes install is saved to your
local machine in the output directory.  By default, this directory is ~/.kraken/<cluster name>/.  The file is 'admin.kubeconfig'.

#### Example usage - local kubectl
To see all nodes in your kubernetes cluster
```
kubectl --kubeconfig=~/.kraken/<cluster name>/admin.kubeconfig --cluster=<cluster name> get nodes
```

To see all installed applications accross all namespaces
```
kubectl --kubeconfig=~/.kraken/<cluster name>/admin.kubeconfig --cluster=<cluster name> get pods --all-namespaces
```

#### Example usage - local helm
Helm requires both the admin.kubeconfig file and the saved local helm state.  The helm state directory is also in the output
directory

To see all installed charts
```
KUBECONFIG=~/.kraken/<cluster name>/admin.kubeconfig HELM_HOME=~/.kraken/<cluster name>/.helm helm list 
```

To install the samsung_cnct maintained kafka chart
```
KUBECONFIG=~/.kraken/<cluster name>/admin.kubeconfig HELM_HOME=~/.kraken/<cluster name>/.helm helm install atlas/kafka
```

### Destroying the running cluster
While not something to be done in production, during development when you are done with your cluster (or with a quickstart) its
best to clean up your resources.  To destroy the running cluster from this guide, simply run:
```
./k2cli cluster down ~/k2configs/config.yaml
```

Note:  if you have specified an '--output' directory during the creation command, make sure you specify it here or the cluster
will still be running!

## Note on using environment variables in yaml configuration

k2cli will automatically attempt to expand all ```$VARIABLE_NAME``` strings in your configuration file, pass the variable and value to the k2 docker container, and mount the path (if it's a path to an existing file or folder) into the k2 docker container.

For example:

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

given that ```export MY_SERVICE_ACCOUNT_KEYFILE=/Users/kraken/.ssh/keyfile.json```

Will expand to:

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

and the k2 container would get a ```/Users/kraken/.ssh/keyfile.json:/Users/kraken/.ssh/keyfile.json``` mount and ```K2_SERVICE_ACCOUNT_KEYFILE=/Users/kraken/.ssh/keyfile.json``` environment variable


Environment variables with 'K2' prefix can also bind automatically to config values. For example, given that ```export K2_DEPLOYMENT_CLUSTER=production-cluster```

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

will evaluate to 

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
We accept pull requests from all users and require no CLA.  To ease merging we prefer PRs that are based on the current Master
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
