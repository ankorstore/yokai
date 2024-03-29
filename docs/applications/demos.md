---
icon: material/application-cog-outline
---

# :material-application-cog-outline: Applications demos

> Yokai provides a [showroom](https://github.com/ankorstore/yokai-showroom) for its demo applications.

## :material-application-cog-outline: HTTP application demo

This [demo application](https://github.com/ankorstore/yokai-showroom/tree/main/http-demo) is a simple HTTP REST API (CRUD) to manage [gophers](https://go.dev/blog/gopher).

It provides:

- a Yokai application container, with
  	- the [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module to expose the REST API
  	- the [fxorm](https://github.com/ankorstore/yokai/tree/main/fxorm) module to enable database interactions
- a [MySQL](https://www.mysql.com/) container to store the gophers

Available on [:fontawesome-brands-github: GitHub](https://github.com/ankorstore/yokai-showroom/tree/main/http-demo).

## :material-application-cog-outline: Worker application demo

This [demo application](https://github.com/ankorstore/yokai-showroom/tree/main/worker-demo) provides is a simple worker example subscribing to [Pub/Sub](https://cloud.google.com/pubsub).

It provides:

- a Yokai application container, with:
	- the [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module to expose a Pub/Sub publication endpoint
	- the [fxworker](https://github.com/ankorstore/yokai/tree/main/fxworker) module to provide a worker running a Pub/Sub subscriber
- a Pub/Sub emulator container

Available on [:fontawesome-brands-github: GitHub](https://github.com/ankorstore/yokai-showroom/tree/main/worker-demo).