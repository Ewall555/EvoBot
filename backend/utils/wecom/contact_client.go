package wecom

import "EvoBot/backend/utils/wecom/client"

type WecomContactClientt interface {
	UserListID(Cursor string) ([]*client.DeptUser, error)
}

func NewWecomContactClient(conf client.WecomConfig) (WecomContactClientt, error) {
	contactClient := client.NewWork(conf).GetAddressList()
	return client.NewWorkAddressList(client.AddressList{
		AddressList: contactClient,
	}), nil
}
