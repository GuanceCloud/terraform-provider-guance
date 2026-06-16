## 0.1.0 (June 15, 2026)

### FEATURES
* [guance_notify_object] Add resource for managing alert notification objects.
* [guance_alert_policy_notice_date] Add resource for managing custom alert policy notice dates.
* [guance_alert_policy] Add resource for managing alert routing, aggregation, silence, escalation, permissions, checker bindings, and notification targets.
* [guance_mute] Add resource for managing alert policy, checker, tag, and custom mute rules.
* [guance_monitor] Add resource for managing monitor/checker rules with structured Terraform fields.
* [data-sources] Add `guance_notify_object`, `guance_alert_policy_notice_date`, `guance_alert_policy`, `guance_mute`, `guance_monitor`, and `guance_monitors`.

### IMPROVEMENTS
* [docs] Add alert, monitor, and data source examples.
* [docs] Add generated Terraform Registry documentation for alert and monitor resources.
* [docs] Add local documentation generation checks with `make docs` and `make check-docs`.
* [examples] Mark `guance_slo` and `guance_synthetics_test` examples as implementation references because those resources are not registered in this release.

### BUGFIXES
* [guance_alert_policy] Preserve `false`, `0`, empty string, and empty-list updates when Forethought OpenAPI treats omitted fields as keeping existing values.
* [guance_alert_policy] Allow `checker_uuids` and `security_rule_uuids` to be cleared.
* [guance_alert_policy] Fix drift for alert target schedules, `df_source`, upgrade target durations, and nested alert option durations.
* [guance_notify_object] Preserve empty permission updates.
* [guance_mute] Page through all mute list results when reading by UUID instead of stopping after 20 pages.
* [guance_mute] Stabilize update and read behavior for clearable fields, repeated mute windows, declarations, tags, filters, notification targets, and notification messages.
* [guance_monitors] Page through all monitor list results instead of returning only the first 100 matches.
* [guance_monitor] Validate `extend` JSON during create and update instead of silently omitting invalid payloads.
* [guance_monitor] Stabilize read/update behavior for permissions, tags, alert policy bindings, and backend-expanded `extend` payloads.

### NOTES
* [guance_monitor] Clearing `secret` by changing a non-empty value to `secret = ""` depends on a pending Forethought OpenAPI fix. The current OpenAPI implementation accepts the update request but keeps the old secret value.

**Full Changelog**: https://github.com/GuanceCloud/terraform-provider-guance/compare/v0.0.9...v0.1.0
