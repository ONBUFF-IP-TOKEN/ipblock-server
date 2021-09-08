package internalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi/commonapi_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/labstack/echo"
)

// 입찰 정보 삭제 (보증금 납부부터 다시 해야함)
func (o *InternalAPI) DeleteAucBidRemove(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidRemove()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.DeleteAucBidRemove(params, ctx)
}

// 입찰 보증금 반환 리스트
func (o *InternalAPI) GetAucBidDepositRefund(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidDepositRefundList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucBidDepositRefund(params, ctx)
}

// 입찰 하기
func (o *InternalAPI) PostAucBidSubmit(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidSubmit()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi_auc.PostAucBidSubmitDummy(params, ctx)
}
