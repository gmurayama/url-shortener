# Webservice Template in Go

Template for web services in Go using [Gin](https://gin-gonic.com/). It provide a basic features with useful features pre-configured from the get-go.

## Overview

### Features

- Observability with [OpenTelemetry](https://opentelemetry.io/docs/languages/go/getting-started/) for tracing and [Prometheus](https://github.com/prometheus/client_golang?tab=readme-ov-file#instrumenting-applications) for metrics
- HTTP web server using [Gin](https://gin-gonic.com/)
- Configuration using environment variables
- Dockerfile with multi-arch build

### Structures

```
webservice-template-golang/
├─ cmd/                   # entrypoint
├─ config/
├─ internal/
│  ├─ application/        # bussiness layer
│  ├─ commons/            # common functions that can be used in all layers
│  ├─ gateways/           # the entrypoint implementation (e.g. the API handlers)
│  ├─ infrastructure/     # adapters (i.e. implementation for abstraction used in _gateways_ or _application_)
│  ├─ .../
├─ docker-compose.yaml    # prometheus, grafana and jaeger services (development purposes)
├─ prometheus.yaml        # basic config for prometheus
```

## Getting Started

```
$ docker-compose up -d
$ make help/api # print config
$ make run/api  # listening on :8080, internal server (pprof and metrics) on :9090
```

