# Variables
variable "regions" {
  description = "Regions for the synthetic tests"
  type        = list(string)
  default     = []
}

# HTTP
resource "guance_synthetics_test" "http_test" {
  type    = "http"
  regions = var.regions
  tags    = ["http", "production"]

  task = {
    name            = "TF HTTP Test"
    url             = "https://www.example.com"
    method          = "GET"
    frequency       = "5m"
    status          = "stop"
    schedule_type   = "frequency"

    advance_options = {
      request_options = {
        headers = {
          "User-Agent" = "Mozilla/5.0"
        }
        auth = {
          username = "test"
          password = "test"
        }
      }
      request_body = {
        body_type = "application/json"
        body      = "{\"key\": \"value\"}"
      }
      certificate = {
        ignore_server_certificate_error = true
      }
      request_timeout = "30s"
    }

    success_when_logic = "and"
    success_when = [
      {
        body = [
          {
            contains = "Example Domain"
          }
        ]
        response_time = [
          {
            target = "1s"
          }
        ]
        header = {
          "Content-Type" = [
            jsonencode({
              is = "application/json"
            }),
            
          ]
        }
      }
    ]
  }
}

# TCP
resource "guance_synthetics_test" "tcp_test" {
  type    = "tcp"
  regions = var.regions
  tags    = ["tcp", "production"]

  task = {
    name              = "TF TCP Service Test"
    host              = "example.com"
    port              = "8080"
    message           = "hello"
    frequency         = "5m"
    timeout           = "5s"
    status            = "ok"
    enable_traceroute = true
    desc              = "TCP test for example.com:8080"

    success_when = [
      {
        response_time = [{
          is_contain_dns = true,
          target = "10ms"
        }]
        response_message = [
          {
            is = "hello"
          }
        ]
        hops = [
          {
            op     = "eq"
            target = 20
          }
        ]
      }
    ]
  }
}

# WebSocket
resource "guance_synthetics_test" "websocket_test" {
  type    = "websocket"
  regions = var.regions
  tags    = ["websocket", "production"]

  task = {
    name      = "TF WebSocket Test"
    url       = "wss://example.com/ws"
    frequency = "15m"
    status    = "ok"
    message   = "ping"

    success_when = [
      {
        response_time = [{
          is_contain_dns = true,
          target = "10ms"
        }]
        response_message = [
          {
            is = "pong"
          }
        ]
        header = {
          status = [
            jsonencode({
              is = "ok"
            })
          ]
          token = [
            jsonencode({
              contains = "aaa"
            }),
            jsonencode({
              contains = "bbb"
            })
          ]
        }
      }
    ]

    advance_options = {
      request_options = {
        timeout = "30s"
        headers = {
          "x-token" = "aaaa"
        }
      }
      auth = {
        username = "test"
        password = "test"
      }
    }
  }
}

# ICMP
resource "guance_synthetics_test" "icmp_test" {
  type    = "icmp"
  regions = var.regions
  tags    = ["icmp", "production"]

  task = {
    name              = "TF ICMP Test"
    host              = "example.com"
    frequency         = "30m"
    status            = "ok"
    enable_traceroute = false
    packet_count      = 5
    timeout           = "3s"

    success_when = [
      {
        response_time = [{
          func = "avg",
          op = "leq",
          target = "100ms"
        }]
        packet_loss_percent = [
          {
            op     = "leq"
            target = 10
          }
        ]
        hops = [
          {
            op     = "eq"
            target = 20
          }
        ]
        packets = [
          {
            op     = "geq"
            target = 4
          }
        ]
      }
    ]
  }
}

# Multi-step
resource "guance_synthetics_test" "multi_step_test" {
  type    = "multi"
  regions = var.regions
  tags    = ["multi", "production"]

  task = {
    name      = "TF Multi-step Test"
    frequency = "15m"
    status    = "ok"
    desc      = "Multi-step test with HTTP and wait steps"

    steps = [
      {
        type          = "http"
        task          = jsonencode({
          "name": "step1", 
          "method": "GET", 
          "url": "https://api.example.com/resource", 
          "post_mode": "script", 
          "post_script": "vars[\"token\"] = \"token_value\""
        })
        allow_failure = false
        extracted_vars = [
          {
            name   = "TOKEN"
            field  = "token"
            secure = false
          }
        ]
      },
      {
        type   = "wait"
        value  = 3
      },
      {
        type          = "http"
        task          = jsonencode({
          "name": "step3", 
          "method": "POST", 
          "url": "https://api.example.com/resource", 
          "post_mode": "script", 
          "post_script": "result[\"is_failed\"]=false"
        })
        allow_failure = false
      }
    ]
  }
}