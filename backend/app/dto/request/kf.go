package request

import "EvoBot/backend/app/dto"

type KF struct {
	dto.KF
	StaffList []string `json:"staff_list"`
}
