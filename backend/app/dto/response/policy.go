package response

import "EvoBot/backend/app/dto"

type Policy struct {
	UUID string `json:"uuid"`
	dto.Policy
}
