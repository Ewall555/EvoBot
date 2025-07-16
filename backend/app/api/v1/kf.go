package v1

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags kf
// @Summary Add kf
// @Description Add a new kf
// @Accept json
// @Produce json
// @Param req body request.KF true "kf Information"
// @Success 200 {object} dto.Response
// @Router /kf [post]
func (u *BaseApi) KFAdd(ctx *gin.Context) {
	var req request.KF
	var kfurl string
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if req.KFID == "" || req.KFName == "" || req.KFPlatform == "" {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	switch req.KFPlatform {
	case "wecom":
		url, err := wecomLogic.AddContactWay(req.KFID)
		if err != nil {
			helper.ErrResponseWithErr(ctx, constant.CodeErrBadRequest, err)
			return
		}
		kfurl = url
		if req.StaffList != nil {
			receptionistlist, err := wecomLogic.ReceptionistList(req.KFID)
			if err != nil {
				helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
				return
			}
			if err := wecomLogic.CheckReceptionist(req.StaffList, receptionistlist); err != nil {
				helper.ErrResponseWithErr(ctx, constant.CodeErrBadRequest, err)
				return
			}
		}
	default:
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	err := kFLogic.KFAdd(req)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, kfurl)
}

// @Tags kf
// @Summary Update kf
// @Description Update the kf by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "kf UUID"
// @Param req body request.KF true "kf Information"
// @Success 200 {object} dto.Response
// @Router /kf/{uuid} [patch]
func (u *BaseApi) KFUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	var req request.KF
	var kfurl string
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if req.KFID != "" {
		switch req.KFPlatform {
		case "wecom":
			url, err := wecomLogic.AddContactWay(req.KFID)
			if err != nil {
				helper.ErrResponseWithErr(ctx, constant.CodeErrBadRequest, err)
				return
			}
			kfurl = url
			if req.StaffList != nil {
				receptionistlist, err := wecomLogic.ReceptionistList(req.KFID)
				if err != nil {
					helper.ErrResponse(ctx, constant.CodeErrInternalServer)
					return
				}
				if err := wecomLogic.CheckReceptionist(req.StaffList, receptionistlist); err != nil {
					helper.ErrResponseWithErr(ctx, constant.CodeErrBadRequest, err)
					return
				}
			}
		default:
			helper.ErrResponse(ctx, constant.CodeErrInternalServer)
			return
		}
	}
	if err := kFLogic.KFUpdate(uuid, req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, kfurl)
}

// @Tags kf
// @Summary Delete KF
// @Description Delete the kf by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "kf UUID"
// @Success 200 {object} dto.Response
// @Router /kf/{uuid} [delete]
func (u *BaseApi) KFDel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if err := kFLogic.KFDel(uuid); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags kf
// @Summary Get KF
// @Description Get the kf information by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "kf UUID"
// @Success 200 {object} dto.Response
// @Router /kf/{uuid} [get]
func (u *BaseApi) KFGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	kf, err := kFLogic.KFGet(uuid)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, kf)
}

// @Tags kf
// @Summary List KF
// @Description Retrieve a list of all kf
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response
// @Router /kf [get]
func (u *BaseApi) KFList(ctx *gin.Context) {
	kfs, err := kFLogic.KFList()
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, kfs)
}
