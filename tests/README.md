# End-to-end tests

## Introduction

The end-to-end tests are a set of tests that check basic functionality. Initially, the tests are intended to be used 
only locally, i.e. not as part of a CI pipeline. At first, the tests will not be fully automated, i.e. require 
some manual setup, which includes setting up a (local) cluster.

The basic pattern of these tests is to install the mesh in the cluster, wait until all rollouts are ready, install
 sample applications, wait until all rollouts are ready, possibly test the sample application, delete the sample application, 
repeat for all sample applications, and finally delete the mesh.

## Prerequisites

* A local cluster, e.g. minikube. The tests were tested using `minikube start --cpus 4 --memory 8192 --kubernetes-version=v1.14.1`
* `kubectl`
* `jq`, a tool for processing JSON inputs.
* `base64`, used to encode kubeconfig. 
* BATS, the Bash Automated Testing System: https://github.com/bats-core/bats-core
* gRPCurl, a command-line tool to interact with gRPC servers: https://github.com/fullstorydev/grpcurl
* The environment variable `MESHERY_ADAPTER_ADDR` is defined and specifies the IP address or DNS name of the meshery
 adapter, e.g. `localhost`. For some configurations, this might be set up in the Makefile.
* The adapter is running.

## Implementation notes

As BATS is executing the test files alphabetically, the first two letters in the filenames are used to enforce the correct sequence.   

## Running the tests

* Use `make` with one of the `e2e-` targets.
* The target `e2e-tests` does not define the environment variable `MESHERY_ADAPTER_ADDR`. Define it before using
 this target.