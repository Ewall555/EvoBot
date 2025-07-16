package response

import "EvoBot/backend/app/dto"

type KF struct {
	UUID string `json:"uuid"`
	dto.KF
	Staffs []Staff `json:"staffs"`
}
