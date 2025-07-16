package v1

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags wecom_config
// @Summary List Wecom configuration
// @Description Retrieve the current Wecom configuration
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "Current Wecom configuration"
// @Router /wecom/config [get]
func (b *BaseApi) WecomConfigList(ctx *gin.Context) {
	conf, err := wecomLogic.ConfigList()
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, conf)
}

// @Tags wecom_config
// @Summary Get Wecom configuration
// @Description Retrieve the current Wecom configuration
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "Current Wecom configuration"
// @Router /wecom/config/{uuid} [get]
func (b *BaseApi) WecomConfigGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	conf, err := wecomLogic.ConfigGet(uuid)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, conf)
}

// @Tags wecom_config
// @Summary Update Wecom configuration
// @Description Update the Wecom configuration with the provided data
// @Accept json
// @Produce json
// @Param data body request.WecomConfigApp true "Wecom configuration data"
// @Success 200 {object} dto.Response "Success response"
// @Router /wecom/config/{uuid} [patch]
func (b *BaseApi) WecomConfigUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	var req request.WecomConfigApp
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := wecomLogic.ConfigUpdate(uuid, req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags wecom_callback
// @Summary WeCom Callback Message Handling
// @Description Handles callback messages from WeCom (Enterprise WeChat). This endpoint receives and parses the callback data sent by WeCom.
// @Accept json
// @Produce json
// @Success 200 {object} string "Parsed WeCom callback message"
// @Router /wecom/callback [get,post]
func (b *BaseApi) WecomHandle(ctx *gin.Context) {
	// https://developer.work.weixin.qq.com/document/path/90930#32-%E6%94%AF%E6%8C%81http-post%E8%AF%B7%E6%B1%82%E6%8E%A5%E6%94%B6%E4%B8%9A%E5%8A%A1%E6%95%B0%E6%8D%AE
	// 正确响应本次请求
	// 企业微信服务器在五秒内收不到响应会断掉连接，并且重新发起请求，总共重试三次。
	// 仅针对网络连接失败或者网络请求超时情况重试，建议开发者接受回调后立即应答，业务异步处理。
	if ctx.Request.Method == http.MethodGet {
		var options dto.SignatureOptions
		if err := ctx.ShouldBindQuery(&options); err != nil {
			helper.ErrResponse(ctx, constant.CodeErrBadRequest)
			return
		}
		echo, err := wecomLogic.VerifyURL(options)
		if err != nil {
			helper.ErrResponse(ctx, constant.CodeErrInternalServer)
			return
		}
		ctx.String(constant.CodeSuccess, echo)
		return
	}
	body, err := ctx.GetRawData()
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	wecomLogic.CallbackHandler(body)
	helper.SuccessWithOutData(ctx)
}

// @Tags wecom_receptionist
// @Summary WeCom receptionist add
// @Description receptionist add
// @Accept json
// @Produce json
// @Param data body request.ReceptionistOptions true "Wecom configuration data"
// @Success 200
// @Router /wecom/receptionist/{kfid} [post]
func (b *BaseApi) WecomReceptionistAdd(ctx *gin.Context) {
	kfid := ctx.Param("kfid")
	var options request.ReceptionistOptions
	if err := ctx.ShouldBindJSON(&options); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := wecomLogic.ReceptionistAdd(kfid, options); err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags wecom_receptionist
// @Summary WeCom receptionist delete
// @Description receptionist delete
// @Accept json
// @Produce json
// @Param data body request.ReceptionistOptions true "Wecom configuration data"
// @Success 200
// @Router /wecom/receptionist/{kfid} [delete]
func (b *BaseApi) WecomReceptionistDel(ctx *gin.Context) {
	kfid := ctx.Param("kfid")
	var options request.ReceptionistOptions
	if err := ctx.ShouldBindJSON(&options); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := wecomLogic.ReceptionistDel(kfid, options); err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags wecom_receptionist
// @Summary WeCom receptionist list
// @Description receptionist list
// @Accept json
// @Produce json
// @Success 200 {object} []client.ReceptionistList
// @Router /wecom/receptionist/{kfid} [get]
func (b *BaseApi) WecomReceptionistList(ctx *gin.Context) {
	kfid := ctx.Param("kfid")
	list, err := wecomLogic.ReceptionistList(kfid)
	if err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithData(ctx, list)
}

// @Tags wecom_account
// @Summary WeCom account list
// @Description account list
// @Accept json
// @Produce json
// @Success 200 {object} []client.AccountInfoSchema
// @Router /wecom/account [get]
func (b *BaseApi) WecomAccountList(ctx *gin.Context) {
	list, err := wecomLogic.AccountList()
	if err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithData(ctx, list)
}

// @Tags wecom_account
// @Summary WeCom account get
// @Description account get
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /wecom/account/{kfid} [get]
func (b *BaseApi) WecomAddContactWay(ctx *gin.Context) {
	kfid := ctx.Param("kfid")
	url, err := wecomLogic.AddContactWay(kfid)
	if err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithData(ctx, url)
}
