# k8snet
Web server that ouputs a graph of kubernetes cluster network.

## Ouput Example:
![alt text](https://raw.githubusercontent.com/yeitany/k8snet/master/docs/images/output_example.png)

## Installation:
Possible to run locally using binary or from source.

helm chart installation (soon)

## Usage

### GET /graph
Request Parameters

|Request Parameter|Allowed Values|Default|
|--|--|--|
|format|svg,png,jpg,dot|png
|layout|circo,dot,fdp,neato,osage,patchwork,sfdp,twopi|circo
