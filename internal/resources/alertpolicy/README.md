Alert policy is a set of rules that define when to trigger an alert. You can create alert policies for your data sources, and set up alert targets to receive alerts.

Guance Cloud supports alert policy management for the results of monitor checks, by sending alert notification emails or group message notifications, so that you can know about the abnormal data situation of the monitoring in time, find problems, and solve problems.

Relationships:

```mermaid
graph LR

A[Monitor] --> B[Alert Policy] --> C[Notification]
```

Notes:

1. When a monitor is created, an alert policy must be selected, and the default is selected by default;
2. When a certain alert policy is deleted, the monitor under the deleted alert policy will automatically be classified into the default.