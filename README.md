# k8snet
Discover your kubernetes cluster topology using simple web server

## Output Example:
![alt text](https://raw.githubusercontent.com/yeitany/k8snet/master/docs/images/output_example.png)

## Installation:
From source

## Usage

### GET /graph
Request Parameters

|Request Parameter|Allowed Values|Default|
|--|--|--|
|format|svg,png,jpg,dot|png
|layout|circo,dot,fdp,neato,osage,patchwork,sfdp,twopi|circo
|targets|cluster namespaces|empty
