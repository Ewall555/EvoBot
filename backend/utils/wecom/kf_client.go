package wecom

import (
	"EvoBot/backend/utils/wecom/client"
)

func NewWecomKFClient(conf client.WecomConfig) (*client.WecomKF, error) {
	kFClient, err := client.NewWork(conf).GetKF()
	if err != nil {
		return nil, err
	}
	return client.NewWecomKF(kFClient), nil
}
