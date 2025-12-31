# URL Shortener

> Case study for System Design. More details about the architecture on [SD.md](SD.md)

URL Shortener creates a short URL (alias) for any URL. When accessed, it redirects to the original long URL.

## Getting Started

```
$ docker-compose up -d
$ make help/api # print config
$ make run/api  # :7000, internal server on :7001
```

## Build

```
$ make docker-build-push/multi-arch DOCKER_TAG=registry.example.com/url-shortener:0.1.0
```
