# k2cli
CLI for K2

Generated [docs](docs/k2cli.md)

## Build

```
git clone https://github.com/samsung-cnct/k2cli.git
cd k2cli
go build
```

## Run

```
./k2cli
```
## Use

Bring up a cluster configure by ~/.kraken/mycluster.yml, output results to ~/output/, time out after 600 seconds, use k2 image quay.io/myorg/k2:latest

```
./k2cli cluster up ~/.kraken/mycluster.yml --output ~/output --timeout 600 --image quay.io/myorg/k2:latest
```

Bring up a cluster configured by ~/.kraken/krakenCuster.yaml, output results to ~/.kraken/[cluster name]

```
./k2cli cluster up
```

Destroy cluster configured by ~/.kraken/mycluster.yml

```
./k2cli cluster down ~/.kraken/mycluster.yml
```

Get all nodes of cluster configured by ~/.kraken/mycluster.yml 

```
./k2cli tool kubectl --config ~/.kraken/mycluster.yml get nodes
```
Generate sensible AWS defaults config at ~/myconfig/

```
./k2cli generate ~/myconfig/
```

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
        serviceAccountKeyFile: "$K2_SERVICE_ACCOUNT_KEYFILE"
        
...
```

given that ```export K2_SERVICE_ACCOUNT_KEYFILE=/Users/kraken/.ssh/keyfile.json```

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

and the k2 container would get a /Users/kraken/.ssh/keyfile.json:/Users/kraken/.ssh/keyfile.json mount and K2_SERVICE_ACCOUNT_KEYFILE=/Users/kraken/.ssh/keyfile.json environment variable