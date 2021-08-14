# Blackbox Traefik SD

Generate targets for Blackbox Exporter using Traefik.

## Config

| Environment Variable | Description                                        | Default                             |
|----------------------|----------------------------------------------------|-------------------------------------|
| `TRAEFIK_URL`        | Traefik API url (with scheme, http:// or https://) | ``                                  |
| `TRAEFIK_USERNAME`   | Traefik API basic auth username (if required)      | ``                                  |
| `TRAEFIK_PASSWORD`   | Traefik API basic auth password (if required)      | ``                                  |
| `LOG_LEVEL`          | The level of log verbosity                         | `Info`                              |
| `TARGETS_FILE`       | The file to output                                 | `/blackbox-traefik-sd/targets.json` |
| `INTERVAL`           | How often to build the targets file                | `600`                               |

## Usage

Enable the Traefik API following [the official doc](https://doc.traefik.io/traefik/operations/api/))

Start this application, passing in the required environment variables. To start the application in Docker:

```bash
$ docker run -d \
  -e TRAEFIK_URL=https://traefik.example.com \
  -e INTERVAL=600 \
  -v /apps/blackbox-traefik-sd:/config \
  -e TARGETS_FILE=/config/targets.json \
  ghcr.io/calvinbui/homer-service-discovery
```

The application will generate a JSON file to the path specified in the environment variable `TARGETS_FILE`.

Update your Prometheus config to use the generated targets file along with Blackbox Exporter:

```yaml
scrape_configs:
  - job_name: blackbox
    metrics_path: /probe
    params:
      module: [http]
    file_sd_configs:
      - files:
        - /blackbox-traefik-sd/targets.json
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: blackbox_exporter:9115
```

## Caveats

If a route is removed from Traefik (i.e. a Docker container is removed), this application will also remove it from being monitored. Therefore it is recommended to set the `INTERVAL` environment variable to twice the amount of time it'll take before your rules/alerts trigger.

## To Do

- Support config Labels for Prometheus
- Integrate with Docker labels to configure labels and scheme (currently on https://)
- Integrate with Traefik Services to get more information

## Thanks

- https://github.com/containeroo/SyncFlaer: for their Traefik code
