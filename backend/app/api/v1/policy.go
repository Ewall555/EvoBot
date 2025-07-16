package v1

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/app/dto"
	"EvoBot/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags policy
// @Summary Add policy
// @Description Add a new policy
// @Accept json
// @Produce json
// @Param req body dto.Policy true "policy Information"
// @Success 200 {object} dto.Response
// @Router /policy [post]
func (u *BaseApi) PolicyAdd(ctx *gin.Context) {
	var req dto.Policy
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if req.PolicyName == "" || req.MaxCount == 0 || req.Repeat == 0 || len(req.WorkTimes) == 0 {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	validRepeats := []int{1, 2, 3, 4, 5}
	isValid := false
	for _, valid := range validRepeats {
		if req.Repeat == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if req.Repeat == 1 && req.Week == "" {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if req.Week != "" {
		if len(req.Week) != 7 {
			helper.ErrResponse(ctx, constant.CodeErrBadRequest)
			return
		}
		for _, char := range req.Week {
			if char != '0' && char != '1' {
				helper.ErrResponse(ctx, constant.CodeErrBadRequest)
				return
			}
		}
	}
	if err := policyLogic.PolicyAdd(req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags policy
// @Summary Update policy
// @Description Update the policy by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "policy UUID"
// @Param req body dto.Policy true "policy Information"
// @Success 200 {object} dto.Response
// @Router /policy/{uuid} [patch]
func (u *BaseApi) PolicyUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	var req dto.Policy
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if req.Repeat != 0 {
		validRepeats := []int{1, 2, 3, 4, 5}
		isValid := false
		for _, valid := range validRepeats {
			if req.Repeat == valid {
				isValid = true
				break
			}
		}
		if !isValid {
			helper.ErrResponse(ctx, constant.CodeErrBadRequest)
			return
		}
		if req.Repeat == 1 && req.Week == "" {
			helper.ErrResponse(ctx, constant.CodeErrBadRequest)
			return
		}
		if req.Week != "" {
			if len(req.Week) != 7 {
				helper.ErrResponse(ctx, constant.CodeErrBadRequest)
				return
			}
			for _, char := range req.Week {
				if char != '0' && char != '1' {
					helper.ErrResponse(ctx, constant.CodeErrBadRequest)
					return
				}
			}
		}
	}
	if err := policyLogic.PolicyUpdate(uuid, req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags policy
// @Summary Delete policy
// @Description Delete the policy by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "policy UUID"
// @Success 200 {object} dto.Response
// @Router /policy/{uuid} [delete]
func (u *BaseApi) PolicyDel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if err := policyLogic.PolicyDel(uuid); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags policy
// @Summary Get policy
// @Description Get the policy information by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "policy UUID"
// @Success 200 {object} dto.Response
// @Router /policy/{uuid} [get]
func (u *BaseApi) PolicyGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	strategy, err := policyLogic.PolicyGet(uuid)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, strategy)
}

// PolicyList retrieves all policy
// @Tags policy
// @Summary List policy
// @Description Retrieve a list of all policy
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response
// @Router /policy [get]
func (u *BaseApi) PolicyList(ctx *gin.Context) {
	strategies, err := policyLogic.PolicyList()
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, strategies)
}
