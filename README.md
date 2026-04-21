# User Info Agent

A custom agent built on the [eclipse-cfm/cfm](https://github.com/eclipse-cfm/cfm) framework. It connects to a NATS
message broker, listens for orchestration tasks of type `user-info-activity`, fetches a random user from a dummy
REST endpoint, and publishes the selected user as output on the activity context.

_Note that his is just for demonstration purposes and does have any practical use or value!_

## TL;DR – how to deploy

Perform the following steps in sequence, assuming you have [JAD](https://github.com/metaform/jad) up and running:

```shell
docker build -t userinfoagent:latest .
kind load docker-image userinfoagent:latest --name jad
kubectl apply -f agent.userinfo.yaml
```
Then execute all `POST` requests in the [Bruno collection](./bruno) in sequence.


## How it works

The agent is structured around three layers:

- **`cmd/server/main.go`** — entry point. Starts the agent and blocks until a shutdown signal (SIGTERM/SIGINT) is
  received.
- **`launcher/launcher.go`** — wires up the CFM NATS agent with its configuration, resolves dependencies (HTTP client),
  and hands control to the framework.
- **`activity/activity.go`** — implements the business logic. On `ProcessDeploy`, it calls the configured remote URL,
  parses the JSON response as a list of `User` objects, picks one at random, and stores it as `selectedUser` in the
  activity output. `ProcessDispose` is a no-op.
- **`model/types.go`** — defines the `User`, `Address`, and `Company` structs that map the remote API response.

The agent registers itself on NATS under the service name `cfm.agent.hello-world` and handles activities of type
`user-info-activity`. The default remote URL is `https://jsonplaceholder.typicode.com/users`.

## Configuration

The agent reads its configuration from a file named `agent.userinfo.env` (YAML) mounted at `/etc/appname/` in
Kubernetes. The following keys are supported:

| Key                | Description                                       | Example                                      |
|--------------------|---------------------------------------------------|----------------------------------------------|
| `uri`              | NATS server URI                                   | `nats://localhost:4222`                      |
| `bucket`           | NATS JetStream bucket name                        | `cfm-bucket`                                 |
| `stream`           | NATS JetStream stream name                        | `cfm-stream`                                 |
| `agent.remote.url` | URL of the REST endpoint that returns a user list | `https://jsonplaceholder.typicode.com/users` |

## Building the Docker image

The `Dockerfile` uses a two-stage build: a Go builder stage compiles the binary, and a distroless runtime image runs it.

```bash
docker build -t userinfoagent:latest .
```

To verify that the image was created:

```bash
docker images userinfoagent
```

## Deploying to Kubernetes

The file `agent.userinfo.yaml` contains a `Deployment` and a `ConfigMap` for the `edc-v` namespace. This namespace was
already created by [JAD](https://github.com/metaform/jad).

**Prerequisites:**

- The `edc-v` namespace exists in the cluster.
- A `telemetry-config` ConfigMap is already present in the namespace (created by JAD).
- A NATS instance is reachable at `nats.edc-v.svc.cluster.local:4222`.
- The image `userinfoagent:latest` is available to the cluster (e.g. loaded into a local registry or `minikube`).

**Loading the image into a local cluster (KinD example):**

```bash
kind load docker-image userinfoagent:latest --name jad
```

**Apply the manifest:**

```bash
kubectl apply -f agent.userinfo.yaml
```

**Verify the deployment:**

```bash
kubectl -n edc-v get pods -l app=user-info-agent
kubectl -n edc-v logs -l app=user-info-agent
```

## Registering the agent with the orchestrator

The `bruno/` directory contains a [Bruno](https://www.usebruno.com/) API collection to interact with the CFM process
manager (`pmBaseUrl`).

1. **Register the activity definition** — tells the process manager that `user-info-activity` exists:

   ```
   POST {{pmBaseUrl}}/api/v1alpha1/activity-definitions
   ```

   Use the request in `bruno/Create custom activity definition.yml`.

2. **Update the orchestration definition** — defines a deployment pipeline that includes `user-info-agent` as the final
   step after onboarding:

   ```
   POST {{pmBaseUrl}}/api/v1alpha1/orchestration-definitions
   ```

   Use the request in `bruno/Create Orchestration Definition.yml`.
   Please note that since the API endpoint overwrites the existing orchestration definition, the new orchestration
   definition must also include all other activities!

3. **Inspect registered definitions** via the GET requests in `Get All Activity Definitions.yml` and
   `Get All Orchestration Definitions.yml`.
