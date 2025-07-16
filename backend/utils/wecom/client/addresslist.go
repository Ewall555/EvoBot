package client

import "github.com/silenceper/wechat/v2/work/addresslist"

type AddressList struct {
	AddressList *addresslist.Client
}

func NewWorkAddressList(workaddresslist AddressList) *AddressList {
	return &workaddresslist
}

func (a *AddressList) UserListID(Cursor string) ([]*DeptUser, error) {
	result, err := a.AddressList.UserListID(&addresslist.UserListIDRequest{
		Cursor: Cursor,
		Limit:  1000,
	})
	if err != nil {
		return []*DeptUser{}, err
	}
	allDeptUsers := convertToDeptUserPtr(result.DeptUser)
	for result.NextCursor != "" {
		resultmore, err := a.AddressList.UserListID(&addresslist.UserListIDRequest{
			Cursor: result.NextCursor,
			Limit:  1000,
		})
		if err != nil {
			return []*DeptUser{}, err
		}
		result.NextCursor = resultmore.NextCursor
		allDeptUsers = append(allDeptUsers, convertToDeptUserPtr(resultmore.DeptUser)...)
	}

	return allDeptUsers, nil
}

func convertToDeptUserPtr(users []*addresslist.DeptUser) []*DeptUser {
	var deptUsers []*DeptUser
	for _, user := range users {
		deptUsers = append(deptUsers, &DeptUser{
			UserID:     user.UserID,
			Department: user.Department,
		})
	}
	return deptUsers
}
