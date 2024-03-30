---
icon: material/folder-eye-outline
---

# :material-folder-eye-outline: Demo - worker application

> Yokai's [showroom](https://github.com/ankorstore/yokai-showroom) provides a [worker demo application](https://github.com/ankorstore/yokai-showroom/tree/main/worker-demo).

## Overview

This [worker demo application](https://github.com/ankorstore/yokai-showroom/tree/main/worker-demo) is a simple subscriber to [Pub/Sub](https://cloud.google.com/pubsub).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [worker](../modules/fxworker.md) module to offer a subscriber worker using the [fxgcppubsub](https://github.com/ankorstore/yokai-contrib/tree/main/fxgcppubsub) contrib module
- a [Pub/Sub emulator](https://github.com/marcelcorso/gcloud-pubsub-emulator) container, with preconfigured topic and subscription
- a [Pub/Sub emulator UI](https://github.com/echocode-io/gcp-pubsub-emulator-ui) container, preconfigured to work with the emulator container
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

### Layout

This demo application is following the [recommended project layout](https://go.dev/doc/modules/layout):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `worker/`: workers
	- `bootstrap.go`: bootstrap
	- `register.go`: dependencies registration

### Makefile

This demo application provides a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make test   # run tests
make lint   # run linter
```

## Usage

### Start the application

To start the application, simply run:

```shell
make fresh
```

After a short moment, the application will offer:

- [http://localhost:8081](http://localhost:8081): application core dashboard
- [http://localhost:8680](http://localhost:8680): pub/sub emulator UI
- [http://localhost:16686](http://localhost:16686): jaeger UI

### Message publication

You can use the Pub/Sub emulator UI to publish a message to the preconfigured topic:

[http://localhost:8680/project/demo-project/topic/demo-topic](http://localhost:8680/project/demo-project/topic/demo-topic)

### Message subscription

Check your application logs by running:

```shell
make logs
```

You will see the [SubscribeWorker](https://github.com/ankorstore/yokai-showroom/blob/main/worker-demo/internal/worker/subscribe.go) subscribed to Pub/Sub in action, logging the received
messages.
