# Blackbox Traefik SD

Generate targets for Blackbox Exporter using Traefik.

## Config

| Environment Variable | Description                                        | Default                 |
|----------------------|----------------------------------------------------|-------------------------|
| `TRAEFIK_URL`        | Traefik API url (with scheme, http:// or https://) | ``                      |
| `TRAEFIK_USERNAME`   | Traefik API basic auth username (if required)      | ``                      |
| `TRAEFIK_PASSWORD`   | Traefik API basic auth password (if required)      | ``                      |
| `LOG_LEVEL`          | The level of log verbosity                         | `Info`                  |
| `OUTPUT_DIR`         | The folder to output all target JSON files         | `/blackbox-traefik-sd/` |
| `INTERVAL`           | How often to build the targets file                | `600`                   |

## Usage

Enable the Traefik API following [the official doc](https://doc.traefik.io/traefik/operations/api/))

Start this application, passing in the required environment variables. To start the application in Docker:

```bash
$ docker run -d \
  -e TRAEFIK_URL=https://traefik.example.com \
  -e INTERVAL=600 \
  -v /apps/blackbox-traefik-sd:/config \
  -e OUTPUT_DIR=/config/targets.json \
  ghcr.io/calvinbui/homer-service-discovery
```

The application will generate JSON files to the path specified in the environment variable `OUTPUT_DIR`.

Update your Prometheus config to use the generated targets file along with Blackbox Exporter:

```yaml
scrape_configs:
  - job_name: blackbox
    metrics_path: /probe
    params:
      module: [http]
    file_sd_configs:
      - files:
        - /blackbox-traefik-sd/*.json
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: blackbox_exporter:9115
```

## Caveats

If a route is removed from Traefik (i.e. a Docker container is removed), you will have to delete the target JSON file manually from Prometheus under the `OUTPUT_DIR`.

## To Do

- Support config Labels for Prometheus
- Integrate with Docker labels to configure labels and scheme (currently on https://)
- Integrate with Traefik Services to get more information
- Delete any unknown targets after a period of time

## Thanks

- https://github.com/containeroo/SyncFlaer: for their Traefik code
