package v1

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags staff
// @Summary Add Staff
// @Description Add a new staff member
// @Accept json
// @Produce json
// @Param req body request.Staff true "Staff Information"
// @Success 200 {object} dto.Response
// @Router /staff [post]
func (u *BaseApi) StaffAdd(ctx *gin.Context) {
	var req request.Staff
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if req.StaffID == "" || req.StaffName == "" {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if req.PolicyList == nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	err := staffLogic.StaffAdd(req)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags staff
// @Summary Update Staff
// @Description Update the staff member by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "Staff UUID"
// @Param req body request.Staff true "Staff Information"
// @Success 200 {object} dto.Response
// @Router /staff/{uuid} [patch]
func (u *BaseApi) StaffUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	var req request.Staff
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := staffLogic.StaffUpdate(uuid, req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags staff
// @Summary Delete Staff
// @Description Delete the staff member by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "Staff UUID"
// @Success 200 {object} dto.Response
// @Router /staff/{uuid} [delete]
func (u *BaseApi) StaffDel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if err := staffLogic.StaffDel(uuid); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags staff
// @Summary Get Staff
// @Description Get the staff information by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "Staff UUID"
// @Success 200 {object} dto.Response
// @Router /staff/{uuid} [get]
func (u *BaseApi) StaffGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	staff, err := staffLogic.StaffGet(uuid)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, staff)
}

// @Tags staff
// @Summary List Staff
// @Description Retrieve a list of all staff members
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response
// @Router /staff [get]
func (u *BaseApi) StaffList(ctx *gin.Context) {
	staffs, err := staffLogic.StaffList()
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, staffs)
}
