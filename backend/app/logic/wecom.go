package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/app/model"
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"EvoBot/backend/i18n"
	utils "EvoBot/backend/utils/redis"
	"EvoBot/backend/utils/wecom"
	wecomclient "EvoBot/backend/utils/wecom/client"
	"fmt"
	"strings"

	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type WecomLogic struct {
	wecomkf *wecomclient.WecomKF
}

type IWecomLogic interface {
	VerifyURL(req dto.SignatureOptions) (string, error)
	CallbackHandler(encryptedMsg []byte) (err error)

	ConfigList() ([]response.WecomConfigApp, error)
	ConfigUpdate(uuid string, req request.WecomConfigApp) error
	ConfigGet(uuid string) (response.WecomConfigApp, error)

	ReceptionistAdd(kfid string, req request.ReceptionistOptions) error
	ReceptionistDel(kfid string, req request.ReceptionistOptions) error
	ReceptionistList(kfid string) ([]wecomclient.ReceptionistList, error)
	CheckReceptionist(stafflist []string, receptionistlist []wecomclient.ReceptionistList) error

	AccountList() ([]wecomclient.AccountInfoSchema, error)
	AddContactWay(kfid string) (string, error)
}

func NewIWecomLogic() IWecomLogic {
	return &WecomLogic{}
}

func (u *WecomLogic) ConfigUpdate(uuid string, req request.WecomConfigApp) error {
	conf, err := wecomRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return err
	}
	conf.UUID = uuid
	conf.CorpID = req.CorpID
	conf.Secret = req.Secret
	conf.AgentID = req.AgentID
	conf.Token = req.Token
	conf.EncodingAESKey = req.EncodingAESKey
	err = wecomRepo.Update(&conf)
	if err != nil {
		global.ZAPLOG.Error("failed to update Wecom config", zap.Error(err))
		return err
	}
	u.wecomkf, err = wecom.NewWecomKFClient(wecomclient.WecomConfig{
		CorpID:         conf.CorpID,
		Token:          conf.Token,
		EncodingAESKey: conf.EncodingAESKey,
		Secret:         conf.Secret,
		AgentID:        conf.AgentID,
	})
	if err != nil {
		global.ZAPLOG.Error("failed to initialize WecomKFClient", zap.Error(err))
		return err
	}
	return nil
}

func (u *WecomLogic) ConfigGet(uuid string) (response.WecomConfigApp, error) {
	var res response.WecomConfigApp
	conf, err := wecomRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return res, err
	}
	res.UUID = conf.UUID
	res.Type = conf.Type
	res.CorpID = conf.CorpID
	res.Token = conf.Token
	res.EncodingAESKey = conf.EncodingAESKey
	res.AgentID = conf.AgentID
	return res, nil
}

func (u *WecomLogic) ConfigList() ([]response.WecomConfigApp, error) {
	var resp []response.WecomConfigApp
	configs, err := wecomRepo.List()
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return resp, err
	}
	for _, config := range configs {
		var res response.WecomConfigApp
		res.UUID = config.UUID
		res.Type = config.Type
		res.CorpID = config.CorpID
		res.Token = config.Token
		res.EncodingAESKey = config.EncodingAESKey
		res.AgentID = config.AgentID
		resp = append(resp, res)
	}
	return resp, nil
}

func (u *WecomLogic) VerifyURL(req dto.SignatureOptions) (echo string, err error) {
	if err = u.initializeWecomClient(); err != nil {
		return
	}
	return u.wecomkf.VerifyURL(wecomclient.SignatureOptions{SignatureOptions: req.SignatureOptions.SignatureOptions})
}

func (u *WecomLogic) CallbackHandler(encryptedMsg []byte) (err error) {
	if err = u.initializeWecomClient(); err != nil {
		return
	}
	callbackmessage, err := u.wecomkf.KFClient.GetCallbackMessage(encryptedMsg)
	if err != nil {
		global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "GetCallbackMessage"), zap.Error(err))
		return
	}
	switch callbackmessage.MsgType {
	case wecomclient.WecomCallbackMsgTypeEvent:
		switch callbackmessage.Event {
		case wecomclient.WecomCallbackEventKFMsgOrEvent:
			go func() {
				var options wecomclient.SyncMsgOptions
				options.Token = callbackmessage.Token
				options.OpenKfID = callbackmessage.OpenKfID
				if err := u.SyncMsg(options); err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SyncMsgs"), zap.Error(err))
					return
				}
			}()
		default:
			global.ZAPLOG.Error(i18n.T("wecom.unknow_event_type"), zap.Error(err))
		}
	default:
		global.ZAPLOG.Error(i18n.T("wecom.unknow_msgtype"), zap.Error(err))
	}
	return
}

func (u *WecomLogic) SyncMsg(options wecomclient.SyncMsgOptions) (err error) {
	ctx := context.Background()
	wecomCursorLockKey := wecomclient.KeyWecomCursorLockPrefix + options.OpenKfID
	wecomCursorLockVal, ok := utils.AcquireRedisLockWithRetry(ctx, global.RDS, wecomCursorLockKey, 30*time.Second, 500*time.Millisecond, 30*time.Second)
	if !ok {
		global.ZAPLOG.Warn("无法获取同步锁，跳过本轮")
		return nil
	}
	defer utils.ReleaseRedisLock(ctx, global.RDS, wecomCursorLockKey, wecomCursorLockVal)

	wecomCursorKey := wecomclient.KeyWecomCursorPrefix + options.OpenKfID
	options.Cursor, err = global.RDS.Get(ctx, wecomCursorKey).Result()
	if err != nil && err != redis.Nil {
		global.ZAPLOG.Error("Get Wecom Cursor error", zap.Error(err))
		return err
	}
	for {
		resp, err := u.wecomkf.SyncMsgs(options)
		if err != nil {
			break
		}
		for _, msg := range resp.MsgList {
			go func() {
				if err := u.processMessage(wecomclient.Message{
					Message: msg,
				}); err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "processMessage"), zap.Error(err))
				}
			}()
		}
		if err = global.RDS.Set(ctx, wecomCursorKey, resp.NextCursor, 0).Err(); err != nil {
			break
		}
		if resp.HasMore == 0 {
			break
		}
		options.Cursor = resp.NextCursor
	}
	return
}

func (u *WecomLogic) ReceptionistAdd(kfid string, req request.ReceptionistOptions) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	err := u.wecomkf.ReceptionistAdd(wecomclient.ReceptionistOptions{
		OpenKFID:   kfid,
		UserIDList: req.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *WecomLogic) ReceptionistDel(kfid string, req request.ReceptionistOptions) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	err := u.wecomkf.ReceptionistDel(wecomclient.ReceptionistOptions{
		OpenKFID:   kfid,
		UserIDList: req.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *WecomLogic) ReceptionistList(kfid string) ([]wecomclient.ReceptionistList, error) {
	if err := u.initializeWecomClient(); err != nil {
		return nil, err
	}
	list, err := u.wecomkf.ReceptionistList(kfid)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (u *WecomLogic) CheckReceptionist(stafflist []string, receptionistlist []wecomclient.ReceptionistList) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	var invalidUserIDs []string
	receptionistIDSet := make(map[string]struct{}, len(receptionistlist))
	for _, receptionist := range receptionistlist {
		receptionistIDSet[receptionist.UserID] = struct{}{}
	}
	staffs, err := staffRepo.List(commonRepo.WithUUIDsIn(stafflist))
	if err != nil {
		return err
	}
	for _, staff := range staffs {
		if _, exists := receptionistIDSet[staff.StaffID]; !exists {
			invalidUserIDs = append(invalidUserIDs, staff.StaffID)
		}
	}
	if len(invalidUserIDs) > 0 {
		return fmt.Errorf("ErrStaffIDLIST: %v", invalidUserIDs)
	}
	return nil
}

func (u *WecomLogic) AccountList() ([]wecomclient.AccountInfoSchema, error) {
	if err := u.initializeWecomClient(); err != nil {
		return nil, err
	}
	accountinfolist, err := u.wecomkf.AccountList()
	if err != nil {
		return nil, err
	}
	return accountinfolist, nil
}

func (u *WecomLogic) AddContactWay(kfid string) (string, error) {
	if err := u.initializeWecomClient(); err != nil {
		return "", err
	}
	url, err := u.wecomkf.AddContactWay(kfid)
	if err != nil {
		return "", err
	}
	return url, nil
}

// Helper
func (u *WecomLogic) initializeWecomClient() (err error) {
	if u.wecomkf != nil {
		return
	}
	conf, err := wecomRepo.Get(wecomRepo.WithType("app"))
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return
	}
	if conf.CorpID == "" || conf.AgentID == "" || conf.EncodingAESKey == "" || conf.Secret == "" || conf.Token == "" {
		global.ZAPLOG.Error("wecom config is nil")
		return
	}
	u.wecomkf, err = wecom.NewWecomKFClient(wecomclient.WecomConfig{
		CorpID:         conf.CorpID,
		Token:          conf.Token,
		EncodingAESKey: conf.EncodingAESKey,
		Secret:         conf.Secret,
		AgentID:        conf.AgentID,
	})
	if err != nil {
		global.ZAPLOG.Error("failed to initialize WecomKFClient", zap.Error(err))
		return
	}
	return
}

func (u *WecomLogic) processMessage(msg wecomclient.Message) (err error) {
	err = kHRepo.FirstOrCreatebyKHID(model.KH{KHID: msg.ExternalUserID})
	if err != nil {
		global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "FirstOrCreate"), zap.Error(err))
		return
	}
	switch msg.MsgType {
	case wecomclient.WecomMsgTypeText:
		textMessage, err := msg.GetTextMessage()
		if err != nil {
			return err
		}
		global.ZAPLOG.Debug("", zap.Any("", textMessage))
		return u.handleTextMessage(wecomclient.Text{Text: textMessage})
	case wecomclient.WecomMsgTypeEvent:
		switch msg.EventType {
		case wecomclient.WecomMsgTypeEnterSessionEvent:
			enterSessionEvent, err := msg.GetEnterSessionEvent()
			if err != nil {
				return err
			}
			global.ZAPLOG.Debug("", zap.Any("", enterSessionEvent))
			return u.handleEnterSessionEvent(wecomclient.EnterSessionEvent{EnterSessionEvent: enterSessionEvent})
		case wecomclient.WecomMsgTypeMsgSendFailEvent:
			msgSendFailEvent, err := msg.GetMsgSendFailEvent()
			if err != nil {
				return err
			}
			global.ZAPLOG.Debug("", zap.Any("", msgSendFailEvent))
			return err
		case wecomclient.WecomMsgTypeSessionStatusChangeEvent:
			sessionStatusChangeEvent, err := msg.GetSessionStatusChangeEvent()
			if err != nil {
				return err
			}
			global.ZAPLOG.Debug("", zap.Any("", sessionStatusChangeEvent))
			return u.handleSessionStatusChangeEvent(wecomclient.SessionStatusChangeEvent{SessionStatusChangeEvent: sessionStatusChangeEvent})
		}
	default:

	}
	return
}

func (u *WecomLogic) handleTextMessage(textMessage wecomclient.Text) (err error) {
	var options wecomclient.ServiceStateGetOptions
	options.OpenKFID = textMessage.OpenKFID
	options.ExternalUserID = textMessage.ExternalUserID
	serviceStateInfo, err := u.wecomkf.ServiceStateGet(options)
	if err != nil {
		global.ZAPLOG.Error(serviceStateInfo.Error(), zap.Error(err))
		return
	}
	switch serviceStateInfo.ServiceState {
	case wecomclient.SessionStatusNew:
		return u.handleNewSession(textMessage)
	case wecomclient.SessionStatusHandled:
		return u.handleBotSession(textMessage)
	case wecomclient.SessionStatusWaiting:

	case wecomclient.SessionStatusInProgress:
		return u.handleInProgress(textMessage)
	case wecomclient.SessionStatusEndedOrNotStarted:

	default:

	}
	return
}

func (u *WecomLogic) handleNewSession(textMessage wecomclient.Text) (err error) {
	var options wecomclient.ServiceStateTransOptions
	options.OpenKFID = textMessage.OpenKFID
	options.ExternalUserID = textMessage.ExternalUserID
	options.ServiceState = wecomclient.SessionStatusHandled
	serviceStateRespInfo, err := u.wecomkf.ServiceStateTrans(options)
	if err != nil {
		global.ZAPLOG.Error(serviceStateRespInfo.Error(), zap.Error(err))
		return
	}
	return u.handleBotSession(textMessage)
}

func (u *WecomLogic) handleBotSession(textMessage wecomclient.Text) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(textMessage.OpenKFID))
	if err != nil {
		return
	}
	switch kFInfo.Status {
	case constant.KFStatusRobotToHuman:
		switch {
		case textMessage.Text.Text.Content == "企业版购买咨询":
			var sendTextMsgOptions wecomclient.SendTextMsgOptions
			sendTextMsgOptions.Touser = textMessage.ExternalUserID
			sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
			sendTextMsgOptions.Text.Content = "正在为你转接人工..."
			err = u.wecomkf.SendTextMsg(sendTextMsgOptions)
			if err != nil {
				global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
			}
			err = u.handleTransferToStaff(textMessage, kFInfo, "售前")
			if err != nil {
				global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "handleTransferToStaff"), zap.Error(err))
			}
			return kHRepo.UpdatebyKHID(model.KH{VerifyStatus: constant.KHStatusUnprocessed, KHID: textMessage.ExternalUserID})
		case textMessage.Text.Text.Content == "售后服务支持":
			var sendTextMsgOptions wecomclient.SendTextMsgOptions
			sendTextMsgOptions.Touser = textMessage.ExternalUserID
			sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
			sendTextMsgOptions.Text.Content = "您好，感谢联系飞致云，很高兴为您服务！请输入产品序列号、提供公司全称或简称，以便确认您所使用产品信息。产品序列号由字母和数字组成，类似“JSN1234567X”，可登陆产品，在“系统设置-许可证”页面，进行查看。（例如：杭州飞致云 或者 JSN1234567X）"
			err = u.wecomkf.SendTextMsg(sendTextMsgOptions)
			if err != nil {
				global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
			}
			verifyKey := fmt.Sprintf("verify_attempt:%s", textMessage.ExternalUserID)
			global.RDS.Set(context.Background(), verifyKey, 0, 900*time.Second)
			return kHRepo.UpdatebyKHID(model.KH{VerifyStatus: constant.KHStatusVerification, KHID: textMessage.ExternalUserID})
		case textMessage.Text.Text.Content == "返回和智能助手对话":
			u.wecomkf.SendMenuMsg(wecomclient.SendMenuMsgOptions{
				BaseSendMsgOptions: wecomclient.BaseSendMsgOptions{
					OpenKfid: textMessage.OpenKFID,
					Touser:   textMessage.ExternalUserID,
				},
				MenuMsgOptions: parseMenuText(kFInfo.BotWelcomeMsg),
			})
			return kHRepo.UpdatebyKHID(model.KH{VerifyStatus: constant.KHStatusUnprocessed, KHID: textMessage.ExternalUserID})
		}
		kHInfo, _ := kHRepo.Get(kHRepo.WithKHID(textMessage.ExternalUserID))
		switch kHInfo.VerifyStatus {
		case 0, constant.KHStatusUnprocessed:
			if strings.Contains(kFInfo.TransferKeywords, textMessage.Text.Text.Content) {
				selectRoleMsg := "#H 您好，请选择需要的咨询:\n#TXT -------------------------\n#CLK 企业版购买咨询\n#TXT -------------------------\n#CLK 售后服务支持\n#TXT -------------------------\n#CLK 返回和智能助手对话\n#T -------------------------"
				err = u.wecomkf.SendMenuMsg(wecomclient.SendMenuMsgOptions{
					BaseSendMsgOptions: wecomclient.BaseSendMsgOptions{
						Touser:         textMessage.ExternalUserID,
						OpenKfid:       textMessage.OpenKFID,
						ForceImmediate: true,
					},
					MenuMsgOptions: parseMenuText(selectRoleMsg),
				})
				if err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendMenuMsg"), zap.Error(err))
				}
				return
			}
			return u.handleBotReply(textMessage)
		case constant.KHStatusVerification:
			const (
				verifyMaxAttempts = 3
			)
			verifyKey := fmt.Sprintf("verify_attempt:%s", textMessage.ExternalUserID)
			ctx := context.Background()
			attempts, err := global.RDS.Get(ctx, verifyKey).Int()
			if err == redis.Nil {
				kHRepo.UpdatebyKHID(model.KH{
					VerifyStatus: constant.KHStatusUnprocessed,
					KHID:         textMessage.ExternalUserID,
				})
				return u.handleBotSession(textMessage)
			}
			if err != nil {
				return err
			}
			var sendTextMsgOptions wecomclient.SendTextMsgOptions
			sendTextMsgOptions.Touser = textMessage.ExternalUserID
			sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
			if attempts >= verifyMaxAttempts {
				sendTextMsgOptions.Text.Content = "验证失败次数已达上限，已返回智能助手的对话。"
				if err := u.wecomkf.SendTextMsg(sendTextMsgOptions); err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
					return err
				}
				_ = global.RDS.Del(ctx, verifyKey).Err()
				u.wecomkf.SendMenuMsg(wecomclient.SendMenuMsgOptions{
					BaseSendMsgOptions: wecomclient.BaseSendMsgOptions{
						OpenKfid: textMessage.OpenKFID,
						Touser:   textMessage.ExternalUserID,
					},
					MenuMsgOptions: parseMenuText(kFInfo.BotWelcomeMsg),
				})
				return kHRepo.UpdatebyKHID(model.KH{VerifyStatus: constant.KHStatusUnprocessed, KHID: textMessage.ExternalUserID})
			}
		case constant.KHStatusUserInfoConfirm:
			switch textMessage.Text.Text.Content {
			case "正确":
				var sendTextMsgOptions wecomclient.SendTextMsgOptions
				sendTextMsgOptions.Touser = textMessage.ExternalUserID
				sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
				sendTextMsgOptions.Text.Content = "正在为你转接人工..."
				err = u.wecomkf.SendTextMsg(sendTextMsgOptions)
				if err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
				}
				err = u.handleTransferToStaff(textMessage, kFInfo, "售后")
				if err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "handleTransferToStaff"), zap.Error(err))
				}
				return kHRepo.UpdatebyKHID(model.KH{VerifyStatus: constant.KHStatusUnprocessed, KHID: textMessage.ExternalUserID})
			case "不正确":
				var sendTextMsgOptions wecomclient.SendTextMsgOptions
				sendTextMsgOptions.Touser = textMessage.ExternalUserID
				sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
				sendTextMsgOptions.Text.Content = "已返回和智能助手的对话，可重新发起转人工请求。"
				err = u.wecomkf.SendTextMsg(sendTextMsgOptions)
				if err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
				}
				u.wecomkf.SendMenuMsg(wecomclient.SendMenuMsgOptions{
					BaseSendMsgOptions: wecomclient.BaseSendMsgOptions{
						OpenKfid: textMessage.OpenKFID,
						Touser:   textMessage.ExternalUserID,
					},
					MenuMsgOptions: parseMenuText(kFInfo.BotWelcomeMsg),
				})
				return kHRepo.UpdatebyKHID(model.KH{VerifyStatus: constant.KHStatusUnprocessed, KHID: textMessage.ExternalUserID})
			default:
				var sendTextMsgOptions wecomclient.SendTextMsgOptions
				sendTextMsgOptions.Touser = textMessage.ExternalUserID
				sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
				sendTextMsgOptions.Text.Content = "已退出转人工流程，返回和智能助手对话"
				err = u.wecomkf.SendTextMsg(sendTextMsgOptions)
				if err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
				}
				err = kHRepo.UpdatebyKHID(model.KH{VerifyStatus: constant.KHStatusUnprocessed, KHID: textMessage.ExternalUserID})
				if err != nil {
					return
				}
				return u.handleBotReply(textMessage)
			}
		default:
			return u.handleBotReply(textMessage)
		}
	case constant.KFStatusOnlyRobot:
		return u.handleBotReply(textMessage)
	case constant.KFStatusOnlyHuman:

	default:

	}
	return
}

func (u *WecomLogic) handleBotReply(textMessage wecomclient.Text) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(textMessage.OpenKFID))
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(kFInfo.BotTimeout)*time.Second)
	defer cancel()
	resultChan := make(chan string)
	errorChan := make(chan error)
	var sendTextMsgOptions wecomclient.SendTextMsgOptions
	sendTextMsgOptions.Touser = textMessage.ExternalUserID
	sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
	go func() {
		llmappLogic := NewILLMAppLogic()
		fullContent, err := llmappLogic.ChatMessage(textMessage.ExternalUserID, kFInfo.BotID, textMessage.Text.Text.Content)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- MarkdownToText(fullContent)
		}
	}()
	go func(ctx context.Context, resultChan chan string, errorChan chan error) {
		select {
		case <-ctx.Done():
			sendTextMsgOptions.Text.Content = kFInfo.BotTimeoutMsg
			err = u.wecomkf.SendTextMsg(sendTextMsgOptions)
			if err != nil {
				global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
			}
			select {
			case fullContent := <-resultChan:
				sendTextMsgOptions.Text.Content = fullContent
			case err = <-errorChan:
				sendTextMsgOptions.Text.Content = err.Error()
			}
		case fullContent := <-resultChan:
			sendTextMsgOptions.Text.Content = fullContent
		case err = <-errorChan:
			sendTextMsgOptions.Text.Content = err.Error()
		}
		if err := u.wecomkf.SendTextMsg(sendTextMsgOptions); err != nil {
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
			return
		}
	}(ctx, resultChan, errorChan)
	chatkey := constant.KeyWecomBotStaffPrefix + kFInfo.BotID + ":" + textMessage.ExternalUserID
	remainingTTL, err := global.RDS.TTL(context.Background(), chatkey).Result()
	if err != nil {
		global.ZAPLOG.Debug("TTL check failed", zap.Error(err))
	}
	if remainingTTL >= 0 {
		global.ZAPLOG.Debug("刷新会话超时时间")
		if err = global.RDS.Set(context.Background(), chatkey, 1, time.Duration(kFInfo.ChatTimeout)*time.Second).Err(); err != nil {
			global.ZAPLOG.Error("redis set error", zap.Error(err))
			return
		}
		return
	}
	global.ZAPLOG.Debug("初始化会话")
	if err = global.RDS.Set(context.Background(), chatkey, 1, time.Duration(kFInfo.ChatTimeout)*time.Second).Err(); err != nil {
		global.ZAPLOG.Error("redis set error", zap.Error(err))
		return
	}
	go func() {
		ticker := time.NewTicker(time.Duration(kFInfo.ChatTimeout) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				remainingTTL, err := global.RDS.TTL(context.Background(), chatkey).Result()
				if err != nil {
					global.ZAPLOG.Debug("TTL check failed", zap.Error(err))
					continue
				}
				if remainingTTL == -2 {
					var serviceStateTransOptions wecomclient.ServiceStateTransOptions
					serviceStateTransOptions.OpenKFID = textMessage.OpenKFID
					serviceStateTransOptions.ExternalUserID = textMessage.ExternalUserID
					serviceStateTransOptions.ServiceState = wecomclient.SessionStatusEndedOrNotStarted
					serviceStateRespInfo, err := u.wecomkf.ServiceStateTrans(serviceStateTransOptions)
					if err != nil {
						global.ZAPLOG.Error(serviceStateRespInfo.Error(), zap.Error(err))
						continue
					}
					if err := kHRepo.ClearChatIDsByKHIDAndBotID(textMessage.ExternalUserID, kFInfo.BotID); err != nil {
						global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "ClearChatIDsByKHIDAndBotID"), zap.Error(err))
						return
					}
					return
				}
				if remainingTTL > 0 {
					ticker.Stop()
					ticker = time.NewTicker(remainingTTL)
					continue
				}
				return
			case <-context.Background().Done():
				return
			}
		}
	}()
	return
}

func (u *WecomLogic) handleInProgress(textMessage wecomclient.Text) (err error) {
	if textMessage.Origin != wecomclient.MessageTypeCustomer {
		return
	}
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(textMessage.OpenKFID))
	if err != nil {
		return
	}
	chatkey := constant.KeyWecomKFStaffPrefix + textMessage.OpenKFID + ":" + textMessage.ExternalUserID
	if err = global.RDS.Set(context.Background(), chatkey, 1, time.Duration(kFInfo.ChatTimeout)*time.Second).Err(); err != nil {
		global.ZAPLOG.Error("redis set error", zap.Error(err))
		return
	}
	return
}

func (u *WecomLogic) handleEnterSessionEvent(enterSessionEvent wecomclient.EnterSessionEvent) (err error) {
	var options wecomclient.ServiceStateGetOptions
	options.OpenKFID = enterSessionEvent.OpenKFID
	options.ExternalUserID = enterSessionEvent.ExternalUserID
	serviceStateInfo, err := u.wecomkf.ServiceStateGet(options)
	if err != nil {
		global.ZAPLOG.Error(serviceStateInfo.Error(), zap.Error(err))
		return
	}
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(enterSessionEvent.OpenKFID))
	if err != nil {
		return err
	}
	if serviceStateInfo.ServiceState == wecomclient.SessionStatusHandled {
		return u.wecomkf.SendMenuMsg(wecomclient.SendMenuMsgOptions{
			BaseSendMsgOptions: wecomclient.BaseSendMsgOptions{
				OpenKfid: enterSessionEvent.OpenKFID,
				Touser:   enterSessionEvent.ExternalUserID,
			},
			MenuMsgOptions: parseMenuText(kFInfo.BotWelcomeMsg),
		})
	}
	if enterSessionEvent.Event.WelcomeCode == "" {
		return
	}
	return u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
		Credential:     enterSessionEvent.Event.WelcomeCode,
		MenuMsgOptions: parseMenuText(kFInfo.BotWelcomeMsg),
	})
}

// 处理会话状态变更事件
func (u *WecomLogic) handleSessionStatusChangeEvent(sessionStatusChangeEvent wecomclient.SessionStatusChangeEvent) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(sessionStatusChangeEvent.OpenKFID))
	if err != nil {
		return err
	}
	switch sessionStatusChangeEvent.Event.ChangeType {
	case wecomclient.WecomEventChangeTypeJoinSession:
	case wecomclient.WecomEventChangeTypeTransferSession:
		isStaffWorkByStaffID(sessionStatusChangeEvent.Event.NewReceptionistUserID)
		var serviceStateTransOptions wecomclient.ServiceStateTransOptions
		serviceStateTransOptions.OpenKFID = sessionStatusChangeEvent.OpenKFID
		serviceStateTransOptions.ExternalUserID = sessionStatusChangeEvent.ExternalUserID
		serviceStateTransOptions.ServicerUserID = sessionStatusChangeEvent.Event.NewReceptionistUserID
		staffweightkey := constant.KeyStaffWeightPrefix + sessionStatusChangeEvent.Event.OldReceptionistUserID
		if err := global.RDS.Incr(context.Background(), staffweightkey).Err(); err != nil {
			global.ZAPLOG.Error("Incr", zap.Error(err))
		}
		return u.handleServiceStateTransInProgress(serviceStateTransOptions)
	case wecomclient.WecomEventChangeTypeEndSession:
		if err := u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
			Credential:     sessionStatusChangeEvent.Event.MsgCode,
			MenuMsgOptions: parseMenuText(kFInfo.ChatendMsg),
		}); err != nil {
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendMenuMsgOnEvent"), zap.Error(err))
		}
		chatkey := constant.KeyWecomKFStaffPrefix + sessionStatusChangeEvent.OpenKFID + ":" + sessionStatusChangeEvent.ExternalUserID
		if err = global.RDS.Del(context.Background(), chatkey).Err(); err != nil {
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "RDS.Del chatkey"), zap.Error(err))
			return
		}
		return
	case wecomclient.WecomEventChangeTypeRejoinSession:
		global.ZAPLOG.Debug("WecomEventChangeTypeRejoinSession")
		isStaffWorkByStaffID(sessionStatusChangeEvent.Event.NewReceptionistUserID)
		var serviceStateTransOptions wecomclient.ServiceStateTransOptions
		serviceStateTransOptions.OpenKFID = sessionStatusChangeEvent.OpenKFID
		serviceStateTransOptions.ExternalUserID = sessionStatusChangeEvent.ExternalUserID
		serviceStateTransOptions.ServicerUserID = sessionStatusChangeEvent.Event.NewReceptionistUserID
		return u.handleServiceStateTransInProgress(serviceStateTransOptions)
	default:
		global.ZAPLOG.Debug(i18n.T("wecom.unknow_msgtype"))
	}
	return
}

func (u *WecomLogic) handleServiceStateTransInProgress(serviceStateTransOptions wecomclient.ServiceStateTransOptions) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(serviceStateTransOptions.OpenKFID))
	if err != nil {
		return
	}
	var options wecomclient.ServiceStateGetOptions
	options.OpenKFID = serviceStateTransOptions.OpenKFID
	options.ExternalUserID = serviceStateTransOptions.ExternalUserID
	serviceStateInfo, err := u.wecomkf.ServiceStateGet(options)
	if err != nil {
		global.ZAPLOG.Error(serviceStateInfo.Error(), zap.Error(err))
		return
	}
	if serviceStateInfo.ServiceState != wecomclient.SessionStatusInProgress {
		serviceStateTransOptions.ServiceState = wecomclient.SessionStatusInProgress
		Info, err := u.wecomkf.ServiceStateTrans(serviceStateTransOptions)
		if err != nil {
			global.ZAPLOG.Error(Info.Error(), zap.Error(err))
			return err
		}
		if Info.MsgCode != "" {
			if err = u.wecomkf.SendTextMsgOnEvent(wecomclient.SendTextMsgOnEventOptions{
				Message:    kFInfo.StaffWelcomeMsg,
				Credential: Info.MsgCode,
			}); err != nil {
				return err
			}
		}
	}

	global.ZAPLOG.Debug("更新客户的上一次接待人员为选出的接待人员")
	if err := kHRepo.UpdatebyKHID(model.KH{
		KHID:    serviceStateTransOptions.ExternalUserID,
		StaffID: serviceStateTransOptions.ServicerUserID,
	}); err != nil {
		global.ZAPLOG.Error("更新客户的上一次接待人员失败", zap.Error(err))
	}

	global.ZAPLOG.Debug("降低该接待人员空闲权重")
	ctx := context.Background()
	staffweightkey := constant.KeyStaffWeightPrefix + serviceStateTransOptions.ServicerUserID
	if err = global.RDS.Decr(ctx, staffweightkey).Err(); err != nil {
		return
	}

	global.ZAPLOG.Debug("缓存会话超时时间")
	chatkey := constant.KeyWecomKFStaffPrefix + serviceStateTransOptions.OpenKFID + ":" + serviceStateTransOptions.ExternalUserID
	remainingTTL, err := global.RDS.TTL(ctx, chatkey).Result()
	if err != nil {
		global.ZAPLOG.Debug("TTL check failed", zap.Error(err))
	}

	if remainingTTL >= 0 {
		return
	}

	if err = global.RDS.Set(ctx, chatkey, 1, time.Duration(kFInfo.ChatTimeout)*time.Second).Err(); err != nil {
		global.ZAPLOG.Error("redis set error", zap.Error(err))
		return
	}

	global.ZAPLOG.Debug("handleChatTimeoutTask")
	go func() {
		ticker := time.NewTicker(time.Duration(kFInfo.ChatTimeout) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				remainingTTL, err := global.RDS.TTL(ctx, chatkey).Result()
				if err != nil {
					global.ZAPLOG.Debug("TTL check failed", zap.Error(err))
					continue
				}
				if remainingTTL == -2 {
					var options wecomclient.ServiceStateGetOptions
					options.OpenKFID = serviceStateTransOptions.OpenKFID
					options.ExternalUserID = serviceStateTransOptions.ExternalUserID
					serviceStateInfo, err := u.wecomkf.ServiceStateGet(options)
					if err != nil {
						global.ZAPLOG.Error(serviceStateInfo.Error(), zap.Error(err))
						return
					}
					if serviceStateInfo.ServiceState != wecomclient.SessionStatusEndedOrNotStarted {
						serviceStateTransOptions.ServiceState = wecomclient.SessionStatusEndedOrNotStarted
						serviceStateRespInfo, err := u.wecomkf.ServiceStateTrans(serviceStateTransOptions)
						if err != nil {
							global.ZAPLOG.Error(serviceStateRespInfo.Error(), zap.Error(err))
							continue
						}
						global.ZAPLOG.Debug("会话超时，已变更状态")
						kFInfo, err := kFRepo.Get(kFRepo.WithKFID(serviceStateTransOptions.OpenKFID))
						if err != nil {
							return
						}
						if err := u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
							Credential:     serviceStateRespInfo.MsgCode,
							MenuMsgOptions: parseMenuText(kFInfo.ChatendMsg),
						}); err != nil {
							global.ZAPLOG.Error("SendTextMsgOnEvent", zap.Error(err))
						}
					}
					kHInfo, err := kHRepo.Get(kHRepo.WithKHID(serviceStateTransOptions.ExternalUserID))
					if err != nil {
						global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "kHRepo.Get"), zap.Error(err))
						return
					}
					staffweightkey := constant.KeyStaffWeightPrefix + kHInfo.StaffID
					if err := global.RDS.Incr(ctx, staffweightkey).Err(); err != nil {
						global.ZAPLOG.Error("Incr", zap.Error(err))
					}
					return
				}
				if remainingTTL >= 0 {
					ticker.Stop()
					ticker = time.NewTicker(remainingTTL)
					continue
				}
				return
			case <-ctx.Done():
				return
			}
		}
	}()
	return
}

func (u *WecomLogic) handleTransferToStaff(textMessage wecomclient.Text, kFInfo model.KF, role string) (err error) {
	if kFInfo.ReceivePriority == 1 {
		kHInfo, err := kHRepo.Get(kHRepo.WithKHID(textMessage.ExternalUserID))
		if err != nil {
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "Create"), zap.Error(err))
			return err
		}
		if kHInfo.StaffID != "" {
			for _, staffInfo := range kFInfo.Staffs {
				if staffInfo.StaffID == kHInfo.StaffID {
					if staffInfo.Role != role {
						global.ZAPLOG.Debug("服务人员角色不匹配，跳过", zap.String("staffid", staffInfo.StaffID), zap.String("角色", staffInfo.Role))
						continue
					}
					isStaffWork, err := isStaffWorkByStaffID(kHInfo.StaffID)
					if err != nil {
						global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "isStaffWorkByStaffID"), zap.Error(err))
						return err
					}
					if isStaffWork {
						var serviceStateTransOptions wecomclient.ServiceStateTransOptions
						serviceStateTransOptions.OpenKFID = textMessage.OpenKFID
						serviceStateTransOptions.ExternalUserID = textMessage.ExternalUserID
						serviceStateTransOptions.ServicerUserID = kHInfo.StaffID
						return u.handleServiceStateTransInProgress(serviceStateTransOptions)
					}
				}
			}
		}
	}
	switch kFInfo.ReceiveRule {
	case constant.KFReceiveRuleRoundRobin:
		global.ZAPLOG.Debug("ToDo")
	case constant.KFReceiveRuleIdle:
		var staffIDs []string
		for _, staffinfo := range kFInfo.Staffs {
			global.ZAPLOG.Debug("staffinfo", zap.Any("", staffinfo))
			isStaffWork, err := isStaffWorkByStaffID(staffinfo.StaffID)
			if err != nil {
				return err
			}
			if isStaffWork && staffinfo.Role == role {
				staffIDs = append(staffIDs, staffinfo.StaffID)
			}
		}
		selectedstaff, _ := getHighestWeightStaff(staffIDs)
		if selectedstaff != "" {
			var serviceStateTransOptions wecomclient.ServiceStateTransOptions
			serviceStateTransOptions.OpenKFID = textMessage.OpenKFID
			serviceStateTransOptions.ExternalUserID = textMessage.ExternalUserID
			serviceStateTransOptions.ServicerUserID = selectedstaff
			return u.handleServiceStateTransInProgress(serviceStateTransOptions)
		}
		global.ZAPLOG.Debug("没有找到空闲的接待人员")
		var sendTextMsgOptions wecomclient.SendTextMsgOptions
		sendTextMsgOptions.Touser = textMessage.ExternalUserID
		sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
		sendTextMsgOptions.Text.Content = kFInfo.UnmannedMsg
		return u.wecomkf.SendTextMsg(sendTextMsgOptions)
	default:

	}
	return
}

func parseMenuText(text string) wecomclient.MenuMsgOptions {
	var resp wecomclient.MenuMsgOptions
	lines := strings.Split(text, "\n")
	menuList := []wecomclient.MenuItem{}
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "#H "):
			if resp.HeadContent != "" {
				resp.HeadContent += "\n"
			}
			resp.HeadContent += strings.TrimPrefix(line, "#H ")
		case strings.HasPrefix(line, "#T "):
			if resp.TailContent != "" {
				resp.TailContent += "\n"
			}
			resp.TailContent += strings.TrimPrefix(line, "#T ")
		case strings.HasPrefix(line, "#TXT "):
			menuList = append(menuList, wecomclient.MenuItem{
				Type: "text",
				Text: &struct {
					Content   string `json:"content"`
					NoNewLine int    `json:"no_newline,omitempty"`
				}{
					Content:   strings.TrimPrefix(line, "#TXT "),
					NoNewLine: 0,
				},
			})
		case strings.HasPrefix(line, "#CLK "):
			menuList = append(menuList, wecomclient.MenuItem{
				Type: "click",
				Click: &struct {
					ID      string `json:"id,omitempty"`
					Content string `json:"content"`
				}{
					ID:      fmt.Sprintf("%d", idx+101),
					Content: strings.TrimPrefix(line, "#CLK "),
				},
			})
		case strings.HasPrefix(line, "#VIEW "):
			viewText := strings.TrimPrefix(line, "#VIEW ")
			openParen := strings.LastIndex(viewText, "(")
			closeParen := strings.LastIndex(viewText, ")")
			if openParen > -1 && closeParen > openParen {
				content := viewText[:openParen]
				url := viewText[openParen+1 : closeParen]
				menuList = append(menuList, wecomclient.MenuItem{
					Type: "view",
					View: &struct {
						URL     string `json:"url"`
						Content string `json:"content"`
					}{
						URL:     url,
						Content: content,
					},
				})
			}
		case strings.HasPrefix(line, "#MINIPROGRAM "):
			mpText := strings.TrimPrefix(line, "#MINIPROGRAM ")
			openParen := strings.LastIndex(mpText, "(")
			closeParen := strings.LastIndex(mpText, ")")
			if openParen > -1 && closeParen > openParen {
				content := mpText[:openParen-1]
				params := mpText[openParen+1 : closeParen]
				paramParts := strings.Split(params, ",")
				if len(paramParts) == 2 {
					appid := strings.TrimSpace(paramParts[0])
					pagepath := strings.TrimSpace(paramParts[1])
					menuList = append(menuList, wecomclient.MenuItem{
						Type: "miniprogram",
						MiniProgram: &struct {
							AppID    string `json:"appid"`
							PagePath string `json:"pagepath"`
							Content  string `json:"content"`
						}{
							AppID:    appid,
							PagePath: pagepath,
							Content:  content,
						},
					})
				}
			}
		}
	}
	if menuList != nil {
		resp.List = menuList
	}
	return resp
}
