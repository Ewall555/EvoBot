package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/app/model"
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"EvoBot/backend/utils/wecom"
	wecomclient "EvoBot/backend/utils/wecom/client"
	"fmt"
	"strings"

	"context"
	"time"

	"go.uber.org/zap"
)

type WecomLogic struct {
	wecomkf              wecom.WecomKFClient
	wecomexternalcontact wecom.WecomExternalContactClient
}

type IWecomLogic interface {
	VerifyURL(req dto.SignatureOptions) (string, error)
	Handle(body []byte) error

	ConfigList() ([]response.WecomConfigApp, error)
	ConfigUpdate(uuid string, req request.WecomConfigApp) error
	ConfigGet(uuid string) (response.WecomConfigApp, error)

	ReceptionistAdd(kfid string, req request.ReceptionistOptions) error
	ReceptionistDel(kfid string, req request.ReceptionistOptions) error
	ReceptionistList(kfid string) ([]wecomclient.ReceptionistList, error)
	CheckReceptionist(stafflist []string, receptionistlist []wecomclient.ReceptionistList) error

	AccountList() ([]wecomclient.AccountInfoSchema, error)
	AddContactWay(kfid string) (string, error)

	AddGroupMsgTemplate(req wecomclient.AddMsgTemplateRequest) (string, []string, error)
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
	conf.CorpID = req.CorpID
	conf.Secret = req.Secret
	conf.AgentID = req.AgentID
	conf.Token = req.Token
	conf.EncodingAESKey = req.EncodingAESKey
	err = wecomRepo.Update(context.Background(), &conf)
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

func (u *WecomLogic) VerifyURL(req dto.SignatureOptions) (string, error) {
	if err := u.initializeWecomClient(); err != nil {
		return "", err
	}
	echo, err := u.wecomkf.VerifyURL(wecomclient.SignatureOptions{
		Signature: req.Signature,
		TimeStamp: req.TimeStamp,
		Nonce:     req.Nonce,
		EchoStr:   req.EchoStr,
	})
	if err != nil {
		global.ZAPLOG.Error("failed to check url", zap.Error(err))
	}
	return echo, nil
}

func (u *WecomLogic) Handle(body []byte) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	streamCallback := func(msginfo *wecomclient.MessageInfo) {
		if msginfo.MessageID != "" || msginfo.SendTime != 0 || msginfo.Origin != 0 || msginfo.KFID != "" || msginfo.KHID != "" || msginfo.StaffID != "" || msginfo.Message != "" || msginfo.MessageType != "" || msginfo.Credential != "" || msginfo.ChatState != 0 {
			global.ZAPLOG.Debug("MessageInfo 内容:",
				zap.Any("MessageID", msginfo.MessageID),
				zap.Any("SendTime", msginfo.SendTime),
				zap.Any("Origin", msginfo.Origin),
				zap.Any("KFID", msginfo.KFID),
				zap.Any("KHID", msginfo.KHID),
				zap.Any("StaffID", msginfo.StaffID),
				zap.Any("Message", msginfo.Message),
				zap.Any("MessageType", msginfo.MessageType),
				zap.Any("Credential", msginfo.Credential),
				zap.Any("ChatState", msginfo.ChatState),
			)
		}
		u.processMessage(*msginfo)
	}
	_, err := u.wecomkf.SyncMsg(body, streamCallback)
	if err != nil {
		global.ZAPLOG.Error("failed SyncMsg", zap.Error(err))
		return err
	}

	return nil
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

func (u *WecomLogic) AddGroupMsgTemplate(req wecomclient.AddMsgTemplateRequest) (string, []string, error) {
	if err := u.initializeWecomClient(); err != nil {
		return "", nil, err
	}
	msgID, failList, err := u.wecomexternalcontact.AddGroupMsgTemplate(wecomclient.AddMsgTemplateRequest{
		Sender: req.Sender,
		Text: wecomclient.MsgText{
			Content: req.Text.Content,
		},
		AllowSelect: req.AllowSelect,
		ChatIDList:  req.ChatIDList,
	})
	if err != nil {
		return "", nil, err
	}
	return msgID, failList, nil
}

// Helper
func (u *WecomLogic) initializeWecomClient() error {
	if u.wecomkf != nil {
		return nil
	}
	conf, err := wecomRepo.Get(wecomRepo.WithType("app"))
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return err
	}
	if conf.CorpID == "" || conf.AgentID == "" || conf.EncodingAESKey == "" || conf.Secret == "" || conf.Token == "" {
		global.ZAPLOG.Error("wecom config is nil", zap.Error(err))
		return fmt.Errorf("incomplete Wecom config")
	}
	config := wecomclient.WecomConfig{
		CorpID:         conf.CorpID,
		Token:          conf.Token,
		EncodingAESKey: conf.EncodingAESKey,
		Secret:         conf.Secret,
		AgentID:        conf.AgentID,
	}
	u.wecomkf, err = wecom.NewWecomKFClient(config)
	if err != nil {
		global.ZAPLOG.Error("failed to initialize WecomKFClient", zap.Error(err))
		return err
	}
	u.wecomexternalcontact, err = wecom.NewWecomExternalContactClient(config)
	if err != nil {
		global.ZAPLOG.Error("failed to initialize WecomExternalContactClient", zap.Error(err))
		return err
	}
	return nil
}

func (u *WecomLogic) processMessage(msginfo wecomclient.MessageInfo) error {
	kfinfo, err := kFRepo.Get(kFRepo.WithKFID(msginfo.KFID))
	if err != nil {
		return err
	}
	switch msginfo.MessageType {
	case wecomclient.WecomEventTypeEnterSession:
		if msginfo.Credential != "" {
			return u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
				Credential:     msginfo.Credential,
				MenuMsgOptions: parseMenuText(kfinfo.BotWelcomeMsg),
			})
		}
		if msginfo.ChatState == wecomclient.SessionStatusHandled {
			return u.wecomkf.SendMenuMsg(wecomclient.SendMenuMsgOptions{
				KFID:           msginfo.KFID,
				KHID:           msginfo.KHID,
				MenuMsgOptions: parseMenuText(kfinfo.BotWelcomeMsg),
			})
		}
	case wecomclient.WecomEventTypeSessionStatusChange:
		switch msginfo.Message {
		case wecomclient.WecomEventChangeTypeJoinSession:
		case wecomclient.WecomEventChangeTypeTransferSession:
			khinfo, err := kHRepo.Get(kHRepo.WithKHID(msginfo.KHID))
			if err != nil {
				global.ZAPLOG.Info("数据库没有找到该客户信息")
				if err := kHRepo.Create(model.KH{KHID: msginfo.KHID}); err != nil {
					return err
				}
				global.ZAPLOG.Info("录入客户信息:", zap.String("khid", msginfo.KHID))
			}
			chatkey := constant.KeyWecomKHStaffPrefix + msginfo.KHID + ":" + khinfo.StaffID
			if err := global.RDS.Del(context.Background(), chatkey).Err(); err != nil {
				global.ZAPLOG.Error("redis del error", zap.Error(err))
				return err
			}
			isStaffWorkByStaffID(msginfo.StaffID)
			return u.handleSuccessfulTransfer(msginfo, msginfo.StaffID, kfinfo)
		case wecomclient.WecomEventChangeTypeEndSession:
			if err := u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
				Credential:     msginfo.Credential,
				MenuMsgOptions: parseMenuText(kfinfo.ChatendMsg),
			}); err != nil {
				global.ZAPLOG.Error("SendMenuMsgOnEvent", zap.Error(err))
			}
			chatkey := constant.KeyWecomKHStaffPrefix + msginfo.KHID + ":" + msginfo.StaffID
			if err := global.RDS.Del(context.Background(), chatkey).Err(); err != nil {
				global.ZAPLOG.Error("redis del error", zap.Error(err))
				return err
			}
			global.ZAPLOG.Info("结束会话监控", zap.String("chatkey:", chatkey))
		case wecomclient.WecomEventChangeTypeRejoinSession:
			isStaffWorkByStaffID(msginfo.StaffID)
			return u.handleSuccessfulTransfer(msginfo, msginfo.StaffID, kfinfo)
		default:
			global.ZAPLOG.Info("未实现的消息类型")
		}
	case wecomclient.WecomMsgTypeText:
		switch msginfo.ChatState {
		case wecomclient.SessionStatusHandled:
			return u.handleBotMessage(msginfo, kfinfo)
		case wecomclient.SessionStatusWaiting:
			global.ZAPLOG.Info("未实现的聊天状态处理")
		case wecomclient.SessionStatusInProgress:
			chatkey := constant.KeyWecomKHStaffPrefix + msginfo.KHID + ":" + msginfo.StaffID
			if err = global.RDS.Set(context.Background(), chatkey, 1, time.Duration(kfinfo.ChatTimeout)*time.Second).Err(); err != nil {
				global.ZAPLOG.Error("redis set error", zap.Error(err))
				return err
			}
		case wecomclient.SessionStatusEndedOrNotStarted:
			global.ZAPLOG.Info("未实现的聊天状态处理")
		default:
			global.ZAPLOG.Info("未实现的聊天状态处理")
		}
	default:
		global.ZAPLOG.Info("未实现的消息类型")
	}
	return nil
}

func (u *WecomLogic) handleSuccessfulTransfer(msginfo wecomclient.MessageInfo, staffid string, kfinfo model.KF) error {
	if msginfo.ChatState != wecomclient.SessionStatusInProgress {
		global.ZAPLOG.Info("变更微信客服会话状态由选出的接待人员接待")
		staffcredential, err := u.wecomkf.ServiceStateTrans(wecomclient.ServiceStateTransOptions{
			OpenKFID:       msginfo.KFID,
			ExternalUserID: msginfo.KHID,
			ServicerUserID: staffid,
		}, wecomclient.SessionStatusInProgress)
		if err != nil {
			global.ZAPLOG.Error("变更微信客服会话状态失败", zap.Error(err))
			currentTime := time.Now()
			formattedTime := currentTime.Format("2006-01-02 15:04:05")
			errormessage := "未知原因出现异常请联系工作人员或可尝试继续与 AI 对话，当前时间：" + formattedTime
			return u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
				KFID:    msginfo.KFID,
				KHID:    msginfo.KHID,
				Message: errormessage,
			})
		}
		if err = u.wecomkf.SendTextMsgOnEvent(wecomclient.SendTextMsgOnEventOptions{
			Message:    kfinfo.StaffWelcomeMsg,
			Credential: staffcredential,
		}); err != nil {
			return err
		}
	}
	global.ZAPLOG.Info("更新客户的上一次接待人员为选出的接待人员")
	if err := kHRepo.UpdatebyKHID(model.KH{
		KHID:    msginfo.KHID,
		StaffID: staffid,
	}); err != nil {
		global.ZAPLOG.Error("更新客户的上一次接待人员失败", zap.Error(err))
	}
	global.ZAPLOG.Info("降低该接待人员空闲权重")
	ctx := context.Background()
	staffweightkey := constant.KeyStaffWeightPrefix + staffid
	if err := global.RDS.Decr(ctx, staffweightkey).Err(); err != nil {
		return err
	}
	global.ZAPLOG.Info("缓存会话超时时间")
	chatkey := constant.KeyWecomKHStaffPrefix + msginfo.KHID + ":" + staffid
	if err := global.RDS.Set(ctx, chatkey, 1, time.Duration(kfinfo.ChatTimeout)*time.Second).Err(); err != nil {
		global.ZAPLOG.Error("redis set error", zap.Error(err))
		return err
	}
	global.ZAPLOG.Info("开始监控会话任务")
	go u.monitorChat(ctx, kfinfo, msginfo, staffweightkey, chatkey)
	return nil
}

func (u *WecomLogic) monitorChat(ctx context.Context, kfinfo model.KF, msginfo wecomclient.MessageInfo, staffweightkey, chatkey string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !u.hasKHReplied(ctx, chatkey) {
				u.handleChatTimeout(ctx, kfinfo, msginfo)
				if err := global.RDS.Incr(ctx, staffweightkey).Err(); err != nil {
					global.ZAPLOG.Error("Incr", zap.Error(err))
				} else {
					global.ZAPLOG.Info("结束会话变更权重", zap.String("staffweightkey:", staffweightkey))
				}
				return
			}
		case <-ctx.Done():
			global.ZAPLOG.Info("监控已停止")
			return
		}
	}
}

func (u *WecomLogic) handleChatTimeout(ctx context.Context, kfinfo model.KF, msginfo wecomclient.MessageInfo) {
	statusinfo, err := u.wecomkf.ServiceStateGet(wecomclient.ServiceStateGetOptions{
		OpenKFID:       kfinfo.KFID,
		ExternalUserID: msginfo.KHID,
	})
	if err != nil {
		global.ZAPLOG.Error("获取会话失败", zap.Error(err))
		return
	}
	if statusinfo == wecomclient.SessionStatusEndedOrNotStarted {
		return
	}
	var cursor uint64
	chatkey := constant.KeyWecomKHStaffPrefix + msginfo.KHID
	result, _, err := global.RDS.Scan(ctx, cursor, chatkey+"*", 10).Result()
	if err != nil {
		return
	}
	if len(result) > 0 {
		return
	}
	endcredential, err := u.wecomkf.ServiceStateTrans(wecomclient.ServiceStateTransOptions{
		OpenKFID:       kfinfo.KFID,
		ExternalUserID: msginfo.KHID,
	}, wecomclient.SessionStatusEndedOrNotStarted)
	if err != nil {
		global.ZAPLOG.Error("结束会话失败", zap.Error(err))
		return
	}
	global.ZAPLOG.Info("会话超时，已变更状态")
	if err := u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
		Credential:     endcredential,
		MenuMsgOptions: parseMenuText(kfinfo.ChatendMsg),
	}); err != nil {
		global.ZAPLOG.Error("SendTextMsgOnEvent", zap.Error(err))
	}
}

func (u *WecomLogic) hasKHReplied(ctx context.Context, chatkey string) bool {
	exists, err := global.RDS.Exists(ctx, chatkey).Result()
	if err != nil {
		global.ZAPLOG.Error("检查客户回复缓存失败", zap.Error(err))
		return false
	}
	if exists > 0 {
		return true
	}
	return false
}

func (u *WecomLogic) handleTransferToStaff(msginfo wecomclient.MessageInfo, kfinfo model.KF) error {
	if kfinfo.ReceivePriority == 1 {
		khinfo, err := kHRepo.Get(kHRepo.WithKHID(msginfo.KHID))
		if err != nil {
			global.ZAPLOG.Info("数据库没有找到该客户信息")
			if err := kHRepo.Create(model.KH{KHID: msginfo.KHID}); err != nil {
				return err
			}
			global.ZAPLOG.Info("录入客户信息:", zap.String("khid", msginfo.KHID))
		}
		if khinfo.StaffID != "" {
			for _, staffinfo := range kfinfo.Staffs {
				if staffinfo.StaffID == khinfo.StaffID {
					isStaffWork, err := isStaffWorkByStaffID(khinfo.StaffID)
					if err != nil {
						return err
					}
					if isStaffWork {
						global.ZAPLOG.Info("选出的接待人员", zap.String("StaffID", khinfo.StaffID))
						return u.handleSuccessfulTransfer(msginfo, khinfo.StaffID, kfinfo)
					}
				}
			}
			global.ZAPLOG.Info("该接待人员不在客服应用中", zap.String("KFName", kfinfo.KFName))
		}
		global.ZAPLOG.Info("该客户无法应用上一次接待人员", zap.String("KHID", msginfo.KHID))
	}
	switch kfinfo.ReceiveRule {
	case constant.KFReceiveRuleRoundRobin:
		global.ZAPLOG.Info("轮流接待模式", zap.String("KFName", kfinfo.KFName))
		global.ZAPLOG.Info("未实现功能")
	case constant.KFReceiveRuleIdle:
		global.ZAPLOG.Info("空闲接待模式", zap.String("KFName", kfinfo.KFName))
		var staffIDs []string
		for _, staffinfo := range kfinfo.Staffs {
			isStaffWork, err := isStaffWorkByStaffID(staffinfo.StaffID)
			if err != nil {
				return err
			}
			if isStaffWork {
				staffIDs = append(staffIDs, staffinfo.StaffID)
			}
		}
		selectedstaff, staffweightvalue := getHighestWeightStaff(staffIDs)
		if selectedstaff != "" {
			global.ZAPLOG.Info("选出的接待人员", zap.String("StaffID", selectedstaff), zap.Int("当前权重", staffweightvalue))
			return u.handleSuccessfulTransfer(msginfo, selectedstaff, kfinfo)
		}
		global.ZAPLOG.Info("没有找到空闲的接待人员")
		return u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
			KFID:    msginfo.KFID,
			KHID:    msginfo.KHID,
			Message: kfinfo.UnmannedMsg,
		})
	default:
		global.ZAPLOG.Info("默认接待模式未配置", zap.String("KFName", kfinfo.KFName))
		return nil
	}
	return nil
}

func (u *WecomLogic) handleBotMessage(msginfo wecomclient.MessageInfo, kfinfo model.KF) error {
	switch kfinfo.Status {
	case constant.KFStatusRobotToHuman:
		keywords := strings.Split(kfinfo.TransferKeywords, ";")
		for _, keyword := range keywords {
			if msginfo.Message == keyword {
				global.ZAPLOG.Info("客户触发客服转人工关键字", zap.String("keyword", keyword))
				return u.handleTransferToStaff(msginfo, kfinfo)
			}
		}
		return u.handleBotReply(msginfo, kfinfo)
	case constant.KFStatusOnlyRobot:
		global.ZAPLOG.Info("客服为仅机器人模式，无法转接人工", zap.String("KFName", kfinfo.KFName))
		return u.handleBotReply(msginfo, kfinfo)
	case constant.KFStatusOnlyHuman:
		global.ZAPLOG.Info("客服为仅人工模式，机器人无法回复消息", zap.String("KFName", kfinfo.KFName))
		return nil
	default:
		global.ZAPLOG.Error("该客服模式不存在")
		return nil
	}
}

func (u *WecomLogic) handleBotReply(msginfo wecomclient.MessageInfo, kfinfo model.KF) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(kfinfo.BotTimeout)*time.Second)
	defer cancel()
	resultChan := make(chan string)
	errorChan := make(chan error)
	go func() {
		message := msginfo.Message + kfinfo.BotPrompt
		llmappLogic := NewILLMAppLogic()
		fullContent, err := llmappLogic.ChatMessage(msginfo.KHID, kfinfo.BotID, message)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- fullContent
		}
	}()
	select {
	case <-ctx.Done():
		err := u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
			KFID:    msginfo.KFID,
			KHID:    msginfo.KHID,
			Message: kfinfo.BotTimeoutMsg,
		})
		if err != nil {
			global.ZAPLOG.Error("failed to send wait message", zap.Error(err))
		}
		global.ZAPLOG.Info("bot超时消息:", zap.String("content", kfinfo.BotTimeoutMsg))
		select {
		case fullContent := <-resultChan:
			err = u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
				KFID:    msginfo.KFID,
				KHID:    msginfo.KHID,
				Message: MarkdownToText(fullContent),
			})
			if err != nil {
				return err
			}
			global.ZAPLOG.Info("bot消息:", zap.String("fullContent", fullContent))
			return nil
		case err := <-errorChan:
			global.ZAPLOG.Error("failed ChatMessage after timeout", zap.Error(err))
			return err
		}
	case fullContent := <-resultChan:
		err := u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
			KFID:    msginfo.KFID,
			KHID:    msginfo.KHID,
			Message: MarkdownToText(fullContent),
		})
		if err != nil {
			return err
		}
		global.ZAPLOG.Info("bot原消息:", zap.String("fullContent", fullContent))
		return nil
	case err := <-errorChan:
		global.ZAPLOG.Error("failed ChatMessage", zap.Error(err))
		return err
	}
}

func parseMenuText(text string) wecomclient.MenuMsgOptions {
	lines := strings.Split(text, "\n")
	headContent := ""
	tailContent := ""
	menuList := []wecomclient.MenuItem{}
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "#H "):
			headContent = strings.TrimPrefix(line, "#H ")
		case strings.HasPrefix(line, "#T "):
			tailContent = strings.TrimPrefix(line, "#T ")
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
	var resp wecomclient.MenuMsgOptions
	if headContent != "" {
		resp.HeadContent = headContent
	}
	if menuList != nil {
		resp.MenuList = menuList
	}
	if tailContent != "" {
		resp.TailContent = tailContent
	}
	return resp
}
