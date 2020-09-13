# End-to-end tests

## Introduction

A basic set of simple end-to-end tests that check basic functionality. Initially, the tests are intended to be used 
only locally, i.e. not as part of the CI pipeline. At first, the tests will not be fully automated, i.e. require 
some manual setup, e.g. setting up a local cluster.

The basic pattern of a test is to install the mesh in the cluster, wait until all pods are ready, install sample 
applications, wait until all pods are ready, possibly test the sample application, delete the sample application, 
repeat for all sample applications, and finally delete the mesh.

## Prerequisites

* A local cluster, e.g. microk8s. The tests were tested using microk8s channel=1.16/stable.
* kubectl
* BATS: https://github.com/bats-core/bats-core
* gRPCurl: https://github.com/fullstorydev/grpcurl
* The environment variable `MESHERY_ADAPTER_ADDR` is defined specifies the IP address or DNS name of the meshery adapter, e.g. `localhost`

## Running the tests

* use `make` with one of the `e2e-` targets
* the target `e2e-tests` does not define the environment variable `MESHERY_ADAPTER_ADDR`, define it before running this target 