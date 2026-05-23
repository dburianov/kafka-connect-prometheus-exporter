# Kafka connect exporter

A [Prometheus](https://prometheus.io/) exporter that collects [Kafka connect](https://docs.confluent.io/current/connect/index.html) metrics.

### Usage

```sh
$ ./kafka_connect_exporter -h
Usage of ./kafka_connect_exporter:
  -listen-address string
        Address on which to expose metrics. (default ":8080")
  -pass string
        Kafka connect basic auth password. Used along with -user, ignored when empty string. (default "")
  -scrape-uri string
        URI on which to scrape kafka connect. (default "http://127.0.0.1:8080")
  -telemetry-path string
        Path under which to expose metrics. (default "/metrics")
  -tls.ca-file string
        The optional certificate authority file for Kafka TLS client authentication
  -tls.cert-file string
        The optional certificate file for Kafka client authentication
  -tls.enabled
        Connect to Kafka using TLS
  -tls.insecure-skip-tls-verify
        If true, the server's certificate will not be checked for validity
  -tls.key-file string
        The optional key file for Kafka client authentication
  -tls.server-name string
        Used to verify the hostname on the returned certificates unless tls.insecure-skip-tls-verify is given. The kafka server's name should be given
  -user string
        Kafka connect basic auth username. Used along with -pass, ignored when empty string. (default "")
  -version
        show version and exit
```

## TLS Configuration

The exporter supports TLS configuration for secure connections to Kafka Connect:

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `-tls.enabled` | `TLS_ENABLED` | Enable TLS for connections (boolean) |
| `-tls.server-name` | `TLS_SERVER_NAME` | Server name for certificate verification |
| `-tls.ca-file` | `TLS_CA_FILE` | Path to CA certificate file for client authentication |
| `-tls.cert-file` | `TLS_CERT_FILE` | Path to client certificate file |
| `-tls.key-file` | `TLS_KEY_FILE` | Path to client key file |
| `-tls.insecure-skip-tls-verify` | `TLS_INSECURE_SKIP_VERIFY` | Skip certificate verification (boolean) |

### TLS Examples

**HTTPS with self-signed certificate:**
```sh
./kafka_connect_exporter \
  -tls.enabled \
  -tls.ca-file=/path/to/ca.crt \
  -tls.insecure-skip-tls-verify \
  -scrape-uri=https://kafka-connect:8083
```

**Mutual TLS (mTLS) with client certificate:**
```sh
./kafka_connect_exporter \
  -tls.enabled \
  -tls.ca-file=/path/to/ca.crt \
  -tls.cert-file=/path/to/client.crt \
  -tls.key-file=/path/to/client.key \
  -scrape-uri=https://kafka-connect:8083
```

**Skip certificate verification (not recommended for production):**
```sh
./kafka_connect_exporter \
  -tls.enabled \
  -tls.insecure-skip-tls-verify \
  -scrape-uri=https://kafka-connect:8083
```

## Metrics

```
# HELP kafka_connect_connector_state_running is the connector running?
# TYPE kafka_connect_connector_state_running gauge
kafka_connect_connector_state_running{connector="test-changesets",state="running",worker="kafka-connect:8083"} 1
# HELP kafka_connect_connector_tasks_state the state of tasks. 0-failed, 1-running, 2-unassigned, 3-paused
# TYPE kafka_connect_connector_tasks_state gauge
kafka_connect_connector_tasks_state{connector="test-changesets",state="running",worker_id="kafka-connect:8083"} 1
# HELP kafka_connect_connectors_count number of deployed connectors
# TYPE kafka_connect_connectors_count gauge
kafka_connect_connectors_count 1
# HELP kafka_connect_up was the last scrape of kafka connect successful?
# TYPE kafka_connect_up gauge
kafka_connect_up 1
```

## Helm Chart

The TLS configuration can also be set via Helm values:

```yaml
configuration:
  env:
    TLS_ENABLED: "true"
    TLS_SERVER_NAME: "kafka.example.com"
    TLS_CA_FILE: "/path/to/ca.crt"
    TLS_CERT_FILE: "/path/to/client.crt"
    TLS_KEY_FILE: "/path/to/client.key"
    TLS_INSECURE_SKIP_VERIFY: "false"