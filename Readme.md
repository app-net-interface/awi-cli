

### Prerequisite

go version go1.20.5

### Building

Run `make build` on the root directory. awi executable should be created on the home directory.


### Usage:
```
./awi --help
CLI for connecting networking and application resources through Cisco Catalyst WAN

Usage:
  awi [command]

Available Commands:
  completion     Generate the autocompletion script for the specified shell
  create         Create resources
  delete         Delete resources
  generate-token Update config with generated tokens
  get            get resource
  help           Help about any command
  list           List resources

Flags:
  -c, --config string   Configuration file in YAML format (default "config.yaml")
  -h, --help            help for awi

Use "awi [command] --help" for more information about a command.

```

### Examples

#### Connecting two VPCs with matching Ids across any cloud

``` yaml
apiVersion: awi.cisco.awi/v1
kind: InterNetworkDomainConnection # Connection across network domains
metadata:
  name: "aws-infra-vpc-to-sandbox-vpc"       #generate an appropriate name
spec:
  source:
    metadata:
      name: "Infra VPC" #source network domain name
      description: ""
    networkDomain:
      selector: #Select network name based on the below selection criteria
        matchId:
          id: "vpc-067cfa335f9a2e657"
  destination:
    metadata:
      name: "Sandbox VPC" #Destination network domain name
      description: ""
    networkDomain:
      selector: #Select network name based on the below selection criteria
        matchId:
          id: "vpc-003643f14c9e5a38d"
```

#### Connecting two VPCs - source vpc labeld as "infra" and destination vpc labled as "sandbox" (across any cloud)

``` yaml
apiVersion: awi.cisco.awi/v1
kind: InterNetworkDomainConnection
metadata:
  name: "aws-infra-vpcs-to-sandbox-vpcs-labels"
  labels:
    awi_watching: true
spec:
  source:
    metadata:
      name: "Infra VPCs"
      description: ""
    networkDomain:
      selector:
        matchLabels:
          name: "infra"
  destination:
    metadata:
      name: "Sandbox VPCs"
      description: ""
    networkDomain:
      selector:
        matchLabels:
          env: "sandbox"
```

## Contributing

Thank you for interest in contributing! Please refer to our
[contributing guide](CONTRIBUTING.md).

## License

awi-infra-guard is released under the Apache 2.0 license. See
[LICENSE](./LICENSE).

awi-infra-guard is also made possible thanks to
[third party open source projects](NOTICE).
