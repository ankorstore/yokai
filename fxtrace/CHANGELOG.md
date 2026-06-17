# Changelog

## [1.3.0](https://github.com/ankorstore/yokai/compare/fxtrace/v1.2.0...fxtrace/v1.3.0) (2026-05-22)


### Bug Fixes

* **fxtrace:** Treat tracer provider `ForceFlush`/`Shutdown` as best-effort during `OnStop`. Errors (e.g. OTLP collector overloaded, `sending_queue is full`) are now logged and swallowed instead of propagated to `fx.App.Stop()`, so a saturated collector no longer turns a graceful pod shutdown into a non-zero exit. Both calls are bounded by a short internal timeout so a hanging exporter cannot consume the entire pod termination grace period.


### Features

* **fxtrace:** `FxTraceParam` now requires a `*log.Logger` so the module can log suppressed best-effort shutdown errors. Consumers wiring `FxTraceModule` outside `fxcore` must also register `fxlog.FxLogModule`.

## [1.2.0](https://github.com/ankorstore/yokai/compare/fxtrace/v1.1.0...fxtrace/v1.2.0) (2024-03-14)


### Features

* **fxtrace:** Updated dependencies ([#149](https://github.com/ankorstore/yokai/issues/149)) ([cbafdb7](https://github.com/ankorstore/yokai/commit/cbafdb7d5ddef34ce63f680eafe28d61ed4c3df3))

## [1.1.0](https://github.com/ankorstore/yokai/compare/fxtrace/v1.0.0...fxtrace/v1.1.0) (2024-01-11)


### Features

* **fxtrace:** Updated module name ([#30](https://github.com/ankorstore/yokai/issues/30)) ([e440bdd](https://github.com/ankorstore/yokai/commit/e440bdd815bf7642b1694e9d96b6d0d061d85efe))

## 1.0.0 (2024-01-11)


### Features

* **fxtrace:** Provided module ([#28](https://github.com/ankorstore/yokai/issues/28)) ([6757f8e](https://github.com/ankorstore/yokai/commit/6757f8e909d6399580a7cf3c4764532bedf8afd4))
