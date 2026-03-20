# Notify Object Resource

This resource manages notify objects in Guance Cloud. Notify objects are used to define how alerts are delivered, such as through DingTalk, WeChat, email, etc.

## Example Usage

```hcl
resource "guance_notify_object" "example" {
  type = "dingTalkRobot"
  name = "Example DingTalk Robot"
  opt_set = jsonencode({
    webhook = "https://oapi.dingtalk.com/robot/send?access_token=example"
    secret  = "example_secret"
  })
  open_permission_set = false
  permission_set = ["wsAdmin"]
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The type of notify object. Valid values are:
  * `dingTalkRobot`
  * `HTTPRequest`
  * `wechatRobot`
  * `mailGroup`
  * `feishuRobot`
  * `sms`
  * `vms`
  * `simpleHTTPRequest`
  * `slackIncomingWebhook`
  * `teamsWorkflowWebhook`
  * `googleChatWebhook`

* `name` - (Required) The name of the notify object.

* `opt_set` - (Required) The configuration options for the notify object, encoded as a JSON string. The structure depends on the `type`:
  * For `dingTalkRobot`: `{"webhook": "string", "secret": "string"}`
  * For `HTTPRequest`: `{"url": "string"}`
  * For `wechatRobot`: `{"webhook": "string"}`
  * For `mailGroup`: `{"to": ["string"]}`
  * For `feishuRobot`: `{"webhook": "string", "secret": "string"}`
  * For `sms`: `{"to": ["string"]}`
  * For `vms`: `{"to": ["string"]}`
  * For `simpleHTTPRequest`: `{"url": "string"}`
  * For `slackIncomingWebhook`: `{"webhook": "string"}`
  * For `teamsWorkflowWebhook`: `{"webhook": "string"}`
  * For `googleChatWebhook`: `{"webhook": "string"}`

* `open_permission_set` - (Optional) Whether to enable custom permission settings. Defaults to `false`.

* `permission_set` - (Optional) The permission set for the notify object. Only applicable if `open_permission_set` is `true`. Can include roles (e.g., `wsAdmin`), member UUIDs, and team UUIDs.

## Attribute Reference

The following attributes are exported:

* `uuid` - The UUID of the notify object.
* `workspace_uuid` - The UUID of the workspace containing the notify object.
* `create_at` - The timestamp when the notify object was created.
* `update_at` - The timestamp when the notify object was last updated.

## Import

Notify objects can be imported using their UUID:

```bash
terraform import guance_notify_object.example <notify_object_uuid>
```
