package client

import "github.com/silenceper/wechat/v2/work/externalcontact"

type ExternalContact struct {
	ExternalContactClient *externalcontact.Client
}

func NewExternalContact(externalcontact ExternalContact) *ExternalContact {
	return &externalcontact
}

func (e *ExternalContact) GetGroupChatList() ([]string, error) {
	var groupChatLists []string
	resp, err := e.ExternalContactClient.GetGroupChatList(&externalcontact.GroupChatListRequest{
		Limit: 1000,
	})
	if err != nil {
		return nil, err
	}
	for _, groupChatList := range resp.GroupChatList {
		groupChatLists = append(groupChatLists, groupChatList.ChatID)
	}
	return groupChatLists, nil
}

func (e *ExternalContact) GetGroupChatDetail(chatID string) (*externalcontact.GroupChatDetailResponse, error) {
	resp, err := e.ExternalContactClient.GetGroupChatDetail(&externalcontact.GroupChatDetailRequest{
		ChatID:   chatID,
		NeedName: 1,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *ExternalContact) AddGroupMsgTemplate(req AddMsgTemplateRequest) (string, []string, error) {
	text := externalcontact.MsgText{
		Content: req.Text.Content,
	}
	resp, err := e.ExternalContactClient.AddMsgTemplate(&externalcontact.AddMsgTemplateRequest{
		ChatType:    "group",
		Sender:      req.Sender,
		Text:        text,
		AllowSelect: req.AllowSelect,
		ChatIDList:  req.ChatIDList,
	})
	if err != nil {
		return "", nil, err
	}
	return resp.MsgID, resp.FailList, nil
}
