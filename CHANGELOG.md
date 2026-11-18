## 0.1.1

### BUGFIXES
* [guance_pipeline] Keep `extend.logging_index = ""` in state when the API returns `"loggingIndex": ""`, and clear `extend` from state when the remote extend no longer has any managed key, so that these changes show up in the plan.
* [guance_pipeline] Keep `extend` keys that the provider does not manage (for example keys written by the console) on update: the live `extend` is read first and merged, because the modify API replaces the whole object.
* [guance_pipeline] Send `"loggingIndex": ""` when `logging_index` is set to an empty string, and stop sending `"appID": null` and `"measurement": null` when those attributes are not set.
* [guance_pipeline] On update, send the live `loggingIndex` back when neither the prior state nor the plan has `logging_index`, so that an update to other attributes does not remove a remote `"loggingIndex": ""` that the plan never showed.
* [guance_pipeline] Keep `extend.app_id = []` and `extend.measurement = []` as empty lists in state when the API leaves those keys out, so that configured empty lists do not produce a diff on every plan.
* [guance_monitor_json] Ignore the order of tags (`tagInfo`) and alert policies (`alertPolicyNames`) in `checker_json` to prevent spurious updates. Other arrays still compare in order.

### NOTES
* [guance_pipeline] After upgrading, the first plan for a console-created pipeline with `"loggingIndex": ""` can show removal of `extend` if the Terraform configuration omits it. Set `extend = { logging_index = "" }` to keep the empty logging index.

## 0.1.0 (June 16, 2026)

### FEATURES
* [guance_alert_policy] Add alert routing policy resource support.
* [guance_alert_policy] Support checker/security rule bindings, notification targets, aggregation, silence, escalation, permissions, and custom schedules.
* [guance_alert_policy_notice_date] Add custom alert policy notification date resource support.
* [guance_monitor] Expose the structured monitor/checker resource.
* [guance_mute] Add alert policy, checker, tag, and custom mute rule resource support.
* [guance_notify_object] Add alert notification object resource support.
* [data-sources] Add alert policy, notice date, monitor, monitor list, mute, and notify object data sources.

### IMPROVEMENTS
* [docs] Add generated Terraform Registry documentation for alert and monitor resources.
* [docs] Add generated Terraform Registry documentation for the new data sources.
* [examples] Add runnable alert and data source examples, and expand monitor examples.
* [examples] Mark `guance_slo` and `guance_synthetics_test` examples as implementation references.
* [developer] Add `make docs`, `make check-docs`, `make test`, and `make testacc` development targets.
* [developer] Upgrade the provider build target to Go 1.26.2.
* [internal] Add shared Terraform value conversion helpers used by alert and mute resources.

### BUGFIXES
* [guance_alert_policy] Preserve `false`, `0`, empty string, empty-list, and nil-clearing updates.
* [guance_alert_policy] Allow checker and security rule bindings to be cleared without drift.
* [guance_alert_policy] Allow alert target schedules and duration fields to be cleared.
* [guance_alert_policy] Avoid sending empty `df_source` values because the current OpenAPI rejects empty security source updates.
* [guance_alert_policy] Detect remote zero-value changes for alert options during refresh.
* [guance_alert_policy] Preserve configured nested empty values during refresh to avoid follow-up diffs.
* [guance_alert_policy] Page through all alert policy list results for name and notification object lookups.
* [guance_monitor] Validate `extend` JSON during create and update instead of silently omitting invalid payloads.
* [guance_monitor] Detect remote clears for `dashboard_uuid` and `secret` during refresh.
* [guance_monitor] Detect remote `extend` changes during refresh.
* [guance_monitor] Stabilize permissions, tags, alert policy bindings, and backend-expanded `extend` payloads.
* [guance_monitors] Page through all monitor list results.
* [guance_monitors] Clarify checker type filtering and avoid examples that send monitor resource types to the list API.
* [guance_mute] Page through all mute list results when reading by UUID.
* [guance_mute] Page through all mute list results for name-based data source lookups.
* [guance_mute] Stabilize clearable fields, repeated mute windows, declarations, tags, filters, notification targets, and messages.
* [guance_alert_policy_notice_date] Page through all notice date list results for name-based data source lookups.
* [guance_notify_object] Preserve disabled and empty permission updates.
* [guance_notify_object] Page through all notify object list results for name-based data source lookups.

### NOTES
* [guance_monitor] Clearing `secret` with `secret = ""` depends on a pending Forethought OpenAPI fix.
* [guance_alert_policy] Clearing a security target `df_source` depends on a pending Forethought OpenAPI contract clarification.
* [guance_mute] The OpenAPI currently requires non-empty `startTime` and `endTime` values for mute requests.

**Full Changelog**: https://github.com/GuanceCloud/terraform-provider-guance/compare/v0.0.9...v0.1.0
