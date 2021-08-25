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

// 경매 등록
func (o *InternalAPI) PostAucAuctionRegister(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewAucAuctionRegister()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostAucAuctionRegister(params, ctx)
}

// 경매 업데이트
func (o *InternalAPI) PostAucAuctionUpdate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewAucAuctionUpdate()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostAucAuctionUpdate(params, ctx)
}

// 경매 리스트 요청
func (o *InternalAPI) GetAucAuctionList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewAuctionList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucAuctionList(params, c)
}

// 경매 리스트 요청 (경매 상태에 따른)
func (o *InternalAPI) GetAucAuctionListByAucState(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewAuctionListByAucState()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	params.ActiveState = context_auc.Auction_active_state_all
	return commonapi_auc.GetAucAuctionListByAucState(params, c)
}

// 경매 삭제
func (o *InternalAPI) DeleteAucAuctiontRemove(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewRemoveAuction()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.DeleteAucAuctiontRemove(params, ctx)
}

// 경매 종료
func (o *InternalAPI) PostAucAuctionFinish(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewAuctionFinish()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostAucAuctionFinish(params, ctx)
}
