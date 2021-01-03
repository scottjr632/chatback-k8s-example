# Chatback - A K8s Application Example

<!-- toc -->
* [About This Project](#about-this-project)
* [Requirements](#requirements)
* [Getting Started](#getting-started)
    * [Components](#components)
    * [Scripts](#scripts)
    * [Monitoring](#monitoring)
* [Usage](#usage)
<!-- tocstop -->

## About this project

This project is an application that was built to help learn Kubernetes. There are multiple parts to this project and instructions on deploying each one. Each part of this project has already been build out and can be deployed on either a local kubernetes instance (e.g. docker-desktop or mini-kube) or deployed to a cloud provider like AKS (Azure Kubernetes Service).

## Requirements

- [Docker](https://www.docker.com/get-started)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- Kubernetes (either)
    - [K8s for Docker Desktop](https://birthday.play-with-docker.com/kubernetes-docker-desktop/)
    - [Mini-Kube](https://minikube.sigs.k8s.io/docs/start/)
- [Docker Compose (optional)](https://docs.docker.com/compose/install/)

## Getting Started

This is an example chat application that can be deployed to any kubernetes cluster. There are several components to this application that are discussed below. 

### Components

#### Client

The client is the frontend of the application. This is what is visible from the browser. The client is a ReactJS application. More can be read about the client in the [client's README](client/README.md).

#### Server

The server is the main backend for the application. It is written in Golang and uses [Fiber](https://gofiber.io/) to serve requests and handle the websockets. More can be read about the server in the [server's README](server/README.md).

#### Broker

The broker's job is to broker messages from one server to another. Because there can be multiple servers, the broker's job is to ensure that each server get's the message sent from a client. [Broker's README](broker/README.md).

#### Database

PostgreSQL is used for the database and stores the messages sent.

### Scripts

There are several scripts that can help ease the deployment.

### Monitoring

#### Prometheus

#### Grafana

## Usage
