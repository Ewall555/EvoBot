package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"

	"go.uber.org/zap"
)

type KFLogic struct {
}

type IKFLogic interface {
	KFAdd(req request.KF) error
	KFUpdate(uuid string, req request.KF) error
	KFDel(uuid string) error
	KFGet(uuid string) (response.KF, error)
	KFList() ([]response.KF, error)
}

func NewIKFLogic() IKFLogic {
	return &KFLogic{}
}

func (u *KFLogic) KFAdd(req request.KF) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	kf := model.KF{
		KFName:           req.KFName,
		KFID:             req.KFID,
		KFPlatform:       req.KFPlatform,
		BotID:            req.BotID,
		Status:           req.Status,
		ReceivePriority:  req.ReceivePriority,
		ReceiveRule:      req.ReceiveRule,
		ChatTimeout:      req.ChatTimeout,
		BotTimeout:       req.BotTimeout,
		BotTimeoutMsg:    req.BotTimeoutMsg,
		BotWelcomeMsg:    req.BotWelcomeMsg,
		StaffWelcomeMsg:  req.StaffWelcomeMsg,
		UnmannedMsg:      req.UnmannedMsg,
		ChatendMsg:       req.ChatendMsg,
		TransferKeywords: req.TransferKeywords,
	}
	staffs, err := staffRepo.List(commonRepo.WithUUIDsIn(req.StaffList))
	if err != nil {
		global.ZAPLOG.Error("Find Staffs", zap.Error(err))
		return err
	}
	kf.Staffs = staffs
	if err := kFRepo.Create(kf); err != nil {
		global.ZAPLOG.Error("Create KF", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

func (u *KFLogic) KFUpdate(uuid string, req request.KF) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	kf, err := kFRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return err
	}
	kf.UUID = uuid
	kf.KFName = req.KFName
	kf.KFID = req.KFID
	kf.KFPlatform = req.KFPlatform
	kf.BotID = req.BotID
	kf.Status = req.Status
	kf.ReceivePriority = req.ReceivePriority
	kf.ReceiveRule = req.ReceiveRule
	kf.ChatTimeout = req.ChatTimeout
	kf.BotTimeout = req.BotTimeout
	kf.BotTimeoutMsg = req.BotTimeoutMsg
	kf.BotWelcomeMsg = req.BotWelcomeMsg
	kf.StaffWelcomeMsg = req.StaffWelcomeMsg
	kf.UnmannedMsg = req.UnmannedMsg
	kf.ChatendMsg = req.ChatendMsg
	kf.TransferKeywords = req.TransferKeywords
	if req.StaffList != nil {
		staffs, err := staffRepo.List(commonRepo.WithUUIDsIn(req.StaffList))
		if err != nil {
			return err
		}
		kf.Staffs = staffs
	}
	if err := kFRepo.Update(kf); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *KFLogic) KFDel(uuid string) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	if err := kFRepo.Delete(commonRepo.WithByUUID(uuid)); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *KFLogic) KFGet(uuid string) (response.KF, error) {
	kf, err := kFRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return response.KF{}, err
	}
	var staffs []response.Staff
	for _, staff := range kf.Staffs {
		var policies []response.Policy
		for _, policy := range staff.Policies {
			var workTimes []dto.WorkTime
			for _, workTime := range policy.WorkTimes {
				workTimes = append(workTimes, dto.WorkTime{
					StartTime: workTime.StartTime,
					EndTime:   workTime.EndTime,
				})
			}
			var policyres response.Policy
			policyres.UUID = policy.UUID
			policyres.PolicyName = policy.PolicyName
			policyres.Repeat = policy.Repeat
			policyres.Week = policy.Week
			policyres.MaxCount = policy.MaxCount
			policyres.WorkTimes = workTimes
			policies = append(policies, policyres)
		}
		var staffres response.Staff
		staffres.UUID = staff.UUID
		staffres.StaffID = staff.StaffID
		staffres.StaffName = staff.StaffName
		staffres.Number = staff.Number
		staffres.Email = staff.Email
		staffres.Policies = policies
		staffs = append(staffs, staffres)
	}
	var kfres response.KF
	kfres.UUID = kf.UUID
	kfres.KFName = kf.KFName
	kfres.KFID = kf.KFID
	kfres.KFPlatform = kf.KFPlatform
	kfres.BotID = kf.BotID
	kfres.Status = kf.Status
	kfres.ReceivePriority = kf.ReceivePriority
	kfres.ReceiveRule = kf.ReceiveRule
	kfres.ChatTimeout = kf.ChatTimeout
	kfres.BotTimeout = kf.BotTimeout
	kfres.BotTimeoutMsg = kf.BotTimeoutMsg
	kfres.BotWelcomeMsg = kf.BotWelcomeMsg
	kfres.StaffWelcomeMsg = kf.StaffWelcomeMsg
	kfres.UnmannedMsg = kf.UnmannedMsg
	kfres.ChatendMsg = kf.ChatendMsg
	kfres.TransferKeywords = kf.TransferKeywords
	kfres.Staffs = staffs
	return kfres, nil
}

func (u *KFLogic) KFList() ([]response.KF, error) {
	kfs, err := kFRepo.List()
	if err != nil {
		return nil, err
	}
	var kflist []response.KF
	for _, kf := range kfs {
		var staffs []response.Staff
		for _, staff := range kf.Staffs {
			var policies []response.Policy
			for _, policy := range staff.Policies {
				var workTimes []dto.WorkTime
				for _, workTime := range policy.WorkTimes {
					workTimes = append(workTimes, dto.WorkTime{
						StartTime: workTime.StartTime,
						EndTime:   workTime.EndTime,
					})
				}
				var policyres response.Policy
				policyres.UUID = policy.UUID
				policyres.PolicyName = policy.PolicyName
				policyres.Repeat = policy.Repeat
				policyres.Week = policy.Week
				policyres.MaxCount = policy.MaxCount
				policyres.WorkTimes = workTimes
				policies = append(policies, policyres)
			}
			var staffres response.Staff
			staffres.UUID = staff.UUID
			staffres.StaffID = staff.StaffID
			staffres.StaffName = staff.StaffName
			staffres.Number = staff.Number
			staffres.Email = staff.Email
			staffres.Policies = policies
			staffs = append(staffs, staffres)
		}
		var kfres response.KF
		kfres.UUID = kf.UUID
		kfres.KFName = kf.KFName
		kfres.KFID = kf.KFID
		kfres.KFPlatform = kf.KFPlatform
		kfres.BotID = kf.BotID
		kfres.Status = kf.Status
		kfres.ReceivePriority = kf.ReceivePriority
		kfres.ReceiveRule = kf.ReceiveRule
		kfres.ChatTimeout = kf.ChatTimeout
		kfres.BotTimeout = kf.BotTimeout
		kfres.BotTimeoutMsg = kf.BotTimeoutMsg
		kfres.BotWelcomeMsg = kf.BotWelcomeMsg
		kfres.StaffWelcomeMsg = kf.StaffWelcomeMsg
		kfres.UnmannedMsg = kf.UnmannedMsg
		kfres.ChatendMsg = kf.ChatendMsg
		kfres.TransferKeywords = kf.TransferKeywords
		kfres.Staffs = staffs
		kflist = append(kflist, kfres)
	}
	return kflist, nil
}
