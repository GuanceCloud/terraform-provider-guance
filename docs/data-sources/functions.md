---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "guance_functions Data Source - guance"
subcategory: ""
description: |-
  Guance Cloud Function is a function development, management and execution platform. It is simple and easy to use, without the need to build Web services from scratch, without managing servers and other infrastructure, just write code and publish, and simply configure to generate HTTP API interfaces for functions.
---

# guance_functions (Data Source)

Guance Cloud Function is a function development, management and execution platform. It is simple and easy to use, without the need to build Web services from scratch, without managing servers and other infrastructure, just write code and publish, and simply configure to generate HTTP API interfaces for functions.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `max_results` (Number) The max results count of the resource will be returned.
- `type_name` (String) The type name of the resource be queried

### Read-Only

- `id` (String) Identifier of the resource.
- `items` (Attributes List) The list of the resource (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--items"></a>
### Nested Schema for `items`

Optional:

- `description` (String) Description
- `func_id` (String) Function ID
- `title` (String) Title

