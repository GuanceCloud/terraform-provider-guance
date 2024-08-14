package api

import "github.com/GuanceCloud/terraform-provider-guance/internal/consts"

type Membergroup struct {
	UUID         string   `json:"uuid,omitempty"`
	Name         string   `json:"name"`
	AccountUUIDs []string `json:"accountUUIDs"`
}

type MembergroupContent struct {
	UUID string `json:"uuid"`
}

type MembergroupMember struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
}
type MembergroupReadContent struct {
	GroupMembers []MembergroupMember `json:"groupMembers"`
}

func init() {
	apiURLs[consts.TypeNameMemberGroup] = map[string]string{
		ResourceCreate: "/workspace/member_group/add",
		ResourceDelete: "/workspace/member_group/%s/delete",
		ResourceRead:   "/workspace/member_group/get?groupUUID=%s",
		ResourceUpdate: "/workspace/member_group/%s/modify",
	}
}
