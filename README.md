# http-prober

A super minimal server that runs HTTP requests against `https://google.com` and exposes results and duration as Prometheus metrics.

It differs from [blackbox-exporter](https://github.com/prometheus/blackbox_exporter) because it is a stateful application. It will increase the same counters and histograms, instead of reseting on every scrape.

## Development

Build and run the binary locally with:

```
make build
./http-prober
```

Or as a docker image:

```
make docker-build
docker run -p 8080:8080 ghcr.io/gitpod-io/http-prober:main
```

## Releasing

Create a new tag with

```
git tag -a <tag number> -m "Tag message"
```

Push the new tag with

```
git push origin <tag number>
```

Github Actions will [generate binaries](.github/workflows/goreleaser.yml) and [docker images](.github/workflows/container.yml)