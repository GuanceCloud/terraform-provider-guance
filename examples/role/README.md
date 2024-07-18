# Role

Guance Cloud supports creating roles to manage the permissions of users. Role management provides users with an intuitive entry for permission management, supporting the flexible adjustment of the permission scope corresponding to different roles, creating new roles for users, and assigning permission scopes to roles to meet the permission needs of different users.

## Create

The first let me create a resource. We will send the create operation to the resource management service

```terraform
resource "guance_role" "role" {
  name = "tf-test-role1"
  desc = "test role"
  keys = ["snapshot.delete", "workspace.readMember"]
}
```