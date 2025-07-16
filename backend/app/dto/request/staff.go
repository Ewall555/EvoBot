package request

import "EvoBot/backend/app/dto"

type Staff struct {
	dto.Staff
	PolicyList []string `json:"policy_list"`
}
