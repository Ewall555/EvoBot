package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
	"time"

	"go.uber.org/zap"
)

type PolicyLogic struct {
}

type IPolicyLogic interface {
	PolicyAdd(req dto.Policy) error
	PolicyUpdate(uuid string, req dto.Policy) error
	PolicyDel(uuid string) error
	PolicyGet(uuid string) (response.Policy, error)
	PolicyList() ([]response.Policy, error)
}

func NewIPolicyLogic() IPolicyLogic {
	return &PolicyLogic{}
}

func (u *PolicyLogic) PolicyAdd(req dto.Policy) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	policy := model.Policy{
		PolicyName: req.PolicyName,
		Repeat:     req.Repeat,
		Week:       req.Week,
		MaxCount:   req.MaxCount,
	}
	for _, worktime := range req.WorkTimes {
		if _, err := time.Parse("15:04:05", worktime.StartTime); err != nil {
			global.ZAPLOG.Error("time.Parse StartTime", zap.Error(err))
			return err
		}
		if _, err := time.Parse("15:04:05", worktime.EndTime); err != nil {
			global.ZAPLOG.Error("time.Parse EndTime", zap.Error(err))
			return err
		}
		workTime := model.WorkTime{
			StartTime: worktime.StartTime,
			EndTime:   worktime.EndTime,
		}
		policy.WorkTimes = append(policy.WorkTimes, workTime)
	}
	if err := policyRepo.Create(policy); err != nil {
		global.ZAPLOG.Error("policyRepo createPolicy failed", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

func (u *PolicyLogic) PolicyUpdate(uuid string, req dto.Policy) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	policy, err := policyRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return err
	}
	policy.UUID = uuid
	policy.PolicyName = req.PolicyName
	policy.Repeat = req.Repeat
	policy.Week = req.Week
	policy.MaxCount = req.MaxCount
	if len(req.WorkTimes) != 0 {
		if err := policyRepo.DeleteWorkTime(policy.ID); err != nil {
			return err
		}
		for _, worktime := range req.WorkTimes {
			workTime := model.WorkTime{
				PolicyID:  policy.ID,
				StartTime: worktime.StartTime,
				EndTime:   worktime.EndTime,
			}
			if err := policyRepo.CreateWorkTime(workTime); err != nil {
				return err
			}
		}
	}
	if err := policyRepo.Update(policy); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *PolicyLogic) PolicyDel(uuid string) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	if err := policyRepo.Delete(commonRepo.WithByUUID(uuid)); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *PolicyLogic) PolicyGet(uuid string) (response.Policy, error) {
	policy, err := policyRepo.Get(commonRepo.WithByUUID(uuid))
	var resp response.Policy
	if err != nil {
		return resp, err
	}
	var workTimes []dto.WorkTime
	for _, workTime := range policy.WorkTimes {
		workTimes = append(workTimes, dto.WorkTime{
			StartTime: workTime.StartTime,
			EndTime:   workTime.EndTime,
		})
	}
	resp.UUID = policy.UUID
	resp.PolicyName = policy.PolicyName
	resp.Repeat = policy.Repeat
	resp.Week = policy.Week
	resp.MaxCount = policy.MaxCount
	resp.WorkTimes = workTimes
	return resp, nil
}

func (u *PolicyLogic) PolicyList() ([]response.Policy, error) {
	policies, _ := policyRepo.List()
	var resp []response.Policy
	for _, policy := range policies {
		var workTimes []dto.WorkTime
		for _, workTime := range policy.WorkTimes {
			workTimes = append(workTimes, dto.WorkTime{
				StartTime: workTime.StartTime,
				EndTime:   workTime.EndTime,
			})
		}
		var res response.Policy
		res.UUID = policy.UUID
		res.PolicyName = policy.PolicyName
		res.Week = policy.Week
		res.MaxCount = policy.MaxCount
		res.Repeat = policy.Repeat
		res.WorkTimes = workTimes
		resp = append(resp, res)
	}
	return resp, nil
}
