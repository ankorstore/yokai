---
icon: material/application-cog-outline
---

# :material-application-cog-outline: Demo applications

> Yokai provides a [showroom](https://github.com/ankorstore/yokai-showroom) for its demo applications.

## gRPC demo application

This [gRPC demo application](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo) is a simple gRPC API offering a [text transformation service](https://github.com/ankorstore/yokai-showroom/tree/main/grpc-demo/proto/transform.proto).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [gRPC server](../modules/fxgrpcserver.md) module to offer the gRPC API
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

To try it, just follow the [README](https://github.com/ankorstore/yokai-showroom/blob/main/grpc-demo/README.md) instructions.

## HTTP demo application 

This [HTTP demo application](https://github.com/ankorstore/yokai-showroom/tree/main/http-demo) is a simple HTTP REST API (CRUD) to manage [gophers](https://go.dev/blog/gopher).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [HTTP server](../modules/fxhttpserver.md) module to offer the gophers API
- a [MySQL](https://www.mysql.com/) container to store the gophers
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

To try it, just follow the [README](https://github.com/ankorstore/yokai-showroom/blob/main/http-demo/README.md) instructions.

## Worker demo application

This [worker demo application](https://github.com/ankorstore/yokai-showroom/tree/main/worker-demo) provides is a simple worker example subscribing to [Pub/Sub](https://cloud.google.com/pubsub).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [worker](../modules/fxworker.md) module to offer a worker subscribing to Pub/Sub (using the [fxgcppubsub](https://github.com/ankorstore/yokai-contrib/tree/main/fxgcppubsub) contrib module)
- a [Pub/Sub emulator](https://cloud.google.com/pubsub) container, with preconfigured topic and subscription
- a [Pub/Sub emulator UI](https://github.com/echocode-io/gcp-pubsub-emulator-ui) container, preconfigured to work with the emulator container
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

To try it, just follow the [README](https://github.com/ankorstore/yokai-showroom/blob/main/worker-demo/README.md) instructions.
