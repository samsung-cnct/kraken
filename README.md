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

