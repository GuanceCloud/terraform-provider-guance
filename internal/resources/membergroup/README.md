# Member Group

Member group is a collection of members in a workspace, and member groups can be authorized to access the resources in the workspace.

Member group is an abstract concept, it can be a team, or a department, it can help us build a reasonable organizational
structure, optimize the management efficiency and user experience of the observability platform.

Relationships:

```mermaid
graph LR

A[Workspace] --> B[Member]
A --> C[MemberGroup]
```

## Example Usage

```terraform
variable "email" {
  type = string
}

data "guance_members" "demo" {
  search = var.email
}

resource "guance_membergroup" "demo" {
  name          = "oac-demo2"
  account_uuids = data.guance_members.demo.members[*].uuid
}

output "member" {
  value = data.guance_members.demo.members
}
```
