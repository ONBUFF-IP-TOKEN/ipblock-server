package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi/commonapi_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/labstack/echo"
)

// 경매 리스트 요청
func (o *ExternalAPI) GetAucAuctionList(c echo.Context) error {
	params := context_auc.NewAuctionList()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	params.ActiveState = context_auc.Auction_active_state_active
	return commonapi_auc.GetAucAuctionList(params, c)
}

// 경매 정보 요청
func (o *ExternalAPI) GetAucAuction(c echo.Context) error {
	params := context_auc.NewGetAuction()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucAuction(params, c)
}
