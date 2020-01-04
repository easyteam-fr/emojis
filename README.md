# Emoji

Emoji is a sample application to demonstrate how to build a Kubernetes Operator.
The operator monitor the Emoji Custom Resources and publish them in an API. To
know more about this application, check Programmez #237 on
[programmez.com](https://www.programmez.com)

## Design

The application design looks like below:

![Design](img/emoji-design.png)

- An Emoji CRD allow to declare Emojis in Kubernetes
- A controller get all the Emojis from Kubernetes and synchronize them with the
  application. Once done, it publishes the status back to the resource.

## Building and installing the application

To build the application, you must have make, go, kustomize and docker installed
on your instance. Set the `REGISTRY` variable and run `make build-docker` like
in the example below:

```shell
cd app
export REGISTRY=registry:5000
make build-docker
```

To install the application, set the `REGISTRY` variable, run change the image
settings in the `kustomization.yaml` file and run a `kubectl apply` command:

```shell
export REGISTRY=registry:5000
kustomize edit set image emojis-app=${REGISTRY}/emojis-app:latest
kubectl apply -k .
```

To test the application, connect to port 8080 and run a `curl` command like
below; It should return an empty JSON and an HTTP-200 status code:

```shell
kubectl port-forward svc/emojis 8080:8080
curl -v 127.0.0.1:8080
```

## Controller example

This project contains a step by step example of how to build a controller. Each
branch contains a step. To work with it, go to `master` which has no operator
at all, only the application. The main steps are:  

- `feature/operator/01-initialization` contains the code scaffold
- `feature/operator/02-namespace-operator` changes the operator so that it is
  only executed on a single namespace
- `feature/operator/03-crd-and-controller` creates en empty emoji CRD and
  controller
- `feature/operator/04-specify-emoji-type` add expected types to the CRD
- `feature/operator/05-implement-support` implements a status.supported
  properties that tells if the Emoji is supported
- `feature/operator/06-implement-finalizer` implements a finalizer to manage
  resource deletion
- `feature/operator/07-implement-applogic` adds the logic to interact with
  the api

