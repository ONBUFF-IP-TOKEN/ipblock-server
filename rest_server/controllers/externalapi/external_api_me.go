package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi/commonapi_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/labstack/echo"
)

func (o *ExternalAPI) GetAucBidListMe(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewMeBidList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucBidListMe(params, ctx)
}

func (o *ExternalAPI) GetAucBidWinnerListMe(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewMeBidList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucBidWinnerListMe(params, ctx)
}

func (o *ExternalAPI) GetAucNftListMe(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewMeBidList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucBidListMe(params, ctx)
}
