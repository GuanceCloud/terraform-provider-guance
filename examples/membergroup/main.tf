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
