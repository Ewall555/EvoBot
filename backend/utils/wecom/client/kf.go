package client

import (
	"EvoBot/backend/global"
	"sync"
	"time"

	"github.com/silenceper/wechat/v2/work/kf"
	"go.uber.org/zap"
)

type UserMsgQueue struct {
	mu sync.Mutex
}

type WecomKF struct {
	KFClient   *kf.Client
	UserQueues map[string]*UserMsgQueue
	Mu         *sync.Mutex
}

func NewWecomKF(client *kf.Client) *WecomKF {
	return &WecomKF{
		KFClient:   client,
		UserQueues: make(map[string]*UserMsgQueue),
		Mu:         &sync.Mutex{},
	}
}

func (k *WecomKF) VerifyURL(options SignatureOptions) (echo string, err error) {
	return k.KFClient.VerifyURL(options.SignatureOptions)
}

func (k *WecomKF) SyncMsgs(options SyncMsgOptions) (info SyncMsgSchema, err error) {
	info.SyncMsgSchema, err = k.KFClient.SyncMsg(options.SyncMsgOptions)
	return
}

func (k *WecomKF) SendTextMsg(options SendTextMsgOptions) (err error) {
	options.MsgType = "text"
	if options.ForceImmediate {
		for _, content := range splitContent(options.Text.Content) {
			options.Text.Content = content
			info, err := k.KFClient.SendMsg(options)
			if err != nil {
				global.ZAPLOG.Error(info.Error(), zap.Error(err))
				return err
			}
			time.Sleep(1 * time.Second)
		}
		return nil
	}
	k.Mu.Lock()
	userQueue, exists := k.UserQueues[options.Touser]
	if !exists {
		userQueue = &UserMsgQueue{}
		k.UserQueues[options.Touser] = userQueue
	}
	k.Mu.Unlock()

	userQueue.mu.Lock()
	defer userQueue.mu.Unlock()

	for _, content := range splitContent(options.Text.Content) {
		options.Text.Content = content
		info, err := k.KFClient.SendMsg(options)
		if err != nil {
			global.ZAPLOG.Error(info.Error(), zap.Error(err))
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return
}

func (k *WecomKF) SendMenuMsg(options SendMenuMsgOptions) error {
	options.MsgType = "msgmenu"
	global.ZAPLOG.Debug("", zap.Any("", options))
	info, err := k.KFClient.SendMsg(options)
	if err != nil {
		global.ZAPLOG.Error(info.Error(), zap.Error(err))
		return err
	}
	return nil
}

func (k *WecomKF) SendTextMsgOnEvent(info SendTextMsgOnEventOptions) error {
	sendmsg := struct {
		Code    string `json:"code"`
		MsgID   string `json:"msgid,omitempty"`
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{
		Code:    info.Credential,
		MsgType: WecomMsgTypeText,
	}
	sendmsg.Text.Content = info.Message
	if _, err := k.KFClient.SendMsgOnEvent(sendmsg); err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) SendMenuMsgOnEvent(info SendMenuMsgOnEventOptions) error {
	sendmsg := struct {
		Code    string `json:"code"`
		MsgID   string `json:"msgid,omitempty"`
		MsgType string `json:"msgtype"`
		MsgMenu struct {
			HeadContent string     `json:"head_content,omitempty"`
			List        []MenuItem `json:"list"`
			TailContent string     `json:"tail_content,omitempty"`
		} `json:"msgmenu"`
	}{
		Code:    info.Credential,
		MsgType: "msgmenu",
	}
	sendmsg.MsgMenu.HeadContent = info.HeadContent
	sendmsg.MsgMenu.TailContent = info.TailContent
	sendmsg.MsgMenu.List = info.List
	if _, err := k.KFClient.SendMsgOnEvent(sendmsg); err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) ServiceStateGet(options ServiceStateGetOptions) (info ServiceStateGetSchema, err error) {
	info.ServiceStateGetSchema, err = k.KFClient.ServiceStateGet(options.ServiceStateGetOptions)
	return
}

func (k *WecomKF) ServiceStateTrans(options ServiceStateTransOptions) (info ServiceStateTransSchema, err error) {
	info.ServiceStateTransSchema, err = k.KFClient.ServiceStateTrans(options.ServiceStateTransOptions)
	return
}

func (k *WecomKF) ReceptionistList(kfid string) ([]ReceptionistList, error) {
	var receptionistlsit []ReceptionistList
	receptionistinfo, err := k.KFClient.ReceptionistList(kfid)
	if err != nil {
		return nil, err
	}
	for _, receptionist := range receptionistinfo.ReceptionistList {
		receptionistlsit = append(receptionistlsit, ReceptionistList{
			UserID: receptionist.UserID,
			Status: receptionist.Status,
		})
	}
	return receptionistlsit, nil
}

func (k *WecomKF) ReceptionistAdd(options ReceptionistOptions) error {
	_, err := k.KFClient.ReceptionistAdd(kf.ReceptionistOptions{
		OpenKFID:   options.OpenKFID,
		UserIDList: options.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) ReceptionistDel(options ReceptionistOptions) error {
	_, err := k.KFClient.ReceptionistDel(kf.ReceptionistOptions{
		OpenKFID:   options.OpenKFID,
		UserIDList: options.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) AccountAdd(options AccountAddOptions) (string, error) {
	accountres, err := k.KFClient.AccountAdd(kf.AccountAddOptions{
		Name:    options.Name,
		MediaID: options.MediaID,
	})
	if err != nil {
		return "", err
	}
	return accountres.OpenKFID, nil
}

func (k *WecomKF) AccountUpdate(options AccountUpdateOptions) error {
	_, err := k.KFClient.AccountUpdate(kf.AccountUpdateOptions{
		OpenKFID: options.OpenKFID,
		Name:     options.Name,
		MediaID:  options.MediaID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) AccountDel(kfid string) error {
	_, err := k.KFClient.AccountDel(kf.AccountDelOptions{
		OpenKFID: kfid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) AccountList() ([]AccountInfoSchema, error) {
	accountres, err := k.KFClient.AccountList()
	if err != nil {
		return nil, err
	}
	var accountinfo []AccountInfoSchema
	for _, account := range accountres.AccountList {
		accountinfo = append(accountinfo, AccountInfoSchema{
			OpenKFID: account.OpenKFID,
			Name:     account.Name,
			Avatar:   account.Avatar,
		})
	}
	return accountinfo, nil
}

func (k *WecomKF) AddContactWay(kfid string) (string, error) {
	var opt kf.AddContactWayOptions
	opt.OpenKFID = kfid
	opt.Scene = "evobot"
	info, err := k.KFClient.AddContactWay(opt)
	if err != nil {
		return "", err
	}
	return info.URL, nil
}
