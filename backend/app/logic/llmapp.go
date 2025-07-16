package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
	"EvoBot/backend/i18n"
	"EvoBot/backend/utils/llmapp"

	"go.uber.org/zap"
)

type LLMAppLogic struct {
	llmAppClients map[string]llmapp.LLMAppClient
}

type ILLMAppLogic interface {
	ConfigAdd(req dto.LLMAppConfig) error
	ConfigUpdate(uuid string, req dto.LLMAppConfig) error
	ConfigDel(uuid string) error
	ConfigGet(uuid string) (response.LLMAppConfig, error)
	ConfigList() ([]response.LLMAppConfig, error)

	ChatMessage(khid, uuid, message string) (string, error)
}

func NewILLMAppLogic() ILLMAppLogic {
	return &LLMAppLogic{
		llmAppClients: make(map[string]llmapp.LLMAppClient),
	}
}

func (u *LLMAppLogic) ChatMessage(khid, uuid, message string) (string, error) {
	var client llmapp.LLMAppClient
	khinfo, err := kHRepo.Get(kHRepo.WithKHID(khid))
	if err != nil {
		global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "kHRepo.Get"), zap.Error(err))
		return "", err
	}

	var targetChatID *string
	for _, chat := range khinfo.ChatList {
		if chat.BotID == uuid {
			chatIDCopy := chat.ChatID
			targetChatID = &chatIDCopy
			break
		}
	}

	if targetChatID == nil {
		global.ZAPLOG.Info("该客户没有和该机器人的聊天ID", zap.String("KHID", khid), zap.String("BotID", uuid))

		client, err = u.getClient(uuid)
		if err != nil {
			global.ZAPLOG.Error("获取 LLMApp 客户端失败", zap.Error(err))
			return "", err
		}

		newchatid, err := client.ChatOpen()
		if err != nil {
			return "", err
		}

		chatList := model.ChatList{
			KHID:   khinfo.ID,
			BotID:  uuid,
			ChatID: *newchatid,
		}

		if err := kHRepo.CreateChatList(chatList); err != nil {
			return "", err
		}
		if err := kHRepo.UpdatebyKHID(khinfo); err != nil {
			return "", err
		}

		global.ZAPLOG.Info("生成聊天id更新客户信息:",
			zap.String("newchatid", *newchatid),
			zap.String("khid", khid),
			zap.String("botid", uuid))
		targetChatID = newchatid
	} else {
		client, err = u.getClient(uuid)
		if err != nil {
			global.ZAPLOG.Error("获取 LLMApp 客户端失败", zap.Error(err))
			return "", err
		}
	}

	fullContent, err := client.ChatMessage(message, targetChatID)
	if err != nil {
		return "", err
	}
	return fullContent, nil
}

func (u *LLMAppLogic) getClient(uuid string) (llmapp.LLMAppClient, error) {
	if u.llmAppClients == nil {
		u.llmAppClients = make(map[string]llmapp.LLMAppClient)
	}

	if client, ok := u.llmAppClients[uuid]; ok {
		return client, nil
	}

	conf, err := llmappRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return nil, err
	}

	varMap := make(map[string]interface{})
	varMap["base_url"] = conf.BaseURL
	varMap["api_key"] = conf.ApiKey

	client, err := llmapp.NewLLMAppClient(conf.LLMAppType, varMap)
	if err != nil {
		global.ZAPLOG.Error("创建 LLMApp 客户端失败", zap.Error(err))
		return nil, err
	}

	u.llmAppClients[uuid] = client
	return client, nil
}

func (u *LLMAppLogic) ConfigAdd(req dto.LLMAppConfig) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	if err := llmappRepo.Create(model.LLMAppConfig{
		LLMAppType: req.LLMAppType,
		ApiKey:     req.ApiKey,
		BaseURL:    req.BaseURL,
		ConfigName: req.ConfigName,
	}); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *LLMAppLogic) ConfigDel(uuid string) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()

	delete(u.llmAppClients, uuid)

	if err := llmappRepo.Delete(commonRepo.WithByUUID(uuid)); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *LLMAppLogic) ConfigUpdate(uuid string, req dto.LLMAppConfig) error {
	conf, err := llmappRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return err
	}
	conf.UUID = uuid
	conf.LLMAppType = req.LLMAppType
	conf.ConfigName = req.ConfigName
	conf.ApiKey = req.ApiKey
	conf.BaseURL = req.BaseURL
	if err := llmappRepo.Update(conf); err != nil {
		return err
	}

	varMap := make(map[string]interface{})
	varMap["base_url"] = req.BaseURL
	varMap["api_key"] = req.ApiKey
	client, err := llmapp.NewLLMAppClient(req.LLMAppType, varMap)
	if err != nil {
		global.ZAPLOG.Error("更新 LLMApp 客户端失败", zap.Error(err))
		return err
	}
	u.llmAppClients[uuid] = client
	kHRepo.DeleteChatListWithBotID(uuid)
	return nil
}

func (u *LLMAppLogic) ConfigGet(uuid string) (response.LLMAppConfig, error) {
	conf, err := llmappRepo.Get(commonRepo.WithByUUID(uuid))
	var resp response.LLMAppConfig
	if err != nil {
		return resp, err
	}
	resp.UUID = conf.UUID
	resp.LLMAppType = conf.LLMAppType
	resp.ConfigName = conf.ConfigName
	resp.BaseURL = conf.BaseURL
	resp.ApiKey = conf.ApiKey
	return resp, nil
}

func (u *LLMAppLogic) ConfigList() ([]response.LLMAppConfig, error) {
	configs, err := llmappRepo.List()
	var resp []response.LLMAppConfig
	if err != nil {
		global.ZAPLOG.Error("failed to get LLMApp config", zap.Error(err))
		return resp, err
	}
	for _, config := range configs {
		var res response.LLMAppConfig
		res.UUID = config.UUID
		res.LLMAppType = config.LLMAppType
		res.ConfigName = config.ConfigName
		res.BaseURL = config.BaseURL
		res.ApiKey = config.ApiKey
		resp = append(resp, res)
	}

	return resp, nil
}
