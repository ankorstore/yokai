# Showroom

> Yokai provides a showroom for demo application.

## HTTP application demo

This [demo application](https://github.com/ankorstore/yokai-showroom/tree/main/http-demo) is a simple REST API (CRUD) to manage [gophers](https://go.dev/blog/gopher).

It provides:

- a Yokai application container, with the [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module to offer the REST API
- a [MySQL](https://www.mysql.com/) container to store the gophers

Available on [:fontawesome-brands-github: GitHub](https://github.com/ankorstore/yokai-showroom/tree/main/http-demo).

## Worker application template

This [demo application](https://github.com/ankorstore/yokai-showroom/tree/main/worker-demo) provides is a simple worker example subscribing to [Pub/Sub](https://cloud.google.com/pubsub).

It provides:

- a Yokai application container, with:
	- the [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module to offer a Pub/Sub
	  messages publication endpoint
	- the [fxworker](https://github.com/ankorstore/yokai/tree/main/fxworker) module to offer a worker running a Pub/Sub
	  messages subscriber
- a Pub/Sub emulator container

Available on [:fontawesome-brands-github: GitHub](https://github.com/ankorstore/yokai-showroom/tree/main/worker-demo).