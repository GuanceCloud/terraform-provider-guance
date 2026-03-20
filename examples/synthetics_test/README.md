# Synthetics Test Example

This example demonstrates how to use the `guance_synthetics_test` resource to create and manage synthetics tests in Guance Cloud. Synthetics tests are proactive monitoring tasks that regularly test targets (such as websites, APIs, TCP services, etc.) from multiple regions around the world to detect their availability and performance.

## Requirements

- Terraform 1.0+
- Guance Cloud API key

## Usage

1. Configure your Guance Cloud API key in the `provider.tf` file or set the `GUANCE_ACCESS_TOKEN` environment variable.
2. Configure your region in the `provider.tf` file or set the `GUANCE_REGION` environment variable.
3. Modify the `main.tf` file to customize your synthetics test configurations.
4. Run `terraform init` to initialize the provider.
5. Run `terraform plan` to preview the changes.
6. Run `terraform apply` to create the synthetics tests.
7. Run `terraform destroy` to delete the synthetics tests.

## Example Configurations

### HTTP Test

```hcl
resource "guance_synthetics_test" "http_test" {
  type    = "http"
  regions = ["hangzhou", "shanghai"]
  tags    = ["http", "production"]

  task {
    name      = "Production Website Test"
    url       = "https://www.example.com"
    method    = "GET"
    frequency = "5m"
    status    = "ok"

    advance_options {
      request_options {
        follow_redirect = true
        headers = {
          "User-Agent" = "Mozilla/5.0"
        }
      }
    }

    success_when_logic = "and"
    success_when {
      status_code {
        is = "200"
      }
      response_time = "1s"
    }
  }
}
```

### TCP Test

```hcl
resource "guance_synthetics_test" "tcp_test" {
  type    = "tcp"
  regions = ["hangzhou", "shanghai"]
  tags    = ["tcp", "production"]

  task {
    name            = "Production TCP Service Test"
    host            = "example.com"
    port            = "8080"
    frequency       = "10m"
    status          = "ok"
    enable_traceroute = true
    timeout         = "5s"

    success_when {
      response_time = "1s"
    }
  }
}
```

### WebSocket Test

```hcl
resource "guance_synthetics_test" "websocket_test" {
  type    = "websocket"
  regions = ["hangzhou", "shanghai"]
  tags    = ["websocket", "production"]

  task {
    name      = "Production WebSocket Test"
    url       = "wss://example.com/ws"
    frequency = "15m"
    status    = "ok"
    message   = "ping"

    success_when {
      response_time = "500ms"
    }
  }
}
```

### ICMP Test

```hcl
resource "guance_synthetics_test" "icmp_test" {
  type    = "icmp"
  regions = ["hangzhou", "shanghai"]
  tags    = ["icmp", "production"]

  task {
    name            = "Production ICMP Test"
    host            = "example.com"
    frequency       = "20m"
    status          = "ok"
    enable_traceroute = true
    packet_count    = 5
    timeout         = "3s"

    success_when {
      response_time = "500ms"
    }
  }
}
```

## Supported Test Types

- `http`: Test HTTP/HTTPS endpoints
- `tcp`: Test TCP services
- `dns`: Test DNS resolution
- `browser`: Test web pages with browser automation
- `icmp`: Test network connectivity with ICMP pings
- `websocket`: Test WebSocket connections

## Notes

- The `regions` field specifies the regions where the test will be executed. You can find the list of available regions in the Guance Cloud documentation.
- The `frequency` field specifies how often the test will be executed. Valid values are: `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`.
- The `success_when_logic` field specifies the logical relationship between success conditions. Valid values are: `and`, `or`.
- The `advance_options` field allows you to configure advanced settings such as request headers, authentication, and proxy settings.
- The `success_when` field allows you to define custom success conditions based on status codes, response content, and response time.

