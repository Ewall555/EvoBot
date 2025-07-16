package v1

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/app/dto"
	"EvoBot/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags llmapp
// @Summary Get llmapp configuration by UUID
// @Description Retrieve llmapp configuration by its UUID
// @Accept json
// @Produce json
// @Param uuid path string true "UUID of the llmapp configuration"
// @Router /llmapp/config/{uuid} [get]
func (u *BaseApi) LLMAppConfigGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	llmapp, err := llmappLogic.ConfigGet(uuid)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, llmapp)
}

// @Tags llmapp
// @Summary Add new llmapp configuration
// @Description Add a new llmapp configuration
// @Accept json
// @Produce json
// @Param data body dto.LLMAppConfig true "llmapp configuration data"
// @Success 200 {object} dto.Response "Success response"
// @Router /llmapp/config [post]
func (u *BaseApi) LLMAppConfigAdd(ctx *gin.Context) {
	var req dto.LLMAppConfig
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	err := llmappLogic.ConfigAdd(req)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags llmapp
// @Summary Delete llmapp configuration by UUID
// @Description Delete llmapp configuration by its UUID
// @Accept json
// @Produce json
// @Param uuid path string true "UUID of the llmapp configuration"
// @Success 200 {object} dto.Response "Success response"
// @Router /llmapp/config/{uuid} [delete]

func (u *BaseApi) LLMAppConfigDel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if err := llmappLogic.ConfigDel(uuid); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags llmapp_config
// @Summary Update llmapp configuration
// @Description Update the llmapp configuration with the provided data
// @Accept json
// @Produce json
// @Param uuid path string true "UUID of the llmapp configuration"
// @Param data body dto.LLMAppConfig true "llmapp configuration data"
// @Success 200 {object} dto.Response "Success response"
// @Router /llmapp/config/{uuid} [patch]
func (b *BaseApi) LLMAppConfigUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	var req dto.LLMAppConfig
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := llmappLogic.ConfigUpdate(uuid, req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags llmapp_config
// @Summary List llmapp configuration
// @Description Retrieve a list of all llmapp configuration
// @Accept json
// @Produce json
// @Success 200 {object} dto.LLMAppConfig
// @Router /llmapp/config [get]
func (b *BaseApi) LLMAppConfigList(ctx *gin.Context) {
	conf, err := llmappLogic.ConfigList()
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, conf)
}
