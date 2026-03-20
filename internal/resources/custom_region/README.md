# Guance Custom Region Resource

## 核心功能

Guance Custom Region 资源用于在 Guance Cloud 中创建和管理自定义拨测节点。自定义拨测节点是用户自行部署的拨测执行点，可以在特定的网络环境中执行拨测任务。

## 适用场景

- **特定网络环境监控**：在企业内部网络或特定区域部署拨测节点，监控内部服务
- **海外服务监控**：在海外部署拨测节点，监控全球服务的可用性
- **专用网络监控**：在特定网络环境中部署拨测节点，确保服务在该环境中的可用性

## 资源配置参数

### 基本参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
| :--- | :--- | :--- | :--- | :--- |
| `internal` | `bool` | 是 | - | 国家属性（true 为国内，false 为海外） |
| `isp` | `string` | 是 | - | 运营商 |
| `country` | `string` | 是 | - | 国家 |
| `province` | `string` | 否 | - | 省份 |
| `city` | `string` | 否 | - | 城市 |
| `name` | `string` | 否 | - | 拨测节点昵称 |
| `company` | `string` | 否 | - | 企业信息 |
| `keycode` | `string` | 是 | - | 拨测节点 keycode（不可重名） |

## 资源创建/更新/删除的 Terraform 命令及 API 调用示例

### 创建资源

```hcl
resource "guance_custom_region" "example" {
  internal = false
  isp      = "telecom"
  country  = "Afghanistan"
  city     = "Shahrak"
  name     = "test"
  keycode  = "Afghanistan-Shahrak-telecom"
}
```

API 调用示例：

```shell
curl 'https://openapi.guance.com/api/v1/dialing_region/regist' \
-H 'DF-API-KEY: <DF-API-KEY>' \
-H 'Content-Type: application/json;charset=UTF-8' \
--data-raw '{"internal":false,"isp":"telecom","country":"Afghanistan","city":"Shahrak","keycode":"Afghanistan-Shahrak-telecom","name":"test"}' \
--compressed
```

### 更新资源

自定义拨测节点不支持更新操作。如果需要修改配置，需要先删除现有节点，然后创建新节点。

### 删除资源

执行以下命令删除资源：

```bash
terraform destroy -target=guance_custom_region.example
```

API 调用示例：

```shell
curl 'https://openapi.guance.com/api/v1/dialing_region/<region_uuid>/delete' \
-H 'DF-API-KEY: <DF-API-KEY>' \
--compressed
```

## 常见问题及解决方案

### 资源创建失败

1. **权限不足**
   - 问题：API 密钥权限不足
   - 解决方案：确保使用的 API 密钥具有创建和管理自定义拨测节点的权限

2. **参数错误**
   - 问题：配置参数不符合要求
   - 解决方案：检查必填参数是否填写正确，keycode 是否重复

3. **网络问题**
   - 问题：无法连接到 Guance Cloud API
   - 解决方案：检查网络连接，确保能够访问 Guance Cloud API 端点

### 资源删除失败

1. **权限不足**
   - 问题：API 密钥权限不足
   - 解决方案：确保使用的 API 密钥具有删除自定义拨测节点的权限

2. **节点不存在**
   - 问题：尝试删除的节点已不存在
   - 解决方案：检查节点是否存在，如不存在可以忽略此错误

## 示例配置

### 基本示例

```hcl
resource "guance_custom_region" "basic" {
  internal = true
  isp      = "telecom"
  country  = "China"
  province = "Zhejiang"
  city     = "Hangzhou"
  name     = "Hangzhou Telecom"
  keycode  = "china-hangzhou-telecom"
}
```

### 海外节点示例

```hcl
resource "guance_custom_region" "overseas" {
  internal = false
  isp      = "aws"
  country  = "United States"
  province = "California"
  city     = "San Francisco"
  name     = "SF AWS"
  company  = "Example Inc."
  keycode  = "us-sanfrancisco-aws"
}
```
