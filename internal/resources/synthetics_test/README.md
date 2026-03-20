# Guance Synthetics Test Resource

## 核心功能

Guance Synthetics Test 资源用于在 Guance Cloud 中创建和管理云拨测任务。云拨测是一种主动监控服务，可以定期从全球多个地域向指定的目标（如网站、API、TCP 服务等）发起测试，以检测其可用性和性能。

## 适用场景

- **网站可用性监控**：定期检测网站是否可访问，响应时间是否正常
- **API 性能监控**：监控 API 的响应时间、状态码和响应内容
- **TCP 服务监控**：监控 TCP 服务的连通性和响应时间
- **WebSocket 服务监控**：监控 WebSocket 服务的连接和消息传递
- **ICMP 网络监控**：监控网络连通性和延迟
- **DNS 解析监控**：监控 DNS 解析的正确性和响应时间

## 资源配置参数

### 基本参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
| :--- | :--- | :--- | :--- | :--- |
| `type` | `string` | 是 | - | 拨测类型，可选值：`http`, `tcp`, `dns`, `browser`, `icmp`, `websocket` |
| `regions` | `list(string)` | 是 | - | 拨测执行的地域列表 |
| `tags` | `list(string)` | 否 | - | 资源标签列表 |

### 任务配置参数 (`task`)

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
| :--- | :--- | :--- | :--- | :--- |
| `name` | `string` | 是 | - | 任务名称 |
| `frequency` | `string` | 是 | - | 拨测频率，可选值：`1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h` |
| `status` | `string` | 否 | `ok` | 任务状态，可选值：`ok`, `stop` |
| `url` | `string` | 否 | - | 测试目标 URL（HTTP、Browser、WebSocket 类型必填） |
| `method` | `string` | 否 | - | HTTP 请求方法（HTTP 类型必填） |
| `host` | `string` | 否 | - | 测试目标主机（TCP、ICMP 类型必填） |
| `port` | `string` | 否 | - | 测试目标端口（TCP 类型必填） |
| `timeout` | `string` | 否 | - | 测试超时时间（TCP、ICMP 类型可选） |
| `message` | `string` | 否 | - | WebSocket 测试消息（WebSocket 类型必填） |
| `enable_traceroute` | `bool` | 否 | `false` | 是否启用路由跟踪（TCP、ICMP 类型可选） |
| `packet_count` | `number` | 否 | - | ICMP 测试发送的数据包数量（ICMP 类型可选） |
| `post_mode` | `string` | 否 | `default` | 可用性判断模式，可选值：`default`, `script` |
| `post_script` | `string` | 否 | - | 可用性判断脚本内容（post_mode 为 script 时使用） |
| `success_when_logic` | `string` | 否 | `and` | 成功条件逻辑关系，可选值：`and`, `or` |

### 高级选项参数 (`task.advance_options`)

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
| :--- | :--- | :--- | :--- | :--- |
| `request_options.follow_redirect` | `bool` | 否 | `false` | 是否跟随重定向 |
| `request_options.headers` | `map(string)` | 否 | - | HTTP 请求头 |
| `request_options.cookies` | `string` | 否 | - | HTTP 请求 cookies |
| `request_options.auth.username` | `string` | 否 | - | HTTP 认证用户名 |
| `request_options.auth.password` | `string` | 否 | - | HTTP 认证密码 |
| `request_body.body_type` | `string` | 否 | - | HTTP 请求体类型 |
| `request_body.body` | `string` | 否 | - | HTTP 请求体内容 |
| `certificate.ignore_server_certificate_error` | `bool` | 否 | `false` | 是否忽略服务器证书错误 |
| `certificate.private_key` | `string` | 否 | - | 客户端证书私钥 |
| `certificate.certificate` | `string` | 否 | - | 客户端证书 |
| `proxy.url` | `string` | 否 | - | 代理服务器 URL |
| `proxy.headers` | `map(string)` | 否 | - | 代理服务器请求头 |
| `secret.not_save` | `bool` | 否 | `false` | 是否不保存敏感信息 |

### 成功条件参数 (`task.success_when`)

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
| :--- | :--- | :--- | :--- | :--- |
| `body[*].contains` | `string` | 否 | - | 响应体包含指定字符串 |
| `body[*].not_contains` | `string` | 否 | - | 响应体不包含指定字符串 |
| `body[*].is` | `string` | 否 | - | 响应体等于指定字符串 |
| `body[*].is_not` | `string` | 否 | - | 响应体不等于指定字符串 |
| `body[*].match_regex` | `string` | 否 | - | 响应体匹配指定正则表达式 |
| `body[*].not_match_regex` | `string` | 否 | - | 响应体不匹配指定正则表达式 |
| `status_code[*].is` | `string` | 否 | - | 状态码等于指定值 |
| `status_code[*].is_not` | `string` | 否 | - | 状态码不等于指定值 |
| `status_code[*].match_regex` | `string` | 否 | - | 状态码匹配指定正则表达式 |
| `status_code[*].not_match_regex` | `string` | 否 | - | 状态码不匹配指定正则表达式 |
| `response_time` | `string` | 否 | - | 最大响应时间，如 "100ms" |

## 资源创建/更新/删除的 Terraform 命令及 API 调用示例

### 创建资源

```hcl
resource "guance_synthetics_test" "example" {
  type    = "http"
  regions = ["hangzhou", "shanghai"]
  tags    = ["test", "http"]

  task {
    name      = "HTTP Test"
    url       = "https://www.example.com"
    method    = "GET"
    frequency = "1m"
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
      response_time = "500ms"
    }
  }
}
```

API 调用示例：

```shell
curl 'https://openapi.guance.com/api/v1/dialing_task/add' \
-H 'DF-API-KEY: <DF-API-KEY>' \
-H 'Content-Type: application/json;charset=UTF-8' \
--data-raw '{"type":"http","regions":["hangzhou","shanghai"],"task":{"name":"HTTP Test","url":"https://www.example.com","method":"GET","frequency":"1m","status":"ok","advance_options":{"request_options":{"follow_redirect":true,"headers":{"User-Agent":"Mozilla/5.0"}}},"success_when_logic":"and","success_when":[{"status_code":[{"is":"200"}],"response_time":"500ms"}]},"tags":["test","http"]}'
```

### 更新资源

修改 Terraform 配置文件中的资源参数，然后执行：

```bash
terraform apply
```

API 调用示例：

```shell
curl 'https://openapi.guance.com/api/v1/dialing_task/<task_uuid>/modify' \
-H 'DF-API-KEY: <DF-API-KEY>' \
-H 'Content-Type: application/json;charset=UTF-8' \
--data-raw '{"regions":["hangzhou","shanghai"],"task":{"name":"HTTP Test Updated","url":"https://www.example.com","method":"GET","frequency":"5m","status":"ok","advance_options":{"request_options":{"follow_redirect":true,"headers":{"User-Agent":"Mozilla/5.0"}}},"success_when_logic":"and","success_when":[{"status_code":[{"is":"200"}],"response_time":"500ms"}]},"tags":["test","http","updated"]}'
```

### 删除资源

执行以下命令删除资源：

```bash
terraform destroy -target=guance_synthetics_test.example
```

API 调用示例：

```shell
curl 'https://openapi.guance.com/api/v1/dialing_task/delete' \
-H 'DF-API-KEY: <DF-API-KEY>' \
-H 'Content-Type: application/json;charset=UTF-8' \
--data-raw '{"taskUUIDs":["<task_uuid>"]}'
```

## 常见问题及解决方案

### 资源创建失败

1. **权限不足**
   - 问题：API 密钥权限不足
   - 解决方案：确保使用的 API 密钥具有创建和管理云拨测任务的权限

2. **参数错误**
   - 问题：配置参数不符合要求
   - 解决方案：检查必填参数是否填写正确，参数值是否符合要求（如频率、类型等）

3. **地域不存在**
   - 问题：指定的地域不存在或不可用
   - 解决方案：使用有效的地域名称，可参考 Guance Cloud 文档中的地域列表

### 资源更新失败

1. **任务不存在**
   - 问题：尝试更新的任务已被删除
   - 解决方案：检查任务是否存在，如不存在需要重新创建

2. **参数冲突**
   - 问题：更新的参数与现有配置冲突
   - 解决方案：检查参数是否符合要求，特别是类型相关的必填参数

### 资源删除失败

1. **权限不足**
   - 问题：API 密钥权限不足
   - 解决方案：确保使用的 API 密钥具有删除云拨测任务的权限

2. **任务不存在**
   - 问题：尝试删除的任务已不存在
   - 解决方案：检查任务是否存在，如不存在可以忽略此错误

## 示例配置

### HTTP 测试示例

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
          "Authorization" = "Bearer <token>"
        }
      }
      request_body {
        body_type = "application/json"
        body      = "{\"key\": \"value\"}"
      }
    }

    success_when_logic = "and"
    success_when {
      status_code {
        is = "200"
      }
      body {
        contains = "Welcome"
      }
      response_time = "1s"
    }
  }
}
```

### TCP 测试示例

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

### WebSocket 测试示例

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
