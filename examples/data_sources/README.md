# Guance Data Source Examples

This example shows how to look up existing Guance objects with Terraform data sources.

## Usage

Set `GUANCE_ACCESS_TOKEN` and any lookup variables you need, then run:

```shell
terraform init
terraform plan \
  -var='monitor_search=Terraform' \
  -var='monitor_type=simpleCheck'
```

Name-based data sources require the target object to already exist and the name to be unique.
Leave `monitor_type` empty to list all monitor/checker types; when set, it filters by checker type such as `simpleCheck`, not by monitor resource type such as `trigger`.
