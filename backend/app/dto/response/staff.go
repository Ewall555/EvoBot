package response

import (
	"EvoBot/backend/app/dto"
)

type Staff struct {
	UUID string `json:"uuid"`
	dto.Staff
	Policies []Policy `json:"policies"`
}
