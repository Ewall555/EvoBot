package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"

	"go.uber.org/zap"
)

type StaffLogic struct {
}

type IStaffLogic interface {
	StaffAdd(req request.Staff) error
	StaffUpdate(uuid string, req request.Staff) error
	StaffDel(uuid string) error
	StaffGet(uuid string) (response.Staff, error)
	StaffList() ([]response.Staff, error)
}

func NewIStaffLogic() IStaffLogic {
	return &StaffLogic{}
}

func (u *StaffLogic) StaffAdd(req request.Staff) error {
	tx := global.DB.Begin()
	staff := model.Staff{
		StaffID:   req.StaffID,
		StaffName: req.StaffName,
		Number:    req.Number,
		Email:     req.Email,
		Role:      req.Role,
	}
	policies, err := policyRepo.List(commonRepo.WithUUIDsIn(req.PolicyList))
	if err != nil {
		global.ZAPLOG.Error("Find Policies", zap.Error(err))
		tx.Rollback()
		return err
	}
	staff.Policies = policies
	if err := staffRepo.Create(staff); err != nil {
		global.ZAPLOG.Error("Create Staff", zap.Error(err))
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *StaffLogic) StaffUpdate(uuid string, req request.Staff) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	staff, err := staffRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return err
	}
	staff.UUID = uuid
	staff.StaffID = req.StaffID
	staff.StaffName = req.StaffName
	staff.Number = req.Number
	staff.Email = req.Email
	staff.Role = req.Role
	if req.PolicyList != nil {
		policies, err := policyRepo.List(commonRepo.WithUUIDsIn(req.PolicyList))
		if err != nil {
			return err
		}
		staff.Policies = policies
	}
	if err := staffRepo.Update(staff); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *StaffLogic) StaffDel(uuid string) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := staffRepo.Delete(commonRepo.WithByUUID(uuid)); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *StaffLogic) StaffGet(uuid string) (response.Staff, error) {
	var resp response.Staff
	staff, err := staffRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		global.ZAPLOG.Error("faile to get staff withbyid", zap.Error(err))
		return resp, err
	}
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
	resp.UUID = staff.UUID
	resp.StaffID = staff.StaffID
	resp.StaffName = staff.StaffName
	resp.Number = staff.Number
	resp.Email = staff.Email
	resp.Role = staff.Role
	resp.Policies = policies
	return resp, nil
}

func (u *StaffLogic) StaffList() ([]response.Staff, error) {
	staffs, err := staffRepo.List()
	if err != nil {
		return nil, err
	}
	var resp []response.Staff
	for _, staff := range staffs {
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
		staffres.Role = staff.Role
		staffres.Policies = policies
		resp = append(resp, staffres)
	}
	return resp, nil
}
