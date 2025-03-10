package wecom

import (
	"EvoBot/backend/utils/wecom/client"
)

type WecomExternalContactClient interface {
	GetGroupChatList() ([]string, error)
	AddGroupMsgTemplate(req client.AddMsgTemplateRequest) (string, []string, error)
}

func NewWecomExternalContactClient(conf client.WecomConfig) (WecomExternalContactClient, error) {
	externalContactClient := client.NewWork(conf).GetExternalContact()
	return client.NewExternalContact(client.ExternalContact{
		ExternalContactClient: externalContactClient,
	}), nil
}
