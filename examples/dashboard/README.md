# Dashboard

A dashboard is a collection of visualizations that you can use to monitor the health of your systems and applications.
Dashboards are made up of one or more panels, which are the visualizations themselves. Each panel displays a single
metric or a single aggregation of metrics.

Dashboards are a great way to visualize your data and monitor your systems. You can use them to track metrics over time,
and to quickly see how your systems are performing. You can also use them to compare metrics from different systems and
applications.

Guance Cloud's dashboard is used to clearly show the range in which the metric data values are located. It is suitable
for slicing messy data into points.

## Create

The first let me create a resource. We will send the create operation to the resource management service

```terraform
resource "guance_dashboard" "demo" {
  name     = "oac-demo"
  manifest = file("${path.module}/dashboard.json")
}
```
